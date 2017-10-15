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
	"database/sql"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// GraphQLSchema is the GraphQL schema for the MusicBrainz META index.
const GraphQLSchema = `
schema {
  query: Query
}

type Query {
  registered_work(
  title: String,
  iswc:  String,
  composite_type: String,
  record_type: String
  ): [RegisteredWork]!

  publisher_control(
   publisher_sequence_n:              String,
   record_type:                       String,
   transaction_sequence_n:            String,
   record_sequence_n:                 String,
   interested_party_number:           String,
   publisher_name:                    String,
   publisher_unknown_indicator:       String,
   publisher_type:                    String,
   tax_id_number:                     String,
   publisher_ipi_name_number:         String,
   pr_affiliation_society_number:     String,
   submitter_agreement_number:        String,
   pr_ownership_share:                String,
   mr_society:                        String,
   mr_ownership_share:                String,
   sr_society:                        String,
   sr_ownership_share:                String,
   special_agreements_indicator:      String,
   first_recording_refusal_ind:       String,
   publisher_ipi_base_number:         String,
   inter_standard_agreement_code:     String,
   society_assigned_agreement_number: String,
   agreement_type:                    String,
   usa_license_ind:                   String
  ): [PublisherControl]!

	writer_control(
		transaction_sequence_n:   String,
    record_sequence_n:        String,
    interested_party_number:  String,
    writer_last_name:         String,
    writer_first_name:        String,
    writer_unknown_indicator: String,
    writer_designation_code:  String,
    tax_id_number:            String,
    writer_ipi_name:          String,
    writer_ipi_base_number:   String,
    personal_number:          String,
    record_type:              String
  ): [WriterControl]!

 transmission_header(
  sender_type: String,
  sender_id: String,
  record_type: String,
  sender_name: String
  ): [TransmissionHeader]!
}

type TransmissionHeader {
 cid:                      String!
 record_type:              String!
 sender_type:              String!
 sender_id:                String!
 sender_name:              String!
}

type PublisherControl {
 cid:                               String!
 record_type:                       String!
 publisher_sequence_n:              String!
 transaction_sequence_n:            String!
 record_sequence_n:                 String!
 interested_party_number:           String!
 publisher_name:                    String!
 publisher_unknown_indicator:       String!
 publisher_type:                    String!
 tax_id_number:                     String!
 publisher_ipi_name_number:         String!
 pr_affiliation_society_number:     String!
 submitter_agreement_number:        String!
 pr_ownership_share:                String!
 mr_society:                        String!
 mr_ownership_share:                String!
 sr_society:                        String!
 sr_ownership_share:                String!
 special_agreements_indicator:      String!
 first_recording_refusal_ind:       String!
 publisher_ipi_base_number:         String!
 inter_standard_agreement_code:     String!
 society_assigned_agreement_number: String!
 agreement_type:                    String!
 usa_license_ind:                   String!
}

type WriterControl {
	cid:                      String!
  transaction_sequence_n:   String!
  record_sequence_n:        String!
  interested_party_number:  String!
  writer_last_name:         String!
  writer_first_name:        String!
  writer_unknown_indicator: String!
  writer_designation_code:  String!
  tax_id_number:            String!
  writer_ipi_name:          String!
  writer_ipi_base_number:   String!
  personal_number:          String!
  record_type:              String!
}

type RegisteredWork {
 cid:                      String!
 record_type:              String!
 title:                    String!
 language_code:            String!
 submitte_worknumber:      String!
 iswc:                     String!
 copyright_date:           String!
 distribution_category:    String!
 duration:                 String!
 recorded_indicator:       String!
 textmusic_relationship:   String!
 composite_type:           String!
 version_type:             String!
 music_arrangement:        String!
 lyric_adaptation:         String!
 contact_name:             String!
 contact_id:               String!
 work_type:                String!
 grandrights_indicator:    String!
 composite_componentcount: String!
 date_ofpublication:       String!
 exceptional_clause:       String!
 opus_number:              String!
 catalogue_number:         String!
 priority_flag:            String!
 contributors:             [WriterControl]
}
`

// Resolver defines GraphQL resolver functions for the schema contained in
// the GraphQLSchema constant, retrieving data from a META store and SQLite3
// index.
type Resolver struct {
	db    *sql.DB
	store *meta.Store
}

// NewResolver returns a Resolver which retrieves data from the given META
// store and SQLite3 index.
func NewResolver(db *sql.DB, store *meta.Store) *Resolver {
	return &Resolver{db, store}
}

// registeredWorkArgs are the arguments for a GraphQL registeredWork query.
type registeredWorkArgs struct {
	RecordType    *string
	Title         *string
	ISWC          *string
	CompositeType *string
}

type publisherControlArgs struct {
	RecordType                     *string
	PublisherSequenceN             *string
	TransactionSequenceN           *string
	RecordSequenceN                *string
	PublisherSequenceNumber        *string
	InterestedPartyNumber          *string
	PublisherName                  *string
	PublisherUnknownIndicator      *string
	PublisherType                  *string
	TaxIDNumber                    *string
	PublisherIPINameNumber         *string
	PRAffiliationSocietyNumber     *string
	SubmitterAgreementNumber       *string
	PROwnershipShare               *string
	MRSociety                      *string
	MROwnershipShare               *string
	SRSociety                      *string
	SROwnershipShare               *string
	SpecialAgreementsIndicator     *string
	FirstRecordingRefusalInd       *string
	PublisherIPIBaseNumber         *string
	InterStandardAgreementCode     *string
	SocietyAssignedAgreementNumber *string
	AgreementType                  *string
	USALicenseInd                  *string
}

type writerControlArgs struct {
	RecordType             *string
	TransactionSequenceN   *string
	RecordSequenceN        *string
	InterestedPartyNumber  *string
	WriterLastName         *string
	WriterFirstName        *string
	WriterUnknownIndicator *string
	WriterDesignationCode  *string
	TaxIDNumber            *string
	WriterIPIName          *string
	WriterIPIBaseNumber    *string
	PersonalNumber         *string
}

type transmissionHeaderArgs struct {
	RecordType *string
	SenderType *string
	SenderID   *string
	SenderName *string
}

// RegisteredWork is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either an RegisteredWork RecordType, Title ,ISWC,or CompositeType, and loads the
// associated META objects from the META store.
func (g *Resolver) RegisteredWork(args registeredWorkArgs) ([]*registeredWorkResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.Title != nil:
		rows, err = g.db.Query("SELECT object_id FROM registered_work WHERE title = ?", *args.Title)
	case args.RecordType != nil:
		rows, err = g.db.Query("SELECT object_id FROM registered_work WHERE record_type = ?", *args.RecordType)
	case args.CompositeType != nil:
		rows, err = g.db.Query("SELECT object_id FROM registered_work WHERE composite_type = ?", *args.CompositeType)
	case args.ISWC != nil:
		rows, err = g.db.Query("SELECT object_id FROM registered_work WHERE iswc = ?", *args.ISWC)
	default:
		return nil, errors.New("missing title, record_type ,iswc or composite_type argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*registeredWorkResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		objCid, err := cid.Parse(objectID)
		if err != nil {
			return nil, err
		}
		obj, err := g.store.Get(objCid)
		if err != nil {
			return nil, err
		}
		var registeredWork RegisteredWork
		if err := obj.Decode(&registeredWork); err != nil {
			return nil, err
		}
		swrRows, err := g.db.Query("SELECT object_id FROM writer_control WHERE tx_id = ?", objectID)
		for swrRows.Next() {
			var objectID string
			if err := swrRows.Scan(&objectID); err != nil {
				return nil, err
			}
			objCid, err := cid.Parse(objectID)
			if err != nil {
				return nil, err
			}
			obj, err := g.store.Get(objCid)
			if err != nil {
				return nil, err
			}
			var writerControlledbySubmitter WriterControlledbySubmitter
			if err := obj.Decode(&writerControlledbySubmitter); err != nil {
				return nil, err
			}
			registeredWork.Contributors = append(registeredWork.Contributors, &writerControlledbySubmitter)
		}

		resolvers = append(resolvers, &registeredWorkResolver{objectID, &registeredWork})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// PublisherControl is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either a PublihserControl RecordType or publisher_sequence_n and loads the
// associated META objects from the META store.
func (g *Resolver) PublisherControl(args publisherControlArgs) ([]*publisherControlResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.PublisherSequenceN != nil:
		rows, err = g.db.Query("SELECT object_id FROM publisher_control WHERE publisher_sequence_n = ?", *args.PublisherSequenceN)
	case args.RecordType != nil:
		rows, err = g.db.Query("SELECT object_id FROM publisher_control WHERE record_type = ?", *args.RecordType)
	case args.TransactionSequenceN != nil:
		rows, err = g.db.Query("SELECT object_id FROM publisher_control WHERE transaction_sequence_n = ?", *args.TransactionSequenceN)
	default:
		return nil, errors.New("missing record_type or PublisherSequenceN argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*publisherControlResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		cid, err := cid.Parse(objectID)
		if err != nil {
			return nil, err
		}
		obj, err := g.store.Get(cid)
		if err != nil {
			return nil, err
		}
		var publisherControllBySubmitter PublisherControllBySubmitter
		if err := obj.Decode(&publisherControllBySubmitter); err != nil {
			return nil, err
		}
		resolvers = append(resolvers, &publisherControlResolver{objectID, &publisherControllBySubmitter})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// WriterControl is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either a WriterControl RecordType ,writer_first_name,
// writer_last_name or writer_ipi_base_number and loads the associated META
// objects from the META store.
func (g *Resolver) WriterControl(args writerControlArgs) ([]*writerControlResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.WriterFirstName != nil:
		rows, err = g.db.Query("SELECT object_id FROM writer_control WHERE writer_first_name = ?", *args.WriterFirstName)
	case args.RecordType != nil:
		rows, err = g.db.Query("SELECT object_id FROM writer_control WHERE record_type = ?", *args.RecordType)
	case args.WriterLastName != nil:
		rows, err = g.db.Query("SELECT object_id FROM writer_control WHERE writer_last_name = ?", *args.WriterLastName)
	case args.WriterIPIBaseNumber != nil:
		rows, err = g.db.Query("SELECT object_id FROM writer_control WHERE writer_ipi_base_number = ?", *args.WriterIPIBaseNumber)
	case args.WriterUnknownIndicator != nil:
		rows, err = g.db.Query("SELECT object_id FROM writer_control WHERE writer_unknown_indicator = ?", *args.WriterUnknownIndicator)
	default:
		return nil, errors.New("missing writer_first_name,record_type,writer_last_name,writer_unknown_indicator or writer_ipi_base_number argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*writerControlResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		cid, err := cid.Parse(objectID)
		if err != nil {
			return nil, err
		}
		obj, err := g.store.Get(cid)
		if err != nil {
			return nil, err
		}
		var writerControlledbySubmitter WriterControlledbySubmitter
		if err := obj.Decode(&writerControlledbySubmitter); err != nil {
			return nil, err
		}
		resolvers = append(resolvers, &writerControlResolver{objectID, &writerControlledbySubmitter})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resolvers, nil
}

// TransmissionHeader is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either a transmission header RecordType,sender_type,sender_type or sender_name and loads the
// associated META objects from the META store.
func (g *Resolver) TransmissionHeader(args transmissionHeaderArgs) ([]*transmissionHeaderResolver, error) {
	var rows *sql.Rows
	var err error
	switch {
	case args.SenderType != nil:
		rows, err = g.db.Query("SELECT object_id FROM transmission_header WHERE sender_type = ?", *args.SenderType)
	case args.SenderID != nil:
		rows, err = g.db.Query("SELECT object_id FROM transmission_header WHERE sender_id = ?", *args.SenderID)
	case args.RecordType != nil:
		rows, err = g.db.Query("SELECT object_id FROM transmission_header WHERE record_type = ?", *args.RecordType)
	case args.SenderName != nil:
		rows, err = g.db.Query("SELECT object_id FROM transmission_header WHERE sender_name = ?", *args.SenderName)
	default:
		return nil, errors.New("missing record_type,sender_type,sender_id,record_type or sender_name argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*transmissionHeaderResolver
	for rows.Next() {
		var objectID string
		if err := rows.Scan(&objectID); err != nil {
			return nil, err
		}
		cid, err := cid.Parse(objectID)
		if err != nil {
			return nil, err
		}
		obj, err := g.store.Get(cid)
		if err != nil {
			return nil, err
		}
		var transmissionHeader TransmissionHeader
		if err := obj.Decode(&transmissionHeader); err != nil {
			return nil, err
		}
		resolvers = append(resolvers, &transmissionHeaderResolver{objectID, &transmissionHeader})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// transmissionHeaderResolver defines GraphQL resolver functions for transmissionHeader fields.
type transmissionHeaderResolver struct {
	cid                string
	transmissionHeader *TransmissionHeader
}

func (t *transmissionHeaderResolver) Cid() string {
	return t.cid
}

func (t *transmissionHeaderResolver) RecordType() string {
	return t.transmissionHeader.RecordType
}

func (t *transmissionHeaderResolver) SenderID() string {
	return t.transmissionHeader.SenderID
}

func (t *transmissionHeaderResolver) SenderType() string {
	return t.transmissionHeader.SenderType
}

func (t *transmissionHeaderResolver) SenderName() string {
	return t.transmissionHeader.SenderName
}

// publisherControlResolver defines GraphQL resolver functions for publisherControl fields.
type publisherControlResolver struct {
	cid              string
	publisherControl *PublisherControllBySubmitter
}

func (p *publisherControlResolver) Cid() string {
	return p.cid
}

func (p *publisherControlResolver) RecordType() string {
	return p.publisherControl.RecordType
}

func (p *publisherControlResolver) PublisherSequenceN() string {
	return p.publisherControl.PublisherSequenceNumber
}

func (p *publisherControlResolver) TransactionSequenceN() string {
	return p.publisherControl.TransactionSequenceN
}

func (p *publisherControlResolver) RecordSequenceN() string {
	return p.publisherControl.RecordSequenceN
}

func (p *publisherControlResolver) InterestedPartyNumber() string {
	return p.publisherControl.InterestedPartyNumber
}

func (p *publisherControlResolver) MROwnershipShare() string {
	return p.publisherControl.MROwnershipShare
}

func (p *publisherControlResolver) AgreementType() string {
	return p.publisherControl.AgreementType
}

func (p *publisherControlResolver) FirstRecordingRefusalInd() string {
	return p.publisherControl.FirstRecordingRefusalInd
}

func (p *publisherControlResolver) InterStandardAgreementCode() string {
	return p.publisherControl.InterStandardAgreementCode
}

func (p *publisherControlResolver) PublisherIPIBaseNumber() string {
	return p.publisherControl.PublisherIPIBaseNumber
}

func (p *publisherControlResolver) PRAffiliationSocietyNumber() string {
	return p.publisherControl.PRAffiliationSocietyNumber
}

func (p *publisherControlResolver) PROwnershipShare() string {
	return p.publisherControl.PROwnershipShare
}

func (p *publisherControlResolver) PublisherIPINameNumber() string {
	return p.publisherControl.PublisherIPINameNumber
}

func (p *publisherControlResolver) PublisherName() string {
	return p.publisherControl.PublisherName
}

func (p *publisherControlResolver) PublisherUnknownIndicator() string {
	return p.publisherControl.PublisherUnknownIndicator
}

func (p *publisherControlResolver) SROwnershipShare() string {
	return p.publisherControl.SROwnershipShare
}

func (p *publisherControlResolver) SRSociety() string {
	return p.publisherControl.SRSociety
}

func (p *publisherControlResolver) SocietyAssignedAgreementNumber() string {
	return p.publisherControl.SocietyAssignedAgreementNumber
}

func (p *publisherControlResolver) TaxIDNumber() string {
	return p.publisherControl.TaxIDNumber
}

func (p *publisherControlResolver) PublisherType() string {
	return p.publisherControl.PublisherType
}

func (p *publisherControlResolver) SubmitterAgreementNumber() string {
	return p.publisherControl.SubmitterAgreementNumber
}

func (p *publisherControlResolver) SpecialAgreementsIndicator() string {
	return p.publisherControl.SpecialAgreementsIndicator
}

func (p *publisherControlResolver) MRSociety() string {
	return p.publisherControl.MRSociety
}

func (p *publisherControlResolver) USALicenseInd() string {
	return p.publisherControl.USALicenseInd
}

// writerControlResolver defines GraphQL resolver functions for WriterControlledbySubmitter fields.
type writerControlResolver struct {
	cid           string
	writerControl *WriterControlledbySubmitter
}

func (p *writerControlResolver) Cid() string {
	return p.cid
}

func (p *writerControlResolver) WriterUnknownIndicator() string {
	return p.writerControl.WriterUnknownIndicator
}

func (p *writerControlResolver) TransactionSequenceN() string {
	return p.writerControl.TransactionSequenceN
}

func (p *writerControlResolver) WriterDesignationCode() string {
	return p.writerControl.WriterDesignationCode
}

func (p *writerControlResolver) WriterFirstName() string {
	return p.writerControl.WriterFirstName
}

func (p *writerControlResolver) WriterLastName() string {
	return p.writerControl.WriterLastName
}

func (p *writerControlResolver) WriterIPIBaseNumber() string {
	return p.writerControl.WriterIPIBaseNumber
}

func (p *writerControlResolver) WriterIPIName() string {
	return p.writerControl.WriterIPIName
}

func (p *writerControlResolver) RecordSequenceN() string {
	return p.writerControl.RecordSequenceN
}

func (p *writerControlResolver) InterestedPartyNumber() string {
	return p.writerControl.InterestedPartyNumber
}

func (p *writerControlResolver) TaxIDNumber() string {
	return p.writerControl.TaxIDNumber
}

func (p *writerControlResolver) PersonalNumber() string {
	return p.writerControl.PersonalNumber
}

func (p *writerControlResolver) RecordType() string {
	return p.writerControl.RecordType
}

// registeredWorkResolver defines GraphQL resolver functions for registeredWork fields.
type registeredWorkResolver struct {
	cid            string
	registeredWork *RegisteredWork
}

func (r *registeredWorkResolver) Contributors() *[]*writerControlResolver {
	var writerControlResolvers []*writerControlResolver
	for _, c := range r.registeredWork.Contributors {
		writerControlResolvers = append(writerControlResolvers, &writerControlResolver{writerControl: c})
	}
	return &writerControlResolvers
}

func (r *registeredWorkResolver) Cid() string {
	return r.cid
}

func (r *registeredWorkResolver) Title() string {
	return r.registeredWork.Title
}

func (r *registeredWorkResolver) RecordType() string {
	return r.registeredWork.RecordType
}

func (r *registeredWorkResolver) ISWC() string {
	return r.registeredWork.ISWC
}

func (r *registeredWorkResolver) CatalogueNumber() string {
	return r.registeredWork.CatalogueNumber
}

func (r *registeredWorkResolver) CompositeComponentCount() string {
	return r.registeredWork.CompositeComponentCount
}

func (r *registeredWorkResolver) CompositeType() string {
	return r.registeredWork.CompositeType
}

func (r *registeredWorkResolver) ContactID() string {
	return r.registeredWork.ContactID
}

func (r *registeredWorkResolver) ContactName() string {
	return r.registeredWork.ContactName
}

func (r *registeredWorkResolver) CopyRightDate() string {
	return r.registeredWork.CopyRightDate
}

func (r *registeredWorkResolver) DateOfPublication() string {
	return r.registeredWork.DateOfPublication
}

func (r *registeredWorkResolver) DistributionCategory() string {
	return r.registeredWork.DistributionCategory
}

func (r *registeredWorkResolver) Duration() string {
	return r.registeredWork.Duration
}

func (r *registeredWorkResolver) ExceptionalClause() string {
	return r.registeredWork.ExceptionalClause
}

func (r *registeredWorkResolver) GrandRightsIndicator() string {
	return r.registeredWork.GrandRightsIndicator
}
func (r *registeredWorkResolver) LanguageCode() string {
	return r.registeredWork.LanguageCode
}

func (r *registeredWorkResolver) LyricAdaptation() string {
	return r.registeredWork.LyricAdaptation
}

func (r *registeredWorkResolver) MusicArrangement() string {
	return r.registeredWork.MusicArrangement
}
func (r *registeredWorkResolver) OpusNumber() string {
	return r.registeredWork.OpusNumber
}

func (r *registeredWorkResolver) PriorityFlag() string {
	return r.registeredWork.PriorityFlag
}

func (r *registeredWorkResolver) RecordedIndicator() string {
	return r.registeredWork.RecordedIndicator
}
func (r *registeredWorkResolver) SubmitteWorkNumber() string {
	return r.registeredWork.SubmitteWorkNumber
}

func (r *registeredWorkResolver) TextMusicRelationship() string {
	return r.registeredWork.TextMusicRelationship
}

func (r *registeredWorkResolver) VersionType() string {
	return r.registeredWork.VersionType
}
func (r *registeredWorkResolver) WorkType() string {
	return r.registeredWork.WorkType
}
