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

package ern

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/meta-network/go-meta/identity"
	"github.com/meta-network/go-meta/media"
	"github.com/meta-network/go-meta/testutil"
)

func TestImport(t *testing.T) {
	// start Media API
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	index, err := store.OpenIndex("media.meta")
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()
	mediaIndex, err := media.NewIndex(index)
	if err != nil {
		t.Fatal(err)
	}
	idIndex, err := store.OpenIndex("id.meta")
	if err != nil {
		t.Fatal(err)
	}
	defer idIndex.Close()
	idResolver, err := identity.NewResolver(idIndex)
	if err != nil {
		t.Fatal(err)
	}
	api, err := media.NewAPI(media.NewResolver(mediaIndex, idResolver))
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(api)
	defer srv.Close()

	// import ERNs
	source := &media.Source{Name: "test"}
	client := media.NewClient(srv.URL, source)
	importer := NewImporter(client)
	erns := []string{
		"Profile_AudioAlbumMusicOnly.xml",
		"Profile_AudioSingle.xml",
		"Profile_AudioSingle_WithCompoundArtistsAndTerritorialOverride.xml",
	}
	for _, path := range erns {
		f, err := os.Open(filepath.Join("testdata", path))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := importer.ImportERN(f); err != nil {
			t.Fatal(err)
		}
	}

	// check record label
	recordLabel, err := client.RecordLabel(&media.Identifier{
		Type:  "dpid",
		Value: "DPID_OF_THE_SENDER",
	})
	if err != nil {
		t.Fatal(err)
	}
	if recordLabel.Name != "NAME_OF_THE_SENDER" {
		t.Fatalf("expected record label to have name %q, got %q", "NAME_OF_THE_SENDER", recordLabel.Name)
	}

	// check sound recordings
	for isrc, title := range map[string]string{
		"CASE01000001": "Can you feel ...the Monkey Claw!",
		"CASE01000002": "Red top mountain, blown sky high",
		"CASE01000003": "Seige of Antioch",
		"CASE01000004": "Warhammer",
		"CASE01000005": "Iron Horse",
		"CASE01000006": "Yes... I can feel the Monkey Claw!",
		"CASE02000001": "Can you feel ...the Monkey Claw!",
		"CASE03000001": "Can you feel ...the Monkey Claw!",
	} {
		recording, err := client.Recording(&media.Identifier{
			Type:  "isrc",
			Value: isrc,
		})
		if err != nil {
			t.Fatal(err)
		}
		if recording.Title != title {
			t.Fatalf("expected sound recording %s to have title %q, got %q", isrc, title, recording.Title)
		}
	}

	// check releases
	for grid, title := range map[string]string{
		"A1UCASE0100000401X": "A Monkey Claw in a Velvet Glove",
		"A1UCASE0200000001X": "Can you feel ...the Monkey Claw!",
		"A1UCASE0300000001X": "Can you feel ...the Monkey Claw!",
	} {
		release, err := client.Release(&media.Identifier{
			Type:  "grid",
			Value: grid,
		})
		if err != nil {
			t.Fatal(err)
		}
		if release.Title != title {
			t.Fatalf("expected release %s to have title %q, got %q", grid, title, release.Title)
		}
	}

	// check songs
	for grid, title := range map[string]string{
		"A1UCASE0100000001X": "Can you feel ...the Monkey Claw!",
		"A1UCASE0100000002X": "Red top mountain, blown sky high",
		"A1UCASE0100000003X": "Seige of Antioch",
		"A1UCASE0100000004X": "Warhammer",
		"A1UCASE0100000005X": "Iron Horse",
		"A1UCASE0100000006X": "Yes... I can feel the Monkey Claw!",
		"A1UCASE0200000001X": "Can you feel ...the Monkey Claw!",
		"A1UCASE0300000001X": "Can you feel ...the Monkey Claw!",
	} {
		song, err := client.Song(&media.Identifier{
			Type:  "grid",
			Value: grid,
		})
		if err != nil {
			t.Fatal(err)
		}
		if song.Title != title {
			t.Fatalf("expected song %s to have title %q, got %q", grid, title, song.Title)
		}
	}

	// check artists
	for dpid, name := range map[string]string{
		"DPID_OF_THE_ARTIST_1": "Monkey Claw",
		"DPID_OF_THE_ARTIST_2": "Steve Albino",
		"DPID_OF_THE_ARTIST_3": "Bob Black",
	} {
		artist, err := client.Performer(&media.Identifier{
			Type:  "dpid",
			Value: dpid,
		})
		if err != nil {
			t.Fatal(err)
		}
		if artist.Name != name {
			t.Fatalf("expected artist %s to have name %q, got %q", dpid, name, artist.Name)
		}
	}
}
