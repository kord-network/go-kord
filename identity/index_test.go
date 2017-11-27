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

	"github.com/meta-network/go-meta/testutil"
	"github.com/meta-network/go-meta/testutil/index"
)

func TestIndexIdentity(t *testing.T) {
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, testid := testindex.GenerateIdentityIndex(t, ".", store)
	defer index.Close()

	rows, err := index.Query(`SELECT id FROM identity WHERE owner = ?`, testid.Owner.String())
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var id string
	for rows.Next() {
		// if we've already set id then we have a duplicate
		if id != "" {
			t.Fatalf("duplicate entries for owner %q", testid.Owner.String())
		}
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
		if id != testid.ID {
			t.Fatalf("expected id value %q, got %q", testid.ID, id)
		}
	}

	//check we got a result and no db errors
	if id == "" {
		t.Fatalf("owner %q not found", testid.Owner.String())
	} else if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
}

func TestIndexClaim(t *testing.T) {
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, claims := testindex.GenerateClaimIndex(t, ".", store)
	defer index.Close()

	rows, err := index.Query(`SELECT subject FROM claim WHERE signature = ? `, claims[0].Signature)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var subject string
	for rows.Next() {
		// if we've already set id then we have a duplicate
		if subject != "" {
			t.Fatalf("duplicate entries for signature %q", claims[0].Signature)
		}
		if err := rows.Scan(&subject); err != nil {
			t.Fatal(err)
		}

		if subject != claims[0].Subject {
			t.Fatalf("expected value %q, got %q", claims[0].Subject, subject)
		}
	}

	//check we got a result and no db errors
	if subject == "" {
		t.Fatalf("claim subject %q not found", claims[0].Subject)
	} else if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
}
