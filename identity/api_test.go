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
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/meta-network/go-meta/graph"
	"github.com/meta-network/go-meta/testutil"
)

func TestIdentityAPI(t *testing.T) {
	// create a test API
	dpa, err := testutil.NewTestDPA()
	if err != nil {
		t.Fatal(err)
	}
	defer dpa.Cleanup()
	driver := graph.NewDriver("meta-id-test", dpa.DPA, testutil.NewTestENS(), dpa.Dir)
	api, err := NewAPI(driver)
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(api)
	defer srv.Close()

	// create an identity
	client := NewClient(srv.URL + "/graphql")
	identity := newTestIdentity(t)
	if err := client.CreateIdentity(identity); err != nil {
		t.Fatal(err)
	}

	// create a claim
	claim := newTestClaim(t, identity, "test.id", "123")
	if err := client.CreateClaim(claim); err != nil {
		t.Fatal(err)
	}

	// get the claim
	id := identity.ID().String()
	claims, err := client.Claim(&ClaimFilter{
		Graph:    claim.Graph,
		Issuer:   &id,
		Subject:  &id,
		Property: &claim.Property,
		Claim:    &claim.Claim,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(claims) != 1 {
		t.Fatalf("expected 1 claim, got %d", len(claims))
	}
	gotClaim := claims[0]
	if gotClaim.ID() != claim.ID() {
		t.Fatalf("expected claim to have ID %s, got %s", claim.ID().String(), gotClaim.ID().String())
	}
	if !bytes.Equal(gotClaim.Signature, claim.Signature) {
		t.Fatalf("expected claim to have signature %s, got %s", hexutil.Encode(claim.Signature), hexutil.Encode(gotClaim.Signature))
	}
}

var testKey, _ = crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032")

func newTestIdentity(t *testing.T) *Identity {
	identity := &Identity{
		Username: "test",
		Owner:    crypto.PubkeyToAddress(testKey.PublicKey),
	}
	id := identity.ID()
	signature, err := crypto.Sign(id.Hash[:], testKey)
	if err != nil {
		t.Fatal(err)
	}
	identity.Signature = signature
	return identity
}

func newTestClaim(t *testing.T, identity *Identity, property, claim string) *Claim {
	c := &Claim{
		Graph:    identity.Username + ".meta",
		Issuer:   identity.ID(),
		Subject:  identity.ID(),
		Property: property,
		Claim:    claim,
	}
	id := c.ID()
	signature, err := crypto.Sign(id[:], testKey)
	if err != nil {
		t.Fatal(err)
	}
	c.Signature = signature
	return c
}
