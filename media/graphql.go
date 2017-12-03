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
	"github.com/meta-network/go-meta/cwr"
	"github.com/meta-network/go-meta/ern"
	"github.com/meta-network/go-meta/identity"
	"github.com/meta-network/go-meta/musicbrainz"
)

// GraphQLSchema is the GraphQL schema for the META Media index.
//
// It consists of the following main entities:
//
// Account
// MusicPerformer
// MusicComposer
// RecordLabel
// MusicPublisher
// MusicRecording
// MusicWork
// Song
// MusicRelease
//
// Each entity has a set of identifiers which are used to uniquely identify
// them:
//
// META-ID - META Identity
// ISNI    - International Standard Name Identifier
// IPI     - Interested Parties Information
// DPID    - DDEX Party Identifier
// ISRC    - International Standard Recording Code
// ISWC    - International Standard Musical Work Code
// GRid    - Global Song Identifier
// UPC     - Universal Release Code
// EAN     - International Article Number
//
// Other entity properties (like performer name or recording title) are linked
// to these identifiers by various sources in the META network.
const GraphQLSchema = `
schema {
  query: Query
}

type Query {
  account(meta_id: String!): Account!

  performer(dpid: String!): MusicPerformer!

  composer(ipi: String!): MusicComposer!

  label(dpid: String!): RecordLabel!

  publisher(ipi: String!): MusicPublisher!

  recording(isrc: String!): MusicRecording!

  work(iswc: String!): MusicWork!

  song(grid: String!): Song!

  release(upc: String!): MusicRelease!
}

#
# --- Main Entities ---
#
type Account {
    performers: [MusicPerformer]!
    labels:     [RecordLabel]!
    composers:  [MusicComposer]!
}

type MusicPerformer {
  identifiers: [PartyIdentifier]!

  name: StringValue

  recordings: [MusicRecordingLink]!
}

type MusicComposer {
  identifiers: [PartyIdentifier]!

  firstName: StringValue
  lastName:  StringValue
  shares:    MusicShare!

  works: [MusicWorkLink]!
}

type RecordLabel {
  identifiers: [PartyIdentifier]!

  name: StringValue

  releases: [MusicReleaseLink]!
}

type MusicPublisher {
  identifiers: [PartyIdentifier]!

  name: StringValue
  shares: MusicShare!
  works: [MusicWorkLink]!
}

type MusicRecording {
  isrc:       StringValue

  title:      StringValue
  duration :  StringValue

  performers: [MusicPerformerLink]!
  songs:      [SongLink]!
  works:      [MusicWorkLink]!
}

type MusicWork {
  iswc: StringValue

  title: StringValue

  composers:  [MusicComposerLink]!
  publishers: [MusicPublisherLink]!
  recordings: [MusicRecordingLink]!
}

type Song {
  grid: StringValue

  title: StringValue

  releases:   [MusicReleaseLink]!
  recordings: [MusicRecordingLink]!
}

type MusicRelease {
  upc: StringValue

  title: StringValue

  songs:      [SongLink]!
  performer:  MusicPerformerLink!
  labels:     [RecordLabelLink]!
}

type MusicShare {
	performance : StringValue
	mechanical  : StringValue
	synch       : StringValue
}


#
# --- Link Types ---
#
type MusicPerformerLink {
  source:    Source!
  performer: MusicPerformer!
  role:      StringValue
}

type MusicComposerLink {
  source:   Source!
  composer: MusicComposer!
}

type RecordLabelLink {
  source: Source!
  label:  RecordLabel!
}

type MusicPublisherLink {
  source:    Source!
  publisher: MusicPublisher!
}

type MusicRecordingLink {
  source:    Source!
  recording: MusicRecording!
}

type MusicWorkLink {
  source: Source!
  work:   MusicWork!
}

type SongLink {
  source:  Source!
  song: Song!
}

type MusicReleaseLink {
  source:  Source!
  release: MusicRelease!
}

#
# --- Value Types ---
#

type PartyIdentifier {
  type:   String!
  value:  String!
  source: Source!
}

type StringValue {
  value:   String!
  sources: [StringSource]!
}

type StringSource {
  value:  String!
  source: Source!
  score:  String!
}

type Source {
  name: String!
}
`

// Resolver defines GraphQL resolver functions for the schema contained in
// the GraphQLSchema constant, retrieving data from a META store and SQLite3
// index.
type Resolver struct {
	MusicBrainz *musicbrainz.Resolver
	Ern         *ern.Resolver
	Cwr         *cwr.Resolver
	Identity    *identity.Resolver
}

type accountArgs struct {
	MetaID string
}

func (r *Resolver) Account(args accountArgs) (*accountResolver, error) {
	claimsResolver, err := r.Identity.Claim(identity.ClaimArgs{Subject: &args.MetaID})
	if err != nil {
		return nil, err
	}
	aux := make(map[string]string)
	for _, claimsResolver := range claimsResolver {
		aux[claimsResolver.Claim()] = claimsResolver.Signature()
	}

	return &accountResolver{
		resolver: r,
		aux:      aux,
	}, nil
}

type accountResolver struct {
	resolver *Resolver
	aux      map[string]string
}

func (a *accountResolver) Performers() ([]*performerResolver, error) {
	var performers []*performerResolver
	if dpid, ok := a.aux["DPID"]; ok {
		parties, err := a.resolver.Ern.Party(ern.PartyArgs{ID: &dpid})
		if err != nil {
			return nil, err
		}
		if len(parties) > 0 {
			performer := a.resolver.performerResolver(parties)
			if len(performer.parties) > 0 {
				performers = append(performers, performer)
			}
		}
	}
	return performers, nil
}

func (a *accountResolver) Labels() ([]*labelResolver, error) {
	var labels []*labelResolver
	if dpid, ok := a.aux["DPID"]; ok {
		parties, err := a.resolver.Ern.Party(ern.PartyArgs{ID: &dpid})
		if err != nil {
			return nil, err
		}
		if len(parties) > 0 {
			labels = append(labels, a.resolver.labelResolver(parties))
		}
	}
	return labels, nil
}

func (a *accountResolver) Composers() ([]*composerResolver, error) {
	var composers []*composerResolver
	if ipi, ok := a.aux["IPI"]; ok {
		ipiBaseWriters, err := a.resolver.Cwr.WriterControl(cwr.WriterControlArgs{WriterIPIBaseNumber: &ipi})
		if err != nil {
			return nil, err
		}
		ipiNameWriters, err := a.resolver.Cwr.WriterControl(cwr.WriterControlArgs{WriterIPIName: &ipi})
		if err != nil {
			return nil, err
		}

		writers := append(ipiBaseWriters, ipiNameWriters...)
		if len(writers) > 0 {
			composers = append(composers, a.resolver.composerResolver(writers))
		}
	}
	return composers, nil
}

type performerArgs struct {
	DPID string
}

func (r *Resolver) Performer(args performerArgs) (*performerResolver, error) {
	// fetch parties by DPID from ERN
	parties, err := r.Ern.Party(ern.PartyArgs{ID: &args.DPID})
	if err != nil {
		return nil, err
	}
	return r.performerResolver(parties), nil
}

func (r *Resolver) performerResolver(parties []*ern.PartyResolver) *performerResolver {
	performer := &performerResolver{resolver: r}
	for _, party := range parties {
		if party.Type() != "DisplayArtist" {
			continue
		}
		performer.parties = append(performer.parties, party)
	}
	return performer
}

type performerResolver struct {
	resolver *Resolver
	parties  []*ern.PartyResolver
}

func (p *performerResolver) Identifiers() []*partyIdentifierResolver {
	var identifiers []*partyIdentifierResolver
	for _, party := range p.parties {
		identifiers = append(identifiers, &partyIdentifierResolver{
			typ:    partyIdentifierTypeDPID,
			value:  party.PartyId(),
			source: &sourceResolver{name: party.Source()},
		})
	}
	return identifiers
}

func (p *performerResolver) Name() *stringValueResolver {
	name := &stringValueResolver{}
	for _, party := range p.parties {
		name.value = party.Fullname()
		name.sources = append(name.sources, &stringSourceResolver{
			value:  party.Fullname(),
			source: &sourceResolver{name: party.Source()},
			score:  "1",
		})
	}
	return name
}

func (p *performerResolver) Recordings() ([]*recordingLinkResolver, error) {
	var recordings []*recordingLinkResolver
	for _, party := range p.parties {
		soundRecordings, err := party.SoundRecordings()
		if err != nil {
			return nil, err
		}
		for _, soundRecording := range soundRecordings {
			recordings = append(recordings, &recordingLinkResolver{
				resolver: p.resolver,
				source:   &sourceResolver{name: soundRecording.Source()},
				isrc:     soundRecording.SoundRecordingId(),
			})
		}
	}
	return recordings, nil
}

type composerArgs struct {
	IPI string
}

func (r *Resolver) Composer(args composerArgs) (*composerResolver, error) {
	ipiBaseWriters, err := r.Cwr.WriterControl(cwr.WriterControlArgs{WriterIPIBaseNumber: &args.IPI})
	if err != nil {
		return nil, err
	}
	ipiNameWriters, err := r.Cwr.WriterControl(cwr.WriterControlArgs{WriterIPIName: &args.IPI})
	if err != nil {
		return nil, err
	}
	return r.composerResolver(append(ipiBaseWriters, ipiNameWriters...)), nil
}

func (r *Resolver) composerResolver(writers []*cwr.WriterControlResolver) *composerResolver {
	return &composerResolver{
		resolver: r,
		writers:  writers,
	}
}

type composerResolver struct {
	resolver *Resolver
	writers  []*cwr.WriterControlResolver
}

func (c *composerResolver) Identifiers() []*partyIdentifierResolver {
	var identifiers []*partyIdentifierResolver
	for _, writer := range c.writers {
		if ipi := writer.WriterIPIBaseNumber(); ipi != "" {
			identifiers = append(identifiers, &partyIdentifierResolver{
				typ:    partyIdentifierTypeIPI,
				value:  ipi,
				source: &sourceResolver{name: writer.Source()},
			})
		}
		if ipi := writer.WriterIPIName(); ipi != "" {
			identifiers = append(identifiers, &partyIdentifierResolver{
				typ:    partyIdentifierTypeIPI,
				value:  ipi,
				source: &sourceResolver{name: writer.Source()},
			})
		}
	}
	return identifiers
}

func (c *composerResolver) FirstName() *stringValueResolver {
	name := &stringValueResolver{}
	for _, writer := range c.writers {
		name.value = writer.WriterFirstName()
		name.sources = append(name.sources, &stringSourceResolver{
			value:  writer.WriterFirstName(),
			source: &sourceResolver{name: writer.Source()},
			score:  "1",
		})
	}
	return name
}

func (c *composerResolver) LastName() *stringValueResolver {
	name := &stringValueResolver{}
	for _, writer := range c.writers {
		name.value = writer.WriterLastName()
		name.sources = append(name.sources, &stringSourceResolver{
			value:  writer.WriterLastName(),
			source: &sourceResolver{name: writer.Source()},
			score:  "1",
		})
	}
	return name
}

func (c *composerResolver) Shares() (*sharesResolver, error) {
	var shares []*musicShares
	for _, writer := range c.writers {
		shares = append(shares, &musicShares{
			performance: writer.PROwnershipShare(),
			mechanical:  writer.MROwnershipShare(),
			synch:       writer.SROwnershipShare(),
			source:      writer.Source(),
		})
	}
	return &sharesResolver{resolver: c.resolver, shares: shares}, nil
}

func (c *composerResolver) Works() ([]*workLinkResolver, error) {
	return nil, nil
}

type labelArgs struct {
	DPID string
}

func (r *Resolver) Label(args labelArgs) (*labelResolver, error) {
	// fetch parties by DPID from ERN
	parties, err := r.Ern.Party(ern.PartyArgs{ID: &args.DPID})
	if err != nil {
		return nil, err
	}
	return r.labelResolver(parties), nil
}

func (r *Resolver) labelResolver(parties []*ern.PartyResolver) *labelResolver {
	label := &labelResolver{resolver: r}
	for _, party := range parties {
		if party.Type() != "MessageSender" {
			continue
		}
		label.parties = append(label.parties, party)
	}
	return label
}

type labelResolver struct {
	resolver *Resolver
	parties  []*ern.PartyResolver
}

func (l *labelResolver) Identifiers() []*partyIdentifierResolver {
	var identifiers []*partyIdentifierResolver
	for _, party := range l.parties {
		identifiers = append(identifiers, &partyIdentifierResolver{
			typ:    partyIdentifierTypeDPID,
			value:  party.PartyId(),
			source: &sourceResolver{name: party.Source()},
		})
	}
	return identifiers
}

func (l *labelResolver) Name() *stringValueResolver {
	name := &stringValueResolver{}
	for _, party := range l.parties {
		name.value = party.Fullname()
		name.sources = append(name.sources, &stringSourceResolver{
			value:  party.Fullname(),
			source: &sourceResolver{name: party.Source()},
			score:  "1",
		})
	}
	return name
}

func (l *labelResolver) Releases() ([]*releaseLinkResolver, error) {
	var releases []*releaseLinkResolver
	for _, party := range l.parties {
		ernReleases, err := party.Releases()
		if err != nil {
			return nil, err
		}
		for _, ernRelease := range ernReleases {
			if icpn := ernRelease.ReleaseID().ICPN(); icpn != "" {
				releases = append(releases, &releaseLinkResolver{
					resolver: l.resolver,
					source:   &sourceResolver{name: ernRelease.Source()},
					upc:      icpn,
				})
			}
		}
	}
	return releases, nil
}

type publisherArgs struct {
	IPI string
}

func (r *Resolver) Publisher(args publisherArgs) (*publisherResolver, error) {
	return nil, nil
}

type publisherResolver struct {
	resolver   *Resolver
	publishers []*cwr.PublisherControlResolver
}

func (p *publisherResolver) Identifiers() []*partyIdentifierResolver {
	var identifiers []*partyIdentifierResolver
	for _, publisher := range p.publishers {
		if ipi := publisher.PublisherIPIBaseNumber(); ipi != "" {
			identifiers = append(identifiers, &partyIdentifierResolver{
				typ:    partyIdentifierTypeIPI,
				value:  ipi,
				source: &sourceResolver{name: publisher.Source()},
			})
		}
		if ipi := publisher.PublisherIPINameNumber(); ipi != "" {
			identifiers = append(identifiers, &partyIdentifierResolver{
				typ:    partyIdentifierTypeIPI,
				value:  ipi,
				source: &sourceResolver{name: publisher.Source()},
			})
		}
	}
	return identifiers
}

func (p *publisherResolver) Name() *stringValueResolver {
	name := &stringValueResolver{}
	for _, publisher := range p.publishers {
		name.value = publisher.PublisherName()
		name.sources = append(name.sources, &stringSourceResolver{
			value:  publisher.PublisherName(),
			source: &sourceResolver{name: publisher.Source()},
			score:  "1",
		})
	}
	return name
}

func (p *publisherResolver) Shares() (*sharesResolver, error) {
	var shares []*musicShares
	for _, publisher := range p.publishers {
		shares = append(shares, &musicShares{
			performance: publisher.PROwnershipShare(),
			mechanical:  publisher.MROwnershipShare(),
			synch:       publisher.SROwnershipShare(),
			source:      publisher.Source()})
	}
	return &sharesResolver{resolver: p.resolver, shares: shares}, nil
}

func (p *publisherResolver) Works() ([]*workLinkResolver, error) {
	return nil, nil
}

type recordingArgs struct {
	ISRC string
}

func (r *Resolver) Recording(args recordingArgs) (*recordingResolver, error) {
	// fetch sound recordings by ISRC from ERN
	soundRecordings, err := r.Ern.SoundRecording(ern.SoundRecordingArgs{ID: &args.ISRC})
	if err != nil {
		return nil, err
	}
	return &recordingResolver{
		resolver:        r,
		soundRecordings: soundRecordings,
	}, nil
}

type recordingResolver struct {
	resolver        *Resolver
	soundRecordings []*ern.SoundRecordingResolver
}

func (r *recordingResolver) ISRC() *stringValueResolver {
	isrc := &stringValueResolver{}
	for _, soundRecording := range r.soundRecordings {
		isrc.value = soundRecording.SoundRecordingId()
		isrc.sources = append(isrc.sources, &stringSourceResolver{
			value:  soundRecording.SoundRecordingId(),
			source: &sourceResolver{name: soundRecording.Source()},
			score:  "1",
		})
	}
	return isrc
}

func (r *recordingResolver) Title() *stringValueResolver {
	title := &stringValueResolver{}
	for _, soundRecording := range r.soundRecordings {
		title.value = soundRecording.ReferenceTitle()
		title.sources = append(title.sources, &stringSourceResolver{
			value:  soundRecording.ReferenceTitle(),
			source: &sourceResolver{name: soundRecording.Source()},
			score:  "1",
		})
	}
	return title
}

func (r *recordingResolver) Duration() *stringValueResolver {
	duration := &stringValueResolver{}
	for _, soundRecording := range r.soundRecordings {
		duration.value = soundRecording.Duration()
		duration.sources = append(duration.sources, &stringSourceResolver{
			value:  soundRecording.Duration(),
			source: &sourceResolver{name: soundRecording.Source()},
			score:  "1",
		})
	}
	return duration
}

func (r *recordingResolver) Performers() ([]*performerLinkResolver, error) {
	var performers []*performerLinkResolver
	for _, soundRecording := range r.soundRecordings {
		parties, err := soundRecording.Contributors()
		if err != nil {
			return nil, err
		}
		performers = append(performers, &performerLinkResolver{
			resolver: r.resolver,
			source:   &sourceResolver{name: soundRecording.Source()},
			parties:  parties,
		})
	}
	return performers, nil
}

func (r *recordingResolver) Songs() ([]*songLinkResolver, error) {
	var songs []*songLinkResolver
	for _, soundRecording := range r.soundRecordings {
		ernReleases, err := soundRecording.Releases()
		if err != nil {
			return nil, err
		}
		for _, ernRelease := range ernReleases {
			if grid := ernRelease.ReleaseID().GRID(); grid != "" {
				songs = append(songs, &songLinkResolver{
					resolver: r.resolver,
					source:   &sourceResolver{name: ernRelease.Source()},
					grid:     grid,
				})
			}
		}
	}
	return songs, nil
}

func (r *recordingResolver) Works() ([]*workLinkResolver, error) {
	var works []*workLinkResolver
	for _, soundRecording := range r.soundRecordings {
		isrc := soundRecording.SoundRecordingId()
		args := musicbrainz.RecordingWorkLinkArgs{ISRC: &isrc}
		links, err := r.resolver.MusicBrainz.RecordingWorkLink(args)
		if err != nil {
			return nil, err
		}
		for _, link := range links {
			work := &workLinkResolver{
				resolver: r.resolver,
				source:   &sourceResolver{name: link.Source()},
				iswc:     link.ISWC(),
			}
			workResolver, err := work.Work()
			if err != nil {
				return nil, err
			}
			if len(workResolver.cwrWorks) > 0 {
				works = append(works, work)
			}
		}
	}
	return works, nil
}

type workArgs struct {
	ISWC string
}

func (r *Resolver) Work(args workArgs) (*workResolver, error) {
	cwrWorks, err := r.Cwr.RegisteredWork(cwr.RegisteredWorkArgs{ISWC: &args.ISWC})
	if err != nil {
		return nil, err
	}
	return &workResolver{
		resolver: r,
		cwrWorks: cwrWorks,
	}, nil
}

type workResolver struct {
	resolver *Resolver
	cwrWorks []*cwr.RegisteredWorkResolver
}

func (w *workResolver) ISWC() *stringValueResolver {
	iswc := &stringValueResolver{}
	for _, cwrWork := range w.cwrWorks {
		iswc.value = cwrWork.ISWC()
		iswc.sources = append(iswc.sources, &stringSourceResolver{
			value:  cwrWork.ISWC(),
			source: &sourceResolver{name: cwrWork.Source()},
			score:  "1",
		})
	}
	return iswc
}

func (w *workResolver) Title() *stringValueResolver {
	title := &stringValueResolver{}
	for _, work := range w.cwrWorks {
		title.value = work.Title()
		title.sources = append(title.sources, &stringSourceResolver{
			value:  work.Title(),
			source: &sourceResolver{name: work.Source()},
			score:  "1",
		})
	}
	return title
}

func (w *workResolver) Composers() ([]*composerLinkResolver, error) {
	var composers []*composerLinkResolver
	for _, cwrWork := range w.cwrWorks {
		for _, writer := range cwrWork.Contributors() {
			composers = append(composers, &composerLinkResolver{
				resolver: w.resolver,
				source:   &sourceResolver{name: writer.Source()},
				writers:  []*cwr.WriterControlResolver{writer},
			})
		}
	}
	return composers, nil
}

func (w *workResolver) Publishers() ([]*publisherLinkResolver, error) {
	var publishers []*publisherLinkResolver
	for _, cwrWork := range w.cwrWorks {
		for _, publisher := range cwrWork.Controls() {
			publishers = append(publishers, &publisherLinkResolver{
				resolver:   w.resolver,
				source:     &sourceResolver{name: publisher.Source()},
				publishers: []*cwr.PublisherControlResolver{publisher},
			})
		}
	}
	return publishers, nil
}

func (w *workResolver) Recordings() ([]*recordingLinkResolver, error) {
	var recordings []*recordingLinkResolver
	for _, cwrWork := range w.cwrWorks {
		iswc := cwrWork.ISWC()
		args := musicbrainz.RecordingWorkLinkArgs{ISWC: &iswc}
		links, err := w.resolver.MusicBrainz.RecordingWorkLink(args)
		if err != nil {
			return nil, err
		}
		for _, link := range links {
			recordings = append(recordings, &recordingLinkResolver{
				resolver: w.resolver,
				source:   &sourceResolver{name: link.Source()},
				isrc:     link.ISRC(),
			})
		}
	}
	return recordings, nil
}

type songArgs struct {
	GRID string
}

func (r *Resolver) Song(args songArgs) (*songResolver, error) {
	// fetch songs by GRid from ERN
	var withMainRelease = false
	ernReleases, err := r.Ern.Release(ern.ReleaseArgs{ID: &args.GRID, WithMainRelease: &withMainRelease})
	if err != nil {
		return nil, err
	}
	withMainRelease = true
	ernMainRelease, err := r.Ern.Release(ern.ReleaseArgs{ID: &args.GRID, WithMainRelease: &withMainRelease})
	if err != nil {
		return nil, err
	}

	return &songResolver{
		resolver:       r,
		ernReleases:    ernReleases,
		ernMainRelease: ernMainRelease,
	}, nil
}

type songResolver struct {
	resolver       *Resolver
	ernReleases    []*ern.ReleaseResolver
	ernMainRelease []*ern.ReleaseResolver
}

func (r *songResolver) GRID() *stringValueResolver {
	grid := &stringValueResolver{}
	for _, ernRelease := range r.ernReleases {
		grid.value = ernRelease.ReleaseID().GRID()
		grid.sources = append(grid.sources, &stringSourceResolver{
			value:  ernRelease.ReleaseID().GRID(),
			source: &sourceResolver{name: ernRelease.Source()},
			score:  "1",
		})
	}
	return grid
}

func (r *songResolver) Title() *stringValueResolver {
	title := &stringValueResolver{}
	for _, ernRelease := range r.ernReleases {
		title.value = ernRelease.ReferenceTitle()
		title.sources = append(title.sources, &stringSourceResolver{
			value:  ernRelease.ReferenceTitle(),
			source: &sourceResolver{name: ernRelease.Source()},
			score:  "1",
		})
	}
	return title
}

func (r *songResolver) Recordings() ([]*recordingLinkResolver, error) {
	var recordings []*recordingLinkResolver
	for _, ernRelease := range r.ernReleases {
		soundRecordings, err := ernRelease.SoundRecordings()
		if err != nil {
			return nil, err
		}
		for _, soundRecording := range soundRecordings {
			recordings = append(recordings, &recordingLinkResolver{
				resolver: r.resolver,
				source:   &sourceResolver{name: soundRecording.Source()},
				isrc:     soundRecording.SoundRecordingId(),
			})
		}
	}
	return recordings, nil
}

func (r *songResolver) Releases() []*releaseLinkResolver {
	// if any of the ERN songs have an ICPN, return them as a release
	releases := make([]*releaseLinkResolver, 0, len(r.ernReleases))
	for _, ernRelease := range r.ernMainRelease {
		if icpn := ernRelease.ReleaseID().ICPN(); icpn != "" {
			releases = append(releases, &releaseLinkResolver{
				resolver: r.resolver,
				source:   &sourceResolver{name: ernRelease.Source()},
				upc:      icpn,
			})
		}
	}
	return releases
}

type releaseArgs struct {
	UPC string
}

func (r *Resolver) Release(args releaseArgs) (*releaseResolver, error) {
	// fetch songs by UPC from ERN
	var withMainRelease = true
	ernMainRelease, err := r.Ern.Release(ern.ReleaseArgs{ID: &args.UPC, WithMainRelease: &withMainRelease})
	if err != nil {
		return nil, err
	}
	withMainRelease = false
	ernReleases, err := r.Ern.Release(ern.ReleaseArgs{ID: &args.UPC, WithMainRelease: &withMainRelease})
	//ernReleases, err := r.Ern.Release(ern.ReleaseArgs{ID: &args.UPC})
	if err != nil {
		return nil, err
	}
	return &releaseResolver{
		resolver:       r,
		ernReleases:    ernReleases,
		ernMainRelease: ernMainRelease,
	}, nil
}

type releaseResolver struct {
	resolver       *Resolver
	ernReleases    []*ern.ReleaseResolver
	ernMainRelease []*ern.ReleaseResolver
}

func (p *releaseResolver) UPC() *stringValueResolver {
	upc := &stringValueResolver{}
	for _, ernRelease := range p.ernMainRelease {
		upc.value = ernRelease.ReleaseID().ICPN()
		upc.sources = append(upc.sources, &stringSourceResolver{
			value:  ernRelease.ReleaseID().ICPN(),
			source: &sourceResolver{name: ernRelease.Source()},
			score:  "1",
		})
	}
	return upc
}

func (p *releaseResolver) Title() *stringValueResolver {
	title := &stringValueResolver{}
	for _, ernRelease := range p.ernMainRelease {
		title.value = ernRelease.ReferenceTitle()
		title.sources = append(title.sources, &stringSourceResolver{
			value:  ernRelease.ReferenceTitle(),
			source: &sourceResolver{name: ernRelease.Source()},
			score:  "1",
		})
	}
	return title
}

func (p *releaseResolver) Songs() []*songLinkResolver {
	// if any of the ERN songs have a GRid, return them as a song
	songs := make([]*songLinkResolver, 0, len(p.ernReleases))

	for _, ernRelease := range p.ernReleases {
		if grid := ernRelease.ReleaseID().GRID(); grid != "" {
			songs = append(songs, &songLinkResolver{
				resolver: p.resolver,
				source:   &sourceResolver{name: ernRelease.Source()},
				grid:     grid,
			})
		}
	}
	return songs
}

func (p *releaseResolver) Labels() ([]*labelLinkResolver, error) {
	var labels []*labelLinkResolver
	for _, ernRelease := range p.ernMainRelease {
		parties, err := ernRelease.MessageSenders()
		if err != nil {
			return nil, err
		}
		labels = append(labels, &labelLinkResolver{
			resolver: p.resolver,
			source:   &sourceResolver{name: ernRelease.Source()},
			parties:  parties,
		})
	}
	return labels, nil
}

func (p *releaseResolver) Performer() (*performerLinkResolver, error) {
	parties, err := p.ernMainRelease[0].Contributors()
	if err != nil {
		return nil, err
	}
	return &performerLinkResolver{
		resolver: p.resolver,
		source:   &sourceResolver{name: p.ernMainRelease[0].Source()},
		parties:  parties,
	}, nil
}

type performerLinkResolver struct {
	resolver *Resolver
	source   *sourceResolver
	parties  []*ern.PartyResolver
}

func (r *performerLinkResolver) Source() *sourceResolver {
	return r.source
}

func (r *performerLinkResolver) Performer() (*performerResolver, error) {
	return r.resolver.performerResolver(r.parties), nil
}

func (p *performerLinkResolver) Role() *stringValueResolver {
	role := &stringValueResolver{}
	for _, party := range p.parties {
		if partyRole := party.Role(); partyRole != "" {
			role.value = partyRole
			role.sources = append(role.sources, &stringSourceResolver{
				value:  partyRole,
				source: &sourceResolver{name: party.Source()},
				score:  "1",
			})
		}
	}
	return role
}

type labelLinkResolver struct {
	resolver *Resolver
	source   *sourceResolver
	parties  []*ern.PartyResolver
}

func (r *labelLinkResolver) Source() *sourceResolver {
	return r.source
}

func (r *labelLinkResolver) Label() (*labelResolver, error) {
	return r.resolver.labelResolver(r.parties), nil
}

type recordingLinkResolver struct {
	resolver *Resolver
	source   *sourceResolver
	isrc     string
}

func (r *recordingLinkResolver) Source() *sourceResolver {
	return r.source
}

func (r *recordingLinkResolver) Recording() (*recordingResolver, error) {
	return r.resolver.Recording(recordingArgs{ISRC: r.isrc})
}

type songLinkResolver struct {
	resolver *Resolver
	source   *sourceResolver
	grid     string
}

func (r *songLinkResolver) Source() *sourceResolver {
	return r.source
}

func (r *songLinkResolver) Song() (*songResolver, error) {
	return r.resolver.Song(songArgs{GRID: r.grid})
}

type releaseLinkResolver struct {
	resolver *Resolver
	source   *sourceResolver
	upc      string
}

func (p *releaseLinkResolver) Source() *sourceResolver {
	return p.source
}

func (p *releaseLinkResolver) Release() (*releaseResolver, error) {
	return p.resolver.Release(releaseArgs{UPC: p.upc})
}

type workLinkResolver struct {
	resolver *Resolver
	source   *sourceResolver
	iswc     string
}

func (p *workLinkResolver) Source() *sourceResolver {
	return p.source
}

func (p *workLinkResolver) Work() (*workResolver, error) {
	return p.resolver.Work(workArgs{ISWC: p.iswc})
}

type composerLinkResolver struct {
	resolver *Resolver
	source   *sourceResolver
	writers  []*cwr.WriterControlResolver
}

func (c *composerLinkResolver) Source() *sourceResolver {
	return c.source
}

func (c *composerLinkResolver) Composer() *composerResolver {
	return &composerResolver{
		resolver: c.resolver,
		writers:  c.writers,
	}
}

type publisherLinkResolver struct {
	resolver   *Resolver
	source     *sourceResolver
	publishers []*cwr.PublisherControlResolver
}

func (p *publisherLinkResolver) Source() *sourceResolver {
	return p.source
}

func (p *publisherLinkResolver) Publisher() *publisherResolver {
	return &publisherResolver{
		resolver:   p.resolver,
		publishers: p.publishers,
	}
}

type musicShares struct {
	performance string
	mechanical  string
	synch       string
	source      string
}

type sharesResolver struct {
	resolver *Resolver
	shares   []*musicShares
}

func (s *sharesResolver) Performance() *stringValueResolver {
	performance := &stringValueResolver{}
	for _, publisher := range s.shares {
		performance.value = publisher.performance
		performance.sources = append(performance.sources, &stringSourceResolver{
			value:  publisher.performance,
			source: &sourceResolver{name: publisher.source},
			score:  "1",
		})
	}
	return performance
}

func (s *sharesResolver) Mechanical() *stringValueResolver {
	mechanical := &stringValueResolver{}
	for _, publisher := range s.shares {
		mechanical.value = publisher.mechanical
		mechanical.sources = append(mechanical.sources, &stringSourceResolver{
			value:  publisher.mechanical,
			source: &sourceResolver{name: publisher.source},
			score:  "1",
		})
	}
	return mechanical
}

func (s *sharesResolver) Synch() *stringValueResolver {
	synch := &stringValueResolver{}
	for _, publisher := range s.shares {
		synch.value = publisher.synch
		synch.sources = append(synch.sources, &stringSourceResolver{
			value:  publisher.synch,
			source: &sourceResolver{name: publisher.source},
			score:  "1",
		})
	}
	return synch
}

const (
	partyIdentifierTypeISNI string = "ISNI"
	partyIdentifierTypeIPI  string = "IPI"
	partyIdentifierTypeDPID string = "DPID"
)

type partyIdentifierResolver struct {
	typ    string
	value  string
	source *sourceResolver
}

func (p *partyIdentifierResolver) Type() string {
	return p.typ
}

func (p *partyIdentifierResolver) Value() string {
	return p.value
}

func (p *partyIdentifierResolver) Source() *sourceResolver {
	return p.source
}

type stringValueResolver struct {
	value   string
	sources []*stringSourceResolver
}

func (s *stringValueResolver) Value() string {
	return s.value
}

func (s *stringValueResolver) Sources() []*stringSourceResolver {
	return s.sources
}

type stringSourceResolver struct {
	value  string
	source *sourceResolver
	score  string
}

func (s *stringSourceResolver) Value() string {
	return s.value
}

func (s *stringSourceResolver) Source() *sourceResolver {
	return s.source
}

func (s *stringSourceResolver) Score() string {
	return s.score
}

type sourceResolver struct {
	name string
}

func (s *sourceResolver) Name() string {
	return s.name
}
