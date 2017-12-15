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

type Variables map[string]interface{}

type Performer struct {
	Name string `json:"name"`
}

type Contributor struct {
	Name string `json:"name"`
}

type Composer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RecordLabel struct {
	Name string `json:"name"`
}

type Publisher struct {
	Name string `json:"name"`
}

type Recording struct {
	Title    string `json:"title"`
	Duration string `json:"duration"`
}

type Work struct {
	Title string `json:"title"`
}

type Song struct {
	Title    string `json:"title"`
	Duration string `json:"duration"`
}

type Release struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Date  string `json:"date"`
}

type PerformerRecord struct {
	ID     int64
	Name   string
	Source int64
}

type ContributorRecord struct {
	ID     int64
	Name   string
	Source int64
}

type ComposerRecord struct {
	ID        int64
	FirstName string
	LastName  string
	Source    int64
}

type RecordLabelRecord struct {
	ID     int64
	Name   string
	Source int64
}

type PublisherRecord struct {
	ID     int64
	Name   string
	Source int64
}

type RecordingRecord struct {
	ID       int64
	Title    string
	Duration string
	Source   int64
}

type WorkRecord struct {
	ID     int64
	Title  string
	Source int64
}

type SongRecord struct {
	ID       int64
	Title    string
	Duration string
	Source   int64
}

type ReleaseRecord struct {
	ID     int64
	Type   string
	Title  string
	Date   string
	Source int64
}

type Identifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type IdentifierRecord struct {
	ID int64

	Identifier
}

type PerformerRecordingLink struct {
	Performer Identifier `json:"performer"`
	Recording Identifier `json:"recording"`
	Role      string     `json:"role"`
}

type PerformerSongLink struct {
	Performer Identifier `json:"performer"`
	Song      Identifier `json:"song"`
	Role      string     `json:"role"`
}

type PerformerReleaseLink struct {
	Performer Identifier `json:"performer"`
	Release   Identifier `json:"release"`
	Role      string     `json:"role"`
}

type ContributorRecordingLink struct {
	Contributor Identifier `json:"contributor"`
	Recording   Identifier `json:"recording"`
	Role        string     `json:"role"`
}

type ComposerWorkLink struct {
	Composer Identifier `json:"composer"`
	Work     Identifier `json:"work"`
	Role     string     `json:"role"`
}

type RecordLabelRecordingLink struct {
	RecordLabel Identifier `json:"record_label"`
	Recording   Identifier `json:"recording"`
}

type RecordLabelSongLink struct {
	RecordLabel Identifier `json:"record_label"`
	Song        Identifier `json:"song"`
}

type RecordLabelReleaseLink struct {
	RecordLabel Identifier `json:"record_label"`
	Release     Identifier `json:"release"`
}

type PublisherWorkLink struct {
	Publisher Identifier `json:"publisher"`
	Work      Identifier `json:"work"`
}

type SongRecordingLink struct {
	Song      Identifier `json:"song"`
	Recording Identifier `json:"recording"`
}

type ReleaseRecordingLink struct {
	Release   Identifier `json:"release"`
	Recording Identifier `json:"recording"`
}

type RecordingWorkLink struct {
	Recording Identifier `json:"recording"`
	Work      Identifier `json:"work"`
}

type ReleaseSongLink struct {
	Release Identifier `json:"release"`
	Song    Identifier `json:"song"`
}

type PerformerRecordingRecord struct {
	ID        int64
	Performer *IdentifierRecord
	Recording *IdentifierRecord
	Role      string
	Source    int64
}

type PerformerSongRecord struct {
	ID        int64
	Performer *IdentifierRecord
	Song      *IdentifierRecord
	Role      string
	Source    int64
}

type PerformerReleaseRecord struct {
	ID        int64
	Performer *IdentifierRecord
	Release   *IdentifierRecord
	Role      string
	Source    int64
}

type ContributorRecordingRecord struct {
	ID          int64
	Contributor *IdentifierRecord
	Recording   *IdentifierRecord
	Role        string
	Source      int64
}

type ComposerWorkRecord struct {
	ID       int64
	Composer *IdentifierRecord
	Work     *IdentifierRecord
	Role     string
	Source   int64
}

type RecordLabelRecordingRecord struct {
	ID          int64
	RecordLabel *IdentifierRecord
	Recording   *IdentifierRecord
	Source      int64
}

type RecordLabelSongRecord struct {
	ID          int64
	RecordLabel *IdentifierRecord
	Song        *IdentifierRecord
	Source      int64
}

type RecordLabelReleaseRecord struct {
	ID          int64
	RecordLabel *IdentifierRecord
	Release     *IdentifierRecord
	Source      int64
}

type PublisherWorkRecord struct {
	ID        int64
	Publisher *IdentifierRecord
	Work      *IdentifierRecord
	Source    int64
}

type SongRecordingRecord struct {
	ID        int64
	Song      *IdentifierRecord
	Recording *IdentifierRecord
	Source    int64
}

type ReleaseRecordingRecord struct {
	ID        int64
	Release   *IdentifierRecord
	Recording *IdentifierRecord
	Source    int64
}

type RecordingWorkRecord struct {
	ID        int64
	Recording *IdentifierRecord
	Work      *IdentifierRecord
	Source    int64
}

type ReleaseSongRecord struct {
	ID      int64
	Release *IdentifierRecord
	Song    *IdentifierRecord
	Source  int64
}

type Source struct {
	Name string `json:"name"`
}

type SourceRecord struct {
	ID   int64
	Name string
}
