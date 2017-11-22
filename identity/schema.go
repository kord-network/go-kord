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

package identity

import "github.com/meta-network/go-meta/migrate"

// migrations is a set of database migrations to run on a SQLite3 database
// to prepare it for indexing a META stream of IDs.
var migrations = migrate.NewMigrations()

func init() {
	// migration 1 creates indexes for the following id records:
	//

	migrations.Add(1, `
CREATE TABLE identity (
	object_id      text NOT NULL,
	owner          text NOT NULL,
	signature      text NOT NULL,
	id             text NOT NULL
);

CREATE INDEX identity_object_id_idx ON identity (object_id);
CREATE INDEX identity_signature_idx ON identity (signature);
CREATE INDEX identity_owner_idx     ON identity (owner);
CREATE INDEX identity_id_idx        ON identity (id);

CREATE TABLE claim (
	object_id text NOT NULL,
	issuer    text NOT NULL,
	holder    text NOT NULL,
	claim     text NOT NULL,
	signature text NOT NULL,
	id        text NOT NULL
);

CREATE INDEX claim_object_id_idx ON claim (object_id);
CREATE INDEX claim_issuer_idx    ON claim (issuer);
CREATE INDEX claim_holder_idx    ON claim (holder);
CREATE INDEX claim_id_idx        ON claim (id);
CREATE INDEX claim_claim_idx     ON claim (claim);
CREATE INDEX claim_signature_idx ON claim (signature);
`,
	)
}
