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

package musicbrainz

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
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

	// import artists and recording-work links
	f, err := os.Open("testdata/artists.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var artists []*Artist
	dec := json.NewDecoder(f)
	for {
		var artist Artist
		err := dec.Decode(&artist)
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		artists = append(artists, &artist)
	}
	f, err = os.Open("testdata/recording-work-links.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var links []*RecordingWorkLink
	dec = json.NewDecoder(f)
	for {
		var link RecordingWorkLink
		err := dec.Decode(&link)
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		links = append(links, &link)
	}
	source := &media.Source{Name: "test"}
	client := media.NewClient(srv.URL+"/graphql", source)
	importer := NewImporter(client)
	for _, artist := range artists {
		if err := importer.ImportArtist(artist); err != nil {
			t.Fatal(err)
		}
	}
	for _, link := range links {
		if err := importer.ImportRecordingWorkLink(link); err != nil {
			t.Fatal(err)
		}
	}

	// check artists
	query := `
query GetPerformer($identifier: IdentifierInput!) {
  performer(identifier: $identifier) {
    name {
      value
      sources {
	value
      }
    }
  }
}
`
	type performerResponse struct {
		Performer struct {
			Name struct {
				Value   string `json:"value"`
				Sources []struct {
					Value string `json:"value"`
				} `json:"sources"`
			} `json:"name"`
		} `json:"performer"`
	}
	for _, artist := range artists {
		// check getting the artist by ISNI
		for _, isni := range artist.ISNI {
			identifier := &media.Identifier{
				Type:  "isni",
				Value: isni,
			}
			var res performerResponse
			if err := client.Query(query, graphql.Variables{"identifier": identifier}, &res); err != nil {
				t.Fatal(err)
			}
			if res.Performer.Name.Value != artist.Name {
				t.Fatalf("expected performer with ISNI %q to have name %q, got %q", isni, artist.Name, res.Performer.Name.Value)
			}
		}
		// check getting the artist by IPI, also checking that IPI
		// "00435760746" returns two performer names "Future" and
		// "Lmars"
		for _, ipi := range artist.IPI {
			identifier := &media.Identifier{
				Type:  "ipi",
				Value: ipi,
			}
			var res performerResponse
			if err := client.Query(query, graphql.Variables{"identifier": identifier}, &res); err != nil {
				t.Fatal(err)
			}
			if ipi == "00435760746" {
				if len(res.Performer.Name.Sources) != 2 {
					t.Fatalf("expected performer with IPI 00435760746 to have two names, got %d", len(res.Performer.Name.Sources))
				}
				if res.Performer.Name.Sources[0].Value != "Future" {
					t.Fatalf("expected first performer with IPI 00435760746 to have name %q, got %q", ipi, "Future", res.Performer.Name.Sources[0].Value)
				}
				if res.Performer.Name.Sources[1].Value != "Lmars" {
					t.Fatalf("expected second performer with IPI 00435760746 to have name %q, got %q", ipi, "Lmars", res.Performer.Name.Sources[1].Value)
				}
			} else {
				if res.Performer.Name.Value != artist.Name {
					t.Fatalf("expected performer with IPI %q to have name %q, got %q", ipi, artist.Name, res.Performer.Name.Value)
				}
			}
		}
	}

	// check recording-work links
	query = `
query GetRecordingWorks($identifier: IdentifierInput!) {
  recording(identifier: $identifier) {
    title {
      value
    }
    works {
      work {
	title {
	  value
	}
      }
    }
  }
}
`
	type recordingResponse struct {
		Recording struct {
			Title struct {
				Value string `json:"value"`
			} `json:"title"`
			Works []struct {
				Work struct {
					Title struct {
						Value string `json:"value"`
					} `json:"title"`
				} `json:"work"`
			} `json:"works"`
		} `json:"recording"`
	}
	for _, link := range links {
		identifier := &media.Identifier{
			Type:  "isrc",
			Value: link.ISRC,
		}
		var res recordingResponse
		if err := client.Query(query, graphql.Variables{"identifier": identifier}, &res); err != nil {
			t.Fatal(err)
		}
		if res.Recording.Title.Value != link.RecordingTitle {
			t.Fatalf("expected recording with ISRC %q to have title %q, got %q", link.ISRC, link.RecordingTitle, res.Recording.Title.Value)
		}
		if len(res.Recording.Works) != 1 {
			t.Fatalf("expected recording with ISRC %q to have %d works, got %d", link.ISRC, 1, len(res.Recording.Works))
		}
	}
	query = `
query GetWorkRecordings($identifier: IdentifierInput!) {
  work(identifier: $identifier) {
    recordings {
      recording {
	title {
	  value
	}
      }
    }
  }
}
`
	type workResponse struct {
		Work struct {
			Recordings []struct {
				Recording struct {
					Title struct {
						Value string `json:"value"`
					} `json:"title"`
				} `json:"recording"`
			} `json:"recordings"`
		} `json:"work"`
	}

	// group the links by ISWC -> ISRC
	groupedLinks := make(map[string][]string)
	for _, link := range links {
		groupedLinks[link.ISWC] = append(groupedLinks[link.ISWC], link.ISRC)
	}

	for iswc, isrcs := range groupedLinks {
		identifier := &media.Identifier{
			Type:  "iswc",
			Value: iswc,
		}
		var res workResponse
		if err := client.Query(query, graphql.Variables{"identifier": identifier}, &res); err != nil {
			t.Fatal(err)
		}
		if len(res.Work.Recordings) != len(isrcs) {
			t.Fatalf("expected work with ISWC %q to have %d recordings, got %d", iswc, len(isrcs), len(res.Work.Recordings))
		}
	}
}
