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

package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/meta-network/go-meta"
)

// TestCWRCommands tests running the 'meta cwr convert' and
// 'meta cwr index' commands.
func TestCWRCommands(t *testing.T) {
	m := newTestMain(t)

	// check 'meta cwr convert' prints a CID
	stdout := m.run("cwr", "convert", "../../cwr/testdata/testfile.cwr", "../../cwr/CWR-DataApi")
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
	m.runWithStdin(stream, "cwr", "index", db)

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

type testMain struct {
	t     *testing.T
	store *meta.Store
}

func newTestMain(t *testing.T) *testMain {
	return &testMain{
		t:     t,
		store: meta.NewStore(datastore.NewMapDatastore()),
	}
}

func (m *testMain) runWithStdin(stdin io.Reader, args ...string) string {
	var stdout bytes.Buffer
	main := NewMain(m.store, stdin, &stdout)
	if err := main.Run(args...); err != nil {
		m.t.Fatal(err)
	}
	return stdout.String()
}

func (m *testMain) run(args ...string) string {
	return m.runWithStdin(nil, args...)
}
