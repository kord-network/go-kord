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

package identity

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	cid "github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// GraphQLSchema is the GraphQL schema for the MusicBrainz META index.
const GraphQLSchema = `
schema {
  query: Query
  mutation: Mutation
}

type Query {
  identity(id: String,owner: String): [Identity]!

  claim(issuer: String,holder: String,claim:String,signature:String,id:String): [Claim]!
}

type Mutation {
	createIdentity(username: String, owner: String,signature:String): Identity

  createClaim(issuer: String,holder: String,claim: String,signature: String): Claim
}

type Identity {
  id:        String!
  owner:     String!
  signature: String!
}

type Claim {
  issuer:     String!
  holder:     String!
  claim:      String!
  signature : String!
  id :        String!
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

// IdentityArgs are the arguments for a GraphQL identity query.
type IdentityArgs struct {
	Owner *string
	ID    *string
}

// Identity is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either an Identity ID or Owner and loads the
// associated META objects from the META store.
func (r *Resolver) Identity(args IdentityArgs) ([]*IdentityResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.Owner != nil:
		rows, err = r.db.Query("SELECT object_id FROM identity WHERE owner = ?", *args.Owner)
	case args.ID != nil:
		rows, err = r.db.Query("SELECT object_id FROM identity WHERE id = ?", *args.ID)
	default:
		return nil, errors.New("missing owner or id argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*IdentityResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		objCid, err := cid.Parse(objectID)
		if err != nil {
			return nil, err
		}
		obj, err := r.store.Get(objCid)
		if err != nil {
			return nil, err
		}
		var identity Identity
		if err := obj.Decode(&identity); err != nil {
			return nil, err
		}
		resolvers = append(resolvers, &IdentityResolver{objectID, &identity})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// IdentityResolver defines GraphQL resolver functions for Identity fields.
type IdentityResolver struct {
	cid      string
	identity *Identity
}

// Cid resolver
func (i *IdentityResolver) Cid() string {
	return i.cid
}

// Owner resolver
func (i *IdentityResolver) Owner() string {
	return i.identity.Owner
}

// ID resolver
func (i *IdentityResolver) ID() string {
	return i.identity.ID
}

// Signature resolver
func (i *IdentityResolver) Signature() string {
	return i.identity.Sig
}

// ClaimArgs are the arguments for a GraphQL claim query.
type ClaimArgs struct {
	Issuer    *string
	Holder    *string
	Claim     *string
	Signature *string
	ID        *string
}

// Claim is a GraphQL resolver function which retrieves object Claims from the
// SQLite3 index using either Claim ID,Issuer,holder,Claim or Signature and loads the
// associated META objects from the META store.
func (r *Resolver) Claim(args ClaimArgs) ([]*ClaimResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.Issuer != nil:
		rows, err = r.db.Query("SELECT object_id FROM claim WHERE issuer = ?", *args.Issuer)
	case args.Holder != nil:
		rows, err = r.db.Query("SELECT object_id FROM claim WHERE holder = ?", *args.Holder)
	case args.Claim != nil:
		rows, err = r.db.Query("SELECT object_id FROM claim WHERE holder = ?", *args.Claim)
	case args.Signature != nil:
		rows, err = r.db.Query("SELECT object_id FROM claim WHERE signature = ?", *args.Signature)
	case args.ID != nil:
		rows, err = r.db.Query("SELECT object_id FROM claim WHERE id = ?", *args.ID)
	default:
		return nil, errors.New("missing issuer,holder,claim,signature or id argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*ClaimResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		objCid, err := cid.Parse(objectID)
		if err != nil {
			return nil, err
		}
		obj, err := r.store.Get(objCid)
		if err != nil {
			return nil, err
		}
		var claim Claim
		if err := obj.Decode(&claim); err != nil {
			return nil, err
		}
		resolvers = append(resolvers, &ClaimResolver{objectID, &claim})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// ClaimResolver defines GraphQL resolver functions for Claim fields.
type ClaimResolver struct {
	cid   string
	claim *Claim
}

// Cid resolver
func (c *ClaimResolver) Cid() string {
	return c.cid
}

// Issuer resolver
func (c *ClaimResolver) Issuer() string {
	return c.claim.Issuer
}

// Holder resolver
func (c *ClaimResolver) Holder() string {
	return c.claim.Holder
}

// Claim resolver
func (c *ClaimResolver) Claim() string {
	return c.claim.Claim
}

// Signature resolver
func (c *ClaimResolver) Signature() string {
	return c.claim.Signature
}

// ID resolver
func (c *ClaimResolver) ID() string {
	return c.claim.ID
}

// CreateIDArgs are the arguments for a GraphQL CreateID mutation.
type CreateIDArgs struct {
	Username  *string
	Owner     *string
	Signature *string
}

// CreateID is a GraphQL resolver function which create ID object, store and
// index it on SQLite3 index.
func (r *Resolver) CreateIdentity(args *CreateIDArgs) (*IdentityResolver, error) {

	metaid, err := NewIdentity(*args.Username, common.HexToAddress(*args.Owner), common.FromHex(*args.Signature))
	if err != nil {
		return nil, err
	}
	converter := NewConverter(r.store)
	identityCid, err := converter.ConvertIdentity(metaid)
	if err != nil {
		return nil, err
	}
	// create a stream of ID
	writer, err := r.store.StreamWriter("id.meta")
	if err != nil {
		return nil, err
	}
	defer writer.Close()
	if err = writer.Write(identityCid); err != nil {
		return nil, err
	}

	// index the stream of ID
	index, err := r.store.OpenIndex("id.index.meta")
	if err != nil {
		return nil, err
	}
	indexer, err := NewIndexer(index, r.store)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reader, err := r.store.StreamReader("id.meta", meta.StreamLimit(1))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	if err := indexer.Index(ctx, reader); err != nil {
		return nil, err
	}
	return &IdentityResolver{cid: identityCid.String(), identity: metaid}, nil
}

// CreateClaimArgs are the arguments for a GraphQL CreateClaim mutation.
type CreateClaimArgs struct {
	Issuer    *string
	Holder    *string
	Claim     *string
	Signature *string
}

// CreateClaim is a GraphQL resolver function which create claim object, store and
// index it on SQLite3 index.
func (r *Resolver) CreateClaim(args *CreateClaimArgs) (*ClaimResolver, error) {

	claim := NewClaim(*args.Issuer, *args.Holder, *args.Claim, *args.Signature)
	converter := NewConverter(r.store)
	claimCid, err := converter.ConvertClaim(claim)
	if err != nil {
		return nil, err
	}
	// create a stream of ID
	writer, err := r.store.StreamWriter("claim.meta")
	if err != nil {
		return nil, err
	}
	defer writer.Close()
	if err = writer.Write(claimCid); err != nil {
		return nil, err
	}
	// index the stream of ID
	index, err := r.store.OpenIndex("claim.index.meta")
	if err != nil {
		return nil, err
	}
	indexer, err := NewIndexer(index, r.store)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reader, err := r.store.StreamReader("claim.meta", meta.StreamLimit(1))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	if err := indexer.Index(ctx, reader); err != nil {
		return nil, err
	}
	return &ClaimResolver{cid: claimCid.String(), claim: claim}, nil
}
