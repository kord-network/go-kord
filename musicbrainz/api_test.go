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

package musicbrainz

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meta-network/go-meta"
	"github.com/neelance/graphql-go"
)

// TestArtistAPI tests querying an artist index via the GraphQL API.
func TestArtistAPI(t *testing.T) {
	// create a test index of artists
	x, err := newTestIndex()
	if err != nil {
		t.Fatal(err)
	}
	defer x.cleanup()

	// start the API server
	s, err := newTestAPI(x.db, x.store)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	// define a function to execute and assert an artist GraphQL query
	assertQuery := func(artist *Artist, query string, args ...interface{}) {
		data, _ := json.Marshal(map[string]string{"query": fmt.Sprintf(query, args...)})
		req, err := http.NewRequest("POST", s.URL+"/graphql", bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("unexpected HTTP status: %s", res.Status)
		}
		var r graphql.Response
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			t.Fatal(err)
		}
		if len(r.Errors) > 0 {
			t.Fatalf("unexpected errors in API response: %v", r.Errors)
		}
		var a struct {
			Artists []*Artist `json:"artist"`
		}
		if err := json.Unmarshal(r.Data, &a); err != nil {
			t.Fatal(err)
		}
		if len(a.Artists) != 1 {
			t.Fatalf("expected 1 artist, got %d", len(a.Artists))
		}
		if a.Artists[0].Name != artist.Name {
			t.Fatalf("unexpected artist name: expected %q, got %q", artist.Name, a.Artists[0].Name)
		}
	}

	for _, artist := range x.artists {
		// check getting the artist by name
		assertQuery(artist, `{ artist(name:%q) { name } }`, artist.Name)

		// check getting the artist by IPI
		for _, ipi := range artist.IPI {
			assertQuery(artist, `{ artist(ipi:%q) { name } }`, ipi)
		}

		// check getting the artist by ISNI
		for _, isni := range artist.ISNI {
			assertQuery(artist, `{ artist(isni:%q) { name } }`, isni)
		}
	}
}

func newTestAPI(db *sql.DB, store *meta.Store) (*httptest.Server, error) {
	api, err := NewAPI(db, store)
	if err != nil {
		return nil, err
	}
	return httptest.NewServer(api), nil
}