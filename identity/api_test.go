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
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meta-network/go-meta"
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
	index, id := testindex.GenerateIdentityIndex(t, ".", store)
	defer index.Close()

	// start the API server
	s, err := newTestAPI(index.DB, index)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	// define a function to execute and assert an identity GraphQL query
	assertIdentity := func(id *identity.Identity, query string, args ...interface{}) error {
		data, _ := json.Marshal(map[string]string{"query": fmt.Sprintf(query, args...)})
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

	if id.Owner.String() != "" {
		if err := assertIdentity(id,
			`{ identity (owner:%q) {id owner}}`,
			id.Owner.String()); err != nil {
			t.Fatal(err)
		}
	}
}

// TestCreateIdentityAPI tests querying a identity index via the GraphQL API.
func TestCreateIdentityAPI(t *testing.T) {

	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, err := store.OpenIndex("id.index.meta")
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()

	// start the API server
	s, err := newTestAPI(index.DB, index)
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
	testMetaId, err := testindex.GenerateTestMetaId()
	if err != nil {
		t.Fatal(err)
	}

	if err := assertCreateIdentity(testMetaId,
		params{Query: `mutation CreateIdentity($username: String, $owner: String,$signature:String) {
						createIdentity(username: $username, owner: $owner,signature:$signature) {
							id
							owner
						}
					}`,
			Variables: map[string]interface{}{
				"username":  "testid",
				"owner":     testMetaId.Owner,
				"signature": testMetaId.Sig}}); err != nil {
		t.Fatal(err)
	}
}

// TestCreateIdentityAPI tests querying a identity index via the GraphQL API.
func TestCreateClaimAPI(t *testing.T) {

	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, err := store.OpenIndex("claim.index.meta")
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()

	// start the API server
	s, err := newTestAPI(index.DB, index)
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
		if rw.Claim.Claim != claim.Claim {
			return fmt.Errorf("unexpected claim Claim: expected %q ", claim.Claim)
		}
		if rw.Claim.Signature != claim.Signature {
			return fmt.Errorf("unexpected claim Signature: expected %q ", claim.Signature)
		}
		return nil
	}
	testMetaId, err := testindex.GenerateTestMetaId()
	if err != nil {
		t.Fatal(err)
	}
	testClaim := identity.NewClaim(testMetaId.ID, testMetaId.ID, "name", "testname")

	if err := assertCreateClaim(testClaim,
		params{Query: `mutation CreateClaim($issuer: String, $subject: String,$claim: String,$signature: String) {
						createClaim(issuer: $issuer, subject: $subject,claim: $claim,signature: $signature) {
							id
							issuer
							subject
							claim
							signature
						}
					}`,
			Variables: map[string]interface{}{
				"issuer":    testClaim.Issuer,
				"subject":   testClaim.Subject,
				"claim":     "name",
				"signature": "testname"}}); err != nil {
		t.Fatal(err)
	}
}

// TestClaimAPI tests querying a claim index via the GraphQL API.
func TestClaimAPI(t *testing.T) {
	// create a test index of identity
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, claims := testindex.GenerateClaimIndex(t, ".", store)
	defer index.Close()

	// start the API server
	s, err := newTestAPI(index.DB, index)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	// define a function to execute and assert a claim GraphQL query
	assertClaim := func(claim *identity.Claim, query string, args ...interface{}) error {
		data, _ := json.Marshal(map[string]string{"query": fmt.Sprintf(query, args...)})
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
			if r.Signature != claim.Signature && i == len(rw.Claims) {
				return fmt.Errorf("unexpected claim name: expected %q ", claim.Signature)
			}
		}
		return nil
	}

	if claims[0].Signature != "" {
		if err := assertClaim(claims[0],
			`{ claim (subject:%q) {issuer subject claim signature }}`,
			claims[0].Subject); err != nil {
			t.Fatal(err)
		}
	}
}

func newTestAPI(db *sql.DB, index *meta.Index) (*httptest.Server, error) {
	api, err := identity.NewAPI(db, index)
	if err != nil {
		return nil, err
	}
	return httptest.NewServer(api), nil
}
