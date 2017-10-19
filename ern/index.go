// This file is part of the go-meta library.
//
// Copyright (C) 2017 JAAK MUSIC LTD
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// If you have any questions please contact yo@jaak.io

package ern

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
	"github.com/lmars/go-ipld-cbor"
	"github.com/mattn/go-sqlite3"
	"github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/xml"
)

// Indexer is a META indexer which indexes a stream of META objects
// representing DDEX ERNs into a SQLite3 database, getting the
// associated META objects from a META store.
type Indexer struct {
	index *meta.Index
	store *meta.Store
}

// NewIndexer returns an Indexer which updates the indexes in the given SQLite3
// database connection, getting META objects from the given META store.
func NewIndexer(index *meta.Index, store *meta.Store) (*Indexer, error) {
	// migrate the db to ensure it has an up-to-date schema
	if err := migrations.Run(index.DB); err != nil {
		return nil, err
	}

	return &Indexer{
		index: index,
		store: store,
	}, nil
}

// Index indexes a stream of META object links which are expected to
// point at DDEX ERNs.
func (i *Indexer) Index(ctx context.Context, stream *meta.StreamReader) error {
	return i.index.Update(func(tx *sql.Tx) error {
		for {
			select {
			case cid, ok := <-stream.Ch():
				if !ok {
					return stream.Err()
				}
				obj, err := i.store.Get(cid)
				if err != nil {
					return err
				}
				if err := i.indexERN(tx, obj); err != nil {
					return err
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
}

type indexContext struct {
	tx   *sql.Tx
	ern  *meta.Object
	refs map[string]*cid.Cid
}

// indexERN indexes a DDEX ERN based on its MessageHeader, WorkList,
// ResourceList and ReleaseList.
func (i *Indexer) indexERN(tx *sql.Tx, ern *meta.Object) error {
	graph := meta.NewGraph(i.store, ern)
	ctx := &indexContext{
		tx:   tx,
		ern:  ern,
		refs: make(map[string]*cid.Cid),
	}

	// index the properties in order to ensure we index ResourceReferences
	// in the correct order (i.e. get the SoundRecording refs before
	// indexing the Release)
	type indexTask struct {
		field   string
		indexFn func(*indexContext, *meta.Object) error
	}
	for _, t := range []indexTask{
		{"MessageHeader", i.indexMessageHeader},
		{"WorkList", i.indexWorkList},
		{"ResourceList", i.indexResourceList},
		{"ReleaseList", i.indexReleaseList},
	} {
		v, err := graph.Get("NewReleaseMessage", t.field)
		if meta.IsPathNotFound(err) {
			continue
		} else if err != nil {
			return err
		}
		id, ok := v.(*cid.Cid)
		if !ok {
			return fmt.Errorf("unexpected field type for %q, expected *cid.Cid, got %T", t.field, v)
		}
		obj, err := i.store.Get(id)
		if err != nil {
			return err
		}
		if err := t.indexFn(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

// isUniqueErr determines whether an error is a SQLite3 uniqueness error.
func isUniqueErr(err error) bool {
	e, ok := err.(sqlite3.Error)
	if !ok {
		return false
	}
	return e.Code == sqlite3.ErrConstraint && e.ExtendedCode == sqlite3.ErrConstraintUnique
}

// DecodeObj decodes whatever is stored at path into the given value
func DecodeObj(store *meta.Store, obj *meta.Object, path ...string) (v *metaxml.Value, err error) {
	graph := meta.NewGraph(store, obj)

	defer func() {
		if err != nil {
			err = fmt.Errorf("Error decoding %s as metaxml.Value: %s", path, err)
		}
	}()

	x, err := graph.Get(path...)
	if err != nil {
		return nil, err
	}
	id, ok := x.(*cid.Cid)
	if !ok {
		return nil, fmt.Errorf("Expected %s to be *cid.Cid, got %T", path, x)
	}

	obj, err = store.Get(id)
	if err != nil {
		return nil, err
	}
	v = &metaxml.Value{}
	return v, obj.Decode(v)
}

// insertParty inserts the PartyName and PartyId fields of a Party object into
// the party index.
func (i *Indexer) insertParty(tx *sql.Tx, obj *meta.Object) error {
	// explicitly ignore the returned error as it is ok for the PartyId to
	// be missing
	var partyID metaxml.Value
	if v, err := DecodeObj(i.store, obj, "PartyId"); err == nil {
		partyID = *v
	}

	partyName, err := DecodeObj(i.store, obj, "PartyName", "FullName")
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO party (cid, id, name) VALUES ($1, $2, $3)",
		obj.Cid().String(), partyID.Value, partyName.Value,
	)
	if err != nil && !isUniqueErr(err) {
		return err
	}
	return nil
}

// insertParties loads parties from the given field and inserts them into the
// party index.
func (i *Indexer) insertParties(tx *sql.Tx, obj *meta.Object, field string) ([]*cid.Cid, error) {
	ids, err := decodeLinks(i.store, obj, field)
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		party, err := i.store.Get(id)
		if err != nil {
			return nil, err
		}
		if err := i.insertParty(tx, party); err != nil {
			return nil, err
		}
	}
	return ids, nil
}

// indexMessageHeader indexes an ERN MessageHeader based on its MessageId,
// MessageThreadId, MessageSender, MessageRecipient and MessageCreatedDateTime.
func (i *Indexer) indexMessageHeader(ctx *indexContext, obj *meta.Object) error {
	senders, err := i.insertParties(ctx.tx, obj, "MessageSender")
	if err != nil {
		return err
	}
	if len(senders) != 1 {
		return fmt.Errorf("expected 1 sender, got %d", len(senders))
	}
	sender := senders[0]

	recipients, err := i.insertParties(ctx.tx, obj, "MessageRecipient")
	if err != nil {
		return err
	}
	if len(recipients) != 1 {
		return fmt.Errorf("expected 1 recipient, got %d", len(recipients))
	}
	recipient := recipients[0]

	// get the MessageId, MessageThreadId and MessageCreatedDateTime
	// values
	messageID, err := DecodeObj(i.store, obj, "MessageId")
	if err != nil {
		return err
	}
	threadID, err := DecodeObj(i.store, obj, "MessageThreadId")
	if err != nil {
		return err
	}
	created, err := DecodeObj(i.store, obj, "MessageCreatedDateTime")
	if err != nil {
		return err
	}

	// update the ERN index
	_, err = ctx.tx.Exec(
		"INSERT INTO ern (cid, message_id, thread_id, sender_id, recipient_id, created) VALUES ($1, $2, $3, $4, $5, $6)",
		ctx.ern.Cid().String(), messageID.Value, threadID.Value, sender.String(), recipient.String(), created.Value,
	)
	return err
}

func (i *Indexer) indexWorkList(ctx *indexContext, obj *meta.Object) error {
	// TODO: index MusicalWorks
	return nil
}

// indexResourceList indexes an ERN ResourceList based on SoundRecordings.
func (i *Indexer) indexResourceList(ctx *indexContext, obj *meta.Object) error {
	// the SoundRecording property can either be a link if there is only
	// one SoundRecording in the list, or an array of links if there are
	// multiple SoundRecordings in the list, so we need to handle both
	// cases
	cids, err := decodeLinks(i.store, obj, "SoundRecording")
	if err != nil {
		return err
	}

	// load and index each SoundRecording link
	for _, cid := range cids {
		obj, err := i.store.Get(cid)
		if err != nil {
			return err
		}
		if err := i.indexSoundRecording(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

func decodeLinks(store *meta.Store, obj *meta.Object, field string) (cids []*cid.Cid, err error) {
	v, err := obj.Get(field)
	if err != nil {
		if strings.HasSuffix(err.Error(), cbornode.ErrNoSuchLink.Error()) {
			err = nil
		}
		return nil, err
	}
	switch v := v.(type) {
	case *format.Link:
		cids = []*cid.Cid{v.Cid}
	case []interface{}:
		for _, x := range v {
			cid, ok := x.(*cid.Cid)
			if !ok {
				return nil, fmt.Errorf("invalid resource type %T, expected *cid.Cid", x)
			}
			cids = append(cids, cid)
		}
	}
	return
}

// indexSoundRecording indexes an ERN SoundRecording based on its ID (either an
// ISRC, CatalogNumber or ProprietaryId) and its ReferenceTitle.
func (i *Indexer) indexSoundRecording(ctx *indexContext, obj *meta.Object) error {
	// Only *attempt* to load the ISRC, other IDs can be retrieved via GraphQL
	// Default to empty string if not present
	var isrc metaxml.Value
	if v, err := DecodeObj(i.store, obj, "SoundRecordingId", "ISRC"); err == nil {
		isrc = *v
	}

	// Insert the DisplayArtist to party table
	cids, err := decodeLinks(i.store, obj, "SoundRecordingDetailsByTerritory")
	if err != nil {
		return err
	}
	for _, cid := range cids {
		detailsObj, err := i.store.Get(cid)
		if err != nil {
			return err
		}
		artistIDs, err := i.insertParties(ctx.tx, detailsObj, "DisplayArtist")
		if err != nil {
			return err
		}
		for _, artistID := range artistIDs {
			if _, err := ctx.tx.Exec(
				"INSERT INTO sound_recording_contributor (sound_recording_id, party_id) VALUES($1, $2)",
				obj.Cid().String(), artistID.String(),
			); err != nil && !isUniqueErr(err) {
				return err
			}
		}
	}

	// load the ReferenceTitle
	var title metaxml.Value
	if v, err := DecodeObj(i.store, obj, "ReferenceTitle", "TitleText"); err == nil {
		title = *v
	}

	// return an error if there is no ReferenceTitle, SoundRecordingId can be empty
	if title.Value == "" {
		return fmt.Errorf("SoundRecording missing ReferenceTitle")
	}

	// add the ResourceReference to ctx.refs
	if ref, err := DecodeObj(i.store, obj, "ResourceReference"); err == nil {
		ctx.refs[ref.Value] = obj.Cid()
	}

	// update the sound_recording and resource_list indexes
	if _, err := ctx.tx.Exec(
		"INSERT INTO sound_recording (cid, id, title) VALUES ($1, $2, $3)",
		obj.Cid().String(), isrc.Value, title.Value,
	); err != nil {
		return err
	}

	if _, err := ctx.tx.Exec(
		"INSERT INTO resource_list (ern_id, resource_id) VALUES ($1, $2)",
		ctx.ern.Cid().String(), obj.Cid().String(),
	); err != nil {
		return err
	}

	return nil
}

// indexReleaseList indexes the ReleaseList for each Release composite
func (i *Indexer) indexReleaseList(ctx *indexContext, metaObj *meta.Object) error {
	// Much like the resource list, the release propoerty can be
	// a single release, or an array of links.
	cids, err := decodeLinks(i.store, metaObj, "Release")
	if err != nil {
		return err
	}

	// load and index each Release link
	for _, id := range cids {
		rObj, err := i.store.Get(id)
		if err != nil {
			return err
		}
		if err := i.indexRelease(ctx, rObj); err != nil {
			return err
		}
	}
	return nil
}

// indexRelease index each Release composite in the ReelaseList
func (i *Indexer) indexRelease(ctx *indexContext, obj *meta.Object) error {
	// Only *attempt* to load the GRid, other IDs can be retrieved via GraphQL
	var grid metaxml.Value
	if v, err := DecodeObj(i.store, obj, "ReleaseId", "GRid"); err == nil {
		grid = *v
	}

	var icpn metaxml.Value
	if v, err := DecodeObj(i.store, obj, "ReleaseId", "ICPN"); err == nil {
		icpn = *v
	}

	// Insert the DisplayArtist to party table
	cids, err := decodeLinks(i.store, obj, "ReleaseDetailsByTerritory")
	if err != nil {
		return err
	}
	for _, cid := range cids {
		obj, err := i.store.Get(cid)
		if err != nil {
			return err
		}
		_, err = i.insertParties(ctx.tx, obj, "DisplayArtist")
		if err != nil {
			return err
		}
	}

	// load the ReferenceTitle
	var title metaxml.Value
	if v, err := DecodeObj(i.store, obj, "ReferenceTitle", "TitleText"); err == nil {
		title = *v
	}

	// return an error if there is no ReferenceTitle, ReleaseId can be empty
	if title.Value == "" {
		return fmt.Errorf("Release missing ReferenceTitle")
	}

	// update the sound_recording_release index
	if listLink, err := obj.GetLink("ReleaseResourceReferenceList"); err == nil {
		listObj, err := i.store.Get(listLink.Cid)
		if err != nil {
			return err
		}
		refIDs, err := decodeLinks(i.store, listObj, "ReleaseResourceReference")
		if err != nil {
			return err
		}
		for _, refID := range refIDs {
			refObj, err := i.store.Get(refID)
			if err != nil {
				return err
			}
			var ref metaxml.Value
			if err := refObj.Decode(&ref); err != nil {
				return err
			}
			soundRecording, ok := ctx.refs[ref.Value]
			if !ok {
				continue
			}
			_, err = ctx.tx.Exec(
				"INSERT INTO sound_recording_release (sound_recording_id, release_id) VALUES ($1, $2)",
				soundRecording.String(), obj.Cid().String(),
			)
			if err != nil && !isUniqueErr(err) {
				return err
			}
		}
	}

	// update the release and release_list indexes
	for _, id := range []string{grid.Value, icpn.Value} {
		_, err = ctx.tx.Exec(
			"INSERT INTO release (cid, id, title) VALUES ($1, $2, $3)",
			obj.Cid().String(), id, title.Value,
		)
		if err != nil {
			return err
		}
	}

	_, err = ctx.tx.Exec(
		"INSERT INTO release_list (ern_id, release_id) VALUES ($1, $2)",
		ctx.ern.Cid().String(), obj.Cid().String(),
	)
	return err

}
