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
	"testing"

	cid "github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta/testutil"
	"github.com/meta-network/go-meta/testutil/index"
)

func TestIndexIdentity(t *testing.T) {
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, _ := testindex.GenerateIdentityIndex(t, ".", store)
	defer index.Close()

	var owner = "0x970e8128ab834e8eac17ab8e3812f010678cf791"
	rows, err := index.Query(`SELECT object_id FROM identity WHERE owner = ?`, owner)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var id string
	for rows.Next() {
		// if we've already set id then we have a duplicate
		if id != "" {
			t.Fatalf("duplicate entries for sender name %q", owner)
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

		v, err := obj.Get("owner")
		if err != nil {
			t.Fatal(err)
		}
		// check the object has the correct sender_id
		actual, ok := v.(string)
		if !ok {
			t.Fatalf("expected owner value to be string, got %T", v)
		}
		if actual != owner {
			t.Fatalf("expected owner value %q, got %q", owner, actual)
		}
	}

	//check we got a result and no db errors
	if id == "" {
		t.Fatalf("owner %q not found", owner)
	} else if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
}

func TestIndexClaim(t *testing.T) {
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, _ := testindex.GenerateClaimIndex(t, ".", store)
	defer index.Close()

	var value = "DPID_OF_THE_ARTIST_1"
	rows, err := index.Query(`SELECT object_id FROM claim WHERE signature = ? `, value)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var id string
	for rows.Next() {
		// if we've already set id then we have a duplicate
		if id != "" {
			t.Fatalf("duplicate entries for value %q", value)
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

		v, err := obj.Get("signature")
		if err != nil {
			t.Fatal(err)
		}
		// check the object has the correct sender_id
		actual, ok := v.(string)
		if !ok {
			t.Fatalf("expected signature value to be string, got %T", v)
		}
		if actual != value {
			t.Fatalf("expected value %q, got %q", value, actual)
		}
	}

	//check we got a result and no db errors
	if id == "" {
		t.Fatalf("value %q not found", value)
	} else if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
}
