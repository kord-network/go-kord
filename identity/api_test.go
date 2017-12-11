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

package identity_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meta-network/go-meta/identity"
	"github.com/meta-network/go-meta/testutil"
	"github.com/meta-network/go-meta/testutil/index"
	"github.com/neelance/graphql-go"
)

type params struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

// TestIdentityAPI tests querying a identity index via the GraphQL API.
func TestIdentityAPI(t *testing.T) {
	// create a test index of identity
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, err := identity.NewIndex(store)
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()

	// start the API server
	s, err := newTestAPI(index)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	// define a function to execute and assert an identity GraphQL query
	assertIdentity := func(id *identity.Identity, p params) error {
		data, _ := json.Marshal(p)
		req, err := http.NewRequest("POST", s.URL+"/graphql", bytes.NewReader(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected HTTP status: %s", res.Status)
		}

		var r graphql.Response
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return err
		}
		if len(r.Errors) > 0 {
			return fmt.Errorf("unexpected errors in API response: %v", r.Errors)
		}

		var rw struct {
			Identities []*identity.Identity `json:"identity"`
		}
		if err := json.Unmarshal(r.Data, &rw); err != nil {
			return err
		}

		if len(rw.Identities) == 0 {
			return fmt.Errorf("expected Identity, got %d", len(rw.Identities))
		}

		for i, r := range rw.Identities {
			if r.Owner != id.Owner && i == len(rw.Identities) {
				return fmt.Errorf("unexpected identity owner: expected %q ", id.Owner)
			}
		}
		return nil
	}

	query := `
query GetIdentity($filter: IdentityFilter!) {
  identity(filter: $filter) {
    id
    owner
  }
}
`
	identity := testindex.GenerateTestIdentity(t)
	if err := index.CreateIdentity(identity); err != nil {
		t.Fatal(err)
	}
	if err := assertIdentity(
		identity,
		params{
			Query: query,
			Variables: map[string]interface{}{
				"filter": map[string]interface{}{
					"owner": identity.Owner.String(),
				},
			},
		},
	); err != nil {
		t.Fatal(err)
	}
}

// TestCreateIdentityAPI tests querying a identity index via the GraphQL API.
func TestCreateIdentityAPI(t *testing.T) {

	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, err := identity.NewIndex(store)
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()

	// start the API server
	s, err := newTestAPI(index)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	// define a function to execute and assert an identity GraphQL query
	assertCreateIdentity := func(id *identity.Identity, p params) error {
		data, _ := json.Marshal(p)
		req, err := http.NewRequest("POST", s.URL+"/graphql", bytes.NewReader(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected HTTP status: %s", res.Status)
		}

		var r graphql.Response
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return err
		}
		if len(r.Errors) > 0 {
			return fmt.Errorf("unexpected errors in API response: %v", r.Errors)
		}

		var rw struct {
			Identity *identity.Identity `json:"createIdentity"`
		}

		if err := json.Unmarshal(r.Data, &rw); err != nil {
			return err
		}
		if rw.Identity == nil {
			return fmt.Errorf("expected Identity, got nothing")
		}
		if rw.Identity.Owner != id.Owner {
			return fmt.Errorf("unexpected identity owner: expected %q ", id.Owner)
		}
		return nil
	}
	identity := testindex.GenerateTestIdentity(t)

	query := `
mutation CreateIdentity($input: IdentityInput!) {
  createIdentity(input: $input) {
    id
    owner
  }
}
`
	if err := assertCreateIdentity(
		identity,
		params{
			Query: query,
			Variables: map[string]interface{}{
				"input": identity,
			},
		},
	); err != nil {
		t.Fatal(err)
	}
}

// TestCreateIdentityAPI tests querying a identity index via the GraphQL API.
func TestCreateClaimAPI(t *testing.T) {

	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, err := identity.NewIndex(store)
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()

	// start the API server
	s, err := newTestAPI(index)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	// define a function to execute and assert a claim GraphQL mutation
	assertCreateClaim := func(claim *identity.Claim, p params) error {

		data, _ := json.Marshal(p)
		req, err := http.NewRequest("POST", s.URL+"/graphql", bytes.NewReader(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected HTTP status: %s", res.Status)
		}

		var r graphql.Response
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return err
		}
		if len(r.Errors) > 0 {
			return fmt.Errorf("unexpected errors in API response: %v", r.Errors)
		}

		var rw struct {
			Claim *identity.Claim `json:"createclaim"`
		}

		if err := json.Unmarshal(r.Data, &rw); err != nil {
			return err
		}
		if rw.Claim == nil {
			return fmt.Errorf("expected claim, got nothing")
		}
		if rw.Claim.Subject != claim.Subject {
			return fmt.Errorf("unexpected claim subject: expected %q ", claim.Subject)
		}
		if rw.Claim.Issuer != claim.Issuer {
			return fmt.Errorf("unexpected claim Issuer: expected %q ", claim.Issuer)
		}
		if rw.Claim.Property != claim.Property {
			return fmt.Errorf("unexpected claim Property: expected %q ", claim.Property)
		}
		if rw.Claim.Claim != claim.Claim {
			return fmt.Errorf("unexpected claim Claim: expected %q ", claim.Claim)
		}
		if !bytes.Equal(rw.Claim.Signature, claim.Signature) {
			return fmt.Errorf("unexpected claim Signature: expected %q ", claim.Signature)
		}
		return nil
	}
	identity := testindex.GenerateTestIdentity(t)
	if err := index.CreateIdentity(identity); err != nil {
		t.Fatal(err)
	}
	testClaim := testindex.GenerateTestClaim(t, identity, "name", "testname")

	query := `
mutation CreateClaim($input: ClaimInput!) {
  createClaim(input: $input) {
    id
    issuer
    subject
    property
    claim
    signature
  }
}`
	if err := assertCreateClaim(
		testClaim,
		params{
			Query: query,
			Variables: map[string]interface{}{
				"input": testClaim,
			},
		},
	); err != nil {
		t.Fatal(err)
	}
}

// TestClaimAPI tests querying a claim index via the GraphQL API.
func TestClaimAPI(t *testing.T) {
	// create a test index of identity
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, err := identity.NewIndex(store)
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()
	testIdentity := testindex.GenerateTestIdentity(t)
	if err := index.CreateIdentity(testIdentity); err != nil {
		t.Fatal(err)
	}
	claims := make([]*identity.Claim, 0, 2)
	for property, claim := range map[string]string{
		"dpid": "123",
		"ipi":  "xyz",
	} {
		claim := testindex.GenerateTestClaim(t, testIdentity, property, claim)
		if err := index.CreateClaim(claim); err != nil {
			t.Fatal(err)
		}
		claims = append(claims, claim)
	}

	// start the API server
	s, err := newTestAPI(index)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	// define a function to execute and assert a claim GraphQL query
	assertClaim := func(claim *identity.Claim, p params) error {
		data, _ := json.Marshal(p)
		req, err := http.NewRequest("POST", s.URL+"/graphql", bytes.NewReader(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected HTTP status: %s", res.Status)
		}

		var r graphql.Response
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return err
		}
		if len(r.Errors) > 0 {
			return fmt.Errorf("unexpected errors in API response: %v", r.Errors)
		}

		var rw struct {
			Claims []*identity.Claim `json:"claim"`
		}
		if err := json.Unmarshal(r.Data, &rw); err != nil {
			return err
		}

		if len(rw.Claims) == 0 {
			return fmt.Errorf("expected Claim, got %d", len(rw.Claims))
		}

		for i, r := range rw.Claims {
			if !bytes.Equal(r.Signature, claim.Signature) && i == len(rw.Claims) {
				return fmt.Errorf("unexpected claim name: expected %q ", claim.Signature)
			}
		}
		return nil
	}

	query := `
query GetClaim($filter: ClaimFilter!) {
  claim(filter: $filter) {
    id
    issuer
    subject
    claim
    signature
  }
}
`
	for _, claim := range claims {
		if err := assertClaim(
			claim,
			params{
				Query: query,
				Variables: map[string]interface{}{
					"filter": map[string]interface{}{
						"subject": claim.Subject,
					},
				},
			},
		); err != nil {
			t.Fatal(err)
		}
	}
}

func newTestAPI(index *identity.Index) (*httptest.Server, error) {
	api, err := identity.NewAPI(index)
	if err != nil {
		return nil, err
	}
	return httptest.NewServer(api), nil
}
