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

package testindex

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	cid "github.com/ipfs/go-cid"
	meta "github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/cwr"
	"github.com/meta-network/go-meta/ern"
)

func GenerateERNIndex(t *testing.T, dir string, store *meta.Store) (*meta.Index, []*cid.Cid) {
	// convert the test ERNs to META objects
	erns := []string{
		"Profile_AudioAlbumMusicOnly.xml",
		"Profile_AudioSingle.xml",
		"Profile_AudioAlbum_WithBooklet.xml",
		"Profile_AudioSingle_WithCompoundArtistsAndTerritorialOverride.xml",
		"Profile_AudioBook.xml",
	}
	converter := ern.NewConverter(store)
	cids := make([]*cid.Cid, len(erns))
	for i, path := range erns {
		f, err := os.Open(filepath.Join(dir, "testdata", path))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		cid, err := converter.ConvertERN(f, "test")
		if err != nil {
			t.Fatal(err)
		}
		cids[i] = cid
	}

	// index the stream of ERNs
	writer, err := store.StreamWriter("ern.meta")
	if err != nil {
		t.Fatal(err)
	}
	defer writer.Close()
	if err := writer.Write(cids...); err != nil {
		t.Fatal(err)
	}

	index, err := store.OpenIndex("ern.index.meta")
	if err != nil {
		t.Fatal(err)
	}
	indexer, err := ern.NewIndexer(index, store)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	reader, err := store.StreamReader("ern.meta", meta.StreamLimit(len(cids)))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if err := indexer.Index(ctx, reader); err != nil {
		t.Fatal(err)
	}
	return index, cids
}

func GenerateCWRIndex(t *testing.T, dir string, store *meta.Store) (*meta.Index, *cid.Cid) {
	f, err := os.Open(filepath.Join(dir, "testdata", "example_nwr.cwr"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	converter := cwr.NewConverter(store)
	cwrCid, err := converter.ConvertCWR(f, "test")
	if err != nil {
		t.Fatal(err)
	}

	// create a stream of CWRs
	writer, err := store.StreamWriter("cwr.meta")
	if err != nil {
		t.Fatal(err)
	}
	defer writer.Close()
	if err := writer.Write(cwrCid); err != nil {
		t.Fatal(err)
	}

	// index the stream of CWRs
	index, err := store.OpenIndex("cwr.index.meta")
	if err != nil {
		t.Fatal(err)
	}
	indexer, err := cwr.NewIndexer(index, store)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reader, err := store.StreamReader("cwr.meta", meta.StreamLimit(1))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	if err := indexer.Index(ctx, reader); err != nil {
		t.Fatal(err)
	}
	return index, cwrCid
}
