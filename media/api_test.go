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

package media

import (
	"testing"

	"github.com/meta-network/go-meta/cwr"
	"github.com/meta-network/go-meta/ern"
	"github.com/meta-network/go-meta/identity"
	"github.com/meta-network/go-meta/musicbrainz"
	"github.com/meta-network/go-meta/testutil"
	"github.com/meta-network/go-meta/testutil/index"
)

// TestResolver tests resolving META Media API queries.
func TestResolver(t *testing.T) {
	// create test indexes
	store, cleanup := testutil.NewTestStore(t)
	defer cleanup()
	ernIndex, _ := testindex.GenerateERNIndex(t, "../ern", store)
	defer ernIndex.Close()
	cwrIndex, _ := testindex.GenerateCWRIndex(t, "../cwr", store)
	defer cwrIndex.Close()
	musicBrainzIndex, _, _ := testindex.GenerateMusicBrainzIndex(t, "../musicbrainz", store)
	defer musicBrainzIndex.Close()
	claimIndex, _ := testindex.GenerateClaimIndex(t, "../claim", store)
	defer claimIndex.Close()

	// create an ISWC -> ISRC mapping in the MusicBrainz index
	link := &musicbrainz.RecordingWorkLink{
		ISRC: "CASE00000001",
		ISWC: "T1234567890",
	}
	linkObj, err := store.Put(link)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := musicBrainzIndex.Exec(
		`INSERT INTO recording_work (object_id, isrc, iswc) VALUES ($1, $2, $3)`,
		linkObj.Cid().String(), link.ISRC, link.ISWC,
	); err != nil {
		t.Fatal(err)
	}

	identityResolver, err := identity.NewResolver(claimIndex.DB, claimIndex)
	if err != nil {
		t.Fatal(err)
	}

	// create the resolver
	resolver := &Resolver{
		Ern:         ern.NewResolver(ernIndex.DB, store),
		Cwr:         cwr.NewResolver(cwrIndex.DB, store),
		MusicBrainz: musicbrainz.NewResolver(musicBrainzIndex.DB, store),
		Identity:    identityResolver,
	}

	// query account
	testMetaID, err := testindex.GenerateTestMetaId()
	if err != nil {
		t.Fatal(err)
	}
	account, err := resolver.Account(accountArgs{MetaID: testMetaID.ID})
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
	if dpid := accountPerformers[0].Identifiers()[0].Value(); dpid != "DPID_OF_THE_ARTIST_1" {
		t.Fatalf("expected dpid to be %q, got %q", "DPID_OF_THE_ARTIST_1", dpid)
	}
	accountComposers, err := account.Composers()
	if err != nil {
		t.Fatal(err)
	}
	if len(accountComposers) != 1 {
		t.Fatalf("expected account to have 1 composer, got %d", len(accountComposers))
	}
	if ipi := accountComposers[0].Identifiers()[0].Value(); ipi != "123456789ABCD" {
		t.Fatalf("expected ipi to be %q, got %q", "123456789ABCD", ipi)
	}

	// query performers
	performer, err := resolver.Performer(performerArgs{DPID: "DPID_OF_THE_ARTIST_1"})
	if err != nil {
		t.Fatal(err)
	}
	if name := performer.Name().Value(); name != "Monkey Claw" {
		t.Fatalf("expected name to be %q, got %q", "Monkey Claw", name)
	}
	performerRecordings, err := performer.Recordings()
	if err != nil {
		t.Fatal(err)
	}
	if len(performerRecordings) != 1 {
		t.Fatalf("expected performer to have 1 recording, got %d", len(performerRecordings))
	}

	// query composers
	composer, err := resolver.Composer(composerArgs{IPI: "123456789ABCD"})
	if err != nil {
		t.Fatal(err)
	}
	if firstName := composer.FirstName().Value(); firstName != "WRITER_FIRST_NAME" {
		t.Fatalf("expected first name to be %q, got %q", "WRITER_FIRST_NAME", firstName)
	}
	if lastName := composer.LastName().Value(); lastName != "WRITER_LAST_NAME" {
		t.Fatalf("expected last name to be %q, got %q", "WRITER_LAST_NAME", lastName)
	}

	// query labels
	label, err := resolver.Label(labelArgs{DPID: "DPID_OF_THE_SENDER"})
	if err != nil {
		t.Fatal(err)
	}
	if name := label.Name().Value(); name != "NAME_OF_THE_SENDER" {
		t.Fatalf("expected name to be %q, got %q", "NAME_OF_THE_SENDER", name)
	}
	labelProducts, err := label.Releases()
	if err != nil {
		t.Fatal(err)
	}
	if len(labelProducts) != 1 {
		t.Fatalf("expected label to have 1 product, got %d", len(labelProducts))
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
	recordingReleases, err := recording.Songs()
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
	song, err := resolver.Song(songArgs{GRID: "A1UCASE0000000001X"})
	if err != nil {
		t.Fatal(err)
	}
	if title := song.Title().Value(); title != "The Tin Drum" {
		t.Fatalf("expected name to be %q, got %q", "The Tin Drum", title)
	}
	songRecordings, err := song.Recordings()
	if err != nil {
		t.Fatal(err)
	}
	if len(songRecordings) != 7 {
		t.Fatalf("expected song to have 7 recordings, got %d", len(songRecordings))
	}

	// query products
	release, err := resolver.Release(releaseArgs{UPC: "UPC000000001"})
	if err != nil {
		t.Fatal(err)
	}
	if release.Title().Value() != "A Monkey Claw in a Velvet Glove" {
		t.Fatalf("expected name to be %q, got %q", "A Monkey Claw in a Velvet Glove", release.Title().Value())
	}

	Labels, err := release.Labels()
	if err != nil {
		t.Fatal(err)
	}
	if len(Labels) != 1 {
		t.Fatalf("expected release to have 1 label, got %d", len(Labels))
	}

	performerResolver, err := release.Performer()
	if err != nil {
		t.Fatal(err)
	}
	p, err := performerResolver.Performer()
	if err != nil {
		t.Fatal(err)
	}
	if p.Name().Value() != "Monkey Claw" {
		t.Fatalf("expected performer name to be %q got %q", "Monkey Claw", p.Name().Value())
	}
}
