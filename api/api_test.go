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
	metagraph "github.com/meta-network/go-meta/graph"
	"github.com/meta-network/go-meta/testutil"
)

func TestAPI(t *testing.T) {
	dpa, err := testutil.NewTestDPA()
	if err != nil {
		t.Fatal(err)
	}
	defer dpa.Cleanup()
	driver := metagraph.NewDriver("meta", dpa.DPA, testutil.NewTestENS(), dpa.Dir)

	// start server
	srv := httptest.NewServer(NewServer(driver))

	// create a graph
	name := "test.meta"
	client := NewClient(srv.URL, name)
	if err := client.Create(); err != nil {
		t.Fatal(err)
	}

	// write some quads
	qw, err := graph.NewQuadWriter("single", client, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := qw.AddQuad(quad.Make("phrase of the day", "is of course", "Hello World!", nil)); err != nil {
		t.Fatal(err)
	}

	// check the quads were added
	qs, err := driver.Get(name)
	if err != nil {
		t.Fatal(err)
	}
	path := cayley.StartPath(qs, quad.String("phrase of the day")).Out(quad.String("is of course"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	results, err := path.Iterate(ctx).All()
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	v := qs.NameOf(results[0])
	s, ok := v.Native().(string)
	if !ok {
		t.Fatalf("expected string, got %T", v.Native())
	}
	if s != "Hello World!" {
		t.Fatalf(`expected "Hello World!", got %q`, s)
	}
}
