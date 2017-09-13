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

package musicbrainz

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
  artist(
    name: String,
    ipi:  String,
    isni: String
  ): [Artist]!
}

type Artist {
  cid:                    String!
  name:                   String!
  sort_name:              String!
  type:                   String
  gender:                 String
  area:                   String
  begin_date:             String
  end_date:               String
  ipi:                    [String!]
  isni:                   [String!]
  alias:                  [String!]
  mbid:                   String!
  disambiguation_comment: String
  annotation:             [String!]
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

// artistArgs are the arguments for a GraphQL artist query.
type artistArgs struct {
	Name *string
	IPI  *string
	ISNI *string
}

// Artist is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either an artist name, IPI or ISNI, and loads the
// associated META objects from the META store.
func (g *Resolver) Artist(args artistArgs) ([]*artistResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.Name != nil:
		rows, err = g.db.Query("SELECT object_id FROM artist WHERE name = ?", *args.Name)
	case args.IPI != nil:
		rows, err = g.db.Query("SELECT object_id FROM artist_ipi WHERE ipi = ?", *args.IPI)
	case args.ISNI != nil:
		rows, err = g.db.Query("SELECT object_id FROM artist_isni WHERE isni = ?", *args.ISNI)
	default:
		return nil, errors.New("missing name, ipi or isni argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*artistResolver
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
		var artist Artist
		if err := obj.Decode(&artist); err != nil {
			return nil, err
		}
		resolvers = append(resolvers, &artistResolver{objectID, &artist})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// artistResolver defines GraphQL resolver functions for artist fields.
type artistResolver struct {
	cid    string
	artist *Artist
}

func (a *artistResolver) Cid() string {
	return a.cid
}

func (a *artistResolver) Name() string {
	return a.artist.Name
}

func (a *artistResolver) SortName() string {
	return a.artist.SortName
}

func (a *artistResolver) Type() *string {
	return &a.artist.Type
}

func (a *artistResolver) Gender() *string {
	return &a.artist.Gender
}

func (a *artistResolver) Area() *string {
	return &a.artist.Area
}

func (a *artistResolver) BeginDate() *string {
	return &a.artist.BeginDate
}

func (a *artistResolver) EndDate() *string {
	return &a.artist.EndDate
}

func (a *artistResolver) Ipi() *[]string {
	if a.artist.IPI == nil {
		return nil
	}
	return &a.artist.IPI
}

func (a *artistResolver) Isni() *[]string {
	if a.artist.ISNI == nil {
		return nil
	}
	return &a.artist.ISNI
}

func (a *artistResolver) Alias() *[]string {
	if a.artist.Alias == nil {
		return nil
	}
	return &a.artist.Alias
}

func (a *artistResolver) Mbid() string {
	return a.artist.MBID
}

func (a *artistResolver) DisambiguationComment() *string {
	return &a.artist.DisambiguationComment
}

func (a *artistResolver) Annotation() *[]string {
	if a.artist.Annotation == nil {
		return nil
	}
	return &a.artist.Annotation
}
