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

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
	"github.com/mattn/go-sqlite3"
	"github.com/meta-network/go-meta"
)

// Indexer is a META indexer which indexes a stream of META objects
// representing DDEX ERNs into a SQLite3 database, getting the
// associated META objects from a META store.
type Indexer struct {
	db    *sql.DB
	store *meta.Store
}

// NewIndexer returns an Indexer which updates the indexes in the given SQLite3
// database connection, getting META objects from the given META store.
func NewIndexer(db *sql.DB, store *meta.Store) (*Indexer, error) {
	// migrate the db to ensure it has an up-to-date schema
	if err := migrations.Run(db); err != nil {
		return nil, err
	}

	return &Indexer{
		db:    db,
		store: store,
	}, nil
}

// Index indexes a stream of META object links which are expected to
// point at DDEX ERNs.
func (i *Indexer) Index(ctx context.Context, stream chan *cid.Cid) error {
	for {
		select {
		case cid, ok := <-stream:
			if !ok {
				return nil
			}
			obj, err := i.store.Get(cid)
			if err != nil {
				return err
			}
			if err := i.index(obj); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// index indexes a DDEX ERN based on its MessageHeader, WorkList, ResourceList
// and ReleaseList.
func (i *Indexer) index(ern *meta.Object) error {
	graph := meta.NewGraph(i.store, ern)

	for field, indexFn := range map[string]func(*cid.Cid, *meta.Object) error{
		"MessageHeader": i.indexMessageHeader,
		"WorkList":      i.indexWorkList,
		"ResourceList":  i.indexResourceList,
		"ReleaseList":   i.indexReleaseList,
	} {
		v, err := graph.Get("NewReleaseMessage", field)
		if meta.IsPathNotFound(err) {
			continue
		} else if err != nil {
			return err
		}
		id, ok := v.(*cid.Cid)
		if !ok {
			return fmt.Errorf("unexpected field type for %q, expected *cid.Cid, got %T", field, v)
		}
		if err := i.indexProperty(ern.Cid(), id, indexFn); err != nil {
			return err
		}
	}

	return nil
}

// indexProperty indexes a particular ERN property using the provided index
// function.
func (i *Indexer) indexProperty(ernID, cid *cid.Cid, indexFn func(*cid.Cid, *meta.Object) error) error {
	obj, err := i.store.Get(cid)
	if err != nil {
		return err
	}
	return indexFn(ernID, obj)
}

// isUniqueErr determines whether an error is a SQLite3 uniqueness error.
func isUniqueErr(err error) bool {
	e, ok := err.(sqlite3.Error)
	if !ok {
		return false
	}
	return e.Code == sqlite3.ErrConstraint && e.ExtendedCode == sqlite3.ErrConstraintUnique
}

// indexMessageHeader indexes an ERN MessageHeader based on its MessageId,
// MessageThreadId, MessageSender, MessageRecipient and MessageCreatedDateTime.
func (i *Indexer) indexMessageHeader(ernID *cid.Cid, obj *meta.Object) error {
	graph := meta.NewGraph(i.store, obj)

	// decode decodes whatever is stored at path into the given value
	decode := func(v interface{}, path ...string) (err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("error decoding %s into %T: %s", path, v, err)
			}
		}()
		x, err := graph.Get(path...)
		if meta.IsPathNotFound(err) {
			return nil
		} else if err != nil {
			return err
		}
		id, ok := x.(*cid.Cid)
		if !ok {
			return fmt.Errorf("expected %s to be *cid.Cid, got %T", path, x)
		}
		obj, err := i.store.Get(id)
		if err != nil {
			return err
		}
		return obj.Decode(v)
	}

	// insert the MessageSender and MessageRecipient into the party index
	insertParty := func(field string) (*cid.Cid, error) {
		var id *cid.Cid
		link, err := obj.GetLink(field)
		if err == nil {
			id = link.Cid
		} else if !meta.IsPathNotFound(err) {
			return nil, err
		}
		var partyID struct {
			Value string `json:"@value"`
		}
		if err := decode(&partyID, field, "PartyId"); err != nil {
			return nil, err
		}
		var partyName struct {
			Value string `json:"@value"`
		}
		if err := decode(&partyName, field, "PartyName", "FullName"); err != nil {
			return nil, err
		}
		_, err = i.db.Exec(
			"INSERT INTO party (cid, id, name) VALUES ($1, $2, $3)",
			id.String(), partyID.Value, partyName.Value,
		)
		if err != nil && !isUniqueErr(err) {
			return nil, err
		}
		return id, nil
	}
	sender, err := insertParty("MessageSender")
	if err != nil {
		return err
	}
	recipient, err := insertParty("MessageRecipient")
	if err != nil {
		return err
	}

	// get the MessageId, MessageThreadId and MessageCreatedDateTime
	// values
	var messageID struct {
		Value string `json:"@value"`
	}
	if err := decode(&messageID, "MessageId"); err != nil {
		return err
	}
	var threadID struct {
		Value string `json:"@value"`
	}
	if err := decode(&threadID, "MessageThreadId"); err != nil {
		return err
	}
	var created struct {
		Value string `json:"@value"`
	}
	if err := decode(&created, "MessageCreatedDateTime"); err != nil {
		return err
	}

	// update the ERN index
	_, err = i.db.Exec(
		"INSERT INTO ern (cid, message_id, thread_id, sender_id, recipient_id, created) VALUES ($1, $2, $3, $4, $5, $6)",
		ernID.String(), messageID.Value, threadID.Value, sender.String(), recipient.String(), created.Value,
	)
	return err
}

func (i *Indexer) indexWorkList(ernID *cid.Cid, obj *meta.Object) error {
	// TODO: index MusicalWorks
	return nil
}

// indexResourceList indexes an ERN ResourceList based on SoundRecordings.
func (i *Indexer) indexResourceList(ernID *cid.Cid, obj *meta.Object) error {
	// the SoundRecording property can either be a link if there is only
	// one SoundRecording in the list, or an array of links if there are
	// multiple SoundRecordings in the list, so we need to handle both
	// cases
	v, err := obj.Get("SoundRecording")
	if err != nil {
		return err
	}
	var cids []*cid.Cid
	switch v := v.(type) {
	case *format.Link:
		cids = []*cid.Cid{v.Cid}
	case []interface{}:
		for _, x := range v {
			cid, ok := x.(*cid.Cid)
			if !ok {
				return fmt.Errorf("invalid resource type %T, expected *cid.Cid", x)
			}
			cids = append(cids, cid)
		}
	}

	// load and index each SoundRecording link
	for _, cid := range cids {
		obj, err := i.store.Get(cid)
		if err != nil {
			return err
		}
		if err := i.indexSoundRecording(ernID, obj); err != nil {
			return err
		}
	}

	return nil
}

// indexSoundRecording indexes an ERN SoundRecording based on its ID (either an
// ISRC, CatalogNumber or ProprietaryId) and its ReferenceTitle.
func (i *Indexer) indexSoundRecording(ernID *cid.Cid, obj *meta.Object) error {
	graph := meta.NewGraph(i.store, obj)

	// load each potential ID separately
	var ids []string
	for _, field := range []string{"ISRC", "CatalogNumber", "ProprietaryId"} {
		v, err := graph.Get("SoundRecordingId", field, "@value")
		if meta.IsPathNotFound(err) {
			continue
		} else if err != nil {
			return err
		}
		ids = append(ids, v.(string))
	}

	// load the ReferenceTitle
	var title string
	v, err := graph.Get("ReferenceTitle", "TitleText", "@value")
	if err == nil {
		title = v.(string)
	} else if !meta.IsPathNotFound(err) {
		return err
	}

	// return an error if there is neither an ID nor a ReferenceTitle
	if len(ids) == 0 && title == "" {
		return fmt.Errorf("SoundRecording missing both SoundRecordingId and ReferenceTitle")
	}

	// update the sound_recording and resource_list indexes with each ID
	for _, id := range ids {
		_, err := i.db.Exec(
			"INSERT INTO sound_recording (cid, id, title) VALUES ($1, $2, $3)",
			obj.Cid().String(), id, title,
		)
		if err != nil {
			return err
		}

		_, err = i.db.Exec(
			"INSERT INTO resource_list (ern_id, resource_id) VALUES ($1, $2)",
			ernID.String(), obj.Cid().String(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) indexReleaseList(ernID *cid.Cid, obj *meta.Object) error {
	// TODO: index Releases
	return nil
}
