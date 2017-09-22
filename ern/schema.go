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

import "github.com/meta-network/go-meta/migrate"

// migrations is a set of database migrations to run on a SQLite3 database
// to prepare it for indexing a META stream of ERNs.
var migrations = migrate.NewMigrations()

func init() {
	// migration 1 creates indexes and associations for ERN Releases,
	// Parties, MusicalWorks and SoundRecordings
	migrations.Add(1, `
--
-- the ern table is an index of ERN NewReleaseMessage using values from the
-- MessageHeader
--
CREATE TABLE ern (
        -- cid is the CID of the NewReleaseMessage
	cid text NOT NULL,

	-- message_id is the value of the MessageId header
	message_id text,

	-- thread_id is the value of the MessageThreadId header
	thread_id text,

	-- sender_id is the cid of the MessageSender party
	sender_id text,

	-- recipient_id is the cid of the MessageRecipient party
	recipient_id text,

	-- created is the value of the MessageCreatedDateTime header
	created datetime
);
CREATE INDEX ern_cid_idx          ON ern (cid);
CREATE INDEX ern_message_id_idx   ON ern (message_id);
CREATE INDEX ern_thread_id_idx    ON ern (thread_id);
CREATE INDEX ern_sender_id_idx    ON ern (sender_id);
CREATE INDEX ern_recipient_id_idx ON ern (recipient_id);
CREATE INDEX ern_created_idx      ON ern (created);

--
-- the party table is an index of ERN types which have PartyId and PartyName
-- properties (for example Artist, MusicalWorkContributor, MessagingParty etc.)
--
CREATE TABLE party (
	-- cid is the CID of the party
	cid text NOT NULL,

	-- id is the value of PartyId which is either a DPID or an ISNI
	id text,

	-- name is the value of PartyName which is either FullName or
	-- KeyName
	name text
);
CREATE INDEX party_cid_idx       ON party (cid);
CREATE INDEX party_id_idx        ON party (id);
CREATE INDEX party_name_idx      ON party (name);
CREATE UNIQUE INDEX party_unique_idx ON party (id, name);

--
-- the release table is an index of ERN Releases
--
CREATE TABLE release (
	-- cid is the CID of the Release
	cid text NOT NULL,

	-- id is the value of ReleaseId which is either a GRid, ISRC, ICPN,
	-- CatalogNumber or ProprietaryId
	id text NOT NULL,

	-- title is the value of the Release ReferenceTitle
	title text
);
CREATE INDEX release_cid_idx   ON release (cid);
CREATE INDEX release_id_idx    ON release (id);
CREATE INDEX release_title_idx ON release (title);

--
-- the release_list table associates an ERN with one or many Releases through
-- the ReleaseList property
--
CREATE TABLE release_list (
	-- ern_id is the cid of the ERN
	ern_id text NOT NULL,

	-- release_id is the cid of the Release
	release_id text NOT NULL
);
CREATE INDEX release_list_ern_id_idx     ON release_list (ern_id);
CREATE INDEX release_list_release_id_idx ON release_list (release_id);
CREATE UNIQUE INDEX release_list_unique_idx ON release_list (ern_id, release_id);

--
-- the musical_work table is an index of ERN MusicalWorks
--
CREATE TABLE musical_work (
	-- cid is the CID of the MusicalWork
	cid text NOT NULL,

	-- id is the value of MusicalWorkId which is either an ISWC,
	-- OpusNumber, ComposerCatalogNumber or ProprietaryId
	id text,

	-- title is the value of the MusicalWork ReferenceTitle
	title text
);
CREATE INDEX musical_work_cid_idx   ON musical_work (cid);
CREATE INDEX musical_work_id_idx    ON musical_work (id);
CREATE INDEX musical_work_title_idx ON musical_work (title);

--
-- the work_list table associates an ERN with one or many MusicalWorks through
-- the WorkList property
--
CREATE TABLE work_list (
	-- ern_id is the cid of the ERN
	ern_id text NOT NULL,

	-- musical_work_id is the cid of the MusicalWork
	musical_work_id text NOT NULL
);
CREATE INDEX work_list_ern_id_idx          ON work_list (ern_id);
CREATE INDEX work_list_musical_work_id_idx ON work_list (musical_work_id);
CREATE UNIQUE INDEX work_list_unique_idx ON work_list (ern_id, musical_work_id);

--
-- the musical_work_contributor table associates a MusicalWork with one or many
-- parties through the MusicalWorkContributor property
--
CREATE TABLE musical_work_contributor (
	-- musical_work_id is the cid of the MusicalWork
	musical_work_id text NOT NULL,

	-- party_id is the cid of the party
	party_id text NOT NULL
);
CREATE INDEX musical_work_contributor_id_idx    ON musical_work_contributor (musical_work_id);
CREATE INDEX musical_work_contributor_party_idx ON musical_work_contributor (party_id);
CREATE UNIQUE INDEX musical_work_contributor_unique_idx ON musical_work_contributor (musical_work_id, party_id);

--
-- the sound_recording table is an index of ERN SoundRecordings
--
CREATE TABLE sound_recording (
	-- cid is the CID of the SoundRecording
	cid text NOT NULL,

	-- id is the value of SoundRecordingId which is either an ISRC,
	-- CatalogNumber or ProprietaryId
	id text,

	-- title is the value of the SoundRecording ReferenceTitle
	title text
);
CREATE INDEX sound_recording_cid_idx   ON sound_recording (cid);
CREATE INDEX sound_recording_id_idx    ON sound_recording (id);
CREATE INDEX sound_recording_title_idx ON sound_recording (title);

--
-- the resource_list table associates an ERN with one or many SoundRecordings
-- through the ResourceList property
--
CREATE TABLE resource_list (
	-- ern_id is the cid of the ERN
	ern_id text NOT NULL,

	-- resource_id is the cid of the SoundRecording
	resource_id text NOT NULL
);
CREATE INDEX resource_list_ern_id_idx      ON resource_list (ern_id);
CREATE INDEX resource_list_resource_id_idx ON resource_list (resource_id);
CREATE UNIQUE INDEX resource_list_unique_idx ON resource_list (ern_id, resource_id);

--
-- the sound_recording_contributor table associates a SoundRecording with one
-- or many parties through the DisplayArtist and ResourceContributor properties
--
CREATE TABLE sound_recording_contributor (
	-- sound_recording_id is the cid of the SoundRecording
	sound_recording_id text NOT NULL,

	-- party_id is the cid of the party
	party_id text NOT NULL
);
CREATE INDEX sound_recording_contributor_id_idx    ON sound_recording_contributor (sound_recording_id);
CREATE INDEX sound_recording_contributor_party_idx ON sound_recording_contributor (party_id);
CREATE UNIQUE INDEX sound_recording_contributor_unique_idx ON sound_recording_contributor (sound_recording_id, party_id);
`,
	)
}
