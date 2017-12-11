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

package cli

import (
	"context"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/identity"
	"github.com/meta-network/go-meta/media"
	"github.com/meta-network/go-meta/testutil"
)

// TestMediaImportERN tests running the 'meta media import ern' command
func TestMediaImportERN(t *testing.T) {
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	srv := newTestServer(t, store)
	defer srv.Close()
	c := &testCLI{t}

	c.run("media", "import",
		"--url", srv.URL,
		"--source", "test",
		"ern",
		"../media/ern/testdata/Profile_AudioAlbumMusicOnly.xml",
		"../media/ern/testdata/Profile_AudioSingle.xml",
		"../media/ern/testdata/Profile_AudioSingle_WithCompoundArtistsAndTerritorialOverride.xml",
	)

	// check the index was populated
	index, err := store.OpenIndex("media.meta")
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()
	cmd := exec.Command("sqlite3", index.Path(), "SELECT COUNT(*) FROM release")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error checking index: %s: %s", err, out)
	}
	count := strings.TrimSpace(string(out))
	if count != "3" {
		t.Fatalf("expected 3 releases, got %d", count)
	}
}

func newTestServer(t *testing.T, store *meta.Store) *testServer {
	identityIndex, err := identity.NewIndex(store)
	if err != nil {
		t.Fatal(err)
	}
	mediaIndex, err := media.NewIndex(store)
	if err != nil {
		identityIndex.Close()
		t.Fatal(err)
	}
	srv, err := NewServer(store, identityIndex, mediaIndex)
	if err != nil {
		identityIndex.Close()
		mediaIndex.Close()
		t.Fatal(err)
	}
	return &testServer{httptest.NewServer(srv), identityIndex, mediaIndex}
}

type testServer struct {
	*httptest.Server

	identityIndex *identity.Index
	mediaIndex    *media.Index
}

func (t *testServer) Close() {
	t.Server.Close()
	t.identityIndex.Close()
	t.mediaIndex.Close()
}

type testCLI struct {
	t *testing.T
}

func (c *testCLI) run(args ...string) {
	if err := Run(context.Background(), args...); err != nil {
		c.t.Fatal(err)
	}
}
