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
	"github.com/meta-network/go-meta/identity"
)

// GraphQLSchema is the GraphQL schema for the META Media index.
//
// It consists of the following main entities:
//
// Account
// Performer
// Composer
// RecordLabel
// Publisher
// Recording
// Work
// Song
// Release
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
// GRid    - Global Release Identifier
// UPC     - Universal Product Code
// EAN     - International Article Number
//
// Other entity properties (like performer name or recording title) are linked
// to these identifiers by various sources in the META network.
const GraphQLSchema = `
schema {
  query:    Query
  mutation: Mutation
}

type Query {
  account(meta_id: String!): Account!

  performer(identifier: IdentifierInput!): Performer!

  composer(identifier: IdentifierInput!): Composer!

  record_label(identifier: IdentifierInput!): RecordLabel!

  publisher(identifier: IdentifierInput!): Publisher!

  recording(identifier: IdentifierInput!): Recording!

  work(identifier: IdentifierInput!): Work!

  song(identifier: IdentifierInput!): Song!

  release(identifier: IdentifierInput!): Release!
}

type Mutation {
  createPerformer(performer: PerformerInput!): Performer!

  createComposer(composer: ComposerInput!): Composer!

  createRecordLabel(record_label: RecordLabelInput!): RecordLabel!

  createPublisher(publisher: PublisherInput!): Publisher!

  createRecording(recording: RecordingInput!): Recording!

  createWork(work: WorkInput!): Work!

  createSong(song: SongInput!): Song!

  createRelease(release: ReleaseInput!): Release!

  createPerformerRecordingLink(link: PerformerRecordingLinkInput!): PerformerRecordingLink!

  createComposerWorkLink(link: ComposerWorkLinkInput!): ComposerWorkLink!

  createRecordLabelSongLink(link: RecordLabelSongLinkInput!): RecordLabelSongLink!

  createRecordLabelReleaseLink(link: RecordLabelReleaseLinkInput!): RecordLabelReleaseLink!

  createPublisherWorkLink(link: PublisherWorkLinkInput!): PublisherWorkLink!

  createSongRecordingLink(link: SongRecordingLinkInput!): SongRecordingLink!

  createReleaseRecordingLink(link: ReleaseRecordingLinkInput!): ReleaseRecordingLink!

  createRecordingWorkLink(link: RecordingWorkLinkInput!): RecordingWorkLink!

  createReleaseSongLink(link: ReleaseSongLinkInput!): ReleaseSongLink!
}

#
# --- Main Entities ---
#
type Account {
  performers:    [Performer]!
  record_labels: [RecordLabel]!
  composers:     [Composer]!
}

type Performer {
  identifiers: [IdentifierValue]!

  name: StringValue

  recordings: [PerformerRecordingLink]!
}

type Composer {
  identifiers: [IdentifierValue]!

  firstName: StringValue
  lastName:  StringValue

  works: [ComposerWorkLink]!
}

type RecordLabel {
  identifiers: [IdentifierValue]!

  name: StringValue

  songs:    [RecordLabelSongLink]!
  releases: [RecordLabelReleaseLink]!
}

type Publisher {
  identifiers: [IdentifierValue]!

  name: StringValue

  works: [PublisherWorkLink]!
}

type Recording {
  identifiers: [IdentifierValue]!

  title:    StringValue
  duration: StringValue

  performers: [PerformerRecordingLink]!
  releases:   [ReleaseRecordingLink]!
  works:      [RecordingWorkLink]!
}

type Work {
  identifiers: [IdentifierValue]!

  title: StringValue

  composers:  [ComposerWorkLink]!
  publishers: [PublisherWorkLink]!
  recordings: [RecordingWorkLink]!
}

type Song {
  identifiers: [IdentifierValue]!

  title:    StringValue
  duration: StringValue

  recordings:    [SongRecordingLink]!
  releases:      [ReleaseSongLink]!
  record_labels: [RecordLabelSongLink]!
}

type Release {
  identifiers: [IdentifierValue]!

  type:  StringValue
  title: StringValue
  date:  StringValue

  recordings:    [ReleaseRecordingLink]!
  songs:         [ReleaseSongLink]!
  record_labels: [RecordLabelReleaseLink]!
}

#
# --- Mutation Inputs ---
#
input PerformerInput {
  identifier: IdentifierInput!
  name:       String!
  source:     SourceInput!
}

input ComposerInput {
  identifier: IdentifierInput!
  firstName:  String!
  lastName:   String!
  source:     SourceInput!
}

input RecordLabelInput {
  identifier: IdentifierInput!
  name:       String!
  source:     SourceInput!
}

input PublisherInput {
  identifier: IdentifierInput!
  name:       String!
  source:     SourceInput!
}

input RecordingInput {
  identifier: IdentifierInput!
  title:      String!
  duration:   String!
  source:     SourceInput!
}

input WorkInput {
  identifier: IdentifierInput!
  title:      String!
  source:     SourceInput!
}

input SongInput {
  identifier: IdentifierInput!
  title:      String!
  duration:   String!
  source:     SourceInput!
}

input ReleaseInput {
  identifier: IdentifierInput!
  type:       String!
  title:      String!
  date:       String!
  source:     SourceInput!
}

input PerformerRecordingLinkInput {
  performer_id:  IdentifierInput!
  recording_id:  IdentifierInput!
  role:          String!
  source:        SourceInput!
}

input ComposerWorkLinkInput {
  composer_id: IdentifierInput!
  work_id:     IdentifierInput!
  role:        String!
  source:      SourceInput!
}

input RecordLabelSongLinkInput {
  record_label_id: IdentifierInput!
  song_id:         IdentifierInput!
  source:          SourceInput!
}

input RecordLabelReleaseLinkInput {
  record_label_id: IdentifierInput!
  release_id:      IdentifierInput!
  source:          SourceInput!
}

input PublisherWorkLinkInput {
  publisher_id: IdentifierInput!
  work_id:      IdentifierInput!
  source:       SourceInput!
}

input SongRecordingLinkInput {
  song_id:      IdentifierInput!
  recording_id: IdentifierInput!
  source:       SourceInput!
}

input ReleaseRecordingLinkInput {
  release_id:   IdentifierInput!
  recording_id: IdentifierInput!
  source:       SourceInput!
}

input RecordingWorkLinkInput {
  recording_id: IdentifierInput!
  work_id:      IdentifierInput!
  source:       SourceInput!
}

input ReleaseSongLinkInput {
  release_id: IdentifierInput!
  song_id:    IdentifierInput!
  source:     SourceInput!
}

input IdentifierInput {
  type:  String!
  value: String!
}

input SourceInput {
  name: String!
}

#
# --- Link Types ---
#

type PerformerRecordingLink {
  performer: Performer!
  recording: Recording!
  role:      String!
  source:    Source!
}

type ComposerWorkLink {
  composer: Composer!
  work:     Work!
  role:     String!
  source:   Source!
}

type RecordLabelSongLink {
  record_label: RecordLabel!
  song:         Song!
  source:       Source!
}

type RecordLabelReleaseLink {
  record_label: RecordLabel!
  release:      Release!
  source:       Source!
}

type PublisherWorkLink {
  publisher: Publisher!
  work:      Work!
  source:    Source!
}

type SongRecordingLink {
  song:      Song!
  recording: Recording!
  source:    Source!
}

type ReleaseRecordingLink {
  release:   Release!
  recording: Recording!
  source:    Source!
}

type RecordingWorkLink {
  recording: Recording!
  work:      Work!
  source:    Source!
}

type ReleaseSongLink {
  release: Release!
  song:    Song!
  source:  Source!
}

#
# --- Value Types ---
#

type IdentifierValue {
  value:   Identifier!
  sources: [Source]!
}

type Identifier {
  type:   String!
  value:  String!
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
// the GraphQLSchema constant, storing and retrieving data from a Media index.
type Resolver struct {
	mediaIndex    *Index
	identityIndex *identity.Index
}

func NewResolver(mediaIndex *Index, identityIndex *identity.Index) *Resolver {
	return &Resolver{mediaIndex, identityIndex}
}

type AccountArgs struct {
	MetaID string
}

func (r *Resolver) Account(args AccountArgs) (*AccountResolver, error) {
	identity, err := r.identityIndex.Identity(args.MetaID)
	if err != nil {
		return nil, err
	}
	return &AccountResolver{r, identity}, nil
}

type AccountResolver struct {
	resolver *Resolver
	identity *identity.Identity
}

func (a *AccountResolver) Performers() ([]*PerformerResolver, error) {
	identifiers, err := a.identifiers("dpid")
	if err != nil {
		return nil, err
	}
	resolvers := make([]*PerformerResolver, 0, len(identifiers))
	for _, id := range identifiers {
		identifier, err := a.resolver.mediaIndex.Identifier(id)
		if err == nil {
			resolvers = append(resolvers, &PerformerResolver{a.resolver, identifier})
		} else if !isIdentifierNotFound(err) {
			return nil, err
		}
	}
	return resolvers, nil
}

func (a *AccountResolver) RecordLabels() ([]*RecordLabelResolver, error) {
	identifiers, err := a.identifiers("dpid")
	if err != nil {
		return nil, err
	}
	resolvers := make([]*RecordLabelResolver, 0, len(identifiers))
	for _, id := range identifiers {
		identifier, err := a.resolver.mediaIndex.Identifier(id)
		if err == nil {
			resolvers = append(resolvers, &RecordLabelResolver{a.resolver, identifier})
		} else if !isIdentifierNotFound(err) {
			return nil, err
		}
	}
	return resolvers, nil
}

func (a *AccountResolver) Composers() ([]*ComposerResolver, error) {
	identifiers, err := a.identifiers("ipi")
	if err != nil {
		return nil, err
	}
	resolvers := make([]*ComposerResolver, 0, len(identifiers))
	for _, id := range identifiers {
		identifier, err := a.resolver.mediaIndex.Identifier(id)
		if err == nil {
			resolvers = append(resolvers, &ComposerResolver{a.resolver, identifier})
		} else if !isIdentifierNotFound(err) {
			return nil, err
		}
	}
	return resolvers, nil
}

func (a *AccountResolver) identifiers(typ string) ([]*Identifier, error) {
	id := a.identity.ID().String()
	claims, err := a.resolver.identityIndex.Claims(identity.ClaimFilter{
		Subject:  &id,
		Property: &typ,
	})
	if err != nil {
		return nil, err
	}
	identifiers := make(map[Identifier]struct{})
	for _, claim := range claims {
		identifiers[Identifier{Type: typ, Value: claim.Claim}] = struct{}{}
	}
	res := make([]*Identifier, 0, len(identifiers))
	for identifier := range identifiers {
		res = append(res, &identifier)
	}
	return res, nil
}

type IdentifierArgs struct {
	Identifier Identifier
}

type createPerformerArgs struct {
	Performer struct {
		Performer

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreatePerformer(args createPerformerArgs) (*PerformerResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Performer.Performer,
		&args.Performer.Identifier,
		&args.Performer.Source,
	)
	if err != nil {
		return nil, err
	}
	return &PerformerResolver{r, identifier}, nil
}

func (r *Resolver) Performer(args IdentifierArgs) (*PerformerResolver, error) {
	identifier, err := r.mediaIndex.Identifier(&args.Identifier)
	if err != nil {
		return nil, err
	}
	return &PerformerResolver{r, identifier}, nil
}

type PerformerResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (p *PerformerResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{p.resolver, p.identifier}}
}

func (p *PerformerResolver) Name() (*stringValueResolver, error) {
	records, err := p.resolver.mediaIndex.Performers(p.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.Name
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.Name,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (p *PerformerResolver) Recordings() ([]*performerRecordingLinkResolver, error) {
	records, err := p.resolver.mediaIndex.PerformerRecordings(p.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*performerRecordingLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &performerRecordingLinkResolver{
			resolver: p.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createComposerArgs struct {
	Composer struct {
		Composer

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateComposer(args createComposerArgs) (*ComposerResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Composer.Composer,
		&args.Composer.Identifier,
		&args.Composer.Source,
	)
	if err != nil {
		return nil, err
	}
	return &ComposerResolver{r, identifier}, nil
}

func (r *Resolver) Composer(args IdentifierArgs) (*ComposerResolver, error) {
	identifier, err := r.mediaIndex.Identifier(&args.Identifier)
	if err != nil {
		return nil, err
	}
	return &ComposerResolver{r, identifier}, nil
}

type ComposerResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (c *ComposerResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{c.resolver, c.identifier}}
}

func (c *ComposerResolver) FirstName() (*stringValueResolver, error) {
	records, err := c.resolver.mediaIndex.Composers(c.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.FirstName
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.FirstName,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (c *ComposerResolver) LastName() (*stringValueResolver, error) {
	records, err := c.resolver.mediaIndex.Composers(c.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.LastName
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.LastName,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (c *ComposerResolver) Works() ([]*composerWorkLinkResolver, error) {
	records, err := c.resolver.mediaIndex.ComposerWorks(c.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*composerWorkLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &composerWorkLinkResolver{
			resolver: c.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createRecordLabelArgs struct {
	RecordLabel struct {
		RecordLabel

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateRecordLabel(args createRecordLabelArgs) (*RecordLabelResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.RecordLabel.RecordLabel,
		&args.RecordLabel.Identifier,
		&args.RecordLabel.Source,
	)
	if err != nil {
		return nil, err
	}
	return &RecordLabelResolver{r, identifier}, nil
}

func (r *Resolver) RecordLabel(args IdentifierArgs) (*RecordLabelResolver, error) {
	identifier, err := r.mediaIndex.Identifier(&args.Identifier)
	if err != nil {
		return nil, err
	}
	return &RecordLabelResolver{r, identifier}, nil
}

type RecordLabelResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (l *RecordLabelResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{l.resolver, l.identifier}}
}

func (l *RecordLabelResolver) Name() (*stringValueResolver, error) {
	records, err := l.resolver.mediaIndex.RecordLabels(l.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.Name
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.Name,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (l *RecordLabelResolver) Songs() ([]*recordLabelSongLinkResolver, error) {
	records, err := l.resolver.mediaIndex.RecordLabelSongs(l.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*recordLabelSongLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &recordLabelSongLinkResolver{
			resolver: l.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (l *RecordLabelResolver) Releases() ([]*recordLabelReleaseLinkResolver, error) {
	records, err := l.resolver.mediaIndex.RecordLabelReleases(l.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*recordLabelReleaseLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &recordLabelReleaseLinkResolver{
			resolver: l.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createPublisherArgs struct {
	Publisher struct {
		Publisher

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreatePublisher(args createPublisherArgs) (*PublisherResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Publisher.Publisher,
		&args.Publisher.Identifier,
		&args.Publisher.Source,
	)
	if err != nil {
		return nil, err
	}
	return &PublisherResolver{r, identifier}, nil
}

func (r *Resolver) Publisher(args IdentifierArgs) (*PublisherResolver, error) {
	identifier, err := r.mediaIndex.Identifier(&args.Identifier)
	if err != nil {
		return nil, err
	}
	return &PublisherResolver{r, identifier}, nil
}

type PublisherResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (p *PublisherResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{p.resolver, p.identifier}}
}

func (p *PublisherResolver) Name() (*stringValueResolver, error) {
	records, err := p.resolver.mediaIndex.Publishers(p.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.Name
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.Name,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (p *PublisherResolver) Works() ([]*PublisherWorkLinkResolver, error) {
	records, err := p.resolver.mediaIndex.PublisherWorks(p.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*PublisherWorkLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &PublisherWorkLinkResolver{
			resolver: p.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createRecordingArgs struct {
	Recording struct {
		Recording

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateRecording(args createRecordingArgs) (*RecordingResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Recording.Recording,
		&args.Recording.Identifier,
		&args.Recording.Source,
	)
	if err != nil {
		return nil, err
	}
	return &RecordingResolver{r, identifier}, nil
}

func (r *Resolver) Recording(args IdentifierArgs) (*RecordingResolver, error) {
	identifier, err := r.mediaIndex.Identifier(&args.Identifier)
	if err != nil {
		return nil, err
	}
	return &RecordingResolver{r, identifier}, nil
}

type RecordingResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (r *RecordingResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{r.resolver, r.identifier}}
}

func (r *RecordingResolver) Title() (*stringValueResolver, error) {
	records, err := r.resolver.mediaIndex.Recordings(r.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.Title
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.Title,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (r *RecordingResolver) Duration() (*stringValueResolver, error) {
	records, err := r.resolver.mediaIndex.Recordings(r.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.Duration
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.Duration,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (r *RecordingResolver) Performers() ([]*performerRecordingLinkResolver, error) {
	records, err := r.resolver.mediaIndex.RecordingPerformers(r.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*performerRecordingLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &performerRecordingLinkResolver{
			resolver: r.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (r *RecordingResolver) Releases() ([]*releaseRecordingLinkResolver, error) {
	records, err := r.resolver.mediaIndex.RecordingReleases(r.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*releaseRecordingLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &releaseRecordingLinkResolver{
			resolver: r.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (r *RecordingResolver) Works() ([]*recordingWorkLinkResolver, error) {
	records, err := r.resolver.mediaIndex.RecordingWorks(r.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*recordingWorkLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &recordingWorkLinkResolver{
			resolver: r.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createWorkArgs struct {
	Work struct {
		Work

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateWork(args createWorkArgs) (*WorkResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Work.Work,
		&args.Work.Identifier,
		&args.Work.Source,
	)
	if err != nil {
		return nil, err
	}
	return &WorkResolver{r, identifier}, nil
}

func (r *Resolver) Work(args IdentifierArgs) (*WorkResolver, error) {
	identifier, err := r.mediaIndex.Identifier(&args.Identifier)
	if err != nil {
		return nil, err
	}
	return &WorkResolver{r, identifier}, nil
}

type WorkResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (w *WorkResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{w.resolver, w.identifier}}
}

func (w *WorkResolver) Title() (*stringValueResolver, error) {
	records, err := w.resolver.mediaIndex.Works(w.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.Title
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.Title,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (w *WorkResolver) Composers() ([]*composerWorkLinkResolver, error) {
	records, err := w.resolver.mediaIndex.WorkComposers(w.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*composerWorkLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &composerWorkLinkResolver{
			resolver: w.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (w *WorkResolver) Publishers() ([]*PublisherWorkLinkResolver, error) {
	records, err := w.resolver.mediaIndex.WorkPublishers(w.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*PublisherWorkLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &PublisherWorkLinkResolver{
			resolver: w.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (w *WorkResolver) Recordings() ([]*recordingWorkLinkResolver, error) {
	records, err := w.resolver.mediaIndex.WorkRecordings(w.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*recordingWorkLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &recordingWorkLinkResolver{
			resolver: w.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createSongArgs struct {
	Song struct {
		Song

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateSong(args createSongArgs) (*SongResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Song.Song,
		&args.Song.Identifier,
		&args.Song.Source,
	)
	if err != nil {
		return nil, err
	}
	return &SongResolver{r, identifier}, nil
}

func (r *Resolver) Song(args IdentifierArgs) (*SongResolver, error) {
	identifier, err := r.mediaIndex.Identifier(&args.Identifier)
	if err != nil {
		return nil, err
	}
	return &SongResolver{r, identifier}, nil
}

type SongResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (s *SongResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{s.resolver, s.identifier}}
}

func (s *SongResolver) Title() (*stringValueResolver, error) {
	records, err := s.resolver.mediaIndex.Songs(s.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.Title
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.Title,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (s *SongResolver) Duration() (*stringValueResolver, error) {
	records, err := s.resolver.mediaIndex.Songs(s.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.Duration
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.Duration,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (s *SongResolver) Recordings() ([]*songRecordingLinkResolver, error) {
	records, err := s.resolver.mediaIndex.RecordingSongs(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*songRecordingLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &songRecordingLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SongResolver) Releases() ([]*releaseSongLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SongReleases(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*releaseSongLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &releaseSongLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SongResolver) RecordLabels() ([]*recordLabelSongLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SongRecordLabels(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*recordLabelSongLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &recordLabelSongLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createReleaseArgs struct {
	Release struct {
		Release

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateRelease(args createReleaseArgs) (*ReleaseResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Release.Release,
		&args.Release.Identifier,
		&args.Release.Source,
	)
	if err != nil {
		return nil, err
	}
	return &ReleaseResolver{r, identifier}, nil
}

func (r *Resolver) Release(args IdentifierArgs) (*ReleaseResolver, error) {
	identifier, err := r.mediaIndex.Identifier(&args.Identifier)
	if err != nil {
		return nil, err
	}
	return &ReleaseResolver{r, identifier}, nil
}

type ReleaseResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (r *ReleaseResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{r.resolver, r.identifier}}
}

func (r *ReleaseResolver) Type() (*stringValueResolver, error) {
	records, err := r.resolver.mediaIndex.Releases(r.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.Type
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.Type,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (r *ReleaseResolver) Title() (*stringValueResolver, error) {
	records, err := r.resolver.mediaIndex.Releases(r.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.Title
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.Title,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (r *ReleaseResolver) Date() (*stringValueResolver, error) {
	records, err := r.resolver.mediaIndex.Releases(r.identifier)
	if err != nil {
		return nil, err
	}
	resolver := &stringValueResolver{}
	for _, record := range records {
		resolver.value = record.Date
		resolver.sources = append(resolver.sources, &stringSourceResolver{
			value:  record.Date,
			source: &sourceResolver{id: record.Source},
			score:  "1",
		})
	}
	return resolver, nil
}

func (r *ReleaseResolver) Recordings() ([]*releaseRecordingLinkResolver, error) {
	records, err := r.resolver.mediaIndex.ReleaseRecordings(r.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*releaseRecordingLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &releaseRecordingLinkResolver{
			resolver: r.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (r *ReleaseResolver) Songs() ([]*releaseSongLinkResolver, error) {
	records, err := r.resolver.mediaIndex.ReleaseSongs(r.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*releaseSongLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &releaseSongLinkResolver{
			resolver: r.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (r *ReleaseResolver) RecordLabels() ([]*recordLabelReleaseLinkResolver, error) {
	records, err := r.resolver.mediaIndex.ReleaseRecordLabels(r.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*recordLabelReleaseLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &recordLabelReleaseLinkResolver{
			resolver: r.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createPerformerRecordingLinkArgs struct {
	Link struct {
		PerformerID Identifier
		RecordingID Identifier
		Role        string
		Source      Source
	}
}

func (r *Resolver) CreatePerformerRecordingLink(args createPerformerRecordingLinkArgs) (*performerRecordingLinkResolver, error) {
	link := &PerformerRecordingLink{
		Performer: args.Link.PerformerID,
		Recording: args.Link.RecordingID,
		Role:      args.Link.Role,
	}
	record, err := r.mediaIndex.CreatePerformerRecording(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &performerRecordingLinkResolver{r, record}, nil
}

type performerRecordingLinkResolver struct {
	resolver *Resolver
	record   *PerformerRecordingRecord
}

func (p *performerRecordingLinkResolver) Performer() *PerformerResolver {
	return &PerformerResolver{p.resolver, p.record.Performer}
}

func (p *performerRecordingLinkResolver) Recording() *RecordingResolver {
	return &RecordingResolver{p.resolver, p.record.Recording}
}

func (p *performerRecordingLinkResolver) Role() string {
	return p.record.Role
}

func (p *performerRecordingLinkResolver) Source() *sourceResolver {
	return &sourceResolver{p.resolver, p.record.Source}
}

type createComposerWorkLinkArgs struct {
	Link struct {
		ComposerID Identifier
		WorkID     Identifier
		Role       string
		Source     Source
	}
}

func (r *Resolver) CreateComposerWorkLink(args createComposerWorkLinkArgs) (*composerWorkLinkResolver, error) {
	link := &ComposerWorkLink{
		Composer: args.Link.ComposerID,
		Work:     args.Link.WorkID,
		Role:     args.Link.Role,
	}
	record, err := r.mediaIndex.CreateComposerWork(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &composerWorkLinkResolver{r, record}, nil
}

type composerWorkLinkResolver struct {
	resolver *Resolver
	record   *ComposerWorkRecord
}

func (c *composerWorkLinkResolver) Composer() *ComposerResolver {
	return &ComposerResolver{c.resolver, c.record.Composer}
}

func (c *composerWorkLinkResolver) Work() *WorkResolver {
	return &WorkResolver{c.resolver, c.record.Work}
}

func (c *composerWorkLinkResolver) Role() string {
	return c.record.Role
}

func (c *composerWorkLinkResolver) Source() *sourceResolver {
	return &sourceResolver{c.resolver, c.record.Source}
}

type createRecordLabelSongLinkArgs struct {
	Link struct {
		RecordLabelID Identifier
		SongID        Identifier
		Source        Source
	}
}

func (r *Resolver) CreateRecordLabelSongLink(args createRecordLabelSongLinkArgs) (*recordLabelSongLinkResolver, error) {
	link := &RecordLabelSongLink{
		RecordLabel: args.Link.RecordLabelID,
		Song:        args.Link.SongID,
	}
	record, err := r.mediaIndex.CreateRecordLabelSong(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &recordLabelSongLinkResolver{r, record}, nil
}

type recordLabelSongLinkResolver struct {
	resolver *Resolver
	record   *RecordLabelSongRecord
}

func (r *recordLabelSongLinkResolver) RecordLabel() *RecordLabelResolver {
	return &RecordLabelResolver{r.resolver, r.record.RecordLabel}
}

func (r *recordLabelSongLinkResolver) Song() *SongResolver {
	return &SongResolver{r.resolver, r.record.Song}
}

func (r *recordLabelSongLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createRecordLabelReleaseLinkArgs struct {
	Link struct {
		RecordLabelID Identifier
		ReleaseID     Identifier
		Source        Source
	}
}

func (r *Resolver) CreateRecordLabelReleaseLink(args createRecordLabelReleaseLinkArgs) (*recordLabelReleaseLinkResolver, error) {
	link := &RecordLabelReleaseLink{
		RecordLabel: args.Link.RecordLabelID,
		Release:     args.Link.ReleaseID,
	}
	record, err := r.mediaIndex.CreateRecordLabelRelease(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &recordLabelReleaseLinkResolver{r, record}, nil
}

type recordLabelReleaseLinkResolver struct {
	resolver *Resolver
	record   *RecordLabelReleaseRecord
}

func (r *recordLabelReleaseLinkResolver) RecordLabel() *RecordLabelResolver {
	return &RecordLabelResolver{r.resolver, r.record.RecordLabel}
}

func (r *recordLabelReleaseLinkResolver) Release() *ReleaseResolver {
	return &ReleaseResolver{r.resolver, r.record.Release}
}

func (r *recordLabelReleaseLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createPublisherWorkLinkArgs struct {
	Link struct {
		PublisherID Identifier
		WorkID      Identifier
		Source      Source
	}
}

func (r *Resolver) CreatePublisherWorkLink(args createPublisherWorkLinkArgs) (*PublisherWorkLinkResolver, error) {
	link := &PublisherWorkLink{
		Publisher: args.Link.PublisherID,
		Work:      args.Link.WorkID,
	}
	record, err := r.mediaIndex.CreatePublisherWork(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &PublisherWorkLinkResolver{r, record}, nil
}

type PublisherWorkLinkResolver struct {
	resolver *Resolver
	record   *PublisherWorkRecord
}

func (p *PublisherWorkLinkResolver) Publisher() *PublisherResolver {
	return &PublisherResolver{p.resolver, p.record.Publisher}
}

func (p *PublisherWorkLinkResolver) Work() *WorkResolver {
	return &WorkResolver{p.resolver, p.record.Work}
}

func (p *PublisherWorkLinkResolver) Source() *sourceResolver {
	return &sourceResolver{p.resolver, p.record.Source}
}

type createSongRecordingLinkArgs struct {
	Link struct {
		SongID      Identifier
		RecordingID Identifier
		Source      Source
	}
}

func (r *Resolver) CreateSongRecordingLink(args createSongRecordingLinkArgs) (*songRecordingLinkResolver, error) {
	link := &SongRecordingLink{
		Song:      args.Link.SongID,
		Recording: args.Link.RecordingID,
	}
	record, err := r.mediaIndex.CreateSongRecording(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &songRecordingLinkResolver{r, record}, nil
}

type songRecordingLinkResolver struct {
	resolver *Resolver
	record   *SongRecordingRecord
}

func (s *songRecordingLinkResolver) Song() *SongResolver {
	return &SongResolver{s.resolver, s.record.Song}
}

func (s *songRecordingLinkResolver) Recording() *RecordingResolver {
	return &RecordingResolver{s.resolver, s.record.Recording}
}

func (s *songRecordingLinkResolver) Source() *sourceResolver {
	return &sourceResolver{s.resolver, s.record.Source}
}

type createReleaseRecordingLinkArgs struct {
	Link struct {
		ReleaseID   Identifier
		RecordingID Identifier
		Source      Source
	}
}

func (r *Resolver) CreateReleaseRecordingLink(args createReleaseRecordingLinkArgs) (*releaseRecordingLinkResolver, error) {
	link := &ReleaseRecordingLink{
		Release:   args.Link.ReleaseID,
		Recording: args.Link.RecordingID,
	}
	record, err := r.mediaIndex.CreateReleaseRecording(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &releaseRecordingLinkResolver{r, record}, nil
}

type releaseRecordingLinkResolver struct {
	resolver *Resolver
	record   *ReleaseRecordingRecord
}

func (r *releaseRecordingLinkResolver) Release() *ReleaseResolver {
	return &ReleaseResolver{r.resolver, r.record.Release}
}

func (r *releaseRecordingLinkResolver) Recording() *RecordingResolver {
	return &RecordingResolver{r.resolver, r.record.Recording}
}

func (r *releaseRecordingLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createRecordingWorkLinkArgs struct {
	Link struct {
		RecordingID Identifier
		WorkID      Identifier
		Source      Source
	}
}

func (r *Resolver) CreateRecordingWorkLink(args createRecordingWorkLinkArgs) (*recordingWorkLinkResolver, error) {
	link := &RecordingWorkLink{
		Recording: args.Link.RecordingID,
		Work:      args.Link.WorkID,
	}
	record, err := r.mediaIndex.CreateRecordingWork(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &recordingWorkLinkResolver{r, record}, nil
}

type recordingWorkLinkResolver struct {
	resolver *Resolver
	record   *RecordingWorkRecord
}

func (r *recordingWorkLinkResolver) Recording() *RecordingResolver {
	return &RecordingResolver{r.resolver, r.record.Recording}
}

func (r *recordingWorkLinkResolver) Work() *WorkResolver {
	return &WorkResolver{r.resolver, r.record.Work}
}

func (r *recordingWorkLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createReleaseSongLinkArgs struct {
	Link struct {
		ReleaseID Identifier
		SongID    Identifier
		Source    Source
	}
}

func (r *Resolver) CreateReleaseSongLink(args createReleaseSongLinkArgs) (*releaseSongLinkResolver, error) {
	link := &ReleaseSongLink{
		Release: args.Link.ReleaseID,
		Song:    args.Link.SongID,
	}
	record, err := r.mediaIndex.CreateReleaseSong(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &releaseSongLinkResolver{r, record}, nil
}

type releaseSongLinkResolver struct {
	resolver *Resolver
	record   *ReleaseSongRecord
}

func (r *releaseSongLinkResolver) Release() *ReleaseResolver {
	return &ReleaseResolver{r.resolver, r.record.Release}
}

func (r *releaseSongLinkResolver) Song() *SongResolver {
	return &SongResolver{r.resolver, r.record.Song}
}

func (r *releaseSongLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
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
	resolver *Resolver
	id       int64
}

func (s *sourceResolver) Name() (string, error) {
	source, err := s.resolver.mediaIndex.Source(s.id)
	if err != nil {
		return "", err
	}
	return source.Name, nil
}

type identifierValueResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (i *identifierValueResolver) Value() *identifierResolver {
	return &identifierResolver{i.resolver, i.identifier}
}

func (i *identifierValueResolver) Sources() []*sourceResolver {
	return nil
}

type identifierResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (i *identifierResolver) Type() string {
	return i.identifier.Type
}

func (i *identifierResolver) Value() string {
	return i.identifier.Value
}
