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

import "github.com/meta-network/go-meta/migrate"

// migrations is a set of database migrations to run on a SQLite3 database
// to prepare it for indexing a META stream of CWRs.
var migrations = migrate.NewMigrations()

func init() {
	// migration 1 creates indexes for the following artist properties:
	//
	// * Title -
	// * ISWC -
	// * CompositeType -
	// * record_type  -
	//
	migrations.Add(1, `
CREATE TABLE registered_work (
	object_id      text NOT NULL,
	title          text NOT NULL,
	iswc           text NOT NULL,
	composite_type text NOT NULL,
	record_type    text NOT NULL
);

CREATE INDEX registered_work_object_id_idx     ON registered_work (object_id);
CREATE INDEX registered_work_title_idx         ON registered_work (title);
CREATE INDEX registered_work_iswc_idx          ON registered_work (iswc);
CREATE INDEX registered_work_composite_type_idx ON registered_work (composite_type);
CREATE INDEX registered_work_record_type_idx   ON registered_work (record_type);
`,
	)
}
