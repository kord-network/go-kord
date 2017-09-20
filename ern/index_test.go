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

package ern

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/meta-network/go-meta"
	"golang.org/x/net/context"
)

func TestIndex(t *testing.T) {
	// convert the test ERNs to META objects
	erns := []string{
		"Profile_AudioAlbumMusicOnly.xml",
		"Profile_AudioSingle.xml",
		"Profile_AudioAlbum_WithBooklet.xml",
		"Profile_AudioSingle_WithCompoundArtistsAndTerritorialOverride.xml",
		"Profile_AudioBook.xml",
	}
	store := meta.NewStore(datastore.NewMapDatastore())
	converter := NewConverter(store)
	cids := make(map[string]*cid.Cid, len(erns))
	for _, path := range erns {
		f, err := os.Open(filepath.Join("testdata", path))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		cid, err := converter.ConvertERN(f)
		if err != nil {
			t.Fatal(err)
		}
		cids[path] = cid
	}

	// create a stream of ERNs
	stream := make(chan *cid.Cid, len(erns))
	go func() {
		defer close(stream)
		for _, cid := range cids {
			stream <- cid
		}
	}()

	// create a test SQLite3 db
	tmpDir, err := ioutil.TempDir("", "musicbrainz-index-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	db, err := sql.Open("sqlite3", filepath.Join(tmpDir, "index.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// index the stream of ERNs
	indexer, err := NewIndexer(db, store)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := indexer.Index(ctx, stream); err != nil {
		t.Fatal(err)
	}

	// check the MessageSender and MessageRecipient were indexed into the
	// party table
	for _, partyID := range []string{"DPID_OF_THE_SENDER", "DPID_OF_THE_RECIPIENT"} {
		rows, err := db.Query(`SELECT cid FROM party WHERE id = ?`, partyID)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		var id string
		for rows.Next() {
			// if we've already set id then we have a duplicate
			if id != "" {
				t.Fatalf("duplicate entries for PartyId %q", partyID)
			}
			if err := rows.Scan(&id); err != nil {
				t.Fatal(err)
			}

			// check we can get the object from the store
			cid, err := cid.Parse(id)
			if err != nil {
				t.Fatal(err)
			}
			obj, err := store.Get(cid)
			if err != nil {
				t.Fatal(err)
			}

			// check the object has the correct PartyId
			graph := meta.NewGraph(store, obj)
			v, err := graph.Get("PartyId", "@value")
			if err != nil {
				t.Fatal(err)
			}
			actual, ok := v.(string)
			if !ok {
				t.Fatalf("expected PartyId value to be string, got %T", v)
			}
			if actual != partyID {
				t.Fatalf("expected PartyId value %q, got %q", partyID, actual)
			}
		}

		// check we got a result and no db errors
		if id == "" {
			t.Fatalf("party %q not found", partyID)
		} else if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
	}

	// check all the ERNs were indexed into the ern table
	for _, cid := range cids {
		obj, err := store.Get(cid)
		if err != nil {
			t.Fatal(err)
		}
		graph := meta.NewGraph(store, obj)
		var (
			messageID string
			threadID  string
		)
		row := db.QueryRow(`SELECT message_id, thread_id FROM ern WHERE cid = ?`, cid.String())
		if err := row.Scan(&messageID, &threadID); err != nil {
			t.Fatal(err)
		}
		for field, actual := range map[string]string{
			"MessageId":       messageID,
			"MessageThreadId": threadID,
		} {
			v, err := graph.Get("NewReleaseMessage", "MessageHeader", field, "@value")
			if err != nil {
				t.Fatal(err)
			}
			expected, ok := v.(string)
			if !ok {
				t.Fatalf("expected %s to be a string, got %T", field, v)
			}
			if actual != expected {
				t.Fatalf("expected %s to be %q, got %q", field, expected, actual)
			}
		}
	}

	// check SoundRecording objects were indexed
	for isrc, title := range map[string]string{
		"CASE00000001": "Can you feel ...the Monkey Claw!",
		"CASE00000002": "Red top mountain, blown sky high",
		"CASE00000003": "Seige of Antioch",
		"CASE00000004": "Warhammer",
		"CASE00000005": "Iron Horse",
		"CASE00000006": "Yes... I can feel the Monkey Claw!",
	} {
		var id string
		row := db.QueryRow("SELECT cid FROM sound_recording WHERE id = ? AND title = ?", isrc, title)
		if err := row.Scan(&id); err != nil {
			t.Fatal(err)
		}
		var ernID string
		row = db.QueryRow("SELECT ern_id FROM resource_list WHERE resource_id = ?", id)
		if err := row.Scan(&ernID); err != nil {
			t.Fatal(err)
		}
	}
}
