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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
	"github.com/ethereum/go-ethereum/swarm/testutil"
	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// TestCWRCommands tests running the 'meta convert cwr' and
// 'meta index cwr' commands.
func TestCWRCommands(t *testing.T) {
	c, err := newTestCLI(t)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(c.bzz.indexDir)
	defer c.srv.Close()
	// check 'meta convert cwr' prints a CID
	stdout := c.run("convert", "cwr",
		"--source", "test",
		"../cwr/testdata/example_double_nwr.cwr",
		"../cwr/testdata/example_nwr.cwr")
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
		"zdqaWBuxwxhZQj9PBzRsHSp2WEq9pF3tF9rP9KYb7TxgqfXQJ",
		"zdqaWGLaDAkMomHFaooZ7GaxxvTmFCJ1DFCVLaQSvb8v2ewoN",
	}
	if !reflect.DeepEqual(ids, expected) {
		t.Fatalf("unexpected CIDs:\nexpected: %v\ngot:      %v", expected, ids)
	}

	db := filepath.Join(c.bzz.indexDir, "1_index.db")

	// run 'meta index cwr' with the CIDs as stdin
	stream := strings.NewReader(stdout)
	c.bzz.indexDirHash = c.runWithStdin(stream, "index", "cwr", filepath.Base(db), fmt.Sprintf("--bzzapi=%s", c.srv.Server.URL), fmt.Sprintf("--bzzdir=%s", c.bzz.indexDirHash))
	if c.bzz.indexDirHash == "" {
		t.Fatal("No hash returned")
	}
	t.Logf("got hash", c.bzz.indexDirHash)

	// check the index was populated
	c.bzz.GetIndexFile(filepath.Base(db), true)
	cmd := exec.Command("sqlite3", db, "SELECT cwr_id FROM transmission_header ORDER BY cwr_id")
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
	c, err := newTestCLI(t)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(c.bzz.indexDir)
	defer c.srv.Close()

	// check 'meta convert ern' prints multiple CIDs
	stdout := c.run("convert", "ern",
		"--source", "test",
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
		"zdqaWQFfLpAj7Hi7B1DMsqfRzh2gtTWJb6PKzK2TJZgb3gCEM",
		"zdqaWJ4jkU4haHkCfJY8Tz7bNtdi38Kq1bRy4iiU9DUDJqUkB",
		"zdqaWRL5oXkGnwhf9JWzgMdV7H1YUWp6YociP5aWB9EU6qAJr",
		"zdqaWSBmxw7xif1hx4ZXVyC6A6Fr6hcx34JLuf9BGDT7P7utp",
		"zdqaWBoE9GfAc1dsJgW1jqtbHYKHgLcFfw55SKDyxEEGcQWso",
	}
	if !reflect.DeepEqual(ids, expected) {
		t.Fatalf("unexpected CIDs:\nexpected: %v\ngot:      %v", expected, ids)
	}

	db := filepath.Join(c.bzz.indexDir, "2_index.db")
	// run 'meta index ern' with the CIDs as stdin
	stream := strings.NewReader(stdout)
	c.bzz.indexDirHash = c.runWithStdin(stream, "index", "ern", filepath.Base(db), fmt.Sprintf("--bzzapi=%s", c.srv.Server.URL), fmt.Sprintf("--bzzdir=%s", c.bzz.indexDirHash))
	if c.bzz.indexDirHash == "" {
		t.Fatal("No hash returned")
	}

	// check the index was populated
	c.bzz.GetIndexFile(filepath.Base(db), true)
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

// TestERNCommands tests running the 'meta convert eidr' and
// 'meta index eidr' commands.
func TestEIDRCommands(t *testing.T) {
	c, err := newTestCLI(t)
	if err != nil {
		t.Fatal(err)
	}

	// check 'meta convert eidr' outputs expected rows
	stdout := c.run("convert", "eidr",
		"--source", "test",
		"../eidr/testdata/dummy_child.xml",
		"../eidr/testdata/dummy_parent.xml",
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
		"zdqaWUgTjLPoFrMCmBcbPcJNPuHTvUs5fBn3f4iYsniHACo7q",
		"zdqaWTaGbQFm2HEYo2iwHsFanffKDYq49XVk5ggEc6dYaBQkv",
	}
	if !reflect.DeepEqual(ids, expected) {
		t.Fatalf("unexpected CIDs:\nexpected: %v\ngot:      %v", expected, ids)
	}

	db := filepath.Join(c.bzz.indexDir, "3_index.db")

	// run 'meta index eidr' with the CIDs as stdin
	stream := strings.NewReader(stdout)
	c.bzz.indexDirHash = c.runWithStdin(stream, "index", "eidr", filepath.Base(db), fmt.Sprintf("--bzzapi=%s", c.srv.Server.URL), fmt.Sprintf("--bzzdir=%s", c.bzz.indexDirHash))
	if c.bzz.indexDirHash == "" {
		t.Fatal("No hash returned")
	}

	// check the index was populated
	c.bzz.GetIndexFile(filepath.Base(db), true)
	cmd := exec.Command("sqlite3", db, "select count(*) from xobject_baseobject_link x inner join baseobject p, xobject_episode e on p.doi_id = x.parent_doi_id where e.id = x.xobject_id")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error checking index: %s: %s", err, out)
	}
	if strings.TrimSpace(string(out)) != "1" {
		t.Fatalf("baseobject/xobject link count mismatch; expected 1, got %s", out)
	}

	// check if associatedorgs are inserted and linked
	cmd = exec.Command("sqlite3", db, "select count(*) from org o inner join baseobject b on o.base_doi_id = b.doi_id")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error checking index: %s: %s", err, out)
	}
	if strings.TrimSpace(string(out)) != "2" {
		t.Fatalf("associatedorg link count mismatch; expected 2, got %s", out)
	}

	// check if alternateids are inserted and linked
	cmd = exec.Command("sqlite3", db, "select count(*) from alternateid a inner join baseobject b on a.base_doi_id = b.doi_id")
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
	bzz   *SwarmBackend
	srv   *testutil.TestSwarmServer
}

func newTestCLI(t *testing.T) (*testCLI, error) {
	// create a path to store the index and to store the meta objects.
	tmpDir, err := ioutil.TempDir("", "meta-main-test")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(tmpDir)
		}
	}()
	srv := testutil.NewTestSwarmServer(t)

	bzzbackend := &SwarmBackend{
		api:      bzzclient.NewClient(srv.Server.URL),
		indexDir: tmpDir,
	}
	hash, err := bzzbackend.api.UploadDirectory(tmpDir, "", "")
	if err != nil {
		t.Fatal(err)
	}
	err = bzzbackend.OpenIndex(srv.Server.URL, hash)
	if err != nil {
		t.Fatal(err)
	}
	store := meta.NewSwarmDatastore(srv.URL)
	return &testCLI{
		t:     t,
		store: store,
		bzz:   bzzbackend,
		srv:   srv,
	}, nil
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
