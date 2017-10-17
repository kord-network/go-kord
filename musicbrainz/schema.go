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

import "github.com/meta-network/go-meta/migrate"

// migrations is a set of database migrations to run on a SQLite3 database
// to prepare it for indexing a META stream of MusicBrainz artists.
var migrations = migrate.NewMigrations()

func init() {
	// migration 1 creates indexes for the following artist properties:
	//
	// * Name - https://musicbrainz.org/doc/Artist#Name
	// * Type - https://musicbrainz.org/doc/Artist#Type
	// * MBID - https://musicbrainz.org/doc/Artist#MBID
	// * IPI  - https://musicbrainz.org/doc/Artist#IPI_code
	// * ISNI - https://musicbrainz.org/doc/Artist#ISNI_code
	//
	migrations.Add(1, `
CREATE TABLE artist (
	object_id text NOT NULL,
	name      text NOT NULL,
	type      text NOT NULL,
	mbid      text NOT NULL
);

CREATE INDEX artist_object_id_idx ON artist (object_id);
CREATE INDEX artist_name_idx      ON artist (name);
CREATE INDEX artist_type_idx      ON artist (type);
CREATE INDEX artist_mbid_idx      ON artist (mbid);

CREATE TABLE artist_ipi (
	object_id text NOT NULL,
	ipi       text NOT NULL
);

CREATE INDEX artist_ipi_idx ON artist_ipi (ipi);

CREATE TABLE artist_isni (
	object_id text NOT NULL,
	isni      text NOT NULL
);

CREATE INDEX artist_isni_idx ON artist_isni (isni);
`,
	)

	// migration 2 creates an index for ISRC to ISWC links which come
	// from the l_recording_work MusicBrainz table.
	migrations.Add(2, `
CREATE TABLE recording_work (
	object_id text NOT NULL,
	isrc      text NOT NULL,
	iswc      text NOT NULL
);

CREATE INDEX recording_work_object_id_idx ON recording_work (object_id);
CREATE UNIQUE INDEX recording_work_idx ON recording_work (isrc, iswc);
`,
	)
}
