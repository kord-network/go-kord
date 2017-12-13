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

import "github.com/meta-network/go-meta/migrate"

// migrations is a set of database migrations to run on a SQLite3 database
// to prepare it for indexing Media data.
var migrations = migrate.NewMigrations()

func init() {
	// migration 1 creates indexes and associations for:
	// * Performer
	// * Composer
	// * RecordLabel
	// * Publisher
	// * Recording
	// * Work
	// * Song
	// * Release
	migrations.Add(1, `
CREATE TABLE performer (
  id     INTEGER PRIMARY KEY,
  name   TEXT NOT NULL,
  source INTEGER NOT NULL REFERENCES source (id)
);

CREATE UNIQUE INDEX performer_unique_idx ON performer (name, source);

CREATE TABLE composer (
  id         INTEGER PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name  TEXT NOT NULL,
  source     INTEGER NOT NULL REFERENCES source (id)
);

CREATE UNIQUE INDEX composer_unique_idx ON composer (first_name, last_name, source);

CREATE TABLE record_label (
  id     INTEGER PRIMARY KEY,
  name   TEXT NOT NULL,
  source INTEGER NOT NULL REFERENCES source (id)
);

CREATE UNIQUE INDEX record_label_unique_idx ON record_label (name, source);

CREATE TABLE publisher (
  id     INTEGER PRIMARY KEY,
  name   TEXT NOT NULL,
  source INTEGER NOT NULL REFERENCES source (id)
);

CREATE UNIQUE INDEX publisher_unique_idx ON publisher (name, source);

CREATE TABLE recording (
  id       INTEGER PRIMARY KEY,
  title    TEXT NOT NULL,
  duration TEXT NOT NULL,
  source   INTEGER NOT NULL REFERENCES source (id)
);

CREATE UNIQUE INDEX recording_unique_idx ON recording (title, duration, source);

CREATE TABLE work (
  id     INTEGER PRIMARY KEY,
  title  TEXT NOT NULL,
  source INTEGER NOT NULL REFERENCES source (id)
);

CREATE UNIQUE INDEX work_unique_idx ON work (title, source);

CREATE TABLE song (
  id       INTEGER PRIMARY KEY,
  title    TEXT NOT NULL,
  duration TEXT NOT NULL,
  source   INTEGER NOT NULL REFERENCES source (id)
);

CREATE UNIQUE INDEX song_unique_idx ON song (title, duration, source);

CREATE TABLE release (
  id     INTEGER PRIMARY KEY,
  type   TEXT NOT NULL,
  title  TEXT NOT NULL,
  date   TEXT NOT NULL,
  source INTEGER NOT NULL REFERENCES source (id)
);

CREATE UNIQUE INDEX release_unique_idx ON release (type, title, date, source);

CREATE TABLE identifier (
  id    INTEGER PRIMARY KEY,
  type  TEXT NOT NULL CHECK (type != ""),
  value TEXT NOT NULL CHECK (value != "")
);

CREATE UNIQUE INDEX identifier_unique_idx ON identifier (type, value);

CREATE TABLE identifier_assignment (
  id            INTEGER PRIMARY KEY,
  identifier_id INTEGER REFERENCES identifier (id),
  record_type   TEXT NOT NULL,
  record_id     INTEGER NOT NULL,
  source        INTEGER REFERENCES source (id)
);

CREATE INDEX identifier_assignment_identifier_idx ON identifier_assignment (identifier_id);
CREATE INDEX identifier_assignment_record_idx     ON identifier_assignment (record_type, record_id);
CREATE UNIQUE INDEX identifier_assignment_unique_index ON identifier_assignment (identifier_id, record_type, record_id, source);

CREATE TABLE performer_recording (
  id                   INTEGER PRIMARY KEY,
  performer_identifier INTEGER REFERENCES identifier (id),
  recording_identifier INTEGER REFERENCES identifier (id),
  role                 TEXT,
  source               INTEGER REFERENCES source (id)
);

CREATE INDEX performer_recording_performer_idx ON performer_recording (performer_identifier);
CREATE INDEX performer_recording_recording_idx ON performer_recording (recording_identifier);
CREATE UNIQUE INDEX performer_recording_unique_idx ON performer_recording (performer_identifier, recording_identifier, role, source);

CREATE TABLE composer_work (
  id                  INTEGER PRIMARY KEY,
  composer_identifier INTEGER REFERENCES identifier (id),
  work_identifier     INTEGER REFERENCES identifier (id),
  role                TEXT,
  source              INTEGER REFERENCES source (id)
);

CREATE INDEX composer_work_composer_idx ON composer_work (composer_identifier);
CREATE INDEX composer_work_work_idx     ON composer_work (work_identifier);
CREATE UNIQUE INDEX composer_work_unique_idx ON composer_work (composer_identifier, work_identifier, role, source);

CREATE TABLE record_label_song (
  id                      INTEGER PRIMARY KEY,
  record_label_identifier INTEGER REFERENCES identifier (id),
  song_identifier         INTEGER REFERENCES identifier (id),
  source                  INTEGER REFERENCES source (id)
);

CREATE INDEX record_label_song_record_label_idx ON record_label_song (record_label_identifier);
CREATE INDEX record_label_song_song_idx         ON record_label_song (song_identifier);
CREATE UNIQUE INDEX record_label_song_unique_idx ON record_label_song (record_label_identifier, song_identifier, source);

CREATE TABLE record_label_release (
  id                      INTEGER PRIMARY KEY,
  record_label_identifier INTEGER REFERENCES identifier (id),
  release_identifier      INTEGER REFERENCES identifier (id),
  source                  INTEGER REFERENCES source (id)
);

CREATE INDEX record_label_release_record_label_idx ON record_label_release (record_label_identifier);
CREATE INDEX record_label_release_release_idx      ON record_label_release (release_identifier);
CREATE UNIQUE INDEX record_label_release_unique_idx ON record_label_release (record_label_identifier, release_identifier, source);

CREATE TABLE publisher_work (
  id                   INTEGER PRIMARY KEY,
  publisher_identifier INTEGER REFERENCES identifier (id),
  work_identifier      INTEGER REFERENCES identifier (id),
  source               INTEGER REFERENCES source (id)
);

CREATE INDEX publisher_work_publisher_idx ON publisher_work (publisher_identifier);
CREATE INDEX publisher_work_work_idx      ON publisher_work (work_identifier);
CREATE UNIQUE INDEX publisher_work_unique_idx ON publisher_work (publisher_identifier, work_identifier, source);

CREATE TABLE song_recording (
  id                   INTEGER PRIMARY KEY,
  song_identifier      INTEGER REFERENCES identifier (id),
  recording_identifier INTEGER REFERENCES identifier (id),
  source               INTEGER REFERENCES source (id)
);

CREATE INDEX song_recording_song_idx      ON song_recording (song_identifier);
CREATE INDEX song_recording_recording_idx ON song_recording (recording_identifier);
CREATE UNIQUE INDEX song_recording_unique_idx ON song_recording (song_identifier, recording_identifier, source);

CREATE TABLE release_recording (
  id                   INTEGER PRIMARY KEY,
  release_identifier   INTEGER REFERENCES identifier (id),
  recording_identifier INTEGER REFERENCES identifier (id),
  source               INTEGER REFERENCES source (id)
);

CREATE INDEX release_recording_release_idx   ON release_recording (release_identifier);
CREATE INDEX release_recording_recording_idx ON release_recording (recording_identifier);
CREATE UNIQUE INDEX release_recording_unique_idx ON release_recording (release_identifier, recording_identifier, source);

CREATE TABLE recording_work (
  id                   INTEGER PRIMARY KEY,
  recording_identifier INTEGER REFERENCES identifier (id),
  work_identifier      INTEGER REFERENCES identifier (id),
  source               INTEGER REFERENCES source (id)
);

CREATE INDEX recording_work_recording_idx ON recording_work (recording_identifier);
CREATE INDEX recording_work_work_idx      ON recording_work (work_identifier);
CREATE UNIQUE INDEX recording_work_unique_idx ON recording_work (recording_identifier, work_identifier, source);

CREATE TABLE release_song (
  id                 INTEGER PRIMARY KEY,
  release_identifier INTEGER REFERENCES identifier (id),
  song_identifier    INTEGER REFERENCES identifier (id),
  source             INTEGER REFERENCES source (id)
);

CREATE INDEX release_song_release_idx ON release_song (release_identifier);
CREATE INDEX release_song_song_idx    ON release_song (song_identifier);
CREATE UNIQUE INDEX release_song_unique_idx ON release_song (release_identifier, song_identifier, source);

CREATE TABLE source (
  id   INTEGER PRIMARY KEY,
  name TEXT NOT NULL
);
CREATE UNIQUE INDEX source_unique_idx ON source (name);
`,
	)
}
