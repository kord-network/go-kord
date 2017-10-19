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
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/testutil"
)

// TestCWRCommands tests running the 'meta convert cwr' and
// 'meta index cwr' commands.
func TestCWRCommands(t *testing.T) {
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	c := &testCLI{t, store}

	// check 'meta convert cwr' adds CIDs to the cwr.meta stream
	c.run("convert", "cwr",
		"--source", "test",
		"../cwr/testdata/example_double_nwr.cwr",
		"../cwr/testdata/example_nwr.cwr")
	ids := c.readStream("cwr.meta", 2)
	expected := []string{
		"zdqaWBuxwxhZQj9PBzRsHSp2WEq9pF3tF9rP9KYb7TxgqfXQJ",
		"zdqaWGLaDAkMomHFaooZ7GaxxvTmFCJ1DFCVLaQSvb8v2ewoN",
	}
	if !reflect.DeepEqual(ids, expected) {
		t.Fatalf("unexpected CIDs:\nexpected: %v\ngot:      %v", expected, ids)
	}

	// run 'meta index cwr'
	indexName := "cwr.test"
	c.run("index", "cwr", indexName, "--count=2")

	// check the index was populated
	index, err := store.OpenIndex(indexName)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sqlite3", index.Path(), "SELECT cwr_id FROM transmission_header ORDER BY cwr_id")
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

// TestERNCommands tests running the 'meta convert ern' and
// 'meta index ern' commands.
func TestERNCommands(t *testing.T) {
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	c := &testCLI{t, store}

	// check 'meta convert ern' adds CIDs to the ern.meta stream
	c.run("convert", "ern",
		"--source", "test",
		"../ern/testdata/Profile_AudioAlbumMusicOnly.xml",
		"../ern/testdata/Profile_AudioAlbum_WithBooklet.xml",
		"../ern/testdata/Profile_AudioBook.xml",
		"../ern/testdata/Profile_AudioSingle.xml",
		"../ern/testdata/Profile_AudioSingle_WithCompoundArtistsAndTerritorialOverride.xml",
	)
	ids := c.readStream("ern.meta", 5)
	expected := []string{
		"zdqaWMMSZgBtPNTUJj6s9GxWdaNLiF8oEGLDc2uYqG7T9nmqL",
		"zdqaWJ4jkU4haHkCfJY8Tz7bNtdi38Kq1bRy4iiU9DUDJqUkB",
		"zdqaWRL5oXkGnwhf9JWzgMdV7H1YUWp6YociP5aWB9EU6qAJr",
		"zdqaWQnNQ9oTMxVnDWWKagJWAwcfQQuirAmgYAPc9CPnRsFR9",
		"zdqaWBoE9GfAc1dsJgW1jqtbHYKHgLcFfw55SKDyxEEGcQWso",
	}
	if !reflect.DeepEqual(ids, expected) {
		t.Fatalf("unexpected CIDs:\nexpected: %v\ngot:      %v", expected, ids)
	}

	// run 'meta index ern'
	indexName := "ern.test"
	c.run("index", "ern", indexName, "--count=5")

	// check the index was populated
	index, err := store.OpenIndex(indexName)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sqlite3", index.Path(), "SELECT cid FROM ern ORDER BY cid")
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

// TestERNCommands tests running the 'meta convert eidr' and
// 'meta index eidr' commands.
func TestEIDRCommands(t *testing.T) {
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	c := &testCLI{t, store}

	// check 'meta convert eidr' adds CIDs to the eidr.meta stream
	c.run("convert", "eidr",
		"--source", "test",
		"../eidr/testdata/dummy_child.xml",
		"../eidr/testdata/dummy_parent.xml",
	)
	ids := c.readStream("eidr.meta", 2)
	expected := []string{
		"zdqaWUgTjLPoFrMCmBcbPcJNPuHTvUs5fBn3f4iYsniHACo7q",
		"zdqaWTaGbQFm2HEYo2iwHsFanffKDYq49XVk5ggEc6dYaBQkv",
	}
	if !reflect.DeepEqual(ids, expected) {
		t.Fatalf("unexpected CIDs:\nexpected: %v\ngot:      %v", expected, ids)
	}

	// run 'meta index eidr'
	indexName := "eidr.test"
	c.run("index", "eidr", indexName, "--count=2")

	// check the index was populated
	index, err := store.OpenIndex(indexName)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sqlite3", index.Path(), "select count(*) from xobject_baseobject_link x inner join baseobject p, xobject_episode e on p.doi_id = x.parent_doi_id where e.id = x.xobject_id")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error checking index: %s: %s", err, out)
	}
	if strings.TrimSpace(string(out)) != "1" {
		t.Fatalf("baseobject/xobject link count mismatch; expected 1, got %s", out)
	}

	// check if associatedorgs are inserted and linked
	cmd = exec.Command("sqlite3", index.Path(), "select count(*) from org o inner join baseobject b on o.base_doi_id = b.doi_id")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error checking index: %s: %s", err, out)
	}
	if strings.TrimSpace(string(out)) != "2" {
		t.Fatalf("associatedorg link count mismatch; expected 2, got %s", out)
	}

	// check if alternateids are inserted and linked
	cmd = exec.Command("sqlite3", index.Path(), "select count(*) from alternateid a inner join baseobject b on a.base_doi_id = b.doi_id")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error checking index: %s: %s", err, out)
	}
	if strings.TrimSpace(string(out)) != "2" {
		t.Fatalf("alternateid link count mismatch; expected 2, got %s", out)
	}
}

type testCLI struct {
	t     *testing.T
	store *meta.Store
}

func (c *testCLI) run(args ...string) {
	cli := New(c.store, nil, os.Stdout)
	if err := cli.Run(context.Background(), args...); err != nil {
		c.t.Fatal(err)
	}
}

func (c *testCLI) readStream(name string, count int) []string {
	reader, err := c.store.StreamReader(name, meta.StreamLimit(count))
	if err != nil {
		c.t.Fatal(err)
	}
	defer reader.Close()
	var ids []string
	timeout := time.After(10 * time.Second)
	for n := 0; n < count; n++ {
		select {
		case id, ok := <-reader.Ch():
			if !ok {
				c.t.Fatal(reader.Err())
			}
			ids = append(ids, id.String())
		case <-timeout:
			c.t.Fatalf("timed out waiting for stream values")
		}
	}
	return ids
}
