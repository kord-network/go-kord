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

package musicbrainz

// Artist represents a MusicBrainz artist, see
// https://musicbrainz.org/doc/Artist
type Artist struct {
	ID                    int64    `json:"id,omitempty"`
	Name                  string   `json:"name,omitempty"`
	SortName              string   `json:"sortName,omitempty"`
	Type                  string   `json:"type,omitempty"`
	Gender                string   `json:"gender,omitempty"`
	Area                  string   `json:"area,omitempty"`
	BeginDate             string   `json:"begin_date,omitempty"`
	EndDate               string   `json:"end_date,omitempty"`
	IPI                   []string `json:"ipi,omitempty"`
	ISNI                  []string `json:"isni,omitempty"`
	Alias                 []string `json:"alias,omitempty"`
	MBID                  string   `json:"mbid,omitempty"`
	DisambiguationComment string   `json:"disambiguation_comment,omitempty"`
	Annotation            []string `json:"annotation,omitempty"`
}

// RecordingWorkLink is a link between an ISRC and ISWC which comes from the
// l_recording_work MusicBrainz table.
type RecordingWorkLink struct {
	RecordingTitle string `json:"recording_title"`
	WorkTitle      string `json:"work_title"`
	ISRC           string `json:"isrc"`
	ISWC           string `json:"iswc"`
}
