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

package eidr

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

	// import EIDR XML
	source := &media.Source{Name: "test"}
	client := media.NewClient(srv.URL+"/graphql", source)
	importer := NewImporter(client)
	f, err := os.Open(filepath.Join("testdata", "dummy_party.xml"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := importer.ImportParty(f); err != nil {
		t.Fatal(err)
	}
	files := []string{
		"dummy_series.xml",
		"dummy_season.xml",
		"dummy_episode.xml",
	}
	for _, path := range files {
		f, err := os.Open(filepath.Join("testdata", path))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := importer.ImportFullMetadata(f); err != nil {
			t.Fatal(err)
		}
	}

	// check organisation
	identifier := &media.Identifier{
		Type:  "doid",
		Value: "10.5237/FFFF-FFFF",
	}
	query := `
query GetOrganisation($identifier: IdentifierInput!) {
  organisation(identifier: $identifier) {
    name {
      value
    }
    series {
      series {
	name {
	  value
	}
      }
    }
    seasons {
      season {
	name {
	  value
	}
      }
    }
    episodes {
      episode {
	name {
	  value
	}
      }
    }
  }
}
`
	var v struct {
		Organisation struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
			Series []struct {
				Series struct {
					Name struct {
						Value string `json:"value"`
					} `json:"name"`
				} `json:"series"`
			} `json:"series"`
			Seasons []struct {
				Season struct {
					Name struct {
						Value string `json:"value"`
					} `json:"name"`
				} `json:"season"`
			} `json:"seasons"`
			Episodes []struct {
				Episode struct {
					Name struct {
						Value string `json:"value"`
					} `json:"name"`
				} `json:"episode"`
			} `json:"episodes"`
		} `json:"organisation"`
	}
	if err := client.Query(query, graphql.Variables{"identifier": identifier}, &v); err != nil {
		t.Fatal(err)
	}
	if v.Organisation.Name.Value != "Acme, Inc" {
		t.Fatalf("expected organisation to have name %q, got %q", "Acme, Inc", v.Organisation.Name.Value)
	}
	if len(v.Organisation.Series) != 1 {
		t.Fatalf("expected organisation to have 1 series, got %d", len(v.Organisation.Series))
	}
	if name := v.Organisation.Series[0].Series.Name.Value; name != "Foo" {
		t.Fatalf("expected series to have name %q, got %q", "Foo", name)
	}
	if len(v.Organisation.Seasons) != 1 {
		t.Fatalf("expected organisation to have 1 season, got %d", len(v.Organisation.Seasons))
	}
	if name := v.Organisation.Seasons[0].Season.Name.Value; name != "Foo: Season 1" {
		t.Fatalf("expected season to have name %q, got %q", "Foo: Season 1", name)
	}
	if len(v.Organisation.Episodes) != 1 {
		t.Fatalf("expected organisation to have 1 episode, got %d", len(v.Organisation.Episodes))
	}
	if name := v.Organisation.Episodes[0].Episode.Name.Value; name != "Foo: Season 1: Episode 1" {
		t.Fatalf("expected episode to have name %q, got %q", "Foo: Season 1: Episode 1", name)
	}

	// check series
	identifier = &media.Identifier{
		Type:  "doid",
		Value: "10.5240/FEED-BEEF-0123-4567-890A-C",
	}
	query = `
query GetSeries($identifier: IdentifierInput!) {
  series(identifier: $identifier) {
    name {
      value
    }
    seasons {
      season {
	name {
	  value
	}
      }
    }
  }
}
`
	var w struct {
		Series struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
			Seasons []struct {
				Season struct {
					Name struct {
						Value string `json:"value"`
					} `json:"name"`
				} `json:"season"`
			} `json:"seasons"`
		} `json:"series"`
	}
	if err := client.Query(query, graphql.Variables{"identifier": identifier}, &w); err != nil {
		t.Fatal(err)
	}
	if w.Series.Name.Value != "Foo" {
		t.Fatalf("expected series to have name %q, got %q", "Foo", w.Series.Name.Value)
	}
	if len(w.Series.Seasons) != 1 {
		t.Fatalf("expected series to have 1 season, got %d", len(w.Series.Seasons))
	}
	if name := w.Series.Seasons[0].Season.Name.Value; name != "Foo: Season 1" {
		t.Fatalf("expected season to have name %q, got %q", "Foo: Season 1", name)
	}

	// check season
	identifier = &media.Identifier{
		Type:  "doid",
		Value: "10.5240/DEAD-BEEF-0123-4567-890A-B",
	}
	query = `
query GetSeason($identifier: IdentifierInput!) {
  season(identifier: $identifier) {
    name {
      value
    }
    episodes {
      episode {
	name {
	  value
	}
      }
    }
  }
}
`
	var x struct {
		Season struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
			Episodes []struct {
				Episode struct {
					Name struct {
						Value string `json:"value"`
					} `json:"name"`
				} `json:"episode"`
			} `json:"episodes"`
		} `json:"season"`
	}
	if err := client.Query(query, graphql.Variables{"identifier": identifier}, &x); err != nil {
		t.Fatal(err)
	}
	if x.Season.Name.Value != "Foo: Season 1" {
		t.Fatalf("expected season to have name %q, got %q", "Foo: Season 1", x.Season.Name.Value)
	}
	if len(x.Season.Episodes) != 1 {
		t.Fatalf("expected season to have 1 episode, got %d", len(x.Season.Episodes))
	}
	if name := x.Season.Episodes[0].Episode.Name.Value; name != "Foo: Season 1: Episode 1" {
		t.Fatalf("expected episode to have name %q, got %q", "Foo: Season 1: Episode 1", name)
	}
}
