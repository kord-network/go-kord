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

package cwr

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

// TestRegistedWorkAPI tests querying a registered work index via the GraphQL API.
func TestRegistedWorkAPI(t *testing.T) {
	// create a test index of registeredWorks
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

	// define a function to execute and assert an registerWork GraphQL query
	assertQuery := func(record *Record, query string, args ...interface{}) {
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
		var rw struct {
			RegisteredWorks []*RegisteredWork `json:"registered_work"`
		}
		if err := json.Unmarshal(r.Data, &rw); err != nil {
			t.Fatal(err)
		}

		if len(rw.RegisteredWorks) == 0 {
			t.Fatalf("expected registeredwork, got %d", len(rw.RegisteredWorks))
		}

		for i, r := range rw.RegisteredWorks {
			if r.Title != record.Title && i == len(rw.RegisteredWorks) {
				t.Fatalf("unexpected registeredwork title: expected %q ", record.Title)
			}
		}
	}

	for _, record := range x.records {
		// check getting the registeredWork title by iswc
		if record.RecordType == "NWR" ||
			record.RecordType == "REV" {
			if record.ISWC != "" {
				assertQuery(record, `{ registered_work(iswc:%q) { title } }`, record.ISWC)
			}
		}
	}
}

// TestPublisherControlAPI tests querying a publisher_control index via the GraphQL API.
func TestPublisherControlAPI(t *testing.T) {
	// create a test index of registeredWorks
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

	// define a function to execute and assert an record GraphQL query
	assertQuery := func(record *Record, query string, args ...interface{}) {
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
		var rw struct {
			PublisherControlledBySubmitters []*PublisherControllBySubmitter `json:"publisher_control"`
		}
		if err := json.Unmarshal(r.Data, &rw); err != nil {
			t.Fatal(err)
		}

		if len(rw.PublisherControlledBySubmitters) == 0 {
			t.Fatalf("expected spu, got %d", len(rw.PublisherControlledBySubmitters))
		}

		for i, r := range rw.PublisherControlledBySubmitters {
			if r.PublisherSequenceNumber != record.PublisherSequenceNumber && i == len(rw.PublisherControlledBySubmitters) {
				t.Fatalf("unexpected SPU sequenc number : expected %q ", record.PublisherSequenceNumber)
			}
		}
	}

	for _, record := range x.records {
		if record.RecordType == "SPU" {
			if record.PublisherSequenceNumber != "" {
				assertQuery(record, `{ publisher_control(publisher_sequence_n:%q) { publisher_sequence_n } }`, record.PublisherSequenceNumber)
			}
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
