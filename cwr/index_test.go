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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	cid "github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

type testIndex struct {
	db     *sql.DB
	store  *meta.Store
	cwrCid *cid.Cid
	tmpDir string
}

func (t *testIndex) cleanup() {
	if t.db != nil {
		t.db.Close()
	}
	if t.tmpDir != "" {
		os.RemoveAll(t.tmpDir)
	}
}

func TestIndex(t *testing.T) {
	x, err := newTestIndex()
	if err != nil {
		t.Fatal(err)
	}
	defer x.cleanup()

	// check the HDR record was indexed into the transmission_header table
	senderName := "JAAK EXAMPLE SENDER NAME"
	rows, err := x.db.Query(`SELECT object_id FROM transmission_header WHERE sender_name = ?`, senderName)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var id string
	for rows.Next() {
		// if we've already set id then we have a duplicate
		if id != "" {
			t.Fatalf("duplicate entries for sender name %q", senderName)
		}
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}

		// check we can get the object from the store
		cid, err := cid.Parse(id)
		if err != nil {
			t.Fatal(err)
		}
		obj, err := x.store.Get(cid)
		if err != nil {
			t.Fatal(err)
		}

		v, err := obj.Get("sender_name")
		if err != nil {
			t.Fatal(err)
		}
		// check the object has the correct sender_id
		actual, ok := v.(string)
		if !ok {
			t.Fatalf("expected sender name value to be string, got %T", v)
		}
		if actual != senderName {
			t.Fatalf("expected sender name value %q, got %q", senderName, actual)
		}
	}

	//check we got a result and no db errors
	if id == "" {
		t.Fatalf("senderName %q not found", senderName)
	} else if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	check := func(id string, fieldsMap map[string]string) error {

		txID, err := cid.Decode(id)
		if err != nil {
			return err
		}
		obj, err := x.store.Get(txID)
		if err != nil {
			return err
		}
		for field, actual := range fieldsMap {
			expected, err := obj.Get(field)
			if err != nil {
				return err
			}
			if expected != actual {
				return fmt.Errorf("expected %s to be %q, got %q", field, expected, actual)
			}
		}
		return nil
	}
	//check all NWR/REV transactions were index into the registered_work table
	var (
		objectID           string
		title              string
		iswc               string
		publisherSequenceN string
		writerFirstName    string
	)

	rows, err = x.db.Query(`SELECT object_id,title,iswc FROM registered_work WHERE cwr_id = ?`, x.cwrCid.String())
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	rowsCount := 0
	for rows.Next() {
		rowsCount++
		if err = rows.Scan(&objectID, &title, &iswc); err != nil {
			t.Fatal(err)
		}
		if err = check(objectID, map[string]string{
			"title": title,
			"iswc":  iswc,
		}); err != nil {
			t.Fatal(err)
		}
		txID, err := cid.Decode(objectID)
		if err != nil {
			t.Fatal(err)
		}
		//get all publisher control records (SPU) which link to the above tx .
		spuRows, err := x.db.Query(`SELECT object_id,publisher_sequence_n FROM publisher_control WHERE tx_id = ?`, txID.String())
		if err != nil {
			t.Fatal(err)
		}
		defer spuRows.Close()
		for spuRows.Next() {

			if err := spuRows.Scan(&objectID, &publisherSequenceN); err != nil {
				t.Fatal(err)
			}
			if err = check(objectID, map[string]string{
				"publisher_sequence_n": publisherSequenceN,
			}); err != nil {
				t.Fatal(err)
			}
		}
		//get all writer control records (SWR or OWR) which link to the above tx .
		swrRows, err := x.db.Query(`SELECT object_id,writer_first_name FROM writer_control WHERE tx_id = ?`, txID.String())
		if err != nil {
			t.Fatal(err)
		}
		defer swrRows.Close()
		for swrRows.Next() {

			if err := swrRows.Scan(&objectID, &writerFirstName); err != nil {
				t.Fatal(err)
			}
			if err = check(objectID, map[string]string{
				"writer_first_name": writerFirstName,
			}); err != nil {
				t.Fatal(err)
			}
		}
	}
	if rowsCount == 0 {
		t.Fatal("no registered_work found")
	}
}

func newTestIndex() (x *testIndex, err error) {
	// convert the test cwr to META object
	x = &testIndex{}
	defer func() {
		if err != nil {
			x.cleanup()
		}
	}()

	x.tmpDir, err = ioutil.TempDir("", "cwr-index-test")
	if err != nil {
		return nil, err
	}

	x.store = meta.NewMapDatastore()

	converter := NewConverter(x.store)

	f, err := os.Open(filepath.Join("testdata", "example_nwr.cwr"))

	if err != nil {
		return nil, err
	}
	defer f.Close()

	x.cwrCid, err = converter.ConvertCWR(f, "test")
	if err != nil {
		return nil, err
	}

	// create a stream of CWR
	stream := make(chan *cid.Cid)
	go func() {
		defer close(stream)
		stream <- x.cwrCid
	}()

	// create a test SQLite3 db
	x.db, err = sql.Open("sqlite3", filepath.Join(x.tmpDir, "index.db"))
	if err != nil {
		return nil, err
	}
	// index the stream of CWR txs
	indexer, err := NewIndexer(x.db, x.store)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := indexer.Index(ctx, stream); err != nil {
		return nil, err
	}
	return x, nil
}
