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
	"bufio"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/meta-network/go-meta"
)

// TestCWRCommands tests running the 'meta cwr convert' and
// 'meta cwr index' commands.
func TestCWRCommands(t *testing.T) {
	c := newTestCLI(t)

	// check 'meta cwr convert' prints a CID
	stdout := c.run("cwr", "convert", "../cwr/testdata/testfile.cwr", "../cwr/CWR-DataApi")
	id, err := cid.Parse(strings.TrimSpace(stdout))
	if err != nil {
		t.Fatal(err)
	}
	expected := "zdpuAxGP3BcFkXyNp1y59FEdMZmqbgxkLdvFynFrpdQ1jpgK3"
	if id.String() != expected {
		t.Fatalf("unexpected CID, expected %q, got %q", expected, id)
	}

	// create a path to store the index
	tmpDir, err := ioutil.TempDir("", "meta-main-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	db := filepath.Join(tmpDir, "index.db")

	// run 'meta cwr index' with the CID as stdin
	stream := strings.NewReader(stdout)
	c.runWithStdin(stream, "cwr", "index", db)

	// check the index was populated
	cmd := exec.Command("sqlite3", db, "SELECT object_id FROM registered_work")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error checking index: %s: %s", err, out)
	}
	if strings.TrimSpace(string(out)) != expected {
		t.Fatalf("unexpected index output, expected %q, got %q", expected, out)
	}
}

// TestERNCommands tests running the 'meta ern convert' and
// 'meta ern index' commands.
func TestERNCommands(t *testing.T) {
	c := newTestCLI(t)

	// check 'meta ern convert' prints multiple CIDs
	stdout := c.run("ern", "convert",
		"../ern/testdata/Profile_AudioAlbumMusicOnly.xml",
		"../ern/testdata/Profile_AudioAlbum_WithBooklet.xml",
		"../ern/testdata/Profile_AudioBook.xml",
		"../ern/testdata/Profile_AudioSingle.xml",
		"../ern/testdata/Profile_AudioSingle_WithCompoundArtistsAndTerritorialOverride.xml",
	)
	var ids []string
	s := bufio.NewScanner(strings.NewReader(stdout))
	for s.Scan() {
		id, err := cid.Parse(s.Text())
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id.String())
	}
	expected := []string{
		"zdpuB3GGepjkoZLLBTPLAkTvzvLw2zSNDUFtSo7NyogUAbprQ",
		"zdpuAn1okygE1mZxNQbvGqtHPBWpEzfsrgUTuMjvCaBtnAQoN",
		"zdpuAxnsxCEpnNk6mQQXmdAJvUN95dX5WENJXmMuAgwCQ4aKb",
		"zdpuAxVUi1d2eqfPhxJzGjrJEUf1Bbr8JFPXgPFNYXNMHcaWj",
		"zdpuAq11sKUZFqvDYfZhUwvhymofqAbz7cQ5VyBZYpVv453XR",
	}
	if !reflect.DeepEqual(ids, expected) {
		t.Fatalf("unexpected CIDs:\nexpected: %v\ngot:      %v", expected, ids)
	}

	// create a path to store the index
	tmpDir, err := ioutil.TempDir("", "meta-main-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	db := filepath.Join(tmpDir, "index.db")

	// run 'meta ern index' with the CIDs as stdin
	stream := strings.NewReader(stdout)
	c.runWithStdin(stream, "ern", "index", db)

	// check the index was populated
	cmd := exec.Command("sqlite3", db, "SELECT cid FROM ern ORDER BY cid")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error checking index: %s: %s", err, out)
	}
	gotIDs := strings.Split(strings.TrimSpace(string(out)), "\n")
	sort.Strings(expected)
	if !reflect.DeepEqual(gotIDs, expected) {
		t.Fatalf("unexpected index output:\nexpected: %v\ngot:      %q", expected, gotIDs)
	}
}

type testCLI struct {
	t     *testing.T
	store *meta.Store
}

func newTestCLI(t *testing.T) *testCLI {
	return &testCLI{
		t:     t,
		store: meta.NewStore(datastore.NewMapDatastore()),
	}
}

func (c *testCLI) runWithStdin(stdin io.Reader, args ...string) string {
	var stdout bytes.Buffer
	cli := New(c.store, stdin, &stdout)
	if err := cli.Run(context.Background(), args...); err != nil {
		c.t.Fatal(err)
	}
	return stdout.String()
}

func (c *testCLI) run(args ...string) string {
	return c.runWithStdin(nil, args...)
}
