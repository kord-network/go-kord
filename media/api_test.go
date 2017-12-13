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

package media_test

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/meta-network/go-meta/identity"
	"github.com/meta-network/go-meta/media"
	"github.com/meta-network/go-meta/media/ern"
	"github.com/meta-network/go-meta/testutil"
	"github.com/meta-network/go-meta/testutil/index"
)

// TestResolver tests resolving META Media API queries.
func TestResolver(t *testing.T) {
	// create test indexes
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
	client := media.NewClient(srv.URL+"/graphql", &media.Source{Name: "test"})
	importer := ern.NewImporter(client)
	erns := []string{
		"Profile_AudioAlbumMusicOnly.xml",
		"Profile_AudioSingle.xml",
		"Profile_AudioSingle_WithCompoundArtistsAndTerritorialOverride.xml",
	}
	for _, path := range erns {
		f, err := os.Open(filepath.Join("ern", "testdata", path))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := importer.ImportERN(f); err != nil {
			t.Fatal(err)
		}
	}

	// create an identity with a DPID and IPI
	identity := testindex.GenerateTestIdentity(t)
	if err := identityIndex.CreateIdentity(identity); err != nil {
		t.Fatal(err)
	}
	for property, claim := range map[string]string{
		"dpid": "DPID_OF_THE_ARTIST_1",
		"ipi":  "123456789ABCD",
	} {
		claim := testindex.GenerateTestClaim(t, identity, property, claim)
		if err := identityIndex.CreateClaim(claim); err != nil {
			t.Fatal(err)
		}
	}

	/*
		// create an ISWC -> ISRC mapping
		link := &media.RecordingWorkLink{
			Recording: media.Identifier{Type: "isrc", Value: "CASE02000001"},
			Work:      media.Identifier{Type: "iswc", Value: "T1234567890"},
		}
		if err := client.CreateRecordingWorkLink(link); err != nil {
			t.Fatal(err)
		}
	*/

	// query account
	resolver := api.Resolver()
	account, err := resolver.Account(media.AccountArgs{MetaID: identity.ID().String()})
	if err != nil {
		t.Fatal(err)
	}
	accountPerformers, err := account.Performers()
	if err != nil {
		t.Fatal(err)
	}
	if len(accountPerformers) != 1 {
		t.Fatalf("expected account to have 1 performer, got %d", len(accountPerformers))
	}
	if dpid := accountPerformers[0].Identifiers()[0].Value().Value(); dpid != "DPID_OF_THE_ARTIST_1" {
		t.Fatalf("expected dpid to be %q, got %q", "DPID_OF_THE_ARTIST_1", dpid)
	}
	// accountComposers, err := account.Composers()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if len(accountComposers) != 1 {
	// 	t.Fatalf("expected account to have 1 composer, got %d", len(accountComposers))
	// }
	// if ipi := accountComposers[0].Identifiers()[0].Value().Value(); ipi != "123456789ABCD" {
	// 	t.Fatalf("expected ipi to be %q, got %q", "123456789ABCD", ipi)
	// }

	/*
		// query performers
		performer, err := resolver.Performer(media.IdentifierArgs{
			Identifier: media.Identifier{Type: "dpid", Value: "DPID_OF_THE_ARTIST_1"},
		})
		if err != nil {
			t.Fatal(err)
		}
		name, err := performer.Name()
		if err != nil {
			t.Fatal(err)
		}
		if name.Value() != "Monkey Claw" {
			t.Fatalf("expected name to be %q, got %q", "Monkey Claw", name.Value())
		}
		performerRecordings, err := performer.Recordings()
		if err != nil {
			t.Fatal(err)
		}
		if len(performerRecordings) != 1 {
			t.Fatalf("expected performer to have 1 recording, got %d", len(performerRecordings))
		}

		// query composers
		composer, err := resolver.Composer(media.IdentifierArgs{
			Identifier: media.Identifier{Type: "ipi", Value: "123456789ABCD"},
		})
		if err != nil {
			t.Fatal(err)
		}
		firstName, err := composer.FirstName()
		if err != nil {
			t.Fatal(err)
		}
		if firstName.Value() != "WRITER_FIRST_NAME" {
			t.Fatalf("expected first name to be %q, got %q", "WRITER_FIRST_NAME", firstName.Value())
		}
		lastName, err := composer.LastName()
		if err != nil {
			t.Fatal(err)
		}
		if lastName.Value() != "WRITER_LAST_NAME" {
			t.Fatalf("expected last name to be %q, got %q", "WRITER_LAST_NAME", lastName.Value())
		}

		// query labels
		label, err := resolver.Label(media.IdentifierArgs{
			Identifier: media.Identifier{Type: "dpid", Value: "DPID_OF_THE_SENDER"},
		})
		if err != nil {
			t.Fatal(err)
		}
		name, err = label.Name()
		if err != nil {
			t.Fatal(err)
		}
		if name.Value() != "NAME_OF_THE_SENDER" {
			t.Fatalf("expected name to be %q, got %q", "NAME_OF_THE_SENDER", name.Value())
		}
		labelReleases, err := label.Releases()
		if err != nil {
			t.Fatal(err)
		}
		if len(labelReleases) != 1 {
			t.Fatalf("expected label to have 1 release, got %d", len(labelRelease))
		}

		// query recordings
		recording, err := resolver.Recording(recordingArgs{ISRC: "CASE00000001"})
		if err != nil {
			t.Fatal(err)
		}
		if title := recording.Title().Value(); title != "Can you feel ...the Monkey Claw!" {
			t.Fatalf("expected name to be %q, got %q", "Can you feel ...the Monkey Claw!", title)
		}
		recordingPerformers, err := recording.Performers()
		if err != nil {
			t.Fatal(err)
		}
		if len(recordingPerformers) != 5 {
			t.Fatalf("expected recording to have 5 performers, got %d", len(recordingPerformers))
		}
		recordingReleases, err := recording.Releases()
		if err != nil {
			t.Fatal(err)
		}
		if len(recordingReleases) != 10 {
			t.Fatalf("expected recording to have 10 releases, got %d", len(recordingReleases))
		}
		recordingWorks, err := recording.Works()
		if err != nil {
			t.Fatal(err)
		}
		if len(recordingWorks) != 5 {
			t.Fatalf("expected recording to have 5 works, got %d", len(recordingWorks))
		}

		// query works
		work, err := resolver.Work(workArgs{ISWC: "T1234567890"})
		if err != nil {
			t.Fatal(err)
		}
		if title := work.Title().Value(); title != "TOTALY MADE MUSIC UP" {
			t.Fatalf("expected title to be %q, got %q", "TOTALY MADE MUSIC UP", title)
		}
		workRecordings, err := work.Recordings()
		if err != nil {
			t.Fatal(err)
		}
		if len(workRecordings) != 1 {
			t.Fatalf("expected work to have 1 recording, got %d", len(workRecordings))
		}

		// query releases
		release, err := resolver.Release(releaseArgs{GRID: "A1UCASE0000000401X"})
		if err != nil {
			t.Fatal(err)
		}
		if title := release.Title().Value(); title != "A Monkey Claw in a Velvet Glove" {
			t.Fatalf("expected name to be %q, got %q", "A Monkey Claw in a Velvet Glove", title)
		}
		releaseRecordings, err := release.Recordings()
		if err != nil {
			t.Fatal(err)
		}
		if len(releaseRecordings) != 12 {
			t.Fatalf("expected release to have 12 recordings, got %d", len(releaseRecordings))
		}
		releaseProducts := release.Products()
		if len(releaseProducts) != 1 {
			t.Fatalf("expected release to have 1 product, got %d", len(releaseProducts))
		}

		// query products
		product, err := resolver.Product(productArgs{UPC: "UPC000000001"})
		if err != nil {
			t.Fatal(err)
		}
		releases := product.Releases()
		if err != nil {
			t.Fatal(err)
		}
		if len(releases) != 1 {
			t.Fatalf("expected product to have 1 release, got %d", len(releases))
		}
		productRelease, err := releases[0].Release()
		if err != nil {
			t.Fatal(err)
		}
		if grid := productRelease.GRID().Value(); grid != "A1UCASE0000000401X" {
			t.Fatalf("expected product release GRid to be %q, got %q", "A1UCASE0000000401X", grid)
		}
		productLabels, err := product.Labels()
		if err != nil {
			t.Fatal(err)
		}
		if len(productLabels) != 1 {
			t.Fatalf("expected product to have 1 label, got %d", len(productLabels))
		}
		productPerformers, err := product.Performers()
		if err != nil {
			t.Fatal(err)
		}
		if len(productPerformers) != 1 {
			t.Fatalf("expected product to have 1 performer, got %d", len(productPerformers))
		}
	*/
}
