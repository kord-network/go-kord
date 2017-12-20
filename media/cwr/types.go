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
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// HDR is a CWR Transmission Header record
type HDR struct {
	SenderType       string `cwr:"length=2"`
	SenderID         string `cwr:"length=9"`
	SenderName       string `cwr:"length=45"`
	EDIVersion       string `cwr:"length=5"`
	CreationDate     string `cwr:"length=8"`
	CreationTime     string `cwr:"length=6"`
	TransmissionDate string `cwr:"length=8"`
	CharacterSet     string `cwr:"length=15"`
}

// GRH is a CWR Group Header record
type GRH struct {
	TransactionType string `cwr:"length=3"`
	GroupID         string `cwr:"length=5"`
	VersionNumber   string `cwr:"length=5"`
	BatchRequest    string `cwr:"length=10"`
	SubmissionType  string `cwr:"length=2"`
}

// GRT is a CWR Group Trailer record
type GRT struct {
	GroupID          string `cwr:"length=5"`
	TransactionCount string `cwr:"length=8"`
	RecordCount      string `cwr:"length=8"`
	Currency         string `cwr:"length=3"`
	TotalValue       string `cwr:"length=10"`
}

// TRL is a CWR Transmission Trailer record
type TRL struct {
	GroupCount       string `cwr:"length=5"`
	TransactionCount string `cwr:"length=8"`
	RecordCount      string `cwr:"length=8"`
}

// AGR is a CWR Agreement Supporting Work Registration record
type AGR struct {
	TransactionSequenceNumber          string `cwr:"length=8"`
	RecordSequenceNumber               string `cwr:"length=8"`
	SubmitterAgreementNumber           string `cwr:"length=14"`
	InternationalStandardAgreementCode string `cwr:"length=14"`
	AgreementType                      string `cwr:"length=2"`
	AgreementStartDate                 string `cwr:"length=8"`
	AgreementEndDate                   string `cwr:"length=8"`
	RetentionEndDate                   string `cwr:"length=8"`
	PriorRoyaltyStatus                 string `cwr:"length=1"`
	PriorRoyaltyStartDate              string `cwr:"length=8"`
	PostTermCollectionStatus           string `cwr:"length=1"`
	PostTermCollectionEndDate          string `cwr:"length=8"`
	SignatureAgreementDate             string `cwr:"length=8"`
	NumberOfWorks                      string `cwr:"length=5"`
	SalesManufactureClause             string `cwr:"length=1"`
	SharesChange                       string `cwr:"length=1"`
	AdvanceGiven                       string `cwr:"length=1"`
	SocietyAssignedAgreementNumber     string `cwr:"length=14"`
}

// NWR is a CWR New Work Registration record
type NWR struct {
	TransactionSequenceNumber       string `cwr:"length=8"`
	RecordSequenceNumber            string `cwr:"length=8"`
	WorkTitle                       string `cwr:"length=60"`
	LanguageCode                    string `cwr:"length=2"`
	SubmitterWorkNumber             string `cwr:"length=14"`
	ISWC                            string `cwr:"length=11"`
	CopyrightDate                   string `cwr:"length=8"`
	CopyrightNumber                 string `cwr:"length=12"`
	MusicalWorkDistributionCategory string `cwr:"length=3"`
	Duration                        string `cwr:"length=6"`
	RecordedIndicator               string `cwr:"length=1"`
	TextMusicRelationship           string `cwr:"length=3"`
	CompositeType                   string `cwr:"length=3"`
	VersionType                     string `cwr:"length=3"`
	ExcerptType                     string `cwr:"length=3"`
	MusicArrangement                string `cwr:"length=3"`
	LyricAdaptation                 string `cwr:"length=3"`
	ContactName                     string `cwr:"length=30"`
	ContactID                       string `cwr:"length=10"`
	CWRWorkType                     string `cwr:"length=2"`
	GrandRightsInd                  string `cwr:"length=1"`
	CompositeComponentCount         string `cwr:"length=3"`
	PrintedEditionPublicationDate   string `cwr:"length=8"`
	ExceptionalClause               string `cwr:"length=1"`
	OpusNumber                      string `cwr:"length=25"`
	CatalogNumber                   string `cwr:"length=25"`
	PriorityFlag                    string `cwr:"length=1"`
}

// REV is a CWR Revised Registration record
type REV struct {
	TransactionSequenceNumber       string `cwr:"length=8"`
	RecordSequenceNumber            string `cwr:"length=8"`
	WorkTitle                       string `cwr:"length=60"`
	LanguageCode                    string `cwr:"length=2"`
	SubmitterWorkNumber             string `cwr:"length=14"`
	ISWC                            string `cwr:"length=11"`
	CopyrightDate                   string `cwr:"length=8"`
	CopyrightNumber                 string `cwr:"length=12"`
	MusicalWorkDistributionCategory string `cwr:"length=3"`
	Duration                        string `cwr:"length=6"`
	RecordedIndicator               string `cwr:"length=1"`
	TextMusicRelationship           string `cwr:"length=3"`
	CompositeType                   string `cwr:"length=3"`
	VersionType                     string `cwr:"length=3"`
	ExcerptType                     string `cwr:"length=3"`
	MusicArrangement                string `cwr:"length=3"`
	LyricAdaptation                 string `cwr:"length=3"`
	ContactName                     string `cwr:"length=30"`
	ContactID                       string `cwr:"length=10"`
	CWRWorkType                     string `cwr:"length=2"`
	GrandRightsInd                  string `cwr:"length=1"`
	CompositeComponentCount         string `cwr:"length=3"`
	PrintedEditionPublicationDate   string `cwr:"length=8"`
	ExceptionalClause               string `cwr:"length=1"`
	OpusNumber                      string `cwr:"length=25"`
	CatalogNumber                   string `cwr:"length=25"`
	PriorityFlag                    string `cwr:"length=1"`
}

// ISW is a CWR Notification of ISWC assign to a work record
type ISW struct {
	TransactionSequenceNumber       string `cwr:"length=8"`
	RecordSequenceNumber            string `cwr:"length=8"`
	WorkTitle                       string `cwr:"length=60"`
	LanguageCode                    string `cwr:"length=2"`
	SubmitterWorkNumber             string `cwr:"length=14"`
	ISWC                            string `cwr:"length=11"`
	CopyrightDate                   string `cwr:"length=8"`
	CopyrightNumber                 string `cwr:"length=12"`
	MusicalWorkDistributionCategory string `cwr:"length=3"`
	Duration                        string `cwr:"length=6"`
	RecordedIndicator               string `cwr:"length=1"`
	TextMusicRelationship           string `cwr:"length=3"`
	CompositeType                   string `cwr:"length=3"`
	VersionType                     string `cwr:"length=3"`
	ExcerptType                     string `cwr:"length=3"`
	MusicArrangement                string `cwr:"length=3"`
	LyricAdaptation                 string `cwr:"length=3"`
	ContactName                     string `cwr:"length=30"`
	ContactID                       string `cwr:"length=10"`
	CWRWorkType                     string `cwr:"length=2"`
	GrandRightsInd                  string `cwr:"length=1"`
	CompositeComponentCount         string `cwr:"length=3"`
	PrintedEditionPublicationDate   string `cwr:"length=8"`
	ExceptionalClause               string `cwr:"length=1"`
	OpusNumber                      string `cwr:"length=25"`
	CatalogNumber                   string `cwr:"length=25"`
	PriorityFlag                    string `cwr:"length=1"`
}

// EXC is a CWR Existing work which is in Conflict with a Work Registration record
type EXC struct {
	TransactionSequenceNumber       string `cwr:"length=8"`
	RecordSequenceNumber            string `cwr:"length=8"`
	WorkTitle                       string `cwr:"length=60"`
	LanguageCode                    string `cwr:"length=2"`
	SubmitterWorkNumber             string `cwr:"length=14"`
	ISWC                            string `cwr:"length=11"`
	CopyrightDate                   string `cwr:"length=8"`
	CopyrightNumber                 string `cwr:"length=12"`
	MusicalWorkDistributionCategory string `cwr:"length=3"`
	Duration                        string `cwr:"length=6"`
	RecordedIndicator               string `cwr:"length=1"`
	TextMusicRelationship           string `cwr:"length=3"`
	CompositeType                   string `cwr:"length=3"`
	VersionType                     string `cwr:"length=3"`
	ExcerptType                     string `cwr:"length=3"`
	MusicArrangement                string `cwr:"length=3"`
	LyricAdaptation                 string `cwr:"length=3"`
	ContactName                     string `cwr:"length=30"`
	ContactID                       string `cwr:"length=10"`
	CWRWorkType                     string `cwr:"length=2"`
	GrandRightsInd                  string `cwr:"length=1"`
	CompositeComponentCount         string `cwr:"length=3"`
	PrintedEditionPublicationDate   string `cwr:"length=8"`
	ExceptionalClause               string `cwr:"length=1"`
	OpusNumber                      string `cwr:"length=25"`
	CatalogNumber                   string `cwr:"length=25"`
	PriorityFlag                    string `cwr:"length=1"`
}

// ACK is a CWR Acknowledgement of Transaction record
type ACK struct {
	TransactionSequenceNumber         string `cwr:"length=8"`
	RecordSequenceNumber              string `cwr:"length=8"`
	CreationDate                      string `cwr:"length=8"`
	CreationTime                      string `cwr:"length=6"`
	OriginalGroupID                   string `cwr:"length=5"`
	OriginalTransactionSequenceNumber string `cwr:"length=8"`
	OriginalTransactionType           string `cwr:"length=3"`
	CreationTitle                     string `cwr:"length=60"`
	SubmitterCreationNumber           string `cwr:"length=20"`
	RecipientCreationNumber           string `cwr:"length=20"`
	ProcessingDate                    string `cwr:"length=8"`
	TransactionStatus                 string `cwr:"length=2"`
}

// TER is a CWR Territory in Agreement record
type TER struct {
	TransactionSequenceNumber   string `cwr:"length=8"`
	RecordSequenceNumber        string `cwr:"length=8"`
	InclusionExclusionIndicator string `cwr:"length=1"`
	TISNumericCode              string `cwr:"length=4"`
}

// IPA is a CWR Interested Party of Agreement record
type IPA struct {
	TransactionSequenceNumber      string `cwr:"length=8"`
	RecordSequenceNumber           string `cwr:"length=8"`
	AgreementRoleCode              string `cwr:"length=2"`
	IPINameNumber                  string `cwr:"length=11"`
	IPIBaseNumber                  string `cwr:"length=13"`
	InterestedPartyNumber          string `cwr:"length=9"`
	InterestedPartyLastName        string `cwr:"length=45"`
	InterestedPartyWriterFirstName string `cwr:"length=30"`
	PRAffiliationSociety           string `cwr:"length=3"`
	PRShare                        string `cwr:"length=5"`
	MRAffiliationSociety           string `cwr:"length=3"`
	MRShare                        string `cwr:"length=5"`
	SRAffiliationSociety           string `cwr:"length=3"`
	SRShare                        string `cwr:"length=5"`
}

// NPA is a CWR Non-Roman Alphabet Agreement Party Name record
type NPA struct {
	TransactionSequenceNumber      string `cwr:"length=8"`
	RecordSequenceNumber           string `cwr:"length=8"`
	InterestedPartyNumber          string `cwr:"length=9"`
	InterestedPartyLastName        string `cwr:"length=160"`
	InterestedPartyWriterFirstName string `cwr:"length=160"`
	LanguageCode                   string `cwr:"length=2"`
}

// SPU is a CWR Publisher Controlled By Submitter record
type SPU struct {
	TransactionSequenceNumber          string `cwr:"length=8"`
	RecordSequenceNumber               string `cwr:"length=8"`
	PublisherSequenceNumber            string `cwr:"length=2"`
	InterestedPartyNumber              string `cwr:"length=9"`
	PublisherName                      string `cwr:"length=45"`
	PublisherUnknownIndicator          string `cwr:"length=1"`
	PublisherType                      string `cwr:"length=2"`
	TaxIDNumber                        string `cwr:"length=9"`
	PublisherIPINameNumber             string `cwr:"length=11"`
	SubmitterAgreementNumber           string `cwr:"length=14"`
	PRAffiliationSociety               string `cwr:"length=3"`
	PRShare                            string `cwr:"length=5"`
	MRAffiliationSociety               string `cwr:"length=3"`
	MRShare                            string `cwr:"length=5"`
	SRAffiliationSociety               string `cwr:"length=3"`
	SRShare                            string `cwr:"length=5"`
	SpecialAgreementsIndicator         string `cwr:"length=1"`
	FirstRecordingRefusalIndicator     string `cwr:"length=1"`
	Filler                             string `cwr:"length=1"`
	PublisherIPIBaseNumber             string `cwr:"length=13"`
	InternationalStandardAgreementCode string `cwr:"length=14"`
	SocietyAssignedAgreementNumber     string `cwr:"length=14"`
	AgreementType                      string `cwr:"length=2"`
	USALicenseIndicator                string `cwr:"length=1"`
}

// OPU is a CWR Other Publisher record
type OPU struct {
	TransactionSequenceNumber          string `cwr:"length=8"`
	RecordSequenceNumber               string `cwr:"length=8"`
	PublisherSequenceNumber            string `cwr:"length=2"`
	InterestedPartyNumber              string `cwr:"length=9"`
	PublisherName                      string `cwr:"length=45"`
	PublisherUnknownIndicator          string `cwr:"length=1"`
	PublisherType                      string `cwr:"length=2"`
	TaxIDNumber                        string `cwr:"length=9"`
	PublisherIPINameNumber             string `cwr:"length=11"`
	SubmitterAgreementNumber           string `cwr:"length=14"`
	PRAffiliationSociety               string `cwr:"length=3"`
	PRShare                            string `cwr:"length=5"`
	MRAffiliationSociety               string `cwr:"length=3"`
	MRShare                            string `cwr:"length=5"`
	SRAffiliationSociety               string `cwr:"length=3"`
	SRShare                            string `cwr:"length=5"`
	SpecialAgreementsIndicator         string `cwr:"length=1"`
	FirstRecordingRefusalIndicator     string `cwr:"length=1"`
	Filler                             string `cwr:"length=1"`
	PublisherIPIBaseNumber             string `cwr:"length=13"`
	InternationalStandardAgreementCode string `cwr:"length=14"`
	SocietyAssignedAgreementNumber     string `cwr:"length=14"`
	AgreementType                      string `cwr:"length=2"`
	USALicenseIndicator                string `cwr:"length=1"`
}

// NPN is a CWR Non-Roman Alphabet Publisher Name record
type NPN struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	PublisherSequenceNumber   string `cwr:"length=2"`
	InterestedPartyNumber     string `cwr:"length=9"`
	PublisherName             string `cwr:"length=480"`
	LanguageCode              string `cwr:"length=2"`
}

// SPT is a CWR Publisher Territory of Control record
type SPT struct {
	TransactionSequenceNumber   string `cwr:"length=8"`
	RecordSequenceNumber        string `cwr:"length=8"`
	InterestedPartyNumber       string `cwr:"length=9"`
	Constant                    string `cwr:"length=6"`
	PRCollectionShare           string `cwr:"length=5"`
	MRCollectionShare           string `cwr:"length=5"`
	SRCollectionShare           string `cwr:"length=5"`
	InclusionExclusionIndicator string `cwr:"length=1"`
	TISNumericCode              string `cwr:"length=4"`
	SharesChange                string `cwr:"length=1"`
	SequenceNumber              string `cwr:"length=3"`
}

// SWR is a CWR Writer Controlled By Submitter record
type SWR struct {
	TransactionSequenceNumber      string `cwr:"length=8"`
	RecordSequenceNumber           string `cwr:"length=8"`
	InterestedPartyNumber          string `cwr:"length=9"`
	WriterLastName                 string `cwr:"length=45"`
	WriterFirstName                string `cwr:"length=30"`
	WriterUnknownIndicator         string `cwr:"length=1"`
	WriterDesignationCode          string `cwr:"length=2"`
	TaxIDNumber                    string `cwr:"length=9"`
	WriterIPINameNumber            string `cwr:"length=11"`
	PRAffiliationSociety           string `cwr:"length=3"`
	PRShare                        string `cwr:"length=5"`
	MRAffiliationSociety           string `cwr:"length=3"`
	MRShare                        string `cwr:"length=5"`
	SRAffiliationSociety           string `cwr:"length=3"`
	SRShare                        string `cwr:"length=5"`
	ReversionaryIndicator          string `cwr:"length=1"`
	FirstRecordingRefusalIndicator string `cwr:"length=1"`
	WorkForHireIndicator           string `cwr:"length=1"`
	Filler                         string `cwr:"length=1"`
	WriterIPIBaseNumber            string `cwr:"length=13"`
	PersonalNumber                 string `cwr:"length=12"`
	USALicenseIndicator            string `cwr:"length=1"`
}

// OWR is a CWR Other Writer record
type OWR struct {
	TransactionSequenceNumber      string `cwr:"length=8"`
	RecordSequenceNumber           string `cwr:"length=8"`
	InterestedPartyNumber          string `cwr:"length=9"`
	WriterLastName                 string `cwr:"length=45"`
	WriterFirstName                string `cwr:"length=30"`
	WriterUnknownIndicator         string `cwr:"length=1"`
	WriterDesignationCode          string `cwr:"length=2"`
	TaxIDNumber                    string `cwr:"length=9"`
	WriterIPINameNumber            string `cwr:"length=11"`
	PRAffiliationSociety           string `cwr:"length=3"`
	PRShare                        string `cwr:"length=5"`
	MRAffiliationSociety           string `cwr:"length=3"`
	MRShare                        string `cwr:"length=5"`
	SRAffiliationSociety           string `cwr:"length=3"`
	SRShare                        string `cwr:"length=5"`
	ReversionaryIndicator          string `cwr:"length=1"`
	FirstRecordingRefusalIndicator string `cwr:"length=1"`
	WorkForHireIndicator           string `cwr:"length=1"`
	Filler                         string `cwr:"length=1"`
	WriterIPIBaseNumber            string `cwr:"length=13"`
	PersonalNumber                 string `cwr:"length=12"`
	USALicenseIndicator            string `cwr:"length=1"`
}

// NWN is a CWR Non-Roman Alphabet Writer Name record
type NWN struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	InterestedPartyNumber     string `cwr:"length=9"`
	WriterLastName            string `cwr:"length=160"`
	WriterFirstName           string `cwr:"length=160"`
	LanguageCode              string `cwr:"length=2"`
}

// SWT is a CWR Writer Territory of Control record
type SWT struct {
	TransactionSequenceNumber   string `cwr:"length=8"`
	RecordSequenceNumber        string `cwr:"length=8"`
	InterestedPartyNumber       string `cwr:"length=9"`
	PRCollectionShare           string `cwr:"length=5"`
	MRCollectionShare           string `cwr:"length=5"`
	SRCollectionShare           string `cwr:"length=5"`
	InclusionExclusionIndicator string `cwr:"length=1"`
	TISNumericCode              string `cwr:"length=4"`
	SharesChange                string `cwr:"length=1"`
	SequenceNumber              string `cwr:"length=3"`
}

// PWR is a CWR Publisher For Writer record
type PWR struct {
	TransactionSequenceNumber      string `cwr:"length=8"`
	RecordSequenceNumber           string `cwr:"length=8"`
	PublisherInterestedPartyNumber string `cwr:"length=9"`
	PublisherName                  string `cwr:"length=45"`
	SubmitterAgreementNumber       string `cwr:"length=14"`
	SocietyAssignedAgreementNumber string `cwr:"length=14"`
	WriterInterestedPartyNumber    string `cwr:"length=9"`
}

// ALT is a CWR Alternate Title record
type ALT struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	AlternateTitle            string `cwr:"length=60"`
	TitleType                 string `cwr:"length=2"`
	LanguageCode              string `cwr:"length=2"`
}

// NAT is a CWR Non-Roman Alphabet Title record
type NAT struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	Title                     string `cwr:"length=640"`
	TitleType                 string `cwr:"length=2"`
	LanguageCode              string `cwr:"length=2"`
}

// EWT is a CWR Entire Work Title for Excerpts record
type EWT struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	EntireWorkTitle           string `cwr:"length=60"`
	EntireWorkISWC            string `cwr:"length=11"`
	LanguageCode              string `cwr:"length=2"`
	Writer1LastName           string `cwr:"length=45"`
	Writer1FirstName          string `cwr:"length=30"`
	Source                    string `cwr:"length=60"`
	Writer1IPINameNumber      string `cwr:"length=11"`
	Writer1IPIBaseNumber      string `cwr:"length=13"`
	Writer2LastName           string `cwr:"length=45"`
	Writer2FirstName          string `cwr:"length=30"`
	Writer2IPINameNumber      string `cwr:"length=11"`
	Writer2IPIBaseNumber      string `cwr:"length=13"`
	SubmitterWorkNumber       string `cwr:"length=14"`
}

// VER is a CWR Original Work Title for Versions record
type VER struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	OriginalWorkTile          string `cwr:"length=60"`
	OriginalWorkISWC          string `cwr:"length=11"`
	LanguageCode              string `cwr:"length=2"`
	Writer1LastName           string `cwr:"length=45"`
	Writer1FirstName          string `cwr:"length=30"`
	Source                    string `cwr:"length=60"`
	Writer1IPINameNumber      string `cwr:"length=11"`
	Writer1IPIBaseNumber      string `cwr:"length=13"`
	Writer2LastName           string `cwr:"length=45"`
	Writer2FirstName          string `cwr:"length=30"`
	Writer2IPINameNumber      string `cwr:"length=11"`
	Writer2IPIBaseNumber      string `cwr:"length=13"`
	SubmitterWorkNumber       string `cwr:"length=14"`
}

// PER is a CWR Performing Artist record
type PER struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	ArtistLastName            string `cwr:"length=45"`
	ArtistFirstName           string `cwr:"length=30"`
	ArtistIPINameNumber       string `cwr:"length=11"`
	ArtistIPIBaseNumber       string `cwr:"length=13"`
}

// NPR is a CWR Performance Data in non-roman alphabet record
type NPR struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	ArtistName                string `cwr:"length=160"`
	ArtistFirstName           string `cwr:"length=160"`
	ArtistIPINameNumber       string `cwr:"length=11"`
	ArtistIPIBaseNumber       string `cwr:"length=13"`
	LanguageCode              string `cwr:"length=2"`
	PerformanceLanguage       string `cwr:"length=2"`
	PerformanceDialect        string `cwr:"length=3"`
}

// REC is a CWR Recording Detail record
type REC struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	FirstReleaseDate          string `cwr:"length=8"`
	Constant1                 string `cwr:"length=60"`
	FirstReleaseDuration      string `cwr:"length=6"`
	Constant2                 string `cwr:"length=5"`
	FirstAlbumTitle           string `cwr:"length=60"`
	FirstAlbumLabel           string `cwr:"length=60"`
	FirstReleaseCatalogNumber string `cwr:"length=18"`
	EAN                       string `cwr:"length=13"`
	ISRC                      string `cwr:"length=12"`
	RecordingFormat           string `cwr:"length=1"`
	RecordingTechnique        string `cwr:"length=1"`
	MediaType                 string `cwr:"length=1"`
}

// ORN is a CWR Work Origin record
type ORN struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	IntendedPurpose           string `cwr:"length=3"`
	ProductionTitle           string `cwr:"length=60"`
	CDIdentifier              string `cwr:"length=15"`
	CutNumber                 string `cwr:"length=4"`
	Library                   string `cwr:"length=60"`
	BLTVR                     string `cwr:"length=1"`
	VISAN                     string `cwr:"length=8"`
	ISAN                      string `cwr:"length=12"`
	Episode                   string `cwr:"length=4"`
	CheckDigit                string `cwr:"length=1"`
	ProductionNumber          string `cwr:"length=12"`
	EpisodeTitle              string `cwr:"length=60"`
	EpisodeNumber             string `cwr:"length=20"`
	ProductionYear            string `cwr:"length=4"`
	AVISocietyCode            string `cwr:"length=3"`
	AudioVisualNumber         string `cwr:"length=15"`
}

// INS is a CWR Instrumentation Summary record
type INS struct {
	TransactionSequenceNumber   string `cwr:"length=8"`
	RecordSequenceNumber        string `cwr:"length=8"`
	NumberOfVoices              string `cwr:"length=3"`
	StandardInstrumentationType string `cwr:"length=3"`
	InstrumentationDescription  string `cwr:"length=50"`
}

// IND is a CWR Instrumentation Detail record
type IND struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	InstrumentCode            string `cwr:"length=3"`
	NumberOfPlayers           string `cwr:"length=3"`
}

// COM is a CWR Component record
type COM struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	Title                     string `cwr:"length=60"`
	ISWC                      string `cwr:"length=11"`
	SubmitterWorkNumber       string `cwr:"length=14"`
	Duration                  string `cwr:"length=6"`
	Writer1LastName           string `cwr:"length=45"`
	Writer1FirstName          string `cwr:"length=30"`
	Writer1IPINameNumber      string `cwr:"length=11"`
	Writer2LastName           string `cwr:"length=45"`
	Writer2FirstName          string `cwr:"length=30"`
	Writer2IPINameNumber      string `cwr:"length=11"`
	Writer2IPIBaseNumber      string `cwr:"length=13"`
	Writer1IPIBaseNumber      string `cwr:"length=13"`
}

// MSG is a CWR Message record
type MSG struct {
	TransactionSequenceNumber    string `cwr:"length=8"`
	RecordSequenceNumber         string `cwr:"length=8"`
	MessageType                  string `cwr:"length=1"`
	OriginalRecordSequenceNumber string `cwr:"length=8"`
	RecordType                   string `cwr:"length=3"`
	MessageLevel                 string `cwr:"length=1"`
	ValidationNumber             string `cwr:"length=3"`
	MessageText                  string `cwr:"length=150"`
}

// NVT is a CWR Non-Roman Alphabet Original Title for Version record
type NVT struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	Title                     string `cwr:"length=640"`
	LanguageCode              string `cwr:"length=2"`
}

// NET is a CWR Non-Roman Alphabet Entire Work Title for Excerpts record
type NET struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	Title                     string `cwr:"length=640"`
	LanguageCode              string `cwr:"length=2"`
}

// NCT is a CWR Non-Roman Alphabet Title for Components record
type NCT struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	Title                     string `cwr:"length=640"`
	LanguageCode              string `cwr:"length=2"`
}

// NOW is a CWR Non-Roman Alphabet Other Writer Name record
type NOW struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	WriterName                string `cwr:"length=160"`
	WriterFirstName           string `cwr:"length=160"`
	LanguageCode              string `cwr:"length=2"`
	WriterPosition            string `cwr:"length=1"`
}

// ARI is a CWR Additional Related Information record
type ARI struct {
	TransactionSequenceNumber string `cwr:"length=8"`
	RecordSequenceNumber      string `cwr:"length=8"`
	SocietyNumber             string `cwr:"length=3"`
	WorkNumber                string `cwr:"length=14"`
	RightType                 string `cwr:"length=3"`
	SubjectCode               string `cwr:"length=2"`
	Note                      string `cwr:"length=160"`
}

var (
	cwrTypes   map[string]reflect.Type
	cwrLengths map[string][]int
)

func init() {
	cwrTypes = map[string]reflect.Type{
		"HDR": reflect.TypeOf(HDR{}),
		"GRH": reflect.TypeOf(GRH{}),
		"GRT": reflect.TypeOf(GRT{}),
		"TRL": reflect.TypeOf(TRL{}),
		"AGR": reflect.TypeOf(AGR{}),
		"NWR": reflect.TypeOf(NWR{}),
		"REV": reflect.TypeOf(REV{}),
		"ISW": reflect.TypeOf(ISW{}),
		"EXC": reflect.TypeOf(EXC{}),
		"ACK": reflect.TypeOf(ACK{}),
		"TER": reflect.TypeOf(TER{}),
		"IPA": reflect.TypeOf(IPA{}),
		"NPA": reflect.TypeOf(NPA{}),
		"SPU": reflect.TypeOf(SPU{}),
		"OPU": reflect.TypeOf(OPU{}),
		"NPN": reflect.TypeOf(NPN{}),
		"SPT": reflect.TypeOf(SPT{}),
		"SWR": reflect.TypeOf(SWR{}),
		"OWR": reflect.TypeOf(OWR{}),
		"NWN": reflect.TypeOf(NWN{}),
		"SWT": reflect.TypeOf(SWT{}),
		"PWR": reflect.TypeOf(PWR{}),
		"ALT": reflect.TypeOf(ALT{}),
		"NAT": reflect.TypeOf(NAT{}),
		"EWT": reflect.TypeOf(EWT{}),
		"VER": reflect.TypeOf(VER{}),
		"PER": reflect.TypeOf(PER{}),
		"NPR": reflect.TypeOf(NPR{}),
		"REC": reflect.TypeOf(REC{}),
		"ORN": reflect.TypeOf(ORN{}),
		"INS": reflect.TypeOf(INS{}),
		"IND": reflect.TypeOf(IND{}),
		"COM": reflect.TypeOf(COM{}),
		"MSG": reflect.TypeOf(MSG{}),
		"NVT": reflect.TypeOf(NVT{}),
		"NET": reflect.TypeOf(NET{}),
		"NCT": reflect.TypeOf(NCT{}),
		"NOW": reflect.TypeOf(NOW{}),
		"ARI": reflect.TypeOf(ARI{}),
	}
	cwrLengths = make(map[string][]int)
	for name, typ := range cwrTypes {
		cwrLengths[name] = make([]int, typ.NumField())
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			tag := f.Tag.Get("cwr")
			tagVal := strings.SplitN(tag, "=", 2)
			if len(tagVal) != 2 {
				panic(fmt.Sprintf("cwr: invalid tag on %s.%s: %q", typ.Name(), f.Name, tag))
			}
			length, err := strconv.Atoi(tagVal[1])
			if err != nil {
				panic(fmt.Sprintf("cwr: invalid tag on %s.%s: %q (%s)", typ.Name(), f.Name, tag, err))
			}
			cwrLengths[name][i] = length
		}
	}
}
