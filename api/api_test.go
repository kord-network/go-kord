// This file is part of the go-meta library.
//
// Copyright (C) 2018 JAAK MUSIC LTD
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

package api

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	meta "github.com/meta-network/go-meta"
	metasql "github.com/meta-network/go-meta/sql"
	"github.com/meta-network/go-meta/testutil"
)

func TestAPI(t *testing.T) {
	dpa, err := testutil.NewTestDPA()
	if err != nil {
		t.Fatal(err)
	}
	defer dpa.Cleanup()

	storage := meta.NewStorage(dpa.Dir, dpa.DPA, &testutil.ENS{})

	metasql.Register(storage)

	state := meta.NewState(storage)

	signer, err := testutil.NewTestSigner()
	if err != nil {
		t.Fatal(err)
	}

	// create store
	srv := httptest.NewServer(NewServer(state))
	client := NewClient(srv.URL)
	store, err := meta.NewQuadStore(
		signer.Address,
		"test.meta",
		client,
		signer,
		storage,
	)
	if err != nil {
		t.Fatal(err)
	}

	qw, err := graph.NewQuadWriter("single", store, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := qw.AddQuad(quad.Make("phrase of the day", "is of course", "Hello World!", nil)); err != nil {
		t.Fatal(err)
	}

	path := cayley.StartPath(store, quad.String("phrase of the day")).Out(quad.String("is of course"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	results, err := path.Iterate(ctx).All()
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	v := store.NameOf(results[0])
	s, ok := v.Native().(string)
	if !ok {
		t.Fatalf("expected string, got %T", v.Native())
	}
	if s != "Hello World!" {
		t.Fatalf(`expected "Hello World!", got %q`, s)
	}
}
