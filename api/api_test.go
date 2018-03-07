// This file is part of the go-kord library.
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

package api

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kord-network/go-kord/graph"
	"github.com/kord-network/go-kord/testutil"
)

func TestAPI(t *testing.T) {
	// create a test API
	dpa, err := testutil.NewTestDPA()
	if err != nil {
		t.Fatal(err)
	}
	defer dpa.Cleanup()
	registry := testutil.NewTestRegistry()
	driver := graph.NewDriver("kord-id-test", dpa.DPA, registry, dpa.Dir)
	api, err := NewAPI(driver)
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(api)
	defer srv.Close()

	// create a graph
	client := NewClient(srv.URL)
	hash, err := client.CreateGraph(testKordID.Hex())
	if err != nil {
		t.Fatal(err)
	}
	sig, err := crypto.Sign(hash[:], testKey)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.SetGraph(testKordID.Hex(), hash, sig); err != nil {
		t.Fatal(err)
	}

	// create a claim
	claim := newTestClaim(t, "username", "test")
	hash, err = client.CreateClaim(testKordID.Hex(), claim)
	if err != nil {
		t.Fatal(err)
	}
	sig, err = crypto.Sign(hash[:], testKey)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.SetGraph(testKordID.Hex(), hash, sig); err != nil {
		t.Fatal(err)
	}

	// get the claim
	id := testKordID.Hex()
	claims, err := client.Claim(testKordID.Hex(), &ClaimFilter{
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

var (
	testKey, _ = crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032")
	testKordID = NewID(crypto.PubkeyToAddress(testKey.PublicKey))
)

func newTestClaim(t *testing.T, property, claim string) *Claim {
	c := &Claim{
		Issuer:   testKordID,
		Subject:  testKordID,
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
