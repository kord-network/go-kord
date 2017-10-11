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

// Context represents a JSON-LD context.
type Context map[string]string

// Party combines the PartyName and PartyID DDEX complex types.
type Party struct {
	Context   Context `json:"@context"`
	PartyId   string  `json:"partyId, omitempty"`
	PartyName string  `json:"fullName, omitempty"`
}

// SoundRecording is the data object for the ERN SoundRecording composite
type SoundRecording struct {
	Context             Context `json:"@context"`
	ArtistName          string  `json:"fullName, omitempty"`
	GenreText           string  `json:"genreText, omitempty"`
	ParentalWarningType string  `json:"parentalWarningType, omitempty"`
	ReferenceTitle      string  `json:"titleText, omitempty"`
	ResourceReference   string  `json:"resourceReference, omitempty"`
	SoundRecordingId    string  `json:"soundRecordingId, omitempty"`
	SubGenre            string  `json:"subGenre, omitempty"`
	TerritoryCode       string  `json:"territoryCode, omitempty"`
}

// Release is the data object for the ERN Release composite
type Release struct {
	Context      Context `json:"@context"`
	ArtistName   string  `json:"fullName, omitempty"`
	DisplayTitle string  `json:"displayTitle, omitempty"`
	Genre        string  `json:"genre, omitempty"`
	ReleaseType  string  `json:"releaseType, omitempty"`
	ReleaseId    string  `json:"releaseId, omitempty"`
}
