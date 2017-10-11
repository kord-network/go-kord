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
	"database/sql"
	"errors"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-format"
	"github.com/meta-network/go-meta"
)

// Resolver is a general purpose GraphQL resolver function,
// retrieves data from a META store and SQLite3 index
type Resolver struct {
	db    *sql.DB
	store *meta.Store
}

// NewResolver returns a new resolver for returning data from the given
// META store and SQLite3 index
func NewResolver(db *sql.DB, store *meta.Store) *Resolver {
	return &Resolver{db, store}
}

// GraphQLSchema schema definition for Party on the ERN index
const GraphQLSchema = `
	schema {
		query: Query
	}

	type Query {
		party(id: String, name: String): [Party]!
		soundRecording(id: String, title: String): [SoundRecording]!
		release(id: String, title: String): [Release]!
	}

	type Party {
		cid: String!
		partyID: String!
		fullName: String!
	}

	# NEEDS ADDING TO QUERY RESPONSE
	type ResourceContributor {
		party: Party!
		role: String!
	}

	type SoundRecording {
		artistName: String!
		genre: String!
		parentalWarningType: String
		resourceReference: String!
		subGenre: String
		soundRecordingId: String!
		territoryCode: String!
		title: String!
	}

	type Release {
		artistName: String!
		displayTitle: String!
		releaseId: String
		genre: String
		releaseType: String
	}
`

// partyArgs query arguments for Party query
type partyArgs struct {
	Name *string
	ID   *string
}

/**
 *	Party
 */

// partyResolver defines grapQL resolver functions for the Party fields
type partyResolver struct {
	cid   string
	party *Party
}

func (pd *partyResolver) Cid() string {
	return pd.cid
}

func (pd *partyResolver) Fullname() string {
	return pd.party.PartyName
}

func (pd *partyResolver) PartyId() string {
	return pd.party.PartyId
}

// Party is the resolver function to retrieve
// the Party information from the SQLite index
func (g *Resolver) Party(args partyArgs) ([]*partyResolver, error) {

	var rows *sql.Rows
	var err error

	switch {
	case args.Name != nil:
		rows, err = g.db.Query("SELECT cid FROM party WHERE name = ?", *args.Name)
	case args.ID != nil:
		rows, err = g.db.Query("SELECT cid FROM party WHERE id = ?", *args.ID)
	default:
		return nil, errors.New("Missing Name or ID argument in query")
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var resolvers []*partyResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		id, err := cid.Parse(objectID)
		if err != nil {
			return nil, err
		}

		obj, err := g.store.Get(id)
		if err != nil {
			return nil, err
		}

		var DdexPartyId struct {
			Value string `json:"@value"`
		}
		if err := DecodeObj(g.store, obj, &DdexPartyId, "PartyId"); err != nil {
			pid := &DdexPartyId
			pid.Value = ""
		}

		var DdexPartyName struct {
			Value string `json:"@value"`
		}
		if err := DecodeObj(g.store, obj, &DdexPartyName, "PartyName", "FullName"); err != nil {
			return nil, err
		}
		// Not keen on the below, but refinement will take time :)
		party := Party{PartyId: DdexPartyId.Value, PartyName: DdexPartyName.Value}
		resolvers = append(resolvers, &partyResolver{objectID, &party})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

/**
 * SoundRecording
 */

type soundRecordingArgs struct {
	ID    *string
	Title *string
}

type soundRecordingResolver struct {
	cid            string
	soundRecording *SoundRecording
}

func (sr *soundRecordingResolver) Cid() string {
	return sr.cid
}

func (sr *soundRecordingResolver) ArtistName() string {
	return sr.soundRecording.ArtistName
}

func (sr *soundRecordingResolver) Genre() string {
	return sr.soundRecording.GenreText
}

func (sr *soundRecordingResolver) ParentalWarningType() *string {
	return &sr.soundRecording.ParentalWarningType
}

func (sr *soundRecordingResolver) ResourceReference() string {
	return sr.soundRecording.ResourceReference
}

func (sr *soundRecordingResolver) SoundRecordingId() string {
	return sr.soundRecording.SoundRecordingId
}

func (sr *soundRecordingResolver) SubGenre() *string {
	return &sr.soundRecording.SubGenre
}

func (sr *soundRecordingResolver) TerritoryCode() string {
	return sr.soundRecording.TerritoryCode
}

func (sr *soundRecordingResolver) Title() string {
	return sr.soundRecording.ReferenceTitle
}

func (g *Resolver) SoundRecording(args soundRecordingArgs) ([]*soundRecordingResolver, error) {
	var rows *sql.Rows
	var err error

	switch {
	case args.ID != nil:
		rows, err = g.db.Query("SELECT cid FROM sound_recording WHERE id = ?", *args.ID)
	case args.Title != nil:
		rows, err = g.db.Query("SELECT cid FROM sound_recording WHERE title = ?", *args.Title)
	default:
		return nil, errors.New("Missing ID or Title argument in query")
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var response []*soundRecordingResolver

	for rows.Next() {
		var objectID string

		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}

		id, err := cid.Parse(objectID)

		if err != nil {
			return nil, err
		}

		obj, err := g.store.Get(id)

		if err != nil {
			return nil, err
		}

		var ArtistName struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &ArtistName, "SoundRecordingDetailsByTerritory", "DisplayArtist", "PartyName", "FullName"); err != nil {
			return nil, err
		}

		var GenreText struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &GenreText, "SoundRecordingDetailsByTerritory", "Genre", "GenreText"); err != nil {
			return nil, err
		}

		var ParentalWarningType struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &ParentalWarningType, "SoundRecordingDetailsByTerritory", "ParentalWarningType"); err != nil {
			return nil, err
		}

		var ReferenceTitle struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &ReferenceTitle, "ReferenceTitle", "TitleText"); err != nil {
			return nil, err
		}

		var ResourceReference struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &ResourceReference, "ResourceReference"); err != nil {
			return nil, err
		}

		var SoundRecordingId struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &SoundRecordingId, "SoundRecordingId", "ISRC"); err != nil {
			return nil, err
		}

		var SubGenre struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &SubGenre, "SoundRecordingDetailsByTerritory", "Genre", "SubGenre"); err != nil {
			return nil, err
		}

		var TerritoryCode struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &TerritoryCode, "SoundRecordingDetailsByTerritory", "TerritoryCode"); err != nil {
			return nil, err
		}

		var soundRecording SoundRecording
		soundRecording.ArtistName = ArtistName.Value
		soundRecording.GenreText = GenreText.Value
		soundRecording.ParentalWarningType = ParentalWarningType.Value
		soundRecording.ReferenceTitle = ReferenceTitle.Value
		soundRecording.ResourceReference = ResourceReference.Value
		soundRecording.SoundRecordingId = SoundRecordingId.Value
		soundRecording.SubGenre = SubGenre.Value
		soundRecording.TerritoryCode = TerritoryCode.Value

		response = append(response, &soundRecordingResolver{objectID, &soundRecording})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return response, nil
}

/*
* Release
 */

type releaseArgs struct {
	ID    *string
	Title *string
}

type releaseResolver struct {
	cid     string
	release *Release
}

func (rl *releaseResolver) Cid() string {
	return rl.cid
}

func (rl *releaseResolver) ArtistName() string {
	return rl.release.ArtistName
}

func (rl *releaseResolver) DisplayTitle() string {
	return rl.release.DisplayTitle
}

func (rl *releaseResolver) Genre() *string {
	return &rl.release.Genre
}

func (rl *releaseResolver) ReleaseType() *string {
	return &rl.release.ReleaseType
}

func (rl *releaseResolver) ReleaseId() *string {
	return &rl.release.ReleaseId
}

func (g *Resolver) Release(args releaseArgs) ([]*releaseResolver, error) {
	var rows *sql.Rows
	var err error

	switch {
	case args.ID != nil:
		rows, err = g.db.Query("SELECT cid FROM release WHERE id = ?", *args.ID)
	case args.Title != nil:
		rows, err = g.db.Query("SELECT cid FROM release WHERE title = ?", *args.Title)
	default:
		return nil, errors.New("Missing ID or Title argument in query")
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var response []*releaseResolver

	for rows.Next() {
		var objectID string

		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}

		id, err := cid.Parse(objectID)

		if err != nil {
			return nil, err
		}

		obj, err := g.store.Get(id)

		if err != nil {
			return nil, err
		}

		var ArtistName struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &ArtistName, "ReleaseDetailsByTerritory", "DisplayArtist", "PartyName", "FullName"); err != nil {
			return nil, err
		}

		// Title field can be either a single object, or multiple
		// with varying "TitleType"s
		// So the parent object should be loaded and then traversed accordingly
		graph := meta.NewGraph(g.store, obj)
		titles, err := graph.Get("ReleaseDetailsByTerritory", "Title")
		if err != nil {
			return nil, err
		}
		var cids []*cid.Cid
		switch titles := titles.(type) {
		case *format.Link:
			cids = []*cid.Cid{titles.Cid}
		case []interface{}:
			for _, x := range titles {
				cid, ok := x.(*cid.Cid)
				if !ok {
					return nil, fmt.Errorf("invalid resource type %T, expected *cid.Cid", x)
				}
				cids = append(cids, cid)
			}
		}

		var DisplayTitle struct {
			Value string `json:"@value"`
		}
		// Load each title CID and check for the TitleType
		for _, cid := range cids {
			obj, err := g.store.Get(cid)
			if err != nil {
				return nil, err
			}
			tt, err := obj.Get("TitleType")
			if err != nil {
				return nil, err
			}

			if tt.(string) == "DisplayTitle" {
				if err := DecodeObj(g.store, obj, &DisplayTitle, "TitleText"); err != nil {
					return nil, err
				}
			}
		}

		var Genre struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &Genre, "ReleaseDetailsByTerritory", "Genre", "GenreText"); err != nil {
			return nil, err
		}

		var ReleaseType struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &ReleaseType, "ReleaseType"); err != nil {
			return nil, err
		}

		var ReleaseId struct {
			Value string `json:"@value"`
		}

		if err := DecodeObj(g.store, obj, &ReleaseId, "ReleaseId", "ISRC"); err != nil {
			return nil, err
		}

		var release Release
		release.ArtistName = ArtistName.Value
		release.DisplayTitle = DisplayTitle.Value
		release.Genre = Genre.Value
		release.ReleaseType = ReleaseType.Value
		release.ReleaseId = ReleaseId.Value

		response = append(response, &releaseResolver{objectID, &release})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return response, nil
}
