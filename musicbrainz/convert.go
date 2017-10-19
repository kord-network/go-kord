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

	"github.com/meta-network/go-meta"
)

// Converter converts MusicBrainz data stored in a PostgreSQL database to META
// objects.
type Converter struct {
	db    *sql.DB
	store *meta.Store
}

// NewConverter returns a Converter which reads data from the given PostgreSQL
// database connection and stores META objects in the given META store.
func NewConverter(db *sql.DB, store *meta.Store) *Converter {
	return &Converter{
		db:    db,
		store: store,
	}
}

// ConvertArtists reads all artists from the database, converts them to META
// objects, stores them in the META store and sends their CIDs to the given
// stream.
func (c *Converter) ConvertArtists(ctx context.Context, stream *meta.StreamWriter, source string) error {
	// get all artists from the db
	rows, err := c.db.QueryContext(ctx, artistsQuery)
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
		a.Context = ArtistContext
		a.Source = source

		// convert the artist to a META object
		obj, err := c.store.Put(a)
		if err != nil {
			return err
		}

		// send the object's CID to the output stream
		if err := stream.Write(obj.Cid()); err != nil {
			return err
		}
	}

	return rows.Err()
}

// ConvertRecordingWorkLinks loads ISRC to ISWC links from the database,
// converts them to META objects and sends their CIDs to the given stream.
func (c *Converter) ConvertRecordingWorkLinks(ctx context.Context, stream *meta.StreamWriter, source string) error {
	// get all links from the db
	rows, err := c.db.QueryContext(ctx, recordingWorkLinksQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var l RecordingWorkLink
		if err := rows.Scan(&l.ISRC, &l.ISWC); err != nil {
			return err
		}
		l.Source = source

		// convert the link to a META object
		obj, err := c.store.Put(l)
		if err != nil {
			return err
		}

		// send the object's CID to the output stream
		if err := stream.Write(obj.Cid()); err != nil {
			return err
		}
	}

	return rows.Err()
}
