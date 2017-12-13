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

	"github.com/meta-network/go-meta/graphql"
	"github.com/meta-network/go-meta/identity"
	"github.com/meta-network/go-meta/media"
	"github.com/meta-network/go-meta/testutil"
)

func TestImport(t *testing.T) {
	// start Media API
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	mediaIndex, err := media.NewIndex(store)
	if err != nil {
		t.Fatal(err)
	}
	defer mediaIndex.Close()
	identityIndex, err := identity.NewIndex(store)
	if err != nil {
		t.Fatal(err)
	}
	defer identityIndex.Close()
	api, err := media.NewAPI(mediaIndex, identityIndex)
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(api)
	defer srv.Close()

	// import ERNs
	source := &media.Source{Name: "test"}
	client := media.NewClient(srv.URL+"/graphql", source)
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
	identifier := &media.Identifier{
		Type:  "dpid",
		Value: "DPID_OF_THE_SENDER",
	}
	query := `
query GetRecordLabel($identifier: IdentifierInput!) {
  record_label(identifier: $identifier) {
    name {
      value
    }
    songs {
      song {
	title {
	  value
	}
      }
    }
    releases {
      release {
	title {
	  value
	}
      }
    }
  }
}
`
	var v struct {
		RecordLabel struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
			Songs []struct {
				Song struct {
					Title struct {
						Value string `json:"value"`
					} `json:"title"`
				} `json:"song"`
			} `json:"songs"`
			Releases []struct {
				Release struct {
					Title struct {
						Value string `json:"value"`
					} `json:"title"`
				} `json:"release"`
			} `json:"releases"`
		} `json:"record_label"`
	}
	if err := client.Query(query, graphql.Variables{"identifier": identifier}, &v); err != nil {
		t.Fatal(err)
	}
	if v.RecordLabel.Name.Value != "NAME_OF_THE_SENDER" {
		t.Fatalf("expected record label to have name %q, got %q", "NAME_OF_THE_SENDER", v.RecordLabel.Name.Value)
	}
	if len(v.RecordLabel.Releases) != 3 {
		t.Fatalf("expected record label to have 3 releases, got %d", len(v.RecordLabel.Releases))
	}
	if len(v.RecordLabel.Songs) != 8 {
		t.Fatalf("expected record label to have 8 songs, got %d", len(v.RecordLabel.Songs))
	}

	// check sound recordings
	type soundRecording struct {
		isrc       string
		title      string
		performers map[string]string
	}
	for _, x := range []soundRecording{
		{
			isrc:  "CASE01000001",
			title: "Can you feel ...the Monkey Claw!",
			performers: map[string]string{
				"MainArtist": "Monkey Claw",
				"Producer":   "Steve Albino",
				"Composer":   "Bob Black",
			},
		},
		{
			isrc:  "CASE01000002",
			title: "Red top mountain, blown sky high",
			performers: map[string]string{
				"MainArtist": "Monkey Claw",
				"Producer":   "Steve Albino",
				"Composer":   "Bob Black",
			},
		},
		{
			isrc:  "CASE01000003",
			title: "Seige of Antioch",
			performers: map[string]string{
				"MainArtist": "Monkey Claw",
				"Producer":   "Steve Albino",
				"Composer":   "Bob Black",
			},
		},
		{
			isrc:  "CASE01000004",
			title: "Warhammer",
			performers: map[string]string{
				"MainArtist": "Monkey Claw",
				"Producer":   "Steve Albino",
				"Composer":   "Bob Black",
			},
		},
		{
			isrc:  "CASE01000005",
			title: "Iron Horse",
			performers: map[string]string{
				"MainArtist": "Monkey Claw",
				"Producer":   "Steve Albino",
				"Composer":   "Bob Black",
			},
		},
		{
			isrc:  "CASE01000006",
			title: "Yes... I can feel the Monkey Claw!",
			performers: map[string]string{
				"MainArtist": "Monkey Claw",
				"Producer":   "Steve Albino",
				"Composer":   "Bob Black",
			},
		},
		{
			isrc:  "CASE02000001",
			title: "Can you feel ...the Monkey Claw!",
			performers: map[string]string{
				"MainArtist":     "Monkey Claw",
				"FeaturedArtist": "Monkey Claw",
				"Producer":       "Steve Albino",
				"Composer":       "Bob Black",
			},
		},
		{
			isrc:  "CASE03000001",
			title: "Can you feel ...the Monkey Claw!",
			performers: map[string]string{
				"MainArtist": "Monkey Claw",
				"Producer":   "Steve Albino",
				"Composer":   "Bob Black",
			},
		},
	} {
		identifier := &media.Identifier{
			Type:  "isrc",
			Value: x.isrc,
		}
		query = `
query GetRecording($identifier: IdentifierInput!) {
  recording(identifier: $identifier) {
    title {
      value
    }
    performers {
      performer {
	name {
	  value
	}
      }
      role
    }
  }
}
		`
		var v struct {
			Recording struct {
				Title struct {
					Value string `json:"value"`
				} `json:"title"`
				Performers []struct {
					Performer struct {
						Name struct {
							Value string `json:"value"`
						} `json:"name"`
					} `json:"performer"`
					Role string `json:"role"`
				} `json:"performers"`
			} `json:"recording"`
		}
		if err := client.Query(query, graphql.Variables{"identifier": identifier}, &v); err != nil {
			t.Fatal(err)
		}
		if v.Recording.Title.Value != x.title {
			t.Fatalf("expected sound recording %s to have title %q, got %q", x.isrc, x.title, v.Recording.Title.Value)
		}
		if len(v.Recording.Performers) != len(x.performers) {
			t.Fatalf("expected sound recording %s to have %d performers, got %d", x.isrc, len(x.performers), len(v.Recording.Performers))
		}
		for _, p := range v.Recording.Performers {
			name, ok := x.performers[p.Role]
			if !ok {
				t.Fatalf("unexpected performer role: %q", p.Role)
			}
			if p.Performer.Name.Value != name {
				t.Fatalf("expected %s performer to have name %q, got %q", p.Role, name, p.Performer.Name.Value)
			}
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
