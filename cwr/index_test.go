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
	"testing"

	"github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	"github.com/meta-network/go-meta"
)

type testIndex struct {
	db              *sql.DB
	store           *meta.Store
	registeredWorks []*RegisteredWork
	tmpDir          string
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
	x.registeredWorks, err = ParseCWRFile(cwrFileReader, "CWR-DataApi")
	if err != nil {
		return nil, err
	}

	// store the registeredWork in a test store
	x.store = meta.NewStore(datastore.NewMapDatastore())
	cids := make([]*cid.Cid, len(x.registeredWorks))
	for i, registerdWork := range x.registeredWorks {
		obj, err := meta.Encode(registerdWork)
		if err != nil {
			return nil, err
		}
		if err := x.store.Put(obj); err != nil {
			return nil, err
		}
		cids[i] = obj.Cid()
	}

	// create a stream
	stream := make(chan *cid.Cid, len(x.registeredWorks))
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

	// index the RegisteredWork
	indexer, err := NewIndexer(x.db, x.store)
	if err != nil {
		return nil, err
	}
	if err := indexer.IndexRegisteredWorks(context.Background(), stream); err != nil {
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
	for _, registeredWork := range x.registeredWorks {
		// check the title, iswc, composite_type indexes
		rows, err := x.db.Query(
			`SELECT object_id FROM registered_work WHERE title = ? AND iswc = ? AND composite_type = ? AND record_type = ?`,
			registeredWork.Title, registeredWork.ISWC, registeredWork.CompositeType, registeredWork.RecordType,
		)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		var objectID string
		for rows.Next() {
			// if we've already set objectID then we have a duplicate
			if objectID != "" {
				t.Fatalf("duplicate entries for registered work %q", registeredWork.Title)
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
				"title":          registeredWork.Title,
				"iswc":           registeredWork.ISWC,
				"composite_type": registeredWork.CompositeType,
				"record_type":    registeredWork.RecordType,
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
			t.Fatalf("registered work %q not found", registeredWork.Title)
		} else if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
	}
}
