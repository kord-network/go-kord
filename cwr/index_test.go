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
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	"github.com/meta-network/go-meta"
)

type testIndex struct {
	db      *sql.DB
	store   *meta.Store
	records []*Record
	tmpDir  string
}

func (t *testIndex) cleanup() {
	if t.db != nil {
		t.db.Close()
	}
	if t.tmpDir != "" {
		os.RemoveAll(t.tmpDir)
	}
}

func newTestIndex() (x *testIndex, err error) {
	x = &testIndex{}
	defer func() {
		if err != nil {
			x.cleanup()
		}
	}()
	cwrFileReader, err := os.Open("testdata/testfile.cwr")
	if err != nil {
		return nil, err
	}
	x.records, err = ParseCWRFile(cwrFileReader)
	if err != nil {
		return nil, err
	}

	// store the record in a test store
	x.store = meta.NewStore(datastore.NewMapDatastore())
	cids := make([]*cid.Cid, len(x.records))
	for i, record := range x.records {
		obj, err := meta.Encode(record)
		if err != nil {
			return nil, err
		}
		if err := x.store.Put(obj); err != nil {
			return nil, err
		}
		cids[i] = obj.Cid()
	}

	stream := make(chan *cid.Cid, len(x.records))
	go func() {
		defer close(stream)
		for _, cid := range cids {
			stream <- cid
		}
	}()

	// create a test SQLite3 db
	x.tmpDir, err = ioutil.TempDir("", "cwr-index-test")
	if err != nil {
		return nil, err
	}
	x.db, err = sql.Open("sqlite3", filepath.Join(x.tmpDir, "index.db"))
	if err != nil {
		return nil, err
	}

	// index the record
	indexer, err := NewIndexer(x.db, x.store)
	if err != nil {
		return nil, err
	}
	if err := indexer.Index(context.Background(), stream); err != nil {
		return nil, err
	}
	return x, nil

}

// TestIndexRegisteredWork tests indexing a stream of cwr regitered works.
func TestIndexRegisteredWorks(t *testing.T) {
	x, err := newTestIndex()
	if err != nil {
		t.Fatal(err)
	}
	defer x.cleanup()
	// check all the registeredWorks were indexed
	for _, record := range x.records {
		// check the title, iswc, composite_type indexes
		if !strings.HasPrefix(record.RecordType, "NWR") &&
			!strings.HasPrefix(record.RecordType, "REV") {
			continue
		}

		rows, err := x.db.Query(
			`SELECT object_id FROM registered_work WHERE title = ? AND iswc = ? AND composite_type = ? AND record_type = ?`,
			record.Title, record.ISWC, record.CompositeType, record.RecordType,
		)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		var objectID string
		for rows.Next() {
			// if we've already set objectID then we have a duplicate
			if objectID != "" {
				t.Fatalf("duplicate entries for registered work %q", record.Title)
			}
			if err := rows.Scan(&objectID); err != nil {
				t.Fatal(err)
			}

			// check we can get the object from the store
			cid, err := cid.Parse(objectID)
			if err != nil {
				t.Fatal(err)
			}
			obj, err := x.store.Get(cid)
			if err != nil {
				t.Fatal(err)
			}
			// check the object has the correct fields
			for key, expected := range map[string]string{
				"title":          record.Title,
				"iswc":           record.ISWC,
				"composite_type": record.CompositeType,
				"record_type":    record.RecordType,
			} {
				actual, err := obj.GetString(key)
				if err != nil {
					t.Fatal(err)
				}
				if actual != expected {
					t.Fatalf("expected object %s to be %q, got %q", key, expected, actual)
				}
			}
		}

		// check we got an object and no db errors
		if objectID == "" {
			t.Fatalf("registered work %q not found", record.Title)
		} else if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
	}
}

// TestIndexPublisherControl tests indexing a stream of cwr publisher control records.
func TestIndexPublisherControl(t *testing.T) {
	x, err := newTestIndex()
	if err != nil {
		t.Fatal(err)
	}
	defer x.cleanup()
	// check all the publisherControlledBySubmitter were indexed
	for _, record := range x.records {
		// check the publisher_sequence_number indexes
		if !strings.HasPrefix(record.RecordType, "SPU") {
			continue
		}
		rows, err := x.db.Query(
			`SELECT object_id FROM publisher_control WHERE publisher_sequence_n = ? AND record_type = ?`,
			record.PublisherSequenceNumber, record.RecordType,
		)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		var objectID string
		for rows.Next() {
			// if we've already set objectID then we have a duplicate
			if objectID != "" {
				t.Fatalf("duplicate entries for SPU publisher sequence number %q", record.PublisherSequenceNumber)
			}
			if err := rows.Scan(&objectID); err != nil {
				t.Fatal(err)
			}

			// check we can get the object from the store
			cid, err := cid.Parse(objectID)
			if err != nil {
				t.Fatal(err)
			}
			obj, err := x.store.Get(cid)
			if err != nil {
				t.Fatal(err)
			}
			// check the object has the correct fields
			for key, expected := range map[string]string{
				"publisher_sequence_n": record.PublisherSequenceNumber,
				"record_type":          record.RecordType,
			} {
				actual, err := obj.GetString(key)
				if err != nil {
					t.Fatal(err)
				}
				if actual != expected {
					t.Fatalf("expected object %s to be %q, got %q", key, expected, actual)
				}
			}
		}

		// check we got an object and no db errors
		if objectID == "" {
			t.Fatalf("SPU %q not found", record.PublisherSequenceNumber)
		} else if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
	}
}
