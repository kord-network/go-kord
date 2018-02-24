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
  graph(id: String!): Graph!
}

type Mutation {
  createGraph(input: GraphInput!): Graph!
  createClaim(input: ClaimInput!): Claim!
}

type Graph {
  id: String!

  claim(filter: ClaimFilter!): [Claim]!
}

input GraphInput {
  id: String!
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
  graph:     String!
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

type Resolver struct {
	driver *metagraph.Driver
}

func NewResolver(driver *metagraph.Driver) *Resolver {
	return &Resolver{driver}
}

type GraphArgs struct {
	ID string
}

func (r *Resolver) Graph(args GraphArgs) (*GraphResolver, error) {
	qs, err := r.driver.Get(args.ID)
	if err != nil {
		return nil, err
	}
	return &GraphResolver{args.ID, qs}, nil
}

type GraphResolver struct {
	id string
	qs graph.QuadStore
}

func (r *GraphResolver) ID() string {
	return r.id
}

// CreateGraphArgs are the arguments for a GraphQL CreateGraph mutation.
type CreateGraphArgs struct {
	Input GraphInput
}

func (r *Resolver) CreateGraph(ctx context.Context, args CreateGraphArgs) (*GraphResolver, error) {
	hash, err := r.driver.Create(args.Input.ID)
	if err != nil {
		return nil, err
	}
	ctx.Value("swarmHash").(*common.Hash).Set(hash)
	return r.Graph(GraphArgs{ID: args.Input.ID})
}

// ClaimArgs are the arguments for a GraphQL claim query.
type ClaimArgs struct {
	Filter ClaimFilter
}

func (r *GraphResolver) Claim(args ClaimArgs) ([]*ClaimResolver, error) {
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

func (r *Resolver) CreateClaim(ctx context.Context, args CreateClaimArgs) (*ClaimResolver, error) {
	claim := &Claim{
		Issuer:    HexToID(args.Input.Issuer),
		Subject:   HexToID(args.Input.Subject),
		Property:  args.Input.Property,
		Claim:     args.Input.Claim,
		Signature: common.FromHex(args.Input.Signature),
	}

	graph := args.Input.Graph
	if err := r.writeClaim(graph, claim); err != nil {
		return nil, err
	}

	hash, err := r.driver.Commit(graph)
	if err != nil {
		return nil, err
	}
	ctx.Value("swarmHash").(*common.Hash).Set(hash)

	return &ClaimResolver{claim: claim}, nil
}

func (r *Resolver) writeClaim(id string, claim *Claim) error {
	if !verifyClaim(claim) {
		return errors.New("invalid claim signature")
	}
	qs, err := r.driver.Get(id)
	if err != nil {
		return err
	}
	qw, err := graph.NewQuadWriter("single", qs, nil)
	if err != nil {
		return err
	}
	w := graph.NewWriter(qw)

	if _, err := schema.WriteAsQuads(w, claim.Quad()); err != nil {
		return err
	}
	return w.Flush()
}

func verifyClaim(claim *Claim) bool {
	id := claim.ID()
	recoveredPub, err := crypto.Ecrecover(id[:], claim.Signature)
	if err != nil {
		return false
	}
	pubKey := crypto.ToECDSAPub(recoveredPub)
	return crypto.PubkeyToAddress(*pubKey) == claim.Issuer.Address
}
