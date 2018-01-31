// This file is part of the go-meta library.
//
// Copyright (C) 2018 JAAK MUSIC LTD
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
	"errors"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/schema"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

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
  issuer:    String
  subject:   String
  property:  String
  claim:     String
}
`

type Resolver struct {
	qs graph.QuadStore
	qw graph.BatchWriter
}

func NewResolver(qs graph.QuadStore) (*Resolver, error) {
	qw, err := graph.NewQuadWriter("single", qs, nil)
	if err != nil {
		return nil, err
	}
	return &Resolver{qs, graph.NewWriter(qw)}, nil
}

// IdentityArgs are the arguments for a GraphQL identity query.
type IdentityArgs struct {
	Filter IdentityFilter
}

func (r *Resolver) Identity(args IdentityArgs) ([]*IdentityResolver, error) {
	path := path.NewPath(r.qs)
	if v := args.Filter.ID; v != nil {
		path = path.Is(quad.IRI(*v))
	}
	if v := args.Filter.Username; v != nil {
		path = path.Has(quad.IRI("id:username"), quad.StringToValue(*v))
	}
	if v := args.Filter.Owner; v != nil {
		path = path.Has(quad.IRI("id:owner"), quad.StringToValue(*v))
	}
	var identities []identityQuad
	if err := schema.LoadPathTo(context.Background(), r.qs, &identities, path); err != nil {
		return nil, err
	}
	resolvers := make([]*IdentityResolver, len(identities))
	for i, v := range identities {
		resolvers[i] = &IdentityResolver{v.Identity()}
	}
	return resolvers, nil
}

// IdentityResolver defines GraphQL resolver functions for Identity fields.
type IdentityResolver struct {
	identity *Identity
}

func (i *IdentityResolver) ID() string {
	return i.identity.ID().String()
}

func (i *IdentityResolver) Username() string {
	return i.identity.Username
}

func (i *IdentityResolver) Owner() string {
	return i.identity.Owner.String()
}

func (i *IdentityResolver) Signature() string {
	return hexutil.Encode(i.identity.Signature)
}

// CreateIdentityArgs are the arguments for a GraphQL CreateIdentity mutation.
type CreateIdentityArgs struct {
	Input IdentityInput
}

func (r *Resolver) CreateIdentity(args CreateIdentityArgs) (*IdentityResolver, error) {
	identity := &Identity{
		Username:  args.Input.Username,
		Owner:     common.HexToAddress(args.Input.Owner),
		Signature: common.FromHex(args.Input.Signature),
	}
	if !verifyIdentity(identity) {
		return nil, errors.New("identity: invalid identity")
	}
	if _, err := schema.WriteAsQuads(r.qw, identity.Quad()); err != nil {
		return nil, err
	}
	if err := r.qw.Flush(); err != nil {
		return nil, err
	}
	return &IdentityResolver{identity: identity}, nil
}

// ClaimArgs are the arguments for a GraphQL claim query.
type ClaimArgs struct {
	Filter ClaimFilter
}

func (r *Resolver) Claim(args ClaimArgs) ([]*ClaimResolver, error) {
	path := path.NewPath(r.qs)
	if v := args.Filter.Issuer; v != nil {
		path = path.Has(quad.IRI("id:issuer"), quad.IRI(*v))
	}
	if v := args.Filter.Subject; v != nil {
		path = path.Has(quad.IRI("id:subject"), quad.IRI(*v))
	}
	if v := args.Filter.Property; v != nil {
		path = path.Has(quad.IRI("id:property"), quad.StringToValue(*v))
	}
	if v := args.Filter.Claim; v != nil {
		path = path.Has(quad.IRI("id:claim"), quad.StringToValue(*v))
	}
	var claims []claimQuad
	if err := schema.LoadPathTo(context.Background(), r.qs, &claims, path); err != nil {
		return nil, err
	}
	resolvers := make([]*ClaimResolver, len(claims))
	for i, v := range claims {
		resolvers[i] = &ClaimResolver{v.ToClaim()}
	}
	return resolvers, nil
}

// ClaimResolver defines GraphQL resolver functions for Claim fields.
type ClaimResolver struct {
	claim *Claim
}

func (c *ClaimResolver) ID() string {
	return c.claim.ID().String()
}

func (c *ClaimResolver) Issuer() string {
	return c.claim.Issuer.String()
}

func (c *ClaimResolver) Subject() string {
	return c.claim.Subject.String()
}

func (c *ClaimResolver) Property() string {
	return c.claim.Property
}

func (c *ClaimResolver) Claim() string {
	return c.claim.Claim
}

func (c *ClaimResolver) Signature() string {
	return hexutil.Encode(c.claim.Signature)
}

// CreateClaimArgs are the arguments for a GraphQL CreateClaim mutation.
type CreateClaimArgs struct {
	Input ClaimInput
}

func (r *Resolver) CreateClaim(args CreateClaimArgs) (*ClaimResolver, error) {
	var v identityQuad
	if err := schema.LoadTo(context.Background(), r.qs, &v, quad.IRI(args.Input.Issuer)); err != nil {
		return nil, err
	}
	issuer := v.Identity()

	claim := &Claim{
		Issuer:    HexToID(args.Input.Issuer),
		Subject:   HexToID(args.Input.Subject),
		Property:  args.Input.Property,
		Claim:     args.Input.Claim,
		Signature: common.FromHex(args.Input.Signature),
	}
	if !verifyClaim(claim, issuer.Owner) {
		return nil, errors.New("claim: invalid claim")
	}
	if _, err := schema.WriteAsQuads(r.qw, claim.Quad()); err != nil {
		return nil, err
	}
	if err := r.qw.Flush(); err != nil {
		return nil, err
	}
	return &ClaimResolver{claim: claim}, nil
}

func verifyIdentity(identity *Identity) bool {
	id := identity.ID()
	return verifySignature(identity.Owner, id.Hash[:], identity.Signature)
}

func verifyClaim(claim *Claim, owner common.Address) bool {
	id := claim.ID()
	return verifySignature(owner, id[:], claim.Signature)
}

func verifySignature(owner common.Address, msg, signature []byte) bool {
	recoveredPub, err := crypto.Ecrecover(msg, signature)
	if err != nil {
		return false
	}
	pubKey := crypto.ToECDSAPub(recoveredPub)
	return crypto.PubkeyToAddress(*pubKey) == owner
}
