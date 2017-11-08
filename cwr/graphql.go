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

	cid "github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// GraphQLSchema is the GraphQL schema for the MusicBrainz META index.
const GraphQLSchema = `
schema {
  query: Query
}

type Query {
  registered_work(
    title:          String,
    iswc:           String,
    composite_type: String,
    record_type:    String
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
  cid:         String!
  source:      String!
  record_type: String!
  sender_type: String!
  sender_id:   String!
  sender_name: String!
}

type PublisherControl {
  cid:                               String!
  source:                            String!
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
  source:                   String!
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
  pr_ownership_share:       String!
  mr_ownership_share:       String!
  sr_ownership_share:       String!
}

type RegisteredWork {
  cid:                      String!
  source:                   String!
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
  contributors:             [WriterControl]!
  controls:                 [PublisherControl]!
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

// RegisteredWorkArgs are the arguments for a GraphQL registeredWork query.
type RegisteredWorkArgs struct {
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

type WriterControlArgs struct {
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
func (g *Resolver) RegisteredWork(args RegisteredWorkArgs) ([]*RegisteredWorkResolver, error) {
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
	var resolvers []*RegisteredWorkResolver
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
		if err != nil {
			return nil, err
		}
		for swrRows.Next() {
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
		spuRows, err := g.db.Query("SELECT object_id FROM publisher_control WHERE tx_id = ?", objectID)
		if err != nil {
			return nil, err
		}
		for spuRows.Next() {
			if err := spuRows.Scan(&objectID); err != nil {
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
			var publisherControllBySubmitter PublisherControllBySubmitter
			if err := obj.Decode(&publisherControllBySubmitter); err != nil {
				return nil, err
			}
			registeredWork.Controls = append(registeredWork.Controls, &publisherControllBySubmitter)
		}
		resolvers = append(resolvers, &RegisteredWorkResolver{objectID, &registeredWork})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// PublisherControl is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either a PublihserControl RecordType or publisher_sequence_n and loads the
// associated META objects from the META store.
func (g *Resolver) PublisherControl(args publisherControlArgs) ([]*PublisherControlResolver, error) {
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
	var resolvers []*PublisherControlResolver
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
		resolvers = append(resolvers, &PublisherControlResolver{objectID, &publisherControllBySubmitter})
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
func (g *Resolver) WriterControl(args WriterControlArgs) ([]*WriterControlResolver, error) {
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
	case args.WriterIPIName != nil:
		rows, err = g.db.Query("SELECT object_id FROM writer_control WHERE writer_ipi_name = ?", *args.WriterIPIName)
	case args.WriterUnknownIndicator != nil:
		rows, err = g.db.Query("SELECT object_id FROM writer_control WHERE writer_unknown_indicator = ?", *args.WriterUnknownIndicator)
	default:
		return nil, errors.New("missing writer_first_name,record_type,writer_last_name,writer_unknown_indicator, writer_ipi_base_number or writer_ipi_name argument")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var resolvers []*WriterControlResolver
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
		resolvers = append(resolvers, &WriterControlResolver{objectID, &writerControlledbySubmitter})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resolvers, nil
}

// TransmissionHeader is a GraphQL resolver function which retrieves object IDs from the
// SQLite3 index using either a transmission header RecordType,sender_type,sender_type or sender_name and loads the
// associated META objects from the META store.
func (g *Resolver) TransmissionHeader(args transmissionHeaderArgs) ([]*TransmissionHeaderResolver, error) {
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
	var resolvers []*TransmissionHeaderResolver
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
		resolvers = append(resolvers, &TransmissionHeaderResolver{objectID, &transmissionHeader})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resolvers, nil
}

// TransmissionHeaderResolver defines GraphQL resolver functions for transmissionHeader fields.
type TransmissionHeaderResolver struct {
	cid                string
	transmissionHeader *TransmissionHeader
}

// Cid resolver
func (t *TransmissionHeaderResolver) Cid() string {
	return t.cid
}

// Source resolver
func (t *TransmissionHeaderResolver) Source() string {
	return t.transmissionHeader.Source
}

// RecordType resolver
func (t *TransmissionHeaderResolver) RecordType() string {
	return t.transmissionHeader.RecordType
}

// SenderID resolver
func (t *TransmissionHeaderResolver) SenderID() string {
	return t.transmissionHeader.SenderID
}

// SenderType resolver
func (t *TransmissionHeaderResolver) SenderType() string {
	return t.transmissionHeader.SenderType
}

// SenderName resolver
func (t *TransmissionHeaderResolver) SenderName() string {
	return t.transmissionHeader.SenderName
}

// PublisherControlResolver defines GraphQL resolver functions for publisherControl fields.
type PublisherControlResolver struct {
	cid              string
	publisherControl *PublisherControllBySubmitter
}

// Cid resolver
func (p *PublisherControlResolver) Cid() string {
	return p.cid
}

// Source resolver
func (p *PublisherControlResolver) Source() string {
	return p.publisherControl.Source
}

// RecordType resolver
func (p *PublisherControlResolver) RecordType() string {
	return p.publisherControl.RecordType
}

// PublisherSequenceN resolver
func (p *PublisherControlResolver) PublisherSequenceN() string {
	return p.publisherControl.PublisherSequenceNumber
}

// TransactionSequenceN resolver
func (p *PublisherControlResolver) TransactionSequenceN() string {
	return p.publisherControl.TransactionSequenceN
}

// RecordSequenceN resolver
func (p *PublisherControlResolver) RecordSequenceN() string {
	return p.publisherControl.RecordSequenceN
}

// InterestedPartyNumber resolver
func (p *PublisherControlResolver) InterestedPartyNumber() string {
	return p.publisherControl.InterestedPartyNumber
}

// MROwnershipShare resolver
func (p *PublisherControlResolver) MROwnershipShare() string {
	return p.publisherControl.MROwnershipShare
}

// AgreementType resolver
func (p *PublisherControlResolver) AgreementType() string {
	return p.publisherControl.AgreementType
}

// FirstRecordingRefusalInd resolver
func (p *PublisherControlResolver) FirstRecordingRefusalInd() string {
	return p.publisherControl.FirstRecordingRefusalInd
}

// InterStandardAgreementCode resolver
func (p *PublisherControlResolver) InterStandardAgreementCode() string {
	return p.publisherControl.InterStandardAgreementCode
}

// PublisherIPIBaseNumber resolver
func (p *PublisherControlResolver) PublisherIPIBaseNumber() string {
	return p.publisherControl.PublisherIPIBaseNumber
}

// PRAffiliationSocietyNumber resolver
func (p *PublisherControlResolver) PRAffiliationSocietyNumber() string {
	return p.publisherControl.PRAffiliationSocietyNumber
}

// PROwnershipShare resolver
func (p *PublisherControlResolver) PROwnershipShare() string {
	return p.publisherControl.PROwnershipShare
}

// PublisherIPINameNumber resolver
func (p *PublisherControlResolver) PublisherIPINameNumber() string {
	return p.publisherControl.PublisherIPINameNumber
}

// PublisherName resolver
func (p *PublisherControlResolver) PublisherName() string {
	return p.publisherControl.PublisherName
}

// PublisherUnknownIndicator resolver
func (p *PublisherControlResolver) PublisherUnknownIndicator() string {
	return p.publisherControl.PublisherUnknownIndicator
}

// SROwnershipShare resolver
func (p *PublisherControlResolver) SROwnershipShare() string {
	return p.publisherControl.SROwnershipShare
}

// SRSociety resolver
func (p *PublisherControlResolver) SRSociety() string {
	return p.publisherControl.SRSociety
}

// SocietyAssignedAgreementNumber resolver
func (p *PublisherControlResolver) SocietyAssignedAgreementNumber() string {
	return p.publisherControl.SocietyAssignedAgreementNumber
}

// TaxIDNumber resolver
func (p *PublisherControlResolver) TaxIDNumber() string {
	return p.publisherControl.TaxIDNumber
}

// PublisherType resolver
func (p *PublisherControlResolver) PublisherType() string {
	return p.publisherControl.PublisherType
}

// SubmitterAgreementNumber resolver
func (p *PublisherControlResolver) SubmitterAgreementNumber() string {
	return p.publisherControl.SubmitterAgreementNumber
}

// SpecialAgreementsIndicator resolver
func (p *PublisherControlResolver) SpecialAgreementsIndicator() string {
	return p.publisherControl.SpecialAgreementsIndicator
}

// MRSociety resolver
func (p *PublisherControlResolver) MRSociety() string {
	return p.publisherControl.MRSociety
}

// USALicenseInd resolver
func (p *PublisherControlResolver) USALicenseInd() string {
	return p.publisherControl.USALicenseInd
}

// WriterControlResolver defines GraphQL resolver functions for WriterControlledbySubmitter fields.
type WriterControlResolver struct {
	cid           string
	writerControl *WriterControlledbySubmitter
}

// Cid resolver
func (p *WriterControlResolver) Cid() string {
	return p.cid
}

// Source resolver
func (p *WriterControlResolver) Source() string {
	return p.writerControl.Source
}

// WriterUnknownIndicator resolver
func (p *WriterControlResolver) WriterUnknownIndicator() string {
	return p.writerControl.WriterUnknownIndicator
}

// TransactionSequenceN resolver
func (p *WriterControlResolver) TransactionSequenceN() string {
	return p.writerControl.TransactionSequenceN
}

// WriterDesignationCode resolver
func (p *WriterControlResolver) WriterDesignationCode() string {
	return p.writerControl.WriterDesignationCode
}

// WriterFirstName resolver
func (p *WriterControlResolver) WriterFirstName() string {
	return p.writerControl.WriterFirstName
}

// WriterLastName resolver
func (p *WriterControlResolver) WriterLastName() string {
	return p.writerControl.WriterLastName
}

// WriterIPIBaseNumber resolver
func (p *WriterControlResolver) WriterIPIBaseNumber() string {
	return p.writerControl.WriterIPIBaseNumber
}

// WriterIPIName resolver
func (p *WriterControlResolver) WriterIPIName() string {
	return p.writerControl.WriterIPIName
}

// RecordSequenceN resolver
func (p *WriterControlResolver) RecordSequenceN() string {
	return p.writerControl.RecordSequenceN
}

// InterestedPartyNumber resolver
func (p *WriterControlResolver) InterestedPartyNumber() string {
	return p.writerControl.InterestedPartyNumber
}

// TaxIDNumber resolver
func (p *WriterControlResolver) TaxIDNumber() string {
	return p.writerControl.TaxIDNumber
}

// PersonalNumber resolver
func (p *WriterControlResolver) PersonalNumber() string {
	return p.writerControl.PersonalNumber
}

// RecordType resolver
func (p *WriterControlResolver) RecordType() string {
	return p.writerControl.RecordType
}

// PROwnershipShare resolver
func (p *WriterControlResolver) PROwnershipShare() string {
	return p.writerControl.PROwnershipShare
}

// MROwnershipShare resolver
func (p *WriterControlResolver) MROwnershipShare() string {
	return p.writerControl.MROwnershipShare
}

// SROwnershipShare resolver
func (p *WriterControlResolver) SROwnershipShare() string {
	return p.writerControl.SROwnershipShare
}

// RegisteredWorkResolver defines GraphQL resolver functions for registeredWork fields.
type RegisteredWorkResolver struct {
	cid            string
	registeredWork *RegisteredWork
}

// Contributors resolver
func (r *RegisteredWorkResolver) Contributors() []*WriterControlResolver {
	var writerControlResolvers []*WriterControlResolver
	for _, c := range r.registeredWork.Contributors {
		writerControlResolvers = append(writerControlResolvers, &WriterControlResolver{writerControl: c})
	}
	return writerControlResolvers
}

// Controls resolver
func (r *RegisteredWorkResolver) Controls() []*PublisherControlResolver {
	var publisherControlResolvers []*PublisherControlResolver
	for _, c := range r.registeredWork.Controls {
		publisherControlResolvers = append(publisherControlResolvers, &PublisherControlResolver{publisherControl: c})
	}
	return publisherControlResolvers
}

// Cid resolver
func (r *RegisteredWorkResolver) Cid() string {
	return r.cid
}

// Source resolver
func (r *RegisteredWorkResolver) Source() string {
	return r.registeredWork.Source
}

// Title resolver
func (r *RegisteredWorkResolver) Title() string {
	return r.registeredWork.Title
}

// RecordType resolver
func (r *RegisteredWorkResolver) RecordType() string {
	return r.registeredWork.RecordType
}

// ISWC resolver
func (r *RegisteredWorkResolver) ISWC() string {
	return r.registeredWork.ISWC
}

// CatalogueNumber resolver
func (r *RegisteredWorkResolver) CatalogueNumber() string {
	return r.registeredWork.CatalogueNumber
}

// CompositeComponentCount resolver
func (r *RegisteredWorkResolver) CompositeComponentCount() string {
	return r.registeredWork.CompositeComponentCount
}

// CompositeType resolver
func (r *RegisteredWorkResolver) CompositeType() string {
	return r.registeredWork.CompositeType
}

// ContactID resolver
func (r *RegisteredWorkResolver) ContactID() string {
	return r.registeredWork.ContactID
}

// ContactName resolver
func (r *RegisteredWorkResolver) ContactName() string {
	return r.registeredWork.ContactName
}

// CopyRightDate resolver
func (r *RegisteredWorkResolver) CopyRightDate() string {
	return r.registeredWork.CopyRightDate
}

// DateOfPublication resolver
func (r *RegisteredWorkResolver) DateOfPublication() string {
	return r.registeredWork.DateOfPublication
}

// DistributionCategory resolver
func (r *RegisteredWorkResolver) DistributionCategory() string {
	return r.registeredWork.DistributionCategory
}

// Duration resolver
func (r *RegisteredWorkResolver) Duration() string {
	return r.registeredWork.Duration
}

// ExceptionalClause resolver
func (r *RegisteredWorkResolver) ExceptionalClause() string {
	return r.registeredWork.ExceptionalClause
}

// GrandRightsIndicator resolver
func (r *RegisteredWorkResolver) GrandRightsIndicator() string {
	return r.registeredWork.GrandRightsIndicator
}

// LanguageCode resolver
func (r *RegisteredWorkResolver) LanguageCode() string {
	return r.registeredWork.LanguageCode
}

// LyricAdaptation resolver
func (r *RegisteredWorkResolver) LyricAdaptation() string {
	return r.registeredWork.LyricAdaptation
}

// MusicArrangement resolver
func (r *RegisteredWorkResolver) MusicArrangement() string {
	return r.registeredWork.MusicArrangement
}

// OpusNumber resolver
func (r *RegisteredWorkResolver) OpusNumber() string {
	return r.registeredWork.OpusNumber
}

// PriorityFlag resolver
func (r *RegisteredWorkResolver) PriorityFlag() string {
	return r.registeredWork.PriorityFlag
}

// RecordedIndicator resolver
func (r *RegisteredWorkResolver) RecordedIndicator() string {
	return r.registeredWork.RecordedIndicator
}

// SubmitteWorkNumber resolver
func (r *RegisteredWorkResolver) SubmitteWorkNumber() string {
	return r.registeredWork.SubmitteWorkNumber
}

// TextMusicRelationship resolver
func (r *RegisteredWorkResolver) TextMusicRelationship() string {
	return r.registeredWork.TextMusicRelationship
}

// VersionType resolver
func (r *RegisteredWorkResolver) VersionType() string {
	return r.registeredWork.VersionType
}

// WorkType resolver
func (r *RegisteredWorkResolver) WorkType() string {
	return r.registeredWork.WorkType
}
