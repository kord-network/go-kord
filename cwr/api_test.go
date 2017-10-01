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
	"strconv"
	"testing"

	cid "github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
	"github.com/neelance/graphql-go"
)

// TestRegisteredWorkAPI tests querying a registered work(NWR/REV) transacation's records index via the GraphQL API.
func TestRegisteredWorkAPI(t *testing.T) {
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
	assertQueryNWR := func(record *Record, query string, args ...interface{}) error {
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
			RegisteredWorks []*RegisteredWork `json:"registered_work"`
		}
		if err := json.Unmarshal(r.Data, &rw); err != nil {
			return err
		}

		if len(rw.RegisteredWorks) == 0 {
			return fmt.Errorf("expected registeredwork, got %d", len(rw.RegisteredWorks))
		}

		for i, r := range rw.RegisteredWorks {
			if r.Title != record.Title && i == len(rw.RegisteredWorks) {
				return fmt.Errorf("unexpected registeredwork title: expected %q ", record.Title)
			}
		}
		return nil
	}
	// define a function to execute and assert an record GraphQL query
	assertQuerySPU := func(record *Record, query string, args ...interface{}) error {
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
			PublisherControlledBySubmitters []*PublisherControllBySubmitter `json:"publisher_control"`
		}
		if err := json.Unmarshal(r.Data, &rw); err != nil {
			return err
		}

		if len(rw.PublisherControlledBySubmitters) == 0 {
			return fmt.Errorf("expected spu, got %d", len(rw.PublisherControlledBySubmitters))
		}

		for i, r := range rw.PublisherControlledBySubmitters {
			if r.PublisherSequenceNumber != record.PublisherSequenceNumber && i == len(rw.PublisherControlledBySubmitters) {
				return fmt.Errorf("unexpected SPU sequenc number : expected %q ", record.PublisherSequenceNumber)
			}
		}
		return nil
	}
	if err := testTxRecords(x, assertQueryNWR, assertQuerySPU); err != nil {
		t.Fatal(err)
	}
}

//testTxRecords get all cwr NWR/REV transactions and assert queries for its
//records
func testTxRecords(x *testIndex,
	assertQueryNWR func(record *Record, query string, args ...interface{}) error,
	assertQuerySPU func(record *Record, query string, args ...interface{}) error,
) (err error) {

	cwrObj, err := x.store.Get(x.cwrCid)
	if err != nil {
		return err
	}
	graph := meta.NewGraph(x.store, cwrObj)

	v, err := graph.Get("GRH")
	if err != nil {
		return err
	}
	numberOfGroups := len(v.([]interface{}))
	record := &Record{}
	for k := 0; k < numberOfGroups; k++ {
		v, err := graph.Get("GRH", strconv.Itoa(k), "NWR")
		if meta.IsPathNotFound(err) {
			continue
		} else if err != nil {
			return err
		}
		numberOfTx := len(v.([]interface{}))
		for j := 0; j < numberOfTx; j++ {
			v, err := graph.Get("GRH", strconv.Itoa(k), "NWR", strconv.Itoa(j))
			if meta.IsPathNotFound(err) {
				continue
			} else if err != nil {
				return err
			}
			tx, ok := v.(map[string]interface{})
			if !ok {
				return fmt.Errorf("unexpected field type for %q, expected *cid.Cid, got %T", "NWR", v)
			}
			nwrCid, ok := tx["NWR"].(*cid.Cid)
			if !ok {
				nwrCid, ok = tx["REV"].(*cid.Cid)
				if !ok {
					return err
				}
			}
			obj, err := x.store.Get(nwrCid)
			if err != nil {
				return err
			}

			if err := obj.Decode(record); err != nil {
				return err
			}
			if record.ISWC != "" {
				if err := assertQueryNWR(record, `{ registered_work(iswc:%q) { title } }`, record.ISWC); err != nil {
					return err
				}
			}

			for _, spuCid := range tx["SPU"].([]interface{}) {
				obj, err := x.store.Get(spuCid.(*cid.Cid))
				if err != nil {
					return err
				}
				if err := obj.Decode(record); err != nil {
					return err
				}
				if record.RecordType == "SPU" {
					if record.PublisherSequenceNumber != "" {
						if err := assertQuerySPU(record, `{ publisher_control(publisher_sequence_n:%q) { publisher_sequence_n } }`, record.PublisherSequenceNumber); err != nil {
							return err
						}
					}
				}
			}

		}
	}
	return nil
}

func newTestAPI(db *sql.DB, store *meta.Store) (*httptest.Server, error) {
	api, err := NewAPI(db, store)
	if err != nil {
		return nil, err
	}
	return httptest.NewServer(api), nil
}
