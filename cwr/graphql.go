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

package cwr

import (
	"database/sql"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// GraphQLSchema is the GraphQL schema for the MusicBrainz META index.
const GraphQLSchema = `
schema {
  query: Query
}

type Query {
  registered_work(
  title: String,
  iswc:  String,
  composite_type: String,
  record_type: String
  ): [RegisteredWork]!

  publisher_control(
  publisher_sequence_n: String,
  record_type: String
  ): [PublisherControl]!
}


type PublisherControl {
	cid:                      String!
	record_type:              String!
	publisher_sequence_n:     String!
}

type RegisteredWork {
	cid:                      String!
	record_type:              String!
	title:                    String!
	language_code:            String!
	submitte_worknumber:      String!
	iswc:                     String!
	copyright_date:           String!
	distribution_category:    String!
	duration:                 String!
	recorded_indicator:       String!
	textmusic_relationship:   String!
	composite_type:           String!
	version_type:             String!
	music_arrangement:        String!
	lyric_adaptation:         String!
	contact_name:             String!
	contact_id:               String!
	work_type:                String!
	grandrights_indicator:    String!
	composite_componentcount: String!
	date_ofpublication:       String!
	exceptional_clause:       String!
	opus_number:              String!
	catalogue_number:         String!
	priority_flag:            String!
}
`

// Resolver defines GraphQL resolver functions for the schema contained in
// the GraphQLSchema constant, retrieving data from a META store and SQLite3
// index.
type Resolver struct {
	db    *sql.DB
	store *meta.Store
}

// NewResolver returns a Resolver which retrieves data from the given META
// store and SQLite3 index.
func NewResolver(db *sql.DB, store *meta.Store) *Resolver {
	return &Resolver{db, store}
}

// registeredWorkArgs are the arguments for a GraphQL registeredWork query.
type registeredWorkArgs struct {
	RecordType    *string
	Title         *string
	ISWC          *string
	CompositeType *string
}

type publisherControlArgs struct {
	RecordType         *string
	PublisherSequenceN *string
}

// RegisteredWork is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either an RegisteredWork RecordType, Title ,ISWC,or CompositeType, and loads the
// associated META objects from the META store.
func (g *Resolver) RegisteredWork(args registeredWorkArgs) ([]*registeredWorkResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.Title != nil:
		rows, err = g.db.Query("SELECT object_id FROM registered_work WHERE title = ?", *args.Title)
	case args.RecordType != nil:
		rows, err = g.db.Query("SELECT object_id FROM registered_work WHERE record_type = ?", *args.RecordType)
	case args.CompositeType != nil:
		rows, err = g.db.Query("SELECT object_id FROM registered_work WHERE composite_type = ?", *args.CompositeType)
	case args.ISWC != nil:
		rows, err = g.db.Query("SELECT object_id FROM registered_work WHERE iswc = ?", *args.ISWC)
	default:
		return nil, errors.New("missing title, record_type ,iswc or composite_type argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*registeredWorkResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		cid, err := cid.Parse(objectID)
		if err != nil {
			return nil, err
		}
		obj, err := g.store.Get(cid)
		if err != nil {
			return nil, err
		}
		var registeredWork RegisteredWork
		if err := obj.Decode(&registeredWork); err != nil {
			return nil, err
		}
		resolvers = append(resolvers, &registeredWorkResolver{objectID, &registeredWork})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// PublisherControl is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either an PublihserControl RecordType or publisher_sequence_n and loads the
// associated META objects from the META store.
func (g *Resolver) PublisherControl(args publisherControlArgs) ([]*publisherControlResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.PublisherSequenceN != nil:
		rows, err = g.db.Query("SELECT object_id FROM publisher_control WHERE publisher_sequence_n = ?", *args.PublisherSequenceN)
	case args.RecordType != nil:
		rows, err = g.db.Query("SELECT object_id FROM publisher_control WHERE record_type = ?", *args.RecordType)
	default:
		return nil, errors.New("missing record_type or PublisherSequenceN argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*publisherControlResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		cid, err := cid.Parse(objectID)
		if err != nil {
			return nil, err
		}
		obj, err := g.store.Get(cid)
		if err != nil {
			return nil, err
		}
		var publisherControllBySubmitter PublisherControllBySubmitter
		if err := obj.Decode(&publisherControllBySubmitter); err != nil {
			return nil, err
		}
		resolvers = append(resolvers, &publisherControlResolver{objectID, &publisherControllBySubmitter})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// registeredWorkResolver defines GraphQL resolver functions for registeredWork fields.
type publisherControlResolver struct {
	cid              string
	publisherControl *PublisherControllBySubmitter
}

func (p *publisherControlResolver) Cid() string {
	return p.cid
}

func (p *publisherControlResolver) RecordType() string {
	return p.publisherControl.RecordType
}

func (p *publisherControlResolver) PublisherSequenceN() string {
	return p.publisherControl.PublisherSequenceNumber
}

// registeredWorkResolver defines GraphQL resolver functions for registeredWork fields.
type registeredWorkResolver struct {
	cid            string
	registeredWork *RegisteredWork
}

func (r *registeredWorkResolver) Cid() string {
	return r.cid
}

func (r *registeredWorkResolver) Title() string {
	return r.registeredWork.Title
}

func (r *registeredWorkResolver) RecordType() string {
	return r.registeredWork.RecordType
}

func (r *registeredWorkResolver) ISWC() string {
	return r.registeredWork.ISWC
}

func (r *registeredWorkResolver) CatalogueNumber() string {
	return r.registeredWork.CatalogueNumber
}

func (r *registeredWorkResolver) CompositeComponentCount() string {
	return r.registeredWork.CompositeComponentCount
}

func (r *registeredWorkResolver) CompositeType() string {
	return r.registeredWork.CompositeType
}

func (r *registeredWorkResolver) ContactId() string {
	return r.registeredWork.ContactId
}

func (r *registeredWorkResolver) ContactName() string {
	return r.registeredWork.ContactName
}

func (r *registeredWorkResolver) CopyRightDate() string {
	return r.registeredWork.CopyRightDate
}

func (r *registeredWorkResolver) DateOfPublication() string {
	return r.registeredWork.DateOfPublication
}

func (r *registeredWorkResolver) DistributionCategory() string {
	return r.registeredWork.DistributionCategory
}

func (r *registeredWorkResolver) Duration() string {
	return r.registeredWork.Duration
}

func (r *registeredWorkResolver) ExceptionalClause() string {
	return r.registeredWork.ExceptionalClause
}

func (r *registeredWorkResolver) GrandRightsIndicator() string {
	return r.registeredWork.GrandRightsIndicator
}
func (r *registeredWorkResolver) LanguageCode() string {
	return r.registeredWork.LanguageCode
}

func (r *registeredWorkResolver) LyricAdaptation() string {
	return r.registeredWork.LyricAdaptation
}

func (r *registeredWorkResolver) MusicArrangement() string {
	return r.registeredWork.MusicArrangement
}
func (r *registeredWorkResolver) OpusNumber() string {
	return r.registeredWork.OpusNumber
}

func (r *registeredWorkResolver) PriorityFlag() string {
	return r.registeredWork.PriorityFlag
}

func (r *registeredWorkResolver) RecordedIndicator() string {
	return r.registeredWork.RecordedIndicator
}
func (r *registeredWorkResolver) SubmitteWorkNumber() string {
	return r.registeredWork.SubmitteWorkNumber
}

func (r *registeredWorkResolver) TextMusicRelationship() string {
	return r.registeredWork.TextMusicRelationship
}

func (r *registeredWorkResolver) VersionType() string {
	return r.registeredWork.VersionType
}
func (r *registeredWorkResolver) WorkType() string {
	return r.registeredWork.WorkType
}
