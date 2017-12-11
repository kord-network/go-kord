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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// GraphQLSchema is the GraphQL schema for the MusicBrainz META index.
const GraphQLSchema = `
schema {
  query: Query
  mutation: Mutation
}

type Query {
  identity(filter: IdentityFilter!): [Identity]!

  claim(filter: ClaimFilter!): [Claim]!
}

type Mutation {
  createIdentity(input: IdentityInput!): Identity!

  createClaim(input: ClaimInput!): Claim!
}

type Identity {
  id:        String!
  username:  String!
  owner:     String!
  signature: String!
}

input IdentityInput {
  username:  String!
  owner:     String!
  signature: String!
}

input IdentityFilter {
  id:       String
  username: String
  owner:    String
}

type Claim {
  id:        String!
  issuer:    String!
  subject:   String!
  property:  String!
  claim:     String!
  signature: String!
}

input ClaimInput {
  issuer:    String!
  subject:   String!
  property:  String!
  claim:     String!
  signature: String!
}

input ClaimFilter {
  issuer:   String
  subject:  String
  property: String
  claim:    String
}
`

// Resolver defines GraphQL resolver functions for the schema contained in
// the GraphQLSchema constant, retrieving data from a META SQLite3
// index.
type Resolver struct {
	index *Index
}

// NewResolver returns a Resolver which retrieves data from the given META
// SQLite3 index.
func NewResolver(index *Index) *Resolver {
	return &Resolver{index}
}

// IdentityArgs are the arguments for a GraphQL identity query.
type IdentityArgs struct {
	Filter IdentityFilter
}

// Identity is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either an Identity ID or Owner and loads the
// associated META objects from the META index.
func (r *Resolver) Identity(args IdentityArgs) ([]*IdentityResolver, error) {
	identities, err := r.index.Identities(args.Filter)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*IdentityResolver, len(identities))
	for i, identity := range identities {
		resolvers[i] = &IdentityResolver{identity}
	}
	return resolvers, nil
}

// IdentityResolver defines GraphQL resolver functions for Identity fields.
type IdentityResolver struct {
	identity *Identity
}

// ID resolver
func (i *IdentityResolver) ID() string {
	return i.identity.ID().String()
}

// Username resolver
func (i *IdentityResolver) Username() string {
	return i.identity.Username
}

// Owner resolver
func (i *IdentityResolver) Owner() string {
	return i.identity.Owner.String()
}

// Signature resolver
func (i *IdentityResolver) Signature() string {
	return hexutil.Encode(i.identity.Signature)
}

// ClaimArgs are the arguments for a GraphQL claim query.
type ClaimArgs struct {
	Filter ClaimFilter
}

// Claim is a GraphQL resolver function which retrieves object Claims from the
// SQLite3 index using either Claim ID,Issuer,subject,Claim or Signature and loads the
// associated META objects from the META index.
func (r *Resolver) Claim(args ClaimArgs) ([]*ClaimResolver, error) {
	claims, err := r.index.Claims(args.Filter)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*ClaimResolver, len(claims))
	for i, claim := range claims {
		resolvers[i] = &ClaimResolver{claim}
	}
	return resolvers, nil
}

// ClaimResolver defines GraphQL resolver functions for Claim fields.
type ClaimResolver struct {
	claim *Claim
}

// ID resolver
func (c *ClaimResolver) ID() string {
	return c.claim.ID().String()
}

// Issuer resolver
func (c *ClaimResolver) Issuer() string {
	return c.claim.Issuer.String()
}

// Subject resolver
func (c *ClaimResolver) Subject() string {
	return c.claim.Subject.String()
}

// Property resolver
func (c *ClaimResolver) Property() string {
	return c.claim.Property
}

// Claim resolver
func (c *ClaimResolver) Claim() string {
	return c.claim.Claim
}

// Signature resolver
func (c *ClaimResolver) Signature() string {
	return hexutil.Encode(c.claim.Signature)
}

// CreateIdentityArgs are the arguments for a GraphQL CreateIdentity mutation.
type CreateIdentityArgs struct {
	Input struct {
		Username  string
		Owner     string
		Signature string
	}
}

// CreateIdentity is a GraphQL resolver function which create identity and
// index it on SQLite3 index.
func (r *Resolver) CreateIdentity(args CreateIdentityArgs) (*IdentityResolver, error) {
	identity := &Identity{
		Username:  args.Input.Username,
		Owner:     common.HexToAddress(args.Input.Owner),
		Signature: common.FromHex(args.Input.Signature),
	}
	if err := r.index.CreateIdentity(identity); err != nil {
		return nil, err
	}
	return &IdentityResolver{identity: identity}, nil
}

// CreateClaimArgs are the arguments for a GraphQL CreateClaim mutation.
type CreateClaimArgs struct {
	Input struct {
		Issuer    string
		Subject   string
		Property  string
		Claim     string
		Signature string
	}
}

// CreateClaim is a GraphQL resolver function which create claim and
// index it on SQLite3 index.
func (r *Resolver) CreateClaim(args *CreateClaimArgs) (*ClaimResolver, error) {
	claim := &Claim{
		Issuer:    HexToID(args.Input.Issuer),
		Subject:   HexToID(args.Input.Subject),
		Property:  args.Input.Property,
		Claim:     args.Input.Claim,
		Signature: common.FromHex(args.Input.Signature),
	}
	if err := r.index.CreateClaim(claim); err != nil {
		return nil, err
	}
	return &ClaimResolver{claim: claim}, nil
}
