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
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	cid "github.com/ipfs/go-cid"
	meta "github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/cwr"
	"github.com/meta-network/go-meta/ern"
	"github.com/meta-network/go-meta/identity"
	"github.com/meta-network/go-meta/musicbrainz"
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

func GenerateMusicBrainzIndex(t *testing.T, dir string, store *meta.Store) (*meta.Index, []*musicbrainz.Artist, []*musicbrainz.RecordingWorkLink) {
	// load the test data
	f, err := os.Open(filepath.Join(dir, "testdata/artists.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var artists []*musicbrainz.Artist
	dec := json.NewDecoder(f)
	for {
		var artist *musicbrainz.Artist
		err := dec.Decode(&artist)
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		artists = append(artists, artist)
	}
	f, err = os.Open(filepath.Join(dir, "testdata/recording-work-links.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var links []*musicbrainz.RecordingWorkLink
	dec = json.NewDecoder(f)
	for {
		var link musicbrainz.RecordingWorkLink
		err := dec.Decode(&link)
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		links = append(links, &link)
	}

	// store the artists and links in a test store
	artistCids := make([]*cid.Cid, len(artists))
	for i, artist := range artists {
		obj, err := store.Put(artist)
		if err != nil {
			t.Fatal(err)
		}
		artistCids[i] = obj.Cid()
	}
	linkCids := make([]*cid.Cid, len(links))
	for i, link := range links {
		obj, err := store.Put(link)
		if err != nil {
			t.Fatal(err)
		}
		linkCids[i] = obj.Cid()
	}

	// create streams
	streams := func(name string, count int) (*meta.StreamReader, *meta.StreamWriter, error) {
		reader, err := store.StreamReader(name, meta.StreamLimit(count))
		if err != nil {
			return nil, nil, err
		}
		writer, err := store.StreamWriter(name)
		if err != nil {
			reader.Close()
			return nil, nil, err
		}
		return reader, writer, nil
	}
	artistStreamR, artistStreamW, err := streams("artists.musicbrainz.meta", len(artistCids))
	if err != nil {
		t.Fatal(err)
	}
	defer artistStreamR.Close()
	defer artistStreamW.Close()
	if err := artistStreamW.Write(artistCids...); err != nil {
		t.Fatal(err)
	}
	linkStreamR, linkStreamW, err := streams("links.musicbrainz.meta", len(linkCids))
	if err != nil {
		t.Fatal(err)
	}
	defer linkStreamR.Close()
	defer linkStreamW.Close()
	if err := linkStreamW.Write(linkCids...); err != nil {
		t.Fatal(err)
	}

	// index the artists and links
	index, err := store.OpenIndex("musicbrainz.index.meta")
	if err != nil {
		t.Fatal(err)
	}
	indexer, err := musicbrainz.NewIndexer(index, store)
	if err != nil {
		index.Close()
		t.Fatal(err)
	}
	if err := indexer.IndexArtists(context.Background(), artistStreamR); err != nil {
		index.Close()
		t.Fatal(err)
	}
	if err := indexer.IndexRecordingWorkLinks(context.Background(), linkStreamR); err != nil {
		index.Close()
		t.Fatal(err)
	}
	return index, artists, links
}

const (
	testIdentityAddr     = "970e8128ab834e8eac17ab8e3812f010678cf791"
	testIdentityKey      = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	testIdentityUsername = "testid"
)

func GenerateTestIdentity(t *testing.T) *identity.Identity {
	key, err := crypto.HexToECDSA(testIdentityKey)
	if err != nil {
		t.Fatal(err)
	}
	identity := &identity.Identity{
		Username: testIdentityUsername,
		Owner:    common.HexToAddress(testIdentityAddr),
	}
	id := identity.ID()
	signature, err := crypto.Sign(id.Hash[:], key)
	if err != nil {
		t.Fatal(err)
	}
	identity.Signature = signature
	return identity
}

func GenerateTestClaim(t *testing.T, testIdentity *identity.Identity, property, claim string) *identity.Claim {
	key, err := crypto.HexToECDSA(testIdentityKey)
	if err != nil {
		t.Fatal(err)
	}
	c := &identity.Claim{
		Issuer:   testIdentity.ID(),
		Subject:  testIdentity.ID(),
		Property: property,
		Claim:    claim,
	}
	id := c.ID()
	signature, err := crypto.Sign(id[:], key)
	if err != nil {
		t.Fatal(err)
	}
	c.Signature = signature
	return c
}
