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
	"database/sql"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/xml"
)

// Resolver is a general purpose GraphQL resolver function,
// retrieves data from a META store and SQLite3 index
type Resolver struct {
	db    *sql.DB
	store *meta.Store
}

// NewResolver returns a new resolver for returning data from the given
// META store and SQLite3 index
func NewResolver(db *sql.DB, store *meta.Store) *Resolver {
	return &Resolver{db, store}
}

// GraphQLSchema schema definition for Party on the ERN index
const GraphQLSchema = `
schema {
  query: Query
}

type Query {
  party(id: String, name: String): [Party]!

  soundRecording(id: String, title: String): [SoundRecording]!

  release(id: String, title: String): [Release]!
}

type Party {
  cid:      String!
  type:     String!
  role:     String!
  partyID:  String!
  fullName: String!

  soundRecordings: [SoundRecording]!
  releases:        [Release]!
}

type Artist {
  cid:   String!
  party: Party!
  role:  String!
}

type ResourceContributor {
  cid:   String!
  party: Party!
  role:  String!
}

type SoundRecording {
  cid:                String!
  source:             String!
  soundRecordingId:   String!
  soundRecordingType: String!
  referenceTitle:     String!
  duration:           String!
  detailsByTerritory: [SoundRecordingDetailsByTerritory]!
}

type SoundRecordingDetailsByTerritory {
  cid:                 String!
  title:               [Title]!
  displayArtist:       [Artist]!
  displayArtistName:   String!
  labelName:           String!
  territoryCode:       String!
  genre:               String!
  parentalWarningType: String!
}

type Release {
  cid:                String!
  source:             String!
  releaseId:          ReleaseID!
  releaseType:        String!
  referenceTitle:     String!
  detailsByTerritory: [ReleaseDetailsByTerritory]!
}

type ReleaseDetailsByTerritory {
  cid:               String!
  title:             [Title]!
  displayArtist:     [Artist]!
  displayArtistName: String!
  labelName:         String!
  territoryCode:     String!
  releaseDate:       String!
  genre:             String!
}

type Title {
  titleType: String!
  titleText: String!
}

type ReleaseID {
  grid: String!
  icpn: String!
}
`

// partyArgs query arguments for Party query
type PartyArgs struct {
	Name *string
	ID   *string
}

/**
 *	Party
 */

// PartyResolver defines grapQL resolver functions for the Party fields
type PartyResolver struct {
	resolver  *Resolver
	cid       string
	source    string
	typ       string
	role      string
	partyId   string
	partyName string
}

func (pd *PartyResolver) Cid() string {
	return pd.cid
}

func (pd *PartyResolver) Source() string {
	return pd.source
}

func (pd *PartyResolver) Type() string {
	return pd.typ
}

func (pd *PartyResolver) Role() string {
	return pd.role
}

func (pd *PartyResolver) Fullname() string {
	return pd.partyName
}

func (pd *PartyResolver) PartyId() string {
	return pd.partyId
}

func (pd *PartyResolver) SoundRecordings() ([]*SoundRecordingResolver, error) {
	rows, err := pd.resolver.db.Query("SELECT sound_recording_id FROM sound_recording_contributor WHERE party_id = ?", pd.cid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var recordings []*SoundRecordingResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		recording, err := pd.resolver.soundRecordingResolver(objectID)
		if err != nil {
			return nil, err
		}
		recordings = append(recordings, recording)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return recordings, nil
}

func (pd *PartyResolver) Releases() ([]*ReleaseResolver, error) {
	rows, err := pd.resolver.db.Query(
		"SELECT release_id FROM release_list INNER JOIN ern ON ern.cid = release_list.ern_id WHERE ern.sender_id = ?",
		pd.cid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var releases []*ReleaseResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		release, err := pd.resolver.releaseResolver(objectID)
		if err != nil {
			return nil, err
		}
		releases = append(releases, release)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return releases, nil
}

// Party is the resolver function to retrieve
// the Party information from the SQLite index
func (g *Resolver) Party(args PartyArgs) ([]*PartyResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.Name != nil:
		rows, err = g.db.Query("SELECT cid FROM party WHERE name = ?", *args.Name)
	case args.ID != nil:
		rows, err = g.db.Query("SELECT cid FROM party WHERE id = ?", *args.ID)
	default:
		return nil, errors.New("Missing Name or ID argument in query")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resolvers []*PartyResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}

		resolver, err := g.partyResolver(objectID)
		if err != nil {
			return nil, err
		}

		resolvers = append(resolvers, resolver)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

func (r *Resolver) partyResolver(objectID string) (*PartyResolver, error) {
	id, err := cid.Parse(objectID)
	if err != nil {
		return nil, err
	}

	obj, err := r.store.Get(id)
	if err != nil {
		return nil, err
	}

	source, err := obj.GetString("@source")
	if err != nil {
		return nil, err
	}

	typ, err := obj.GetString("@type")
	if err != nil {
		return nil, err
	}

	var artistRole metaxml.Value
	if v, err := DecodeObj(r.store, obj, "ArtistRole"); err == nil {
		artistRole = *v
	}

	var partyId metaxml.Value
	if v, err := DecodeObj(r.store, obj, "PartyId"); err == nil {
		partyId = *v
	}

	partyName, err := DecodeObj(r.store, obj, "PartyName", "FullName")
	if err != nil {
		return nil, err
	}

	return &PartyResolver{
		resolver:  r,
		cid:       objectID,
		source:    source,
		typ:       typ,
		role:      artistRole.Value,
		partyId:   partyId.Value,
		partyName: partyName.Value,
	}, nil
}

/**
 * SoundRecording
 */

type SoundRecordingArgs struct {
	ID    *string
	Title *string
}

type SoundRecordingResolver struct {
	resolver           *Resolver
	cid                string
	source             string
	soundRecordingId   string
	soundRecordingType string
	referenceTitle     string
	duration           string
	detailsByTerritory []*SoundRecordingDetailsResolver
}

func (sr *SoundRecordingResolver) Cid() string {
	return sr.cid
}

func (sr *SoundRecordingResolver) Source() string {
	return sr.source
}

func (sr *SoundRecordingResolver) SoundRecordingId() string {
	return sr.soundRecordingId
}

func (sr *SoundRecordingResolver) SoundRecordingType() string {
	return sr.soundRecordingType
}

func (sr *SoundRecordingResolver) ReferenceTitle() string {
	return sr.referenceTitle
}

func (sr *SoundRecordingResolver) Duration() string {
	return sr.duration
}

func (sr *SoundRecordingResolver) DetailsByTerritory() []*SoundRecordingDetailsResolver {
	return sr.detailsByTerritory
}

func (sr *SoundRecordingResolver) Contributors() ([]*PartyResolver, error) {
	rows, err := sr.resolver.db.Query(
		"SELECT party_id FROM sound_recording_contributor WHERE sound_recording_id = ?",
		sr.cid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var parties []*PartyResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		party, err := sr.resolver.partyResolver(objectID)
		if err != nil {
			return nil, err
		}
		parties = append(parties, party)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return parties, nil
}

func (sr *SoundRecordingResolver) Releases() ([]*ReleaseResolver, error) {
	rows, err := sr.resolver.db.Query(
		"SELECT release_id FROM sound_recording_release WHERE sound_recording_id = ?",
		sr.cid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var releases []*ReleaseResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		release, err := sr.resolver.releaseResolver(objectID)
		if err != nil {
			return nil, err
		}
		releases = append(releases, release)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return releases, nil
}

type SoundRecordingDetailsResolver struct {
	cid                 string
	title               []*TitleResolver
	displayArtist       []*ArtistResolver
	displayArtistName   string
	labelName           string
	territoryCode       string
	genre               string
	parentalWarningType string
}

func (s *SoundRecordingDetailsResolver) Cid() string {
	return s.cid
}

func (s *SoundRecordingDetailsResolver) Title() []*TitleResolver {
	return s.title
}

func (s *SoundRecordingDetailsResolver) DisplayArtist() []*ArtistResolver {
	return s.displayArtist
}

func (s *SoundRecordingDetailsResolver) DisplayArtistName() string {
	return s.displayArtistName
}

func (s *SoundRecordingDetailsResolver) LabelName() string {
	return s.labelName
}

func (s *SoundRecordingDetailsResolver) TerritoryCode() string {
	return s.territoryCode
}

func (s *SoundRecordingDetailsResolver) Genre() string {
	return s.genre
}

func (s *SoundRecordingDetailsResolver) ParentalWarningType() string {
	return s.parentalWarningType
}

type ArtistResolver struct {
	party *PartyResolver
	role  string
}

func (a *ArtistResolver) Cid() string {
	return a.party.Cid()
}

func (a *ArtistResolver) Party() *PartyResolver {
	return a.party
}

func (a *ArtistResolver) Role() string {
	return a.role
}

type TitleResolver struct {
	titleText string
	titleType string
}

func (t *TitleResolver) TitleText() string {
	return t.titleText
}

func (t *TitleResolver) TitleType() string {
	return t.titleType
}

type ContibutorResolver struct {
	party *PartyResolver
	role  string
}

func (c *ContibutorResolver) Party() *PartyResolver {
	return c.party
}

func (c *ContibutorResolver) Role() string {
	return c.role
}

func (g *Resolver) SoundRecording(args SoundRecordingArgs) ([]*SoundRecordingResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.ID != nil:
		rows, err = g.db.Query("SELECT cid FROM sound_recording WHERE id = ?", *args.ID)
	case args.Title != nil:
		rows, err = g.db.Query("SELECT cid FROM sound_recording WHERE title = ?", *args.Title)
	default:
		return nil, errors.New("Missing ID or Title argument in query")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recordings []*SoundRecordingResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}

		recording, err := g.soundRecordingResolver(objectID)
		if err != nil {
			return nil, err
		}

		recordings = append(recordings, recording)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recordings, nil
}

func (r *Resolver) soundRecordingResolver(objectID string) (*SoundRecordingResolver, error) {
	id, err := cid.Parse(objectID)
	if err != nil {
		return nil, err
	}

	obj, err := r.store.Get(id)
	if err != nil {
		return nil, err
	}

	source, err := obj.GetString("@source")
	if err != nil {
		return nil, err
	}

	soundRecordingId, err := DecodeObj(r.store, obj, "SoundRecordingId", "ISRC")
	if err != nil {
		return nil, err
	}

	soundRecordingType, err := DecodeObj(r.store, obj, "SoundRecordingType")
	if err != nil {
		return nil, err
	}

	referenceTitle, err := DecodeObj(r.store, obj, "ReferenceTitle", "TitleText")
	if err != nil {
		return nil, err
	}

	duration, err := DecodeObj(r.store, obj, "Duration")
	if err != nil {
		return nil, err
	}

	recording := &SoundRecordingResolver{
		resolver:           r,
		cid:                objectID,
		source:             source,
		soundRecordingId:   soundRecordingId.Value,
		soundRecordingType: soundRecordingType.Value,
		referenceTitle:     referenceTitle.Value,
		duration:           duration.Value,
	}

	detailIDs, err := decodeLinks(r.store, obj, "SoundRecordingDetailsByTerritory")
	if err != nil {
		return nil, err
	}
	for _, detailID := range detailIDs {
		obj, err := r.store.Get(detailID)
		if err != nil {
			return nil, err
		}
		titleIDs, err := decodeLinks(r.store, obj, "Title")
		if err != nil {
			return nil, err
		}
		var titles []*TitleResolver
		for _, titleID := range titleIDs {
			obj, err := r.store.Get(titleID)
			if err != nil {
				return nil, err
			}
			titleText, err := DecodeObj(r.store, obj, "TitleText")
			if err != nil {
				return nil, err
			}
			var titleType metaxml.Value
			if v, err := DecodeObj(r.store, obj, "TitleType"); err == nil {
				titleType = *v
			}
			titles = append(titles, &TitleResolver{
				titleType: titleType.Value,
				titleText: titleText.Value,
			})
		}
		var artistName metaxml.Value
		if v, err := DecodeObj(r.store, obj, "DisplayArtistName"); err == nil {
			artistName = *v
		}
		var labelName metaxml.Value
		if v, err := DecodeObj(r.store, obj, "LabelName"); err == nil {
			labelName = *v
		}
		var territoryCode metaxml.Value
		if v, err := DecodeObj(r.store, obj, "TerritoryCode"); err == nil {
			territoryCode = *v
		}
		var genre metaxml.Value
		if v, err := DecodeObj(r.store, obj, "Genre", "GenreText"); err == nil {
			genre = *v
		}
		var parentalWarningType metaxml.Value
		if v, err := DecodeObj(r.store, obj, "ParentalWarningType"); err == nil {
			parentalWarningType = *v
		}
		artistIDs, err := decodeLinks(r.store, obj, "DisplayArtist")
		if err != nil {
			return nil, err
		}
		var artists []*ArtistResolver
		for _, artistID := range artistIDs {
			party, err := r.partyResolver(artistID.String())
			if err != nil {
				return nil, err
			}
			var role metaxml.Value
			if v, err := DecodeObj(r.store, obj, "ArtistRole"); err == nil {
				role = *v
			}
			artists = append(artists, &ArtistResolver{
				party: party,
				role:  role.Value,
			})
		}
		recording.detailsByTerritory = append(recording.detailsByTerritory, &SoundRecordingDetailsResolver{
			cid:                 obj.Cid().String(),
			title:               titles,
			displayArtist:       artists,
			displayArtistName:   artistName.Value,
			labelName:           labelName.Value,
			territoryCode:       territoryCode.Value,
			genre:               genre.Value,
			parentalWarningType: parentalWarningType.Value,
		})

	}
	return recording, nil
}

/*
* Release
 */

type ReleaseArgs struct {
	ID              *string
	Title           *string
	WithMainRelease *bool
}

type ReleaseResolver struct {
	resolver           *Resolver
	cid                string
	source             string
	releaseID          *ReleaseIDResolver
	releaseType        string
	referenceTitle     string
	detailsByTerritory []*ReleaseDetailsResolver
}

func (r *ReleaseResolver) Cid() string {
	return r.cid
}

func (r *ReleaseResolver) Source() string {
	return r.source
}

func (r *ReleaseResolver) ReleaseID() *ReleaseIDResolver {
	return r.releaseID
}

func (r *ReleaseResolver) ReleaseType() string {
	return r.releaseType
}

func (r *ReleaseResolver) ReferenceTitle() string {
	return r.referenceTitle
}

func (r *ReleaseResolver) DetailsByTerritory() []*ReleaseDetailsResolver {
	return r.detailsByTerritory
}

func (r *ReleaseResolver) MessageSenders() ([]*PartyResolver, error) {
	rows, err := r.resolver.db.Query(
		"SELECT sender_id FROM ern INNER JOIN release_list ON release_list.ern_id = ern.cid WHERE release_list.release_id = ?",
		r.cid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var parties []*PartyResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		party, err := r.resolver.partyResolver(objectID)
		if err != nil {
			return nil, err
		}
		parties = append(parties, party)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return parties, nil
}

func (r *ReleaseResolver) Contributors() ([]*PartyResolver, error) {
	rows, err := r.resolver.db.Query(
		"SELECT party_id FROM sound_recording_contributor INNER JOIN sound_recording_release ON sound_recording_release.sound_recording_id = sound_recording_contributor.sound_recording_id WHERE sound_recording_release.release_id = ?",
		r.cid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var parties []*PartyResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		party, err := r.resolver.partyResolver(objectID)
		if err != nil {
			return nil, err
		}
		parties = append(parties, party)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return parties, nil
}

func (r *ReleaseResolver) SoundRecordings() ([]*SoundRecordingResolver, error) {
	rows, err := r.resolver.db.Query(
		"SELECT sound_recording_id FROM sound_recording_release WHERE release_id = ?",
		r.cid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var recordings []*SoundRecordingResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		recording, err := r.resolver.soundRecordingResolver(objectID)
		if err != nil {
			return nil, err
		}
		recordings = append(recordings, recording)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return recordings, nil
}

type ReleaseIDResolver struct {
	grid string
	icpn string
}

func (r *ReleaseIDResolver) GRID() string {
	return r.grid
}

func (r *ReleaseIDResolver) ICPN() string {
	return r.icpn
}

type ReleaseDetailsResolver struct {
	cid               string
	title             []*TitleResolver
	displayArtist     []*ArtistResolver
	displayArtistName string
	labelName         string
	territoryCode     string
	releaseDate       string
	genre             string
}

func (r *ReleaseDetailsResolver) Cid() string {
	return r.cid
}

func (r *ReleaseDetailsResolver) Title() []*TitleResolver {
	return r.title
}

func (r *ReleaseDetailsResolver) DisplayArtist() []*ArtistResolver {
	return r.displayArtist
}

func (r *ReleaseDetailsResolver) DisplayArtistName() string {
	return r.displayArtistName
}

func (r *ReleaseDetailsResolver) LabelName() string {
	return r.labelName
}

func (r *ReleaseDetailsResolver) TerritoryCode() string {
	return r.territoryCode
}

func (r *ReleaseDetailsResolver) ReleaseDate() string {
	return r.releaseDate
}

func (r *ReleaseDetailsResolver) Genre() string {
	return r.genre
}

func (g *Resolver) Release(args ReleaseArgs) ([]*ReleaseResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.ID != nil:
		rows, err = g.db.Query("SELECT cid FROM release WHERE id = ?", *args.ID)
	case args.Title != nil:
		rows, err = g.db.Query("SELECT cid FROM release WHERE title = ?", *args.Title)
	default:
		return nil, errors.New("Missing ID or Title argument in query")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var releases []*ReleaseResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		if args.WithMainRelease != nil {
			isMainRelease, err := g.isMainRelease(objectID)
			if err != nil {
				return nil, err
			}
			if !*args.WithMainRelease && isMainRelease {
				continue
			}
			if *args.WithMainRelease && !isMainRelease {
				continue
			}
		}

		release, err := g.releaseResolver(objectID)
		if err != nil {
			return nil, err
		}

		releases = append(releases, release)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return releases, nil
}

func (g *Resolver) isMainRelease(objectID string) (bool, error) {
	objCid, err := cid.Parse(objectID)
	if err != nil {
		return false, err
	}
	obj, err := g.store.Get(objCid)
	if err != nil {
		return false, err
	}

	isMainRelease, err := obj.GetString("IsMainRelease")
	if err != nil {
		return false, nil
	}
	return (isMainRelease == "true"), err
}

func (r *Resolver) releaseResolver(objectID string) (*ReleaseResolver, error) {
	id, err := cid.Parse(objectID)
	if err != nil {
		return nil, err
	}

	obj, err := r.store.Get(id)
	if err != nil {
		return nil, err
	}

	source, err := obj.GetString("@source")
	if err != nil {
		return nil, err
	}

	var grid metaxml.Value
	if v, err := DecodeObj(r.store, obj, "ReleaseId", "GRid"); err == nil {
		grid = *v
	}

	var icpn metaxml.Value
	if v, err := DecodeObj(r.store, obj, "ReleaseId", "ICPN"); err == nil {
		icpn = *v
	}

	releaseType, err := DecodeObj(r.store, obj, "ReleaseType")
	if err != nil {
		return nil, err
	}

	referenceTitle, err := DecodeObj(r.store, obj, "ReferenceTitle", "TitleText")
	if err != nil {
		return nil, err
	}

	release := &ReleaseResolver{
		resolver: r,
		cid:      objectID,
		source:   source,
		releaseID: &ReleaseIDResolver{
			grid: grid.Value,
			icpn: icpn.Value,
		},
		releaseType:    releaseType.Value,
		referenceTitle: referenceTitle.Value,
	}

	detailIDs, err := decodeLinks(r.store, obj, "ReleaseDetailsByTerritory")
	if err != nil {
		return nil, err
	}
	for _, detailID := range detailIDs {
		obj, err := r.store.Get(detailID)
		if err != nil {
			return nil, err
		}
		titleIDs, err := decodeLinks(r.store, obj, "Title")
		if err != nil {
			return nil, err
		}
		var titles []*TitleResolver
		for _, titleID := range titleIDs {
			obj, err := r.store.Get(titleID)
			if err != nil {
				return nil, err
			}
			titleText, err := DecodeObj(r.store, obj, "TitleText")
			if err != nil {
				return nil, err
			}
			var titleType metaxml.Value
			if v, err := DecodeObj(r.store, obj, "TitleType"); err == nil {
				titleType = *v
			}
			titles = append(titles, &TitleResolver{
				titleType: titleType.Value,
				titleText: titleText.Value,
			})
		}
		var artistName metaxml.Value
		if v, err := DecodeObj(r.store, obj, "DisplayArtistName"); err == nil {
			artistName = *v
		}
		var labelName metaxml.Value
		if v, err := DecodeObj(r.store, obj, "LabelName"); err == nil {
			labelName = *v
		}
		var territoryCode metaxml.Value
		if v, err := DecodeObj(r.store, obj, "TerritoryCode"); err == nil {
			territoryCode = *v
		}
		var genre metaxml.Value
		if v, err := DecodeObj(r.store, obj, "Genre", "GenreText"); err == nil {
			genre = *v
		}
		var releaseDate metaxml.Value
		if v, err := DecodeObj(r.store, obj, "ReleaseDate"); err == nil {
			releaseDate = *v
		}
		artistIDs, err := decodeLinks(r.store, obj, "DisplayArtist")
		if err != nil {
			return nil, err
		}
		var artists []*ArtistResolver
		for _, artistID := range artistIDs {
			party, err := r.partyResolver(artistID.String())
			if err != nil {
				return nil, err
			}
			var role metaxml.Value
			if v, err := DecodeObj(r.store, obj, "ArtistRole"); err == nil {
				role = *v
			}
			artists = append(artists, &ArtistResolver{
				party: party,
				role:  role.Value,
			})
		}
		release.detailsByTerritory = append(release.detailsByTerritory, &ReleaseDetailsResolver{
			cid:               obj.Cid().String(),
			title:             titles,
			displayArtist:     artists,
			displayArtistName: artistName.Value,
			labelName:         labelName.Value,
			territoryCode:     territoryCode.Value,
			releaseDate:       releaseDate.Value,
			genre:             genre.Value,
		})

	}
	return release, nil
}
