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
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/common"
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

  claim(issuer: String,subject: String,claim:String,signature:String,id:String): [Claim]!
}

type Mutation {
	createIdentity(username: String, owner: String,signature:String): Identity

  createClaim(issuer: String,subject: String,claim: String,signature: String): Claim
}

type Identity {
  id:        String!
  owner:     String!
  signature: String!
}

type Claim {
  issuer:     String!
  subject:    String!
  claim:      String!
  signature : String!
  id :        String!
}
`

// Resolver defines GraphQL resolver functions for the schema contained in
// the GraphQLSchema constant, retrieving data from a META SQLite3
// index.
type Resolver struct {
	db      *sql.DB
	indexer *Indexer
}

// NewResolver returns a Resolver which retrieves data from the given META
// SQLite3 index.
func NewResolver(db *sql.DB, index *meta.Index) (*Resolver, error) {
	indexer, err := NewIndexer(index)
	if err != nil {
		return nil, err
	}
	return &Resolver{db, indexer}, nil
}

// IdentityArgs are the arguments for a GraphQL identity query.
type IdentityArgs struct {
	Owner *string
	ID    *string
}

// Identity is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either an Identity ID or Owner and loads the
// associated META objects from the META index.
func (r *Resolver) Identity(args IdentityArgs) ([]*IdentityResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.Owner != nil:
		rows, err = r.db.Query("SELECT * FROM identity WHERE owner = ?", *args.Owner)
	case args.ID != nil:
		rows, err = r.db.Query("SELECT * FROM identity WHERE id = ?", *args.ID)
	default:
		return nil, errors.New("missing owner or id argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*IdentityResolver
	for rows.Next() {

		var owner, signature, id string
		if err := rows.Scan(&owner, &signature, &id); err != nil {
			return nil, err
		}

		resolvers = append(resolvers, &IdentityResolver{&Identity{Owner: common.HexToAddress(owner), Sig: signature, ID: id}})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// IdentityResolver defines GraphQL resolver functions for Identity fields.
type IdentityResolver struct {
	identity *Identity
}

// Owner resolver
func (i *IdentityResolver) Owner() string {
	return i.identity.Owner.String()
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
	Subject   *string
	Claim     *string
	Signature *string
	ID        *string
}

// Claim is a GraphQL resolver function which retrieves object Claims from the
// SQLite3 index using either Claim ID,Issuer,subject,Claim or Signature and loads the
// associated META objects from the META index.
func (r *Resolver) Claim(args ClaimArgs) ([]*ClaimResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.Issuer != nil:
		rows, err = r.db.Query("SELECT * FROM claim WHERE issuer = ?", *args.Issuer)
	case args.Subject != nil:
		rows, err = r.db.Query("SELECT * FROM claim WHERE subject = ?", *args.Subject)
	case args.Claim != nil:
		rows, err = r.db.Query("SELECT * FROM claim WHERE subject = ?", *args.Claim)
	case args.Signature != nil:
		rows, err = r.db.Query("SELECT * FROM claim WHERE signature = ?", *args.Signature)
	case args.ID != nil:
		rows, err = r.db.Query("SELECT * FROM claim WHERE id = ?", *args.ID)
	default:
		return nil, errors.New("missing issuer,subject,claim,signature or id argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*ClaimResolver
	for rows.Next() {
		var claim Claim
		if err := rows.Scan(&claim.Issuer, &claim.Subject, &claim.Claim, &claim.Signature, &claim.ID); err != nil {
			return nil, err
		}
		resolvers = append(resolvers, &ClaimResolver{&claim})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// ClaimResolver defines GraphQL resolver functions for Claim fields.
type ClaimResolver struct {
	claim *Claim
}

// Issuer resolver
func (c *ClaimResolver) Issuer() string {
	return c.claim.Issuer
}

// Subject resolver
func (c *ClaimResolver) Subject() string {
	return c.claim.Subject
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

// CreateIdentityArgs are the arguments for a GraphQL CreateIdentity mutation.
type CreateIdentityArgs struct {
	Username  *string
	Owner     *string
	Signature *string
}

// CreateIdentity is a GraphQL resolver function which create identity and
// index it on SQLite3 index.
func (r *Resolver) CreateIdentity(args *CreateIdentityArgs) (*IdentityResolver, error) {
	if args.Owner == nil || args.Username == nil || args.Signature == nil {
		return nil, errors.New("CreateIdentity: one or more argument is nil")
	}
	metaid, err := NewIdentity(*args.Username, common.HexToAddress(*args.Owner), common.FromHex(*args.Signature))
	if err != nil {
		return nil, err
	}
	if err := r.indexer.IndexIdentity(metaid); err != nil {
		return nil, err
	}
	return &IdentityResolver{identity: metaid}, nil
}

// CreateClaimArgs are the arguments for a GraphQL CreateClaim mutation.
type CreateClaimArgs struct {
	Issuer    *string
	Subject   *string
	Claim     *string
	Signature *string
}

// CreateClaim is a GraphQL resolver function which create claim and
// index it on SQLite3 index.
func (r *Resolver) CreateClaim(args *CreateClaimArgs) (*ClaimResolver, error) {

	if args.Issuer == nil || args.Subject == nil || args.Claim == nil || args.Signature == nil {
		return nil, errors.New("CreateClaim: one or more argument is nil")
	}
	claim := NewClaim(*args.Issuer, *args.Subject, *args.Claim, *args.Signature)

	if err := r.indexer.IndexClaim(claim); err != nil {
		return nil, err
	}
	return &ClaimResolver{claim: claim}, nil
}
