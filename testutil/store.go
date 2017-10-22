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

package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/meta-network/go-meta"
)

func NewTestStore(t *testing.T) (*meta.Store, func()) {
	tmpDir, err := ioutil.TempDir("", "meta-test")
	if err != nil {
		t.Fatal(err)
	}
	ensDir := filepath.Join(tmpDir, "ens")
	if err := os.Mkdir(ensDir, 0755); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}
	store, err := meta.NewStore(tmpDir, meta.LocalENS(ensDir))
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}
	return store, func() { os.RemoveAll(tmpDir) }
}
