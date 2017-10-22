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
	"github.com/meta-network/go-meta/testutil"
	"github.com/neelance/graphql-go"
)

// TestAPI tests querying a MusicBrainz index via the GraphQL API.
func TestAPI(t *testing.T) {
	// create a test index
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	x, err := newTestIndex(store)
	if err != nil {
		t.Fatal(err)
	}
	defer x.cleanup()

	// start the API server
	s, err := newTestAPI(x.index.DB, x.store)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	// define a function to execute a GraphQL query
	query := func(q string, args ...interface{}) []byte {
		q = fmt.Sprintf(q, args...)
		data, _ := json.Marshal(map[string]string{"query": q})
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
		return r.Data
	}

	// define a function to assert the response of an artist GraphQL query
	assertArtists := func(res []byte, expected ...*Artist) {
		var a struct {
			Artists []*Artist `json:"artist"`
		}
		if err := json.Unmarshal(res, &a); err != nil {
			t.Fatal(err)
		}
		if len(a.Artists) != len(expected) {
			t.Fatalf("expected query to return %d artists, got %d", len(expected), len(a.Artists))
		}
		for i, artist := range expected {
			if a.Artists[i].Name != artist.Name {
				t.Fatalf("unexpected artist name: expected %q, got %q", artist.Name, a.Artists[i].Name)
			}
		}
	}

	for _, artist := range x.artists {
		// check getting the artist by name
		assertArtists(query(`{ artist(name:%q) { name } }`, artist.Name), artist)

		// check getting the artist by IPI returns just the artist,
		// with the exception of IPI "00435760746" which should
		// return two artists "Future" and "Lmars"
		for _, ipi := range artist.IPI {
			expected := []*Artist{artist}
			if ipi == "00435760746" {
				expected = []*Artist{
					{Name: "Future"},
					{Name: "Lmars"},
				}
			}
			assertArtists(query(`{ artist(ipi:%q) { name } }`, ipi), expected...)
		}

		// check getting the artist by ISNI
		for _, isni := range artist.ISNI {
			assertArtists(query(`{ artist(isni:%q) { name } }`, isni), artist)
		}
	}

	linkQuery := func(q string, args ...interface{}) []*RecordingWorkLink {
		res := query(q, args...)
		var l struct {
			Links []*RecordingWorkLink `json:"recording_work_link"`
		}
		if err := json.Unmarshal(res, &l); err != nil {
			t.Fatal(err)
		}
		return l.Links
	}

	// group the links by ISWC -> ISRC
	links := make(map[string][]string)
	for _, link := range x.links {
		links[link.ISWC] = append(links[link.ISWC], link.ISRC)
	}

	for iswc, isrcs := range links {
		// check getting the link by ISRC
		for _, isrc := range isrcs {
			res := linkQuery(`{ recording_work_link(isrc:%q) { iswc } }`, isrc)
			if len(res) != 1 {
				t.Fatalf("expected ISRC query to return one result, got %d", len(res))
			}
			if res[0].ISWC != iswc {
				t.Fatalf("expected ISWC %q, got %q", iswc, res[0].ISWC)
			}
		}

		// check getting the link by ISWC
		res := linkQuery(`{ recording_work_link(iswc:%q) { isrc } }`, iswc)
		if len(res) != len(isrcs) {
			t.Fatalf("expected %d ISRCs, got %d", len(isrcs), len(res))
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
