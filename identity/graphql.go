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
	_ "github.com/cayleygraph/cayley/writer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	metagraph "github.com/meta-network/go-meta/graph"
)

const GraphQLSchema = `
schema {
  query: Query
  mutation: Mutation
}

type Query {
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

type Claim {
  id:        String!
  graph:     String!
  issuer:    String!
  subject:   String!
  property:  String!
  claim:     String!
  signature: String!
}

input ClaimInput {
  graph:     String!
  issuer:    String!
  subject:   String!
  property:  String!
  claim:     String!
  signature: String!
}

input ClaimFilter {
  graph:     String!
  issuer:    String
  subject:   String
  property:  String
  claim:     String
}
`

type Resolver struct {
	driver *metagraph.Driver
}

func NewResolver(driver *metagraph.Driver) *Resolver {
	return &Resolver{driver}
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
	name := identity.Username + ".meta"
	if err := r.driver.Create(name); err != nil {
		return nil, err
	}
	return &IdentityResolver{identity: identity}, nil
}

// ClaimArgs are the arguments for a GraphQL claim query.
type ClaimArgs struct {
	Filter ClaimFilter
}

func (r *Resolver) Claim(args ClaimArgs) ([]*ClaimResolver, error) {
	qs, err := r.driver.Get(args.Filter.Graph)
	if err != nil {
		return nil, err
	}
	path := path.NewPath(qs)
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
	if err := schema.LoadPathTo(context.Background(), qs, &claims, path); err != nil {
		return nil, err
	}
	resolvers := make([]*ClaimResolver, len(claims))
	for i, v := range claims {
		resolvers[i] = &ClaimResolver{v.ToClaim(args.Filter.Graph)}
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

func (c *ClaimResolver) Graph() string {
	return c.claim.Graph
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
	claim := &Claim{
		Issuer:    HexToID(args.Input.Issuer),
		Subject:   HexToID(args.Input.Subject),
		Property:  args.Input.Property,
		Claim:     args.Input.Claim,
		Signature: common.FromHex(args.Input.Signature),
	}

	qs, err := r.driver.Get(args.Input.Graph)
	if err != nil {
		return nil, err
	}
	qw, err := graph.NewQuadWriter("single", qs, nil)
	if err != nil {
		return nil, err
	}
	w := graph.NewWriter(qw)

	if _, err := schema.WriteAsQuads(w, claim.Quad()); err != nil {
		return nil, err
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}
	return &ClaimResolver{claim: claim}, nil
}

func verifyIdentity(identity *Identity) bool {
	id := identity.ID()
	return verifySignature(identity.Owner, id.Hash[:], identity.Signature)
}

func verifySignature(owner common.Address, msg, signature []byte) bool {
	recoveredPub, err := crypto.Ecrecover(msg, signature)
	if err != nil {
		return false
	}
	pubKey := crypto.ToECDSAPub(recoveredPub)
	return crypto.PubkeyToAddress(*pubKey) == owner
}
