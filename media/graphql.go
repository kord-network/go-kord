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
// Contributor
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

  contributor(identifier: IdentifierInput!): Contributor!

  composer(identifier: IdentifierInput!): Composer!

  record_label(identifier: IdentifierInput!): RecordLabel!

  publisher(identifier: IdentifierInput!): Publisher!

  recording(identifier: IdentifierInput!): Recording!

  work(identifier: IdentifierInput!): Work!

  song(identifier: IdentifierInput!): Song!

  release(identifier: IdentifierInput!): Release!

  organisation(identifier: IdentifierInput!): Organisation!

  series(identifier: IdentifierInput!): Series!

  season(identifier: IdentifierInput!): Season!

  episode(identifier: IdentifierInput!): Episode!

  supplemental(identifier: IdentifierInput!): Supplemental!
}

type Mutation {
  createPerformer(performer: PerformerInput!): Performer!

  createContributor(contributor: ContributorInput!): Contributor!

  createComposer(composer: ComposerInput!): Composer!

  createRecordLabel(record_label: RecordLabelInput!): RecordLabel!

  createPublisher(publisher: PublisherInput!): Publisher!

  createRecording(recording: RecordingInput!): Recording!

  createWork(work: WorkInput!): Work!

  createSong(song: SongInput!): Song!

  createRelease(release: ReleaseInput!): Release!

  createOrganisation(organisation: OrganisationInput!): Organisation!

  createSeries(series: SeriesInput!): Series!

  createSeason(season: SeasonInput!): Season!

  createEpisode(episode: EpisodeInput!): Episode!

  createSupplemental(supplemental: SupplementalInput!): Supplemental!

  createPerformerRecordingLink(link: PerformerRecordingLinkInput!): PerformerRecordingLink!

  createPerformerSongLink(link: PerformerSongLinkInput!): PerformerSongLink!

  createPerformerReleaseLink(link: PerformerReleaseLinkInput!): PerformerReleaseLink!

  createContributorRecordingLink(link: ContributorRecordingLinkInput!): ContributorRecordingLink!

  createComposerWorkLink(link: ComposerWorkLinkInput!): ComposerWorkLink!

  createRecordLabelRecordingLink(link: RecordLabelRecordingLinkInput!): RecordLabelRecordingLink!

  createRecordLabelSongLink(link: RecordLabelSongLinkInput!): RecordLabelSongLink!

  createRecordLabelReleaseLink(link: RecordLabelReleaseLinkInput!): RecordLabelReleaseLink!

  createPublisherWorkLink(link: PublisherWorkLinkInput!): PublisherWorkLink!

  createSongRecordingLink(link: SongRecordingLinkInput!): SongRecordingLink!

  createReleaseRecordingLink(link: ReleaseRecordingLinkInput!): ReleaseRecordingLink!

  createRecordingWorkLink(link: RecordingWorkLinkInput!): RecordingWorkLink!

  createReleaseSongLink(link: ReleaseSongLinkInput!): ReleaseSongLink!

  createOrganisationSeriesLink(link: OrganisationSeriesLinkInput!): OrganisationSeriesLink!

  createOrganisationSeasonLink(link: OrganisationSeasonLinkInput!): OrganisationSeasonLink!

  createOrganisationEpisodeLink(link: OrganisationEpisodeLinkInput!): OrganisationEpisodeLink!

  createOrganisationSupplementalLink(link: OrganisationSupplementalLinkInput!): OrganisationSupplementalLink!

  createSeriesSeasonLink(link: SeriesSeasonLinkInput!): SeriesSeasonLink!

  createSeriesEpisodeLink(link: SeriesEpisodeLinkInput!): SeriesEpisodeLink!

  createSeriesSupplementalLink(link: SeriesSupplementalLinkInput!): SeriesSupplementalLink!

  createSeasonEpisodeLink(link: SeasonEpisodeLinkInput!): SeasonEpisodeLink!

  createSeasonSupplementalLink(link: SeasonSupplementalLinkInput!): SeasonSupplementalLink!

  createEpisodeSupplementalLink(link: EpisodeSupplementalLinkInput!): EpisodeSupplementalLink!
}

#
# --- Main Entities ---
#
type Account {
  performers:    [Performer]!
  record_labels: [RecordLabel]!
  composers:     [Composer]!
  publishers:    [Publisher]!
  organisations: [Organisation]!
}

type Performer {
  identifiers: [IdentifierValue]!

  name: StringValue

  recordings: [PerformerRecordingLink]!
  songs:      [PerformerSongLink]!
  releases:   [PerformerReleaseLink]!
}

type Contributor {
  identifiers: [IdentifierValue]!

  name: StringValue

  recordings: [ContributorRecordingLink]!
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

  recordings: [RecordLabelRecordingLink]!
  songs:      [RecordLabelSongLink]!
  releases:   [RecordLabelReleaseLink]!
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

  performers:    [PerformerRecordingLink]!
  contributors:  [ContributorRecordingLink]!
  songs:         [SongRecordingLink]!
  releases:      [ReleaseRecordingLink]!
  record_labels: [RecordLabelRecordingLink]!
  works:         [RecordingWorkLink]!
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

  performers:    [PerformerSongLink]!
  recordings:    [SongRecordingLink]!
  releases:      [ReleaseSongLink]!
  record_labels: [RecordLabelSongLink]!
}

type Release {
  identifiers: [IdentifierValue]!

  type:  StringValue
  title: StringValue
  date:  StringValue

  performers:    [PerformerReleaseLink]!
  recordings:    [ReleaseRecordingLink]!
  songs:         [ReleaseSongLink]!
  record_labels: [RecordLabelReleaseLink]!
}

type Organisation {
  identifiers: [IdentifierValue]!

  name: StringValue

  series:        [OrganisationSeriesLink]!
  seasons:       [OrganisationSeasonLink]!
  episodes:      [OrganisationEpisodeLink]!
  supplementals: [OrganisationSupplementalLink]!
}

type Series {
  identifiers: [IdentifierValue]!

  name: StringValue

  organisations: [OrganisationSeriesLink]!
  seasons:       [SeriesSeasonLink]!
  episodes:      [SeriesEpisodeLink]!
  supplementals: [SeriesSupplementalLink]!
}

type Season {
  identifiers: [IdentifierValue]!

  name: StringValue

  organisations: [OrganisationSeasonLink]!
  series:        [SeriesSeasonLink]!
  episodes:      [SeasonEpisodeLink]!
  supplementals: [SeasonSupplementalLink]!
}

type Episode {
  identifiers: [IdentifierValue]!

  name: StringValue

  organisations: [OrganisationEpisodeLink]!
  series:        [SeriesEpisodeLink]!
  seasons:       [SeasonEpisodeLink]!
  supplementals: [EpisodeSupplementalLink]!
}

type Supplemental {
  identifiers: [IdentifierValue]!

  name: StringValue

  organisations: [OrganisationSupplementalLink]!
  series:        [SeriesSupplementalLink]!
  seasons:       [SeasonSupplementalLink]!
  episodes:      [EpisodeSupplementalLink]!
}

#
# --- Mutation Inputs ---
#
input PerformerInput {
  identifier: IdentifierInput!
  name:       String!
  source:     SourceInput!
}

input ContributorInput {
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

input OrganisationInput {
  identifier: IdentifierInput!
  name:       String!
  source:     SourceInput!
}

input SeriesInput {
  identifier: IdentifierInput!
  name:       String!
  source:     SourceInput!
}

input SeasonInput {
  identifier: IdentifierInput!
  name:       String!
  source:     SourceInput!
}

input EpisodeInput {
  identifier: IdentifierInput!
  name:       String!
  source:     SourceInput!
}

input SupplementalInput {
  identifier: IdentifierInput!
  name:       String!
  source:     SourceInput!
}

input PerformerRecordingLinkInput {
  performer_id:  IdentifierInput!
  recording_id:  IdentifierInput!
  role:          String!
  source:        SourceInput!
}

input PerformerSongLinkInput {
  performer_id: IdentifierInput!
  song_id:      IdentifierInput!
  role:         String!
  source:       SourceInput!
}

input PerformerReleaseLinkInput {
  performer_id: IdentifierInput!
  release_id:   IdentifierInput!
  role:         String!
  source:       SourceInput!
}

input ContributorRecordingLinkInput {
  contributor_id: IdentifierInput!
  recording_id:   IdentifierInput!
  role:           String!
  source:         SourceInput!
}

input ComposerWorkLinkInput {
  composer_id: IdentifierInput!
  work_id:     IdentifierInput!
  role:        String!
  pr_share:    String!
  mr_share:    String!
  sr_share:    String!
  source:      SourceInput!
}

input RecordLabelRecordingLinkInput {
  record_label_id: IdentifierInput!
  recording_id:    IdentifierInput!
  source:          SourceInput!
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
  role:         String!
  pr_share:     String!
  mr_share:     String!
  sr_share:     String!
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

input OrganisationSeriesLinkInput {
  organisation_id: IdentifierInput!
  series_id:       IdentifierInput!
  source:          SourceInput!
}

input OrganisationSeasonLinkInput {
  organisation_id: IdentifierInput!
  season_id:       IdentifierInput!
  source:          SourceInput!
}

input OrganisationEpisodeLinkInput {
  organisation_id: IdentifierInput!
  episode_id:      IdentifierInput!
  source:          SourceInput!
}

input OrganisationSupplementalLinkInput {
  organisation_id: IdentifierInput!
  supplemental_id: IdentifierInput!
  source:          SourceInput!
}

input SeriesSeasonLinkInput {
  series_id: IdentifierInput!
  season_id: IdentifierInput!
  source:    SourceInput!
}

input SeriesEpisodeLinkInput {
  series_id:  IdentifierInput!
  episode_id: IdentifierInput!
  source:     SourceInput!
}

input SeriesSupplementalLinkInput {
  series_id:       IdentifierInput!
  supplemental_id: IdentifierInput!
  source:          SourceInput!
}

input SeasonEpisodeLinkInput {
  season_id:  IdentifierInput!
  episode_id: IdentifierInput!
  source:     SourceInput!
}

input SeasonSupplementalLinkInput {
  season_id:       IdentifierInput!
  supplemental_id: IdentifierInput!
  source:          SourceInput!
}

input EpisodeSupplementalLinkInput {
  episode_id:      IdentifierInput!
  supplemental_id: IdentifierInput!
  source:          SourceInput!
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

type PerformerSongLink {
  performer: Performer!
  song:      Song!
  role:      String!
  source:    Source!
}

type PerformerReleaseLink {
  performer: Performer!
  release:   Release!
  role:      String!
  source:    Source!
}

type ContributorRecordingLink {
  contributor: Contributor!
  recording:   Recording!
  role:        String!
  source:      Source!
}

type ComposerWorkLink {
  composer: Composer!
  work:     Work!
  role:     String!
  pr_share: String!
  mr_share: String!
  sr_share: String!
  source:   Source!
}

type RecordLabelRecordingLink {
  record_label: RecordLabel!
  recording:    Recording!
  source:       Source!
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
  role:      String!
  pr_share:  String!
  mr_share:  String!
  sr_share:  String!
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

type OrganisationSeriesLink {
  organisation: Organisation!
  series:       Series!
  source:       Source!
}

type OrganisationSeasonLink {
  organisation: Organisation!
  season:       Season!
  source:       Source!
}

type OrganisationEpisodeLink {
  organisation: Organisation!
  episode:      Episode!
  source:       Source!
}

type OrganisationSupplementalLink {
  organisation: Organisation!
  supplemental: Supplemental!
  source:       Source!
}

type SeriesSeasonLink {
  series: Series!
  season: Season!
  source: Source!
}

type SeriesEpisodeLink {
  series:  Series!
  episode: Episode!
  source:  Source!
}

type SeriesSupplementalLink {
  series:       Series!
  supplemental: Supplemental!
  source:       Source!
}

type SeasonEpisodeLink {
  season:  Season!
  episode: Episode!
  source:  Source!
}

type SeasonSupplementalLink {
  season:       Season!
  supplemental: Supplemental!
  source:       Source!
}

type EpisodeSupplementalLink {
  episode:      Episode!
  supplemental: Supplemental!
  source:       Source!
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
		identifier, err := a.resolver.mediaIndex.Identifier("performer", id)
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
		identifier, err := a.resolver.mediaIndex.Identifier("record_label", id)
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
		identifier, err := a.resolver.mediaIndex.Identifier("composer", id)
		if err == nil {
			resolvers = append(resolvers, &ComposerResolver{a.resolver, identifier})
		} else if !isIdentifierNotFound(err) {
			return nil, err
		}
	}
	return resolvers, nil
}

func (a *AccountResolver) Publishers() ([]*PublisherResolver, error) {
	identifiers, err := a.identifiers("ipi")
	if err != nil {
		return nil, err
	}
	resolvers := make([]*PublisherResolver, 0, len(identifiers))
	for _, id := range identifiers {
		identifier, err := a.resolver.mediaIndex.Identifier("publisher", id)
		if err == nil {
			resolvers = append(resolvers, &PublisherResolver{a.resolver, identifier})
		} else if !isIdentifierNotFound(err) {
			return nil, err
		}
	}
	return resolvers, nil
}

func (a *AccountResolver) Organisations() ([]*OrganisationResolver, error) {
	identifiers, err := a.identifiers("doid")
	if err != nil {
		return nil, err
	}
	resolvers := make([]*OrganisationResolver, 0, len(identifiers))
	for _, id := range identifiers {
		identifier, err := a.resolver.mediaIndex.Identifier("organisation", id)
		if err == nil {
			resolvers = append(resolvers, &OrganisationResolver{a.resolver, identifier})
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
	identifier, err := r.mediaIndex.Identifier("performer", &args.Identifier)
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

func (p *PerformerResolver) Songs() ([]*performerSongLinkResolver, error) {
	records, err := p.resolver.mediaIndex.PerformerSongs(p.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*performerSongLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &performerSongLinkResolver{
			resolver: p.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (p *PerformerResolver) Releases() ([]*performerReleaseLinkResolver, error) {
	records, err := p.resolver.mediaIndex.PerformerReleases(p.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*performerReleaseLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &performerReleaseLinkResolver{
			resolver: p.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createContributorArgs struct {
	Contributor struct {
		Contributor

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateContributor(args createContributorArgs) (*ContributorResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Contributor.Contributor,
		&args.Contributor.Identifier,
		&args.Contributor.Source,
	)
	if err != nil {
		return nil, err
	}
	return &ContributorResolver{r, identifier}, nil
}

func (r *Resolver) Contributor(args IdentifierArgs) (*ContributorResolver, error) {
	identifier, err := r.mediaIndex.Identifier("contributor", &args.Identifier)
	if err != nil {
		return nil, err
	}
	return &ContributorResolver{r, identifier}, nil
}

type ContributorResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (p *ContributorResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{p.resolver, p.identifier}}
}

func (p *ContributorResolver) Name() (*stringValueResolver, error) {
	records, err := p.resolver.mediaIndex.Contributors(p.identifier)
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

func (p *ContributorResolver) Recordings() ([]*contributorRecordingLinkResolver, error) {
	records, err := p.resolver.mediaIndex.ContributorRecordings(p.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*contributorRecordingLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &contributorRecordingLinkResolver{
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
	identifier, err := r.mediaIndex.Identifier("composer", &args.Identifier)
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
	identifier, err := r.mediaIndex.Identifier("record_label", &args.Identifier)
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

func (l *RecordLabelResolver) Recordings() ([]*recordLabelRecordingLinkResolver, error) {
	records, err := l.resolver.mediaIndex.RecordLabelRecordings(l.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*recordLabelRecordingLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &recordLabelRecordingLinkResolver{
			resolver: l.resolver,
			record:   record,
		}
	}
	return resolvers, nil
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
	identifier, err := r.mediaIndex.Identifier("publisher", &args.Identifier)
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

func (p *PublisherResolver) Works() ([]*publisherWorkLinkResolver, error) {
	records, err := p.resolver.mediaIndex.PublisherWorks(p.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*publisherWorkLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &publisherWorkLinkResolver{
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
	identifier, err := r.mediaIndex.Identifier("recording", &args.Identifier)
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

func (r *RecordingResolver) Contributors() ([]*contributorRecordingLinkResolver, error) {
	records, err := r.resolver.mediaIndex.RecordingContributors(r.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*contributorRecordingLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &contributorRecordingLinkResolver{
			resolver: r.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (r *RecordingResolver) Songs() ([]*songRecordingLinkResolver, error) {
	records, err := r.resolver.mediaIndex.RecordingSongs(r.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*songRecordingLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &songRecordingLinkResolver{
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

func (r *RecordingResolver) RecordLabels() ([]*recordLabelRecordingLinkResolver, error) {
	records, err := r.resolver.mediaIndex.RecordingRecordLabels(r.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*recordLabelRecordingLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &recordLabelRecordingLinkResolver{
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
	identifier, err := r.mediaIndex.Identifier("work", &args.Identifier)
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

func (w *WorkResolver) Publishers() ([]*publisherWorkLinkResolver, error) {
	records, err := w.resolver.mediaIndex.WorkPublishers(w.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*publisherWorkLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &publisherWorkLinkResolver{
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
	identifier, err := r.mediaIndex.Identifier("song", &args.Identifier)
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

func (s *SongResolver) Performers() ([]*performerSongLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SongPerformers(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*performerSongLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &performerSongLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SongResolver) Recordings() ([]*songRecordingLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SongRecordings(s.identifier)
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
	identifier, err := r.mediaIndex.Identifier("release", &args.Identifier)
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

func (r *ReleaseResolver) Performers() ([]*performerReleaseLinkResolver, error) {
	records, err := r.resolver.mediaIndex.ReleasePerformers(r.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*performerReleaseLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &performerReleaseLinkResolver{
			resolver: r.resolver,
			record:   record,
		}
	}
	return resolvers, nil
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

type createOrganisationArgs struct {
	Organisation struct {
		Organisation

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateOrganisation(args createOrganisationArgs) (*OrganisationResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Organisation.Organisation,
		&args.Organisation.Identifier,
		&args.Organisation.Source,
	)
	if err != nil {
		return nil, err
	}
	return &OrganisationResolver{r, identifier}, nil
}

func (r *Resolver) Organisation(args IdentifierArgs) (*OrganisationResolver, error) {
	identifier, err := r.mediaIndex.Identifier("organisation", &args.Identifier)
	if err != nil {
		return nil, err
	}
	return &OrganisationResolver{r, identifier}, nil
}

type OrganisationResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (o *OrganisationResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{o.resolver, o.identifier}}
}

func (o *OrganisationResolver) Name() (*stringValueResolver, error) {
	records, err := o.resolver.mediaIndex.Organisation(o.identifier)
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

func (o *OrganisationResolver) Series() ([]*organisationSeriesLinkResolver, error) {
	records, err := o.resolver.mediaIndex.OrganisationSeries(o.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*organisationSeriesLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &organisationSeriesLinkResolver{
			resolver: o.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (o *OrganisationResolver) Seasons() ([]*organisationSeasonLinkResolver, error) {
	records, err := o.resolver.mediaIndex.OrganisationSeasons(o.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*organisationSeasonLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &organisationSeasonLinkResolver{
			resolver: o.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (o *OrganisationResolver) Episodes() ([]*organisationEpisodeLinkResolver, error) {
	records, err := o.resolver.mediaIndex.OrganisationEpisodes(o.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*organisationEpisodeLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &organisationEpisodeLinkResolver{
			resolver: o.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (o *OrganisationResolver) Supplementals() ([]*organisationSupplementalLinkResolver, error) {
	records, err := o.resolver.mediaIndex.OrganisationSupplementals(o.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*organisationSupplementalLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &organisationSupplementalLinkResolver{
			resolver: o.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createSeriesArgs struct {
	Series struct {
		Series

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateSeries(args createSeriesArgs) (*SeriesResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Series.Series,
		&args.Series.Identifier,
		&args.Series.Source,
	)
	if err != nil {
		return nil, err
	}
	return &SeriesResolver{r, identifier}, nil
}

func (r *Resolver) Series(args IdentifierArgs) (*SeriesResolver, error) {
	identifier, err := r.mediaIndex.Identifier("series", &args.Identifier)
	if err != nil {
		return nil, err
	}
	return &SeriesResolver{r, identifier}, nil
}

type SeriesResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (s *SeriesResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{s.resolver, s.identifier}}
}

func (s *SeriesResolver) Name() (*stringValueResolver, error) {
	records, err := s.resolver.mediaIndex.Series(s.identifier)
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

func (s *SeriesResolver) Organisations() ([]*organisationSeriesLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SeriesOrganisations(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*organisationSeriesLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &organisationSeriesLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SeriesResolver) Seasons() ([]*seriesSeasonLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SeriesSeasons(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*seriesSeasonLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &seriesSeasonLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SeriesResolver) Episodes() ([]*seriesEpisodeLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SeriesEpisodes(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*seriesEpisodeLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &seriesEpisodeLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SeriesResolver) Supplementals() ([]*seriesSupplementalLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SeriesSupplementals(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*seriesSupplementalLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &seriesSupplementalLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createSeasonArgs struct {
	Season struct {
		Season

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateSeason(args createSeasonArgs) (*SeasonResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Season.Season,
		&args.Season.Identifier,
		&args.Season.Source,
	)
	if err != nil {
		return nil, err
	}
	return &SeasonResolver{r, identifier}, nil
}

func (r *Resolver) Season(args IdentifierArgs) (*SeasonResolver, error) {
	identifier, err := r.mediaIndex.Identifier("season", &args.Identifier)
	if err != nil {
		return nil, err
	}
	return &SeasonResolver{r, identifier}, nil
}

type SeasonResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (s *SeasonResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{s.resolver, s.identifier}}
}

func (s *SeasonResolver) Name() (*stringValueResolver, error) {
	records, err := s.resolver.mediaIndex.Season(s.identifier)
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

func (s *SeasonResolver) Organisations() ([]*organisationSeasonLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SeasonOrganisations(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*organisationSeasonLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &organisationSeasonLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SeasonResolver) Series() ([]*seriesSeasonLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SeasonSeries(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*seriesSeasonLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &seriesSeasonLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SeasonResolver) Episodes() ([]*seasonEpisodeLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SeasonEpisodes(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*seasonEpisodeLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &seasonEpisodeLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SeasonResolver) Supplementals() ([]*seasonSupplementalLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SeasonSupplementals(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*seasonSupplementalLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &seasonSupplementalLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createEpisodeArgs struct {
	Episode struct {
		Episode

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateEpisode(args createEpisodeArgs) (*EpisodeResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Episode.Episode,
		&args.Episode.Identifier,
		&args.Episode.Source,
	)
	if err != nil {
		return nil, err
	}
	return &EpisodeResolver{r, identifier}, nil
}

func (r *Resolver) Episode(args IdentifierArgs) (*EpisodeResolver, error) {
	identifier, err := r.mediaIndex.Identifier("episode", &args.Identifier)
	if err != nil {
		return nil, err
	}
	return &EpisodeResolver{r, identifier}, nil
}

type EpisodeResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (e *EpisodeResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{e.resolver, e.identifier}}
}

func (e *EpisodeResolver) Name() (*stringValueResolver, error) {
	records, err := e.resolver.mediaIndex.Episode(e.identifier)
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

func (e *EpisodeResolver) Organisations() ([]*organisationEpisodeLinkResolver, error) {
	records, err := e.resolver.mediaIndex.EpisodeOrganisations(e.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*organisationEpisodeLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &organisationEpisodeLinkResolver{
			resolver: e.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (e *EpisodeResolver) Series() ([]*seriesEpisodeLinkResolver, error) {
	records, err := e.resolver.mediaIndex.EpisodeSeries(e.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*seriesEpisodeLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &seriesEpisodeLinkResolver{
			resolver: e.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (e *EpisodeResolver) Seasons() ([]*seasonEpisodeLinkResolver, error) {
	records, err := e.resolver.mediaIndex.EpisodeSeasons(e.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*seasonEpisodeLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &seasonEpisodeLinkResolver{
			resolver: e.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (e *EpisodeResolver) Supplementals() ([]*episodeSupplementalLinkResolver, error) {
	records, err := e.resolver.mediaIndex.EpisodeSupplementals(e.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*episodeSupplementalLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &episodeSupplementalLinkResolver{
			resolver: e.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

type createSupplementalArgs struct {
	Supplemental struct {
		Supplemental

		Identifier Identifier
		Source     Source
	}
}

func (r *Resolver) CreateSupplemental(args createSupplementalArgs) (*SupplementalResolver, error) {
	identifier, err := r.mediaIndex.CreateRecord(
		&args.Supplemental.Supplemental,
		&args.Supplemental.Identifier,
		&args.Supplemental.Source,
	)
	if err != nil {
		return nil, err
	}
	return &SupplementalResolver{r, identifier}, nil
}

func (r *Resolver) Supplemental(args IdentifierArgs) (*SupplementalResolver, error) {
	identifier, err := r.mediaIndex.Identifier("supplemental", &args.Identifier)
	if err != nil {
		return nil, err
	}
	return &SupplementalResolver{r, identifier}, nil
}

type SupplementalResolver struct {
	resolver   *Resolver
	identifier *IdentifierRecord
}

func (s *SupplementalResolver) Identifiers() []*identifierValueResolver {
	return []*identifierValueResolver{{s.resolver, s.identifier}}
}

func (s *SupplementalResolver) Name() (*stringValueResolver, error) {
	records, err := s.resolver.mediaIndex.Supplemental(s.identifier)
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

func (s *SupplementalResolver) Organisations() ([]*organisationSupplementalLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SupplementalOrganisations(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*organisationSupplementalLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &organisationSupplementalLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SupplementalResolver) Series() ([]*seriesSupplementalLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SupplementalSeries(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*seriesSupplementalLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &seriesSupplementalLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SupplementalResolver) Episodes() ([]*episodeSupplementalLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SupplementalEpisodes(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*episodeSupplementalLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &episodeSupplementalLinkResolver{
			resolver: s.resolver,
			record:   record,
		}
	}
	return resolvers, nil
}

func (s *SupplementalResolver) Seasons() ([]*seasonSupplementalLinkResolver, error) {
	records, err := s.resolver.mediaIndex.SupplementalSeasons(s.identifier)
	if err != nil {
		return nil, err
	}
	resolvers := make([]*seasonSupplementalLinkResolver, len(records))
	for i, record := range records {
		resolvers[i] = &seasonSupplementalLinkResolver{
			resolver: s.resolver,
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

type createPerformerSongLinkArgs struct {
	Link struct {
		PerformerID Identifier
		SongID      Identifier
		Role        string
		Source      Source
	}
}

func (r *Resolver) CreatePerformerSongLink(args createPerformerSongLinkArgs) (*performerSongLinkResolver, error) {
	link := &PerformerSongLink{
		Performer: args.Link.PerformerID,
		Song:      args.Link.SongID,
		Role:      args.Link.Role,
	}
	record, err := r.mediaIndex.CreatePerformerSong(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &performerSongLinkResolver{r, record}, nil
}

type performerSongLinkResolver struct {
	resolver *Resolver
	record   *PerformerSongRecord
}

func (p *performerSongLinkResolver) Performer() *PerformerResolver {
	return &PerformerResolver{p.resolver, p.record.Performer}
}

func (p *performerSongLinkResolver) Song() *SongResolver {
	return &SongResolver{p.resolver, p.record.Song}
}

func (p *performerSongLinkResolver) Role() string {
	return p.record.Role
}

func (p *performerSongLinkResolver) Source() *sourceResolver {
	return &sourceResolver{p.resolver, p.record.Source}
}

type createPerformerReleaseLinkArgs struct {
	Link struct {
		PerformerID Identifier
		ReleaseID   Identifier
		Role        string
		Source      Source
	}
}

func (r *Resolver) CreatePerformerReleaseLink(args createPerformerReleaseLinkArgs) (*performerReleaseLinkResolver, error) {
	link := &PerformerReleaseLink{
		Performer: args.Link.PerformerID,
		Release:   args.Link.ReleaseID,
		Role:      args.Link.Role,
	}
	record, err := r.mediaIndex.CreatePerformerRelease(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &performerReleaseLinkResolver{r, record}, nil
}

type performerReleaseLinkResolver struct {
	resolver *Resolver
	record   *PerformerReleaseRecord
}

func (p *performerReleaseLinkResolver) Performer() *PerformerResolver {
	return &PerformerResolver{p.resolver, p.record.Performer}
}

func (p *performerReleaseLinkResolver) Release() *ReleaseResolver {
	return &ReleaseResolver{p.resolver, p.record.Release}
}

func (p *performerReleaseLinkResolver) Role() string {
	return p.record.Role
}

func (p *performerReleaseLinkResolver) Source() *sourceResolver {
	return &sourceResolver{p.resolver, p.record.Source}
}

type createContributorRecordingLinkArgs struct {
	Link struct {
		ContributorID Identifier
		RecordingID   Identifier
		Role          string
		Source        Source
	}
}

func (r *Resolver) CreateContributorRecordingLink(args createContributorRecordingLinkArgs) (*contributorRecordingLinkResolver, error) {
	link := &ContributorRecordingLink{
		Contributor: args.Link.ContributorID,
		Recording:   args.Link.RecordingID,
		Role:        args.Link.Role,
	}
	record, err := r.mediaIndex.CreateContributorRecording(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &contributorRecordingLinkResolver{r, record}, nil
}

type contributorRecordingLinkResolver struct {
	resolver *Resolver
	record   *ContributorRecordingRecord
}

func (c *contributorRecordingLinkResolver) Contributor() *ContributorResolver {
	return &ContributorResolver{c.resolver, c.record.Contributor}
}

func (c *contributorRecordingLinkResolver) Recording() *RecordingResolver {
	return &RecordingResolver{c.resolver, c.record.Recording}
}

func (c *contributorRecordingLinkResolver) Role() string {
	return c.record.Role
}

func (c *contributorRecordingLinkResolver) Source() *sourceResolver {
	return &sourceResolver{c.resolver, c.record.Source}
}

type createComposerWorkLinkArgs struct {
	Link struct {
		ComposerID Identifier
		WorkID     Identifier
		Role       string
		PRShare    string
		MRShare    string
		SRShare    string
		Source     Source
	}
}

func (r *Resolver) CreateComposerWorkLink(args createComposerWorkLinkArgs) (*composerWorkLinkResolver, error) {
	link := &ComposerWorkLink{
		Composer: args.Link.ComposerID,
		Work:     args.Link.WorkID,
		Role:     args.Link.Role,
		PerformanceRightsShare:     args.Link.PRShare,
		MechanicalRightsShare:      args.Link.MRShare,
		SynchronizationRightsShare: args.Link.SRShare,
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

func (c *composerWorkLinkResolver) PRShare() string {
	return c.record.PerformanceRightsShare
}

func (c *composerWorkLinkResolver) MRShare() string {
	return c.record.MechanicalRightsShare
}

func (c *composerWorkLinkResolver) SRShare() string {
	return c.record.SynchronizationRightsShare
}

func (c *composerWorkLinkResolver) Source() *sourceResolver {
	return &sourceResolver{c.resolver, c.record.Source}
}

type createRecordLabelRecordingLinkArgs struct {
	Link struct {
		RecordLabelID Identifier
		RecordingID   Identifier
		Source        Source
	}
}

func (r *Resolver) CreateRecordLabelRecordingLink(args createRecordLabelRecordingLinkArgs) (*recordLabelRecordingLinkResolver, error) {
	link := &RecordLabelRecordingLink{
		RecordLabel: args.Link.RecordLabelID,
		Recording:   args.Link.RecordingID,
	}
	record, err := r.mediaIndex.CreateRecordLabelRecording(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &recordLabelRecordingLinkResolver{r, record}, nil
}

type recordLabelRecordingLinkResolver struct {
	resolver *Resolver
	record   *RecordLabelRecordingRecord
}

func (r *recordLabelRecordingLinkResolver) RecordLabel() *RecordLabelResolver {
	return &RecordLabelResolver{r.resolver, r.record.RecordLabel}
}

func (r *recordLabelRecordingLinkResolver) Recording() *RecordingResolver {
	return &RecordingResolver{r.resolver, r.record.Recording}
}

func (r *recordLabelRecordingLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
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
		Role        string
		PRShare     string
		MRShare     string
		SRShare     string
		Source      Source
	}
}

func (r *Resolver) CreatePublisherWorkLink(args createPublisherWorkLinkArgs) (*publisherWorkLinkResolver, error) {
	link := &PublisherWorkLink{
		Publisher: args.Link.PublisherID,
		Work:      args.Link.WorkID,
		Role:      args.Link.Role,
		PerformanceRightsShare:     args.Link.PRShare,
		MechanicalRightsShare:      args.Link.MRShare,
		SynchronizationRightsShare: args.Link.SRShare,
	}
	record, err := r.mediaIndex.CreatePublisherWork(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &publisherWorkLinkResolver{r, record}, nil
}

type publisherWorkLinkResolver struct {
	resolver *Resolver
	record   *PublisherWorkRecord
}

func (p *publisherWorkLinkResolver) Publisher() *PublisherResolver {
	return &PublisherResolver{p.resolver, p.record.Publisher}
}

func (p *publisherWorkLinkResolver) Work() *WorkResolver {
	return &WorkResolver{p.resolver, p.record.Work}
}

func (p *publisherWorkLinkResolver) Source() *sourceResolver {
	return &sourceResolver{p.resolver, p.record.Source}
}

func (p *publisherWorkLinkResolver) Role() string {
	return p.record.Role
}

func (p *publisherWorkLinkResolver) PRShare() string {
	return p.record.PerformanceRightsShare
}

func (p *publisherWorkLinkResolver) MRShare() string {
	return p.record.MechanicalRightsShare
}

func (p *publisherWorkLinkResolver) SRShare() string {
	return p.record.SynchronizationRightsShare
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

type createOrganisationSeriesLinkArgs struct {
	Link struct {
		OrganisationID Identifier
		SeriesID       Identifier
		Source         Source
	}
}

func (r *Resolver) CreateOrganisationSeriesLink(args createOrganisationSeriesLinkArgs) (*organisationSeriesLinkResolver, error) {
	link := &OrganisationSeriesLink{
		Organisation: args.Link.OrganisationID,
		Series:       args.Link.SeriesID,
	}
	record, err := r.mediaIndex.CreateOrganisationSeries(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &organisationSeriesLinkResolver{r, record}, nil
}

type organisationSeriesLinkResolver struct {
	resolver *Resolver
	record   *OrganisationSeriesRecord
}

func (r *organisationSeriesLinkResolver) Organisation() *OrganisationResolver {
	return &OrganisationResolver{r.resolver, r.record.Organisation}
}

func (r *organisationSeriesLinkResolver) Series() *SeriesResolver {
	return &SeriesResolver{r.resolver, r.record.Series}
}

func (r *organisationSeriesLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createOrganisationSeasonLinkArgs struct {
	Link struct {
		OrganisationID Identifier
		SeasonID       Identifier
		Source         Source
	}
}

func (r *Resolver) CreateOrganisationSeasonLink(args createOrganisationSeasonLinkArgs) (*organisationSeasonLinkResolver, error) {
	link := &OrganisationSeasonLink{
		Organisation: args.Link.OrganisationID,
		Season:       args.Link.SeasonID,
	}
	record, err := r.mediaIndex.CreateOrganisationSeason(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &organisationSeasonLinkResolver{r, record}, nil
}

type organisationSeasonLinkResolver struct {
	resolver *Resolver
	record   *OrganisationSeasonRecord
}

func (r *organisationSeasonLinkResolver) Organisation() *OrganisationResolver {
	return &OrganisationResolver{r.resolver, r.record.Organisation}
}

func (r *organisationSeasonLinkResolver) Season() *SeasonResolver {
	return &SeasonResolver{r.resolver, r.record.Season}
}

func (r *organisationSeasonLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createOrganisationEpisodeLinkArgs struct {
	Link struct {
		OrganisationID Identifier
		EpisodeID      Identifier
		Source         Source
	}
}

func (r *Resolver) CreateOrganisationEpisodeLink(args createOrganisationEpisodeLinkArgs) (*organisationEpisodeLinkResolver, error) {
	link := &OrganisationEpisodeLink{
		Organisation: args.Link.OrganisationID,
		Episode:      args.Link.EpisodeID,
	}
	record, err := r.mediaIndex.CreateOrganisationEpisode(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &organisationEpisodeLinkResolver{r, record}, nil
}

type organisationEpisodeLinkResolver struct {
	resolver *Resolver
	record   *OrganisationEpisodeRecord
}

func (r *organisationEpisodeLinkResolver) Organisation() *OrganisationResolver {
	return &OrganisationResolver{r.resolver, r.record.Organisation}
}

func (r *organisationEpisodeLinkResolver) Episode() *EpisodeResolver {
	return &EpisodeResolver{r.resolver, r.record.Episode}
}

func (r *organisationEpisodeLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createOrganisationSupplementalLinkArgs struct {
	Link struct {
		OrganisationID Identifier
		SupplementalID Identifier
		Source         Source
	}
}

func (r *Resolver) CreateOrganisationSupplementalLink(args createOrganisationSupplementalLinkArgs) (*organisationSupplementalLinkResolver, error) {
	link := &OrganisationSupplementalLink{
		Organisation: args.Link.OrganisationID,
		Supplemental: args.Link.SupplementalID,
	}
	record, err := r.mediaIndex.CreateOrganisationSupplemental(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &organisationSupplementalLinkResolver{r, record}, nil
}

type organisationSupplementalLinkResolver struct {
	resolver *Resolver
	record   *OrganisationSupplementalRecord
}

func (r *organisationSupplementalLinkResolver) Organisation() *OrganisationResolver {
	return &OrganisationResolver{r.resolver, r.record.Organisation}
}

func (r *organisationSupplementalLinkResolver) Supplemental() *SupplementalResolver {
	return &SupplementalResolver{r.resolver, r.record.Supplemental}
}

func (r *organisationSupplementalLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createSeriesSeasonLinkArgs struct {
	Link struct {
		SeriesID Identifier
		SeasonID Identifier
		Source   Source
	}
}

func (r *Resolver) CreateSeriesSeasonLink(args createSeriesSeasonLinkArgs) (*seriesSeasonLinkResolver, error) {
	link := &SeriesSeasonLink{
		Series: args.Link.SeriesID,
		Season: args.Link.SeasonID,
	}
	record, err := r.mediaIndex.CreateSeriesSeason(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &seriesSeasonLinkResolver{r, record}, nil
}

type seriesSeasonLinkResolver struct {
	resolver *Resolver
	record   *SeriesSeasonRecord
}

func (r *seriesSeasonLinkResolver) Series() *SeriesResolver {
	return &SeriesResolver{r.resolver, r.record.Series}
}

func (r *seriesSeasonLinkResolver) Season() *SeasonResolver {
	return &SeasonResolver{r.resolver, r.record.Season}
}

func (r *seriesSeasonLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createSeriesEpisodeLinkArgs struct {
	Link struct {
		SeriesID  Identifier
		EpisodeID Identifier
		Source    Source
	}
}

func (r *Resolver) CreateSeriesEpisodeLink(args createSeriesEpisodeLinkArgs) (*seriesEpisodeLinkResolver, error) {
	link := &SeriesEpisodeLink{
		Series:  args.Link.SeriesID,
		Episode: args.Link.EpisodeID,
	}
	record, err := r.mediaIndex.CreateSeriesEpisode(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &seriesEpisodeLinkResolver{r, record}, nil
}

type seriesEpisodeLinkResolver struct {
	resolver *Resolver
	record   *SeriesEpisodeRecord
}

func (r *seriesEpisodeLinkResolver) Series() *SeriesResolver {
	return &SeriesResolver{r.resolver, r.record.Series}
}

func (r *seriesEpisodeLinkResolver) Episode() *EpisodeResolver {
	return &EpisodeResolver{r.resolver, r.record.Episode}
}

func (r *seriesEpisodeLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createSeriesSupplementalLinkArgs struct {
	Link struct {
		SeriesID       Identifier
		SupplementalID Identifier
		Source         Source
	}
}

func (r *Resolver) CreateSeriesSupplementalLink(args createSeriesSupplementalLinkArgs) (*seriesSupplementalLinkResolver, error) {
	link := &SeriesSupplementalLink{
		Series:       args.Link.SeriesID,
		Supplemental: args.Link.SupplementalID,
	}
	record, err := r.mediaIndex.CreateSeriesSupplemental(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &seriesSupplementalLinkResolver{r, record}, nil
}

type seriesSupplementalLinkResolver struct {
	resolver *Resolver
	record   *SeriesSupplementalRecord
}

func (r *seriesSupplementalLinkResolver) Series() *SeriesResolver {
	return &SeriesResolver{r.resolver, r.record.Series}
}

func (r *seriesSupplementalLinkResolver) Supplemental() *SupplementalResolver {
	return &SupplementalResolver{r.resolver, r.record.Supplemental}
}

func (r *seriesSupplementalLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createSeasonEpisodeLinkArgs struct {
	Link struct {
		SeasonID  Identifier
		EpisodeID Identifier
		Source    Source
	}
}

func (r *Resolver) CreateSeasonEpisodeLink(args createSeasonEpisodeLinkArgs) (*seasonEpisodeLinkResolver, error) {
	link := &SeasonEpisodeLink{
		Season:  args.Link.SeasonID,
		Episode: args.Link.EpisodeID,
	}
	record, err := r.mediaIndex.CreateSeasonEpisode(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &seasonEpisodeLinkResolver{r, record}, nil
}

type seasonEpisodeLinkResolver struct {
	resolver *Resolver
	record   *SeasonEpisodeRecord
}

func (r *seasonEpisodeLinkResolver) Season() *SeasonResolver {
	return &SeasonResolver{r.resolver, r.record.Season}
}

func (r *seasonEpisodeLinkResolver) Episode() *EpisodeResolver {
	return &EpisodeResolver{r.resolver, r.record.Episode}
}

func (r *seasonEpisodeLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createSeasonSupplementalLinkArgs struct {
	Link struct {
		SeasonID       Identifier
		SupplementalID Identifier
		Source         Source
	}
}

func (r *Resolver) CreateSeasonSupplementalLink(args createSeasonSupplementalLinkArgs) (*seasonSupplementalLinkResolver, error) {
	link := &SeasonSupplementalLink{
		Season:       args.Link.SeasonID,
		Supplemental: args.Link.SupplementalID,
	}
	record, err := r.mediaIndex.CreateSeasonSupplemental(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &seasonSupplementalLinkResolver{r, record}, nil
}

type seasonSupplementalLinkResolver struct {
	resolver *Resolver
	record   *SeasonSupplementalRecord
}

func (r *seasonSupplementalLinkResolver) Season() *SeasonResolver {
	return &SeasonResolver{r.resolver, r.record.Season}
}

func (r *seasonSupplementalLinkResolver) Supplemental() *SupplementalResolver {
	return &SupplementalResolver{r.resolver, r.record.Supplemental}
}

func (r *seasonSupplementalLinkResolver) Source() *sourceResolver {
	return &sourceResolver{r.resolver, r.record.Source}
}

type createEpisodeSupplementalLinkArgs struct {
	Link struct {
		EpisodeID      Identifier
		SupplementalID Identifier
		Source         Source
	}
}

func (r *Resolver) CreateEpisodeSupplementalLink(args createEpisodeSupplementalLinkArgs) (*episodeSupplementalLinkResolver, error) {
	link := &EpisodeSupplementalLink{
		Episode:      args.Link.EpisodeID,
		Supplemental: args.Link.SupplementalID,
	}
	record, err := r.mediaIndex.CreateEpisodeSupplemental(link, &args.Link.Source)
	if err != nil {
		return nil, err
	}
	return &episodeSupplementalLinkResolver{r, record}, nil
}

type episodeSupplementalLinkResolver struct {
	resolver *Resolver
	record   *EpisodeSupplementalRecord
}

func (r *episodeSupplementalLinkResolver) Episode() *EpisodeResolver {
	return &EpisodeResolver{r.resolver, r.record.Episode}
}

func (r *episodeSupplementalLinkResolver) Supplemental() *SupplementalResolver {
	return &SupplementalResolver{r.resolver, r.record.Supplemental}
}

func (r *episodeSupplementalLinkResolver) Source() *sourceResolver {
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
