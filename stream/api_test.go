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

package stream

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/meta-network/go-meta/testutil"
)

// TestAPI tests reading and writing a META stream via the API.
func TestAPI(t *testing.T) {
	// start the stream API
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	api := NewAPI(store)
	srv := httptest.NewServer(api)
	defer srv.Close()

	// write some values to a stream
	client := NewClient(srv.URL)
	name := "test.meta"
	count := 5
	type obj struct {
		Value int `json:"value"`
	}
	for i := 0; i < count; i++ {
		if _, err := client.WriteStream(name, &obj{i}); err != nil {
			t.Fatal(err)
		}
	}

	// check reading the stream returns the values
	ch := make(chan *obj)
	stream, err := client.ReadStream(name, ch)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	timeout := time.After(10 * time.Second)
	for i := 0; i < count; i++ {
		select {
		case obj, ok := <-ch:
			if !ok {
				t.Fatalf("stream closed unexpectedly: %s", stream.Err())
			}
			if obj.Value != i {
				t.Fatalf("expected object to have value %d, got %d", i, obj.Value)
			}
		case <-timeout:
			t.Fatal("timed out waiting for stream values")
		}
	}
}
