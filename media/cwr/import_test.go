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

package cwr

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

	// import CWRs
	source := &media.Source{Name: "test"}
	client := media.NewClient(srv.URL+"/graphql", source)
	importer := NewImporter(client)
	cwrs := []string{
		"example_nwr.cwr",
		"example_double_nwr.cwr",
	}
	for _, path := range cwrs {
		f, err := os.Open(filepath.Join("testdata", path))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := importer.ImportCWR(f); err != nil {
			t.Fatal(err)
		}
	}

	// check publisher
	identifier := &media.Identifier{
		Type:  "ipi",
		Value: "00123456789",
	}
	query := `
query GetPublisher($identifier: IdentifierInput!) {
  publisher(identifier: $identifier) {
    name {
      value
    }
    works {
      work {
	title {
	  value
	}
      }
      role
    }
  }
}
`
	var v struct {
		Publisher struct {
			Name struct {
				Value string `json:"value"`
			} `json:"name"`
			Works []struct {
				Work struct {
					Title struct {
						Value string `json:"value"`
					} `json:"title"`
				} `json:"work"`
				Role string `json:"role"`
			} `json:"works"`
		} `json:"publisher"`
	}
	if err := client.Query(query, graphql.Variables{"identifier": identifier}, &v); err != nil {
		t.Fatal(err)
	}
	if v.Publisher.Name.Value != "NAME_OF_THE_SENDER" {
		t.Fatalf("expected publisher to have name %q, got %q", "NAME_OF_THE_SENDER", v.Publisher.Name.Value)
	}
	if len(v.Publisher.Works) != 1 {
		t.Fatalf("expected publisher to have 1 work, got %d", len(v.Publisher.Works))
	}
}
