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

// RegisteredWork represents a CWR work registratin , see
// see http://musicmark.com/documents/cwr11-1494_cwr_user_manual_2011-09-23_e_2011-09-23_en.pdf
// NWR or REV record
type RegisteredWork struct {
	RecordType              string `json:"record_type,omitempty"`
	TransactionSequenceN    string `json:"transactionSequenceN,omitempty"`
	RecordSequenceN         string `json:"recordSequenceN,omitempty"`
	Title                   string `json:"title,omitempty"`
	LanguageCode            string `json:"languageCode,omitempty"`
	SubmitteWorkNumber      string `json:"submitterWorkNumber,omitempty"`
	ISWC                    string `json:"iswc,omitempty"`
	CopyRightDate           string `json:"copyRightDate,omitempty"`
	DistributionCategory    string `json:"distributionCategory,omitempty"`
	Duration                string `json:"duration,omitempty"`
	RecordedIndicator       string `json:"recordedIndicator,omitempty"`
	TextMusicRelationship   string `json:"textMusicRelationship,omitempty"`
	CompositeType           string `json:"composite_type,omitempty"`
	VersionType             string `json:"versionType,omitempty"`
	MusicArrangement        string `json:"musicArrangement,omitempty"`
	LyricAdaptation         string `json:"lyricAdaptation,omitempty"`
	ContactName             string `json:"contactName,omitempty"`
	ContactID               string `json:"contactId,omitempty"`
	WorkType                string `json:"workType,omitempty"`
	GrandRightsIndicator    string `json:"grandRightsIndicator,omitempty"`
	CompositeComponentCount string `json:"compositeComponentCount,omitempty"`
	DateOfPublication       string `json:"dateOfPublication,omitempty"`
	ExceptionalClause       string `json:"exceptionalClause,omitempty"`
	OpusNumber              string `json:"opusNumber,omitempty"`
	CatalogueNumber         string `json:"catalogueNumber,omitempty"`
	PriorityFlag            string `json:"priorityFlag,omitempty"`
}

// TransmissionHeader Record - HDR
type TransmissionHeader struct {
	RecordType string `json:"record_type,omitempty"`
	SenderType string `json:"sender_type,omitempty"`
	SenderID   string `json:"sender_id,omitempty"`
	SenderName string `json:"sender_name,omitempty"`
}

// GroupTrailer Record - GRT
type GroupTrailer struct {
	RecordType string `json:"record_type,omitempty"`
	GroupID    string `json:"group_id,omitempty"`
}

// GroupHeader Record - GRH
type GroupHeader struct {
	RecordType      string `json:"record_type,omitempty"`
	TransactionType string `json:"transaction_type,omitempty"`
	GroupID         string `json:"group_id,omitempty"`
}

// TransmissionTrailer Record - TRL
type TransmissionTrailer struct {
	RecordType string `json:"record_type,omitempty"`
}

// PublisherControllBySubmitter Record - SPU
type PublisherControllBySubmitter struct {
	RecordType              string `json:"record_type,omitempty"`
	TransactionSequenceN    string `json:"transactionSequenceN,omitempty"`
	RecordSequenceN         string `json:"recordSequenceN,omitempty"`
	PublisherSequenceNumber string `json:"publisher_sequence_n,omitempty"`
}

// Record - include all CWR records fields
type Record struct {
	RecordType              string `json:"record_type,omitempty"`
	TransactionSequenceN    string `json:"transactionSequenceN,omitempty"`
	RecordSequenceN         string `json:"recordSequenceN,omitempty"`
	Title                   string `json:"title,omitempty"`
	LanguageCode            string `json:"languageCode,omitempty"`
	SubmitteWorkNumber      string `json:"submitterWorkNumber,omitempty"`
	ISWC                    string `json:"iswc,omitempty"`
	CopyRightDate           string `json:"copyRightDate,omitempty"`
	DistributionCategory    string `json:"distributionCategory,omitempty"`
	Duration                string `json:"duration,omitempty"`
	RecordedIndicator       string `json:"recordedIndicator,omitempty"`
	TextMusicRelationship   string `json:"textMusicRelationship,omitempty"`
	CompositeType           string `json:"composite_type,omitempty"`
	VersionType             string `json:"versionType,omitempty"`
	MusicArrangement        string `json:"musicArrangement,omitempty"`
	LyricAdaptation         string `json:"lyricAdaptation,omitempty"`
	ContactName             string `json:"contactName,omitempty"`
	ContactID               string `json:"contactId,omitempty"`
	WorkType                string `json:"workType,omitempty"`
	GrandRightsIndicator    string `json:"grandRightsIndicator,omitempty"`
	CompositeComponentCount string `json:"compositeComponentCount,omitempty"`
	DateOfPublication       string `json:"dateOfPublication,omitempty"`
	ExceptionalClause       string `json:"exceptionalClause,omitempty"`
	OpusNumber              string `json:"opusNumber,omitempty"`
	CatalogueNumber         string `json:"catalogueNumber,omitempty"`
	PriorityFlag            string `json:"priorityFlag,omitempty"`
	PublisherSequenceNumber string `json:"publisher_sequence_n,omitempty"`
	TransactionType         string `json:"transaction_type,omitempty"`
	GroupID                 string `json:"group_id,omitempty"`
	SenderType              string `json:"sender_type,omitempty"`
	SenderID                string `json:"sender_id,omitempty"`
	SenderName              string `json:"sender_name,omitempty"`
}
