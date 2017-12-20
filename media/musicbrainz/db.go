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

import (
	"context"
	"database/sql"
	"strings"
)

func LoadArtists(ctx context.Context, db *sql.DB, fn func(*Artist) error) error {
	rows, err := db.QueryContext(ctx, artistsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		// read the db row into an Artist struct, handling nullable
		// columns
		var (
			a          Artist
			typ        *string
			gender     *string
			area       *string
			beginDate  *string
			endDate    *string
			ipi        []byte
			isni       []byte
			alias      []byte
			annotation []byte
		)
		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.SortName,
			&typ,
			&gender,
			&area,
			&beginDate,
			&endDate,
			&ipi,
			&isni,
			&alias,
			&a.MBID,
			&a.DisambiguationComment,
			&annotation,
		)
		if err != nil {
			return err
		}
		if typ != nil {
			a.Type = *typ
		}
		if gender != nil {
			a.Gender = *gender
		}
		if area != nil {
			a.Area = *area
		}
		if beginDate != nil {
			a.BeginDate = *beginDate
		}
		if endDate != nil {
			a.EndDate = *endDate
		}
		if len(ipi) > 2 {
			a.IPI = strings.Split(string(ipi)[1:len(ipi)-1], ",")
		}
		if len(isni) > 2 {
			a.ISNI = strings.Split(string(isni)[1:len(isni)-1], ",")
		}
		if len(alias) > 2 {
			a.Alias = strings.Split(string(alias)[1:len(alias)-1], ",")
		}
		if len(annotation) > 2 {
			a.Annotation = strings.Split(string(annotation)[1:len(annotation)-1], ",")
		}

		if err := fn(&a); err != nil {
			return err
		}
	}

	return rows.Err()
}

func LoadRecordingWorkLinks(ctx context.Context, db *sql.DB, fn func(*RecordingWorkLink) error) error {
	rows, err := db.QueryContext(ctx, recordingWorkLinksQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var link RecordingWorkLink
		if err := rows.Scan(
			&link.RecordingTitle,
			&link.WorkTitle,
			&link.ISRC,
			&link.ISWC,
		); err != nil {
			return err
		}

		if err := fn(&link); err != nil {
			return err
		}
	}

	return rows.Err()
}

// artistsQuery is a SQL query to load artists from a MusicBrainz PostgreSQL
// database.
var artistsQuery = `
SELECT
  artist.id,
  artist.name,
  artist.sort_name,
  artist_type.name AS type,
  gender.name AS gender,
  area.name AS area,
  CASE
    WHEN artist.begin_date_year IS NULL OR artist.begin_date_year < 0 THEN NULL
    ELSE make_date(
      artist.begin_date_year,
      CASE
	WHEN artist.begin_date_month IS NULL THEN 1
	ELSE artist.begin_date_month
      END,
      CASE
	WHEN artist.begin_date_day IS NULL THEN 1
	ELSE artist.begin_date_day
      END
    )
  END AS begin_date,
  CASE
    WHEN artist.end_date_year IS NULL OR artist.end_date_year < 0 THEN NULL
    ELSE make_date(
      artist.end_date_year,
      CASE
	WHEN artist.end_date_month IS NULL THEN 1
	ELSE artist.end_date_month
      END,
      CASE
	WHEN artist.end_date_day IS NULL THEN 1
	ELSE artist.end_date_day
      END
    )
  END AS end_date,
  ARRAY(SELECT ipi FROM artist_ipi WHERE artist_ipi.artist = artist.id) AS ipi,
  ARRAY(SELECT isni FROM artist_isni WHERE artist_isni.artist = artist.id) AS isni,
  ARRAY(SELECT name FROM artist_alias WHERE artist_alias.artist = artist.id) AS alias,
  artist.gid AS mbid,
  artist.comment AS disambiguation_comment,
  ARRAY(SELECT text FROM annotation LEFT JOIN artist_annotation ON artist_annotation.annotation = annotation.id WHERE artist_annotation.artist = artist.id) AS annotation
FROM artist
LEFT JOIN artist_type ON artist.type = artist_type.id
LEFT JOIN gender ON artist.gender = gender.id
LEFT JOIN area ON artist.area = area.id
`[1:]

var recordingWorkLinksQuery = `
SELECT
  recording.name,
  work.name,
  isrc.isrc,
  iswc.iswc
FROM l_recording_work
INNER JOIN recording ON recording.id = l_recording_work.entity0
INNER JOIN work ON work.id = l_recording_work.entity1
INNER JOIN isrc ON isrc.recording = l_recording_work.entity0
INNER JOIN iswc ON iswc.work = l_recording_work.entity1
`[1:]
