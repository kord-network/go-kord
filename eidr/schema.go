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

package eidr

import (
	"github.com/meta-network/go-meta/migrate"
)

// migrations is a set of database migrations to run on a SQLite3 database
// to prepare it to index EIDR objects
var migrations = migrate.NewMigrations()

func init() {
	// migration 1 creates indexes for eidr objects
	// intended for use as a compound onthology of
	// ReferentType and respective ExtraMetadata derived types
	//
	// baseobject: Properties common to all eidr objects
	// alternateid: 0+ references to external ids pointing to the same object
	// associatedorg: 0+ references to organizations involved in the realization of the media object
	// xobject_*: Respective ExtraMetadata derived types
	//
	// Certain objects have a compulsory parent reference
	migrations.Add(1, `
CREATE TABLE baseobject (
	doi_id TEXT NOT NULL,
	structural_type TEXT NOT NULL,
	referent_type TEXT NOT NULL,
	status TEXT NOT NULL,
	resource_name TEXT NOT NULL,
	resource_name_lang TEXT NOT NULL,
	resource_name_class TEXT
);

CREATE UNIQUE INDEX baseobject_doi_id_name_idx ON baseobject(doi_id, resource_name);

CREATE TABLE org (
	id TEXT,
	idtype TEXT,
	display_name TEXT NOT NULL,
	role TEXT NOT NULL,
	base_doi_id TEXT NOT NULL,
	FOREIGN KEY(base_doi_id) REFERENCES baseobject(doi_id)
);

CREATE INDEX obj_idx ON org(id, idtype, base_doi_id);

CREATE TABLE alternateid (
	id TEXT NOT NULL,
	type TEXT NOT NULL,
	domain TEXT,
	relation TEXT,
	base_doi_id INTEGER NOT NULL,
	FOREIGN KEY(base_doi_id) REFERENCES baseobjecttype(doi_id)
);

CREATE INDEX alternateid_baseobject_idx ON alternateid(id, base_doi_id);

CREATE TABLE xobject_baseobject_link (
	base_doi_id TEXT NOT NULL,
	parent_doi_id TEXT,
	xobject_id INTEGER NOT NULL,
	xobject_type TEXT NOT NULL
);

CREATE UNIQUE INDEX xobject_baseobject_link_idx ON xobject_baseobject_link(base_doi_id, xobject_id, xobject_type);
CREATE UNIQUE INDEX baseobject_parent_child_link_idx ON xobject_baseobject_link(base_doi_id, parent_doi_id);

CREATE TABLE xobject_series (
	id INTEGER NOT NULL PRIMARY KEY,
	series_class TEXT
);

CREATE TABLE xobject_season (
	id INTEGER NOT NULL PRIMARY KEY,
	sequence_number INT
);

CREATE TABLE xobject_episode (
	id INTEGER NOT NULL PRIMARY KEY,
	episode_class TEXT
);

`,
	)
}
