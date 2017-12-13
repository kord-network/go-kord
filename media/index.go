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

import (
	"database/sql"
	"fmt"

	sqlite3 "github.com/mattn/go-sqlite3"
	meta "github.com/meta-network/go-meta"
)

type Index struct {
	*meta.Index
}

func NewIndex(store *meta.Store) (*Index, error) {
	index, err := store.OpenIndex("media.meta")
	if err != nil {
		return nil, err
	}
	// migrate the db to ensure it has an up-to-date schema
	if err := migrations.Run(index.DB); err != nil {
		index.Close()
		return nil, err
	}
	return &Index{index}, nil
}

type errIdentifierNotFound struct {
	identifier *Identifier
}

func (e errIdentifierNotFound) Error() string {
	return fmt.Sprintf("media: identifier not found: type:%q value:%q", e.identifier.Type, e.identifier.Value)
}

func isIdentifierNotFound(err error) bool {
	_, ok := err.(errIdentifierNotFound)
	return ok
}

func (i *Index) Identifier(recordType string, identifier *Identifier) (*IdentifierRecord, error) {
	record := &IdentifierRecord{Identifier: *identifier}
	err := i.QueryRow(
		"SELECT identifier.id FROM identifier JOIN identifier_assignment ON identifier_assignment.identifier_id = identifier.id AND identifier_assignment.record_type = $1 WHERE identifier.type = $2 AND identifier.value = $3",
		recordType, identifier.Type, identifier.Value,
	).Scan(&record.ID)
	if err == sql.ErrNoRows {
		return nil, errIdentifierNotFound{identifier}
	} else if err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Index) Performers(identifier *IdentifierRecord) ([]*PerformerRecord, error) {
	rows, err := i.Query(
		"SELECT id, name, source FROM performer WHERE id IN (SELECT record_id FROM identifier_assignment WHERE record_type = 'performer' AND identifier_id = $1)",
		identifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var performers []*PerformerRecord
	for rows.Next() {
		var performer PerformerRecord
		if err := rows.Scan(
			&performer.ID,
			&performer.Name,
			&performer.Source,
		); err != nil {
			return nil, err
		}
		performers = append(performers, &performer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return performers, nil
}

func (i *Index) Composers(identifier *IdentifierRecord) ([]*ComposerRecord, error) {
	rows, err := i.Query(
		"SELECT id, first_name, last_name, source FROM composer WHERE id IN (SELECT record_id FROM identifier_assignment WHERE record_type = 'composer' AND identifier_id = $1)",
		identifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var composers []*ComposerRecord
	for rows.Next() {
		var composer ComposerRecord
		if err := rows.Scan(
			&composer.ID,
			&composer.FirstName,
			&composer.LastName,
			&composer.Source,
		); err != nil {
			return nil, err
		}
		composers = append(composers, &composer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return composers, nil
}

func (i *Index) RecordLabels(identifier *IdentifierRecord) ([]*RecordLabelRecord, error) {
	rows, err := i.Query(
		"SELECT id, name, source FROM record_label WHERE id IN (SELECT record_id FROM identifier_assignment WHERE record_type = 'record_label' AND identifier_id = $1)",
		identifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var recordLabels []*RecordLabelRecord
	for rows.Next() {
		var recordLabel RecordLabelRecord
		if err := rows.Scan(
			&recordLabel.ID,
			&recordLabel.Name,
			&recordLabel.Source,
		); err != nil {
			return nil, err
		}
		recordLabels = append(recordLabels, &recordLabel)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return recordLabels, nil
}

func (i *Index) Publishers(identifier *IdentifierRecord) ([]*PublisherRecord, error) {
	rows, err := i.Query(
		"SELECT id, name, source FROM publisher WHERE id IN (SELECT record_id FROM identifier_assignment WHERE record_type = 'publisher' AND identifier_id = $1)",
		identifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var publishers []*PublisherRecord
	for rows.Next() {
		var publisher PublisherRecord
		if err := rows.Scan(
			&publisher.ID,
			&publisher.Name,
			&publisher.Source,
		); err != nil {
			return nil, err
		}
		publishers = append(publishers, &publisher)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return publishers, nil
}

func (i *Index) Recordings(identifier *IdentifierRecord) ([]*RecordingRecord, error) {
	rows, err := i.Query(
		"SELECT id, title, duration, source FROM recording WHERE id IN (SELECT record_id FROM identifier_assignment WHERE record_type = 'recording' AND identifier_id = $1)",
		identifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var recordings []*RecordingRecord
	for rows.Next() {
		var recording RecordingRecord
		if err := rows.Scan(
			&recording.ID,
			&recording.Title,
			&recording.Duration,
			&recording.Source,
		); err != nil {
			return nil, err
		}
		recordings = append(recordings, &recording)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return recordings, nil
}

func (i *Index) Works(identifier *IdentifierRecord) ([]*WorkRecord, error) {
	rows, err := i.Query(
		"SELECT id, title, source FROM work WHERE id IN (SELECT record_id FROM identifier_assignment WHERE record_type = 'work' AND identifier_id = $1)",
		identifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var works []*WorkRecord
	for rows.Next() {
		var work WorkRecord
		if err := rows.Scan(
			&work.ID,
			&work.Title,
			&work.Source,
		); err != nil {
			return nil, err
		}
		works = append(works, &work)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return works, nil
}

func (i *Index) Songs(identifier *IdentifierRecord) ([]*SongRecord, error) {
	rows, err := i.Query(
		"SELECT id, title, duration, source FROM song WHERE id IN (SELECT record_id FROM identifier_assignment WHERE record_type = 'song' AND identifier_id = $1)",
		identifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var songs []*SongRecord
	for rows.Next() {
		var song SongRecord
		if err := rows.Scan(
			&song.ID,
			&song.Title,
			&song.Duration,
			&song.Source,
		); err != nil {
			return nil, err
		}
		songs = append(songs, &song)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return songs, nil
}

func (i *Index) Releases(identifier *IdentifierRecord) ([]*ReleaseRecord, error) {
	rows, err := i.Query(
		"SELECT id, type, title, date, source FROM release WHERE id IN (SELECT record_id FROM identifier_assignment WHERE record_type = 'release' AND identifier_id = $1)",
		identifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var releases []*ReleaseRecord
	for rows.Next() {
		var release ReleaseRecord
		if err := rows.Scan(
			&release.ID,
			&release.Type,
			&release.Title,
			&release.Date,
			&release.Source,
		); err != nil {
			return nil, err
		}
		releases = append(releases, &release)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return releases, nil
}

func (i *Index) PerformerRecordings(performerIdentifier *IdentifierRecord) ([]*PerformerRecordingRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.role, a.source FROM performer_recording AS a JOIN identifier AS b ON b.id = a.performer_identifier JOIN identifier AS c ON c.id = a.recording_identifier WHERE a.performer_identifier = $1",
		performerIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*PerformerRecordingRecord
	for rows.Next() {
		record := PerformerRecordingRecord{
			Performer: &IdentifierRecord{},
			Recording: &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Performer.ID,
			&record.Performer.Type,
			&record.Performer.Value,
			&record.Recording.ID,
			&record.Recording.Type,
			&record.Recording.Value,
			&record.Role,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) RecordingPerformers(recordingIdentifier *IdentifierRecord) ([]*PerformerRecordingRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.role, a.source FROM performer_recording AS a JOIN identifier AS b ON b.id = a.performer_identifier JOIN identifier AS c ON c.id = a.recording_identifier WHERE a.recording_identifier = $1",
		recordingIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*PerformerRecordingRecord
	for rows.Next() {
		record := PerformerRecordingRecord{
			Performer: &IdentifierRecord{},
			Recording: &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Performer.ID,
			&record.Performer.Type,
			&record.Performer.Value,
			&record.Recording.ID,
			&record.Recording.Type,
			&record.Recording.Value,
			&record.Role,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) ComposerWorks(composerIdentifier *IdentifierRecord) ([]*ComposerWorkRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.role, a.source FROM composer_work AS a JOIN identifier AS b ON b.id = a.composer_identifier JOIN identifier AS c ON c.id = a.work_identifier WHERE a.composer_identifier = $1",
		composerIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*ComposerWorkRecord
	for rows.Next() {
		record := ComposerWorkRecord{
			Composer: &IdentifierRecord{},
			Work:     &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Composer.ID,
			&record.Composer.Type,
			&record.Composer.Value,
			&record.Work.ID,
			&record.Work.Type,
			&record.Work.Value,
			&record.Role,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) WorkComposers(workIdentifier *IdentifierRecord) ([]*ComposerWorkRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.role, a.source FROM composer_work AS a JOIN identifier AS b ON b.id = a.composer_identifier JOIN identifier AS c ON c.id = a.work_identifier WHERE a.work_identifier = $1",
		workIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*ComposerWorkRecord
	for rows.Next() {
		record := ComposerWorkRecord{
			Composer: &IdentifierRecord{},
			Work:     &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Composer.ID,
			&record.Composer.Type,
			&record.Composer.Value,
			&record.Work.ID,
			&record.Work.Type,
			&record.Work.Value,
			&record.Role,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) RecordLabelSongs(recordLabelIdentifier *IdentifierRecord) ([]*RecordLabelSongRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM record_label_song AS a JOIN identifier AS b ON b.id = a.record_label_identifier JOIN identifier AS c ON c.id = a.song_identifier WHERE a.record_label_identifier = $1",
		recordLabelIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*RecordLabelSongRecord
	for rows.Next() {
		record := RecordLabelSongRecord{
			RecordLabel: &IdentifierRecord{},
			Song:        &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.RecordLabel.ID,
			&record.RecordLabel.Type,
			&record.RecordLabel.Value,
			&record.Song.ID,
			&record.Song.Type,
			&record.Song.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) SongRecordLabels(songIdentifier *IdentifierRecord) ([]*RecordLabelSongRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM record_label_song AS a JOIN identifier AS b ON b.id = a.record_label_identifier JOIN identifier AS c ON c.id = a.song_identifier WHERE a.song_identifier = $1",
		songIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*RecordLabelSongRecord
	for rows.Next() {
		record := RecordLabelSongRecord{
			RecordLabel: &IdentifierRecord{},
			Song:        &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.RecordLabel.ID,
			&record.RecordLabel.Type,
			&record.RecordLabel.Value,
			&record.Song.ID,
			&record.Song.Type,
			&record.Song.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) RecordLabelReleases(recordLabelIdentifier *IdentifierRecord) ([]*RecordLabelReleaseRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM record_label_release AS a JOIN identifier AS b ON b.id = a.record_label_identifier JOIN identifier AS c ON c.id = a.release_identifier WHERE a.record_label_identifier = $1",
		recordLabelIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*RecordLabelReleaseRecord
	for rows.Next() {
		record := RecordLabelReleaseRecord{
			RecordLabel: &IdentifierRecord{},
			Release:     &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.RecordLabel.ID,
			&record.RecordLabel.Type,
			&record.RecordLabel.Value,
			&record.Release.ID,
			&record.Release.Type,
			&record.Release.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) ReleaseRecordLabels(releaseIdentifier *IdentifierRecord) ([]*RecordLabelReleaseRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM record_label_release AS a JOIN identifier AS b ON b.id = a.record_label_identifier JOIN identifier AS c ON c.id = a.release_identifier WHERE a.release_identifier = $1",
		releaseIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*RecordLabelReleaseRecord
	for rows.Next() {
		record := RecordLabelReleaseRecord{
			RecordLabel: &IdentifierRecord{},
			Release:     &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.RecordLabel.ID,
			&record.RecordLabel.Type,
			&record.RecordLabel.Value,
			&record.Release.ID,
			&record.Release.Type,
			&record.Release.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) PublisherWorks(publisherIdentifier *IdentifierRecord) ([]*PublisherWorkRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM publisher_work AS a JOIN identifier AS b ON b.id = a.publisher_identifier JOIN identifier AS c ON c.id = a.work_identifier WHERE a.publisher_identifier = $1",
		publisherIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*PublisherWorkRecord
	for rows.Next() {
		record := PublisherWorkRecord{
			Publisher: &IdentifierRecord{},
			Work:      &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Publisher.ID,
			&record.Publisher.Type,
			&record.Publisher.Value,
			&record.Work.ID,
			&record.Work.Type,
			&record.Work.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) WorkPublishers(workIdentifier *IdentifierRecord) ([]*PublisherWorkRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM publisher_work AS a JOIN identifier AS b ON b.id = a.publisher_identifier JOIN identifier AS c ON c.id = a.work_identifier WHERE a.work_identifier = $1",
		workIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*PublisherWorkRecord
	for rows.Next() {
		record := PublisherWorkRecord{
			Publisher: &IdentifierRecord{},
			Work:      &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Publisher.ID,
			&record.Publisher.Type,
			&record.Publisher.Value,
			&record.Work.ID,
			&record.Work.Type,
			&record.Work.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) SongRecordings(songIdentifier *IdentifierRecord) ([]*SongRecordingRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM song_recording AS a JOIN identifier AS b ON b.id = a.song_identifier JOIN identifier AS c ON c.id = a.recording_identifier WHERE a.song_identifier = $1",
		songIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*SongRecordingRecord
	for rows.Next() {
		record := SongRecordingRecord{
			Song:      &IdentifierRecord{},
			Recording: &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Song.ID,
			&record.Song.Type,
			&record.Song.Value,
			&record.Recording.ID,
			&record.Recording.Type,
			&record.Recording.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) RecordingSongs(recordingIdentifier *IdentifierRecord) ([]*SongRecordingRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM song_recording AS a JOIN identifier AS b ON b.id = a.song_identifier JOIN identifier AS c ON c.id = a.recording_identifier WHERE a.recording_identifier = $1",
		recordingIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*SongRecordingRecord
	for rows.Next() {
		record := SongRecordingRecord{
			Song:      &IdentifierRecord{},
			Recording: &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Song.ID,
			&record.Song.Type,
			&record.Song.Value,
			&record.Recording.ID,
			&record.Recording.Type,
			&record.Recording.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) ReleaseRecordings(releaseIdentifier *IdentifierRecord) ([]*ReleaseRecordingRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM release_recording AS a JOIN identifier AS b ON b.id = a.release_identifier JOIN identifier AS c ON c.id = a.recording_identifier WHERE a.release_identifier = $1",
		releaseIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*ReleaseRecordingRecord
	for rows.Next() {
		record := ReleaseRecordingRecord{
			Release:   &IdentifierRecord{},
			Recording: &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Release.ID,
			&record.Release.Type,
			&record.Release.Value,
			&record.Recording.ID,
			&record.Recording.Type,
			&record.Recording.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) RecordingReleases(recordingIdentifier *IdentifierRecord) ([]*ReleaseRecordingRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM release_recording AS a JOIN identifier AS b ON b.id = a.release_identifier JOIN identifier AS c ON c.id = a.recording_identifier WHERE a.recording_identifier = $1",
		recordingIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*ReleaseRecordingRecord
	for rows.Next() {
		record := ReleaseRecordingRecord{
			Release:   &IdentifierRecord{},
			Recording: &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Release.ID,
			&record.Release.Type,
			&record.Release.Value,
			&record.Recording.ID,
			&record.Recording.Type,
			&record.Recording.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) RecordingWorks(recordingIdentifier *IdentifierRecord) ([]*RecordingWorkRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM recording_work AS a JOIN identifier AS b ON b.id = a.recording_identifier JOIN identifier AS c ON c.id = a.work_identifier WHERE a.recording_identifier = $1",
		recordingIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*RecordingWorkRecord
	for rows.Next() {
		record := RecordingWorkRecord{
			Recording: &IdentifierRecord{},
			Work:      &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Recording.ID,
			&record.Recording.Type,
			&record.Recording.Value,
			&record.Work.ID,
			&record.Work.Type,
			&record.Work.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) WorkRecordings(workIdentifier *IdentifierRecord) ([]*RecordingWorkRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM recording_work AS a JOIN identifier AS b ON b.id = a.recording_identifier JOIN identifier AS c ON c.id = a.work_identifier WHERE a.work_identifier = $1",
		workIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*RecordingWorkRecord
	for rows.Next() {
		record := RecordingWorkRecord{
			Recording: &IdentifierRecord{},
			Work:      &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Recording.ID,
			&record.Recording.Type,
			&record.Recording.Value,
			&record.Work.ID,
			&record.Work.Type,
			&record.Work.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) ReleaseSongs(releaseIdentifier *IdentifierRecord) ([]*ReleaseSongRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM release_song AS a JOIN identifier AS b ON b.id = a.release_identifier JOIN identifier AS c ON c.id = a.song_identifier WHERE a.release_identifier = $1",
		releaseIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*ReleaseSongRecord
	for rows.Next() {
		record := ReleaseSongRecord{
			Release: &IdentifierRecord{},
			Song:    &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Release.ID,
			&record.Release.Type,
			&record.Release.Value,
			&record.Song.ID,
			&record.Song.Type,
			&record.Song.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) SongReleases(songIdentifier *IdentifierRecord) ([]*ReleaseSongRecord, error) {
	rows, err := i.Query(
		"SELECT a.id, b.id, b.type, b.value, c.id, c.type, c.value, a.source FROM release_song AS a JOIN identifier AS b ON b.id = a.release_identifier JOIN identifier AS c ON c.id = a.song_identifier WHERE a.song_identifier = $1",
		songIdentifier.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []*ReleaseSongRecord
	for rows.Next() {
		record := ReleaseSongRecord{
			Release: &IdentifierRecord{},
			Song:    &IdentifierRecord{},
		}
		if err := rows.Scan(
			&record.ID,
			&record.Release.ID,
			&record.Release.Type,
			&record.Release.Value,
			&record.Song.ID,
			&record.Song.Type,
			&record.Song.Value,
			&record.Source,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (i *Index) Source(id int64) (*SourceRecord, error) {
	record := &SourceRecord{ID: id}
	if err := i.QueryRow("SELECT name FROM source WHERE id = $1", id).Scan(&record.Name); err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Index) CreateRecord(record interface{}, identifier *Identifier, source *Source) (*IdentifierRecord, error) {
	var recordType string
	err := i.Update(func(tx *sql.Tx) error {
		source, err := i.createSource(tx, source)
		if err != nil {
			return err
		}
		var insertStmt []interface{}
		var selectStmt []interface{}
		switch v := record.(type) {
		case *Performer:
			recordType = "performer"
			insertStmt = []interface{}{
				"INSERT INTO performer (name, source) VALUES ($1, $2)",
				v.Name, source.ID,
			}
			selectStmt = []interface{}{
				"SELECT id FROM performer WHERE name = $1 AND source = $2",
				v.Name, source.ID,
			}
		case *Composer:
			recordType = "composer"
			insertStmt = []interface{}{
				"INSERT INTO composer (first_name, last_name, source) VALUES ($1, $2, $3)",
				v.FirstName, v.LastName, source.ID,
			}
			selectStmt = []interface{}{
				"SELECT id FROM composer WHERE first_name = $1 AND last_name = $2 AND source = $3",
				v.FirstName, v.LastName, source.ID,
			}
		case *RecordLabel:
			recordType = "record_label"
			insertStmt = []interface{}{
				"INSERT INTO record_label (name, source) VALUES ($1, $2)",
				v.Name, source.ID,
			}
			selectStmt = []interface{}{
				"SELECT id FROM record_label WHERE name = $1 AND source = $2",
				v.Name, source.ID,
			}
		case *Publisher:
			recordType = "publisher"
			insertStmt = []interface{}{
				"INSERT INTO publisher (name, source) VALUES ($1, $2)",
				v.Name, source.ID,
			}
			selectStmt = []interface{}{
				"SELECT id FROM publisher WHERE name = $1 AND source = $2",
				v.Name, source.ID,
			}
		case *Recording:
			recordType = "recording"
			insertStmt = []interface{}{
				"INSERT INTO recording (title, duration, source) VALUES ($1, $2, $3)",
				v.Title, v.Duration, source.ID,
			}
			selectStmt = []interface{}{
				"SELECT id FROM recording WHERE title = $1 AND duration = $2 AND source = $3",
				v.Title, v.Duration, source.ID,
			}
		case *Work:
			recordType = "work"
			insertStmt = []interface{}{
				"INSERT INTO work (title, source) VALUES ($1, $2)",
				v.Title, source.ID,
			}
			selectStmt = []interface{}{
				"SELECT id FROM work WHERE title = $1 AND source = $2",
				v.Title, source.ID,
			}
		case *Song:
			recordType = "song"
			insertStmt = []interface{}{
				"INSERT INTO song (title, duration, source) VALUES ($1, $2, $3)",
				v.Title, v.Duration, source.ID,
			}
			selectStmt = []interface{}{
				"SELECT id FROM song WHERE title = $1 AND duration = $2 AND source = $3",
				v.Title, v.Duration, source.ID,
			}
		case *Release:
			recordType = "release"
			insertStmt = []interface{}{
				"INSERT INTO release (type, title, date, source) VALUES ($1, $2, $3, $4)",
				v.Type, v.Title, v.Date, source.ID,
			}
			selectStmt = []interface{}{
				"SELECT id FROM release WHERE type = $1 AND title = $2 AND date = $3 AND source = $4",
				v.Type, v.Title, v.Date, source.ID,
			}
		}
		var id int64
		res, err := tx.Exec(insertStmt[0].(string), insertStmt[1:]...)
		if err == nil {
			id, err = res.LastInsertId()
			if err != nil {
				return err
			}
		} else if isUniqueErr(err) {
			if err := tx.QueryRow(selectStmt[0].(string), selectStmt[1:]...).Scan(&id); err != nil {
				return err
			}
		} else {
			return err
		}
		if _, err := tx.Exec(
			"INSERT OR IGNORE INTO identifier (type, value) VALUES ($1, $2)",
			identifier.Type, identifier.Value,
		); err != nil {
			return err
		}
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO identifier_assignment (identifier_id, record_type, record_id, source) VALUES ((SELECT id FROM identifier WHERE type = $1 AND value = $2), $3, $4, $5)",
			identifier.Type, identifier.Value, recordType, id, source.ID,
		)
		return err
	})
	if err != nil {
		return nil, err
	}
	return i.Identifier(recordType, identifier)
}

func (i *Index) CreatePerformerRecording(link *PerformerRecordingLink, source *Source) (*PerformerRecordingRecord, error) {
	performer, err := i.Identifier("performer", &link.Performer)
	if err != nil {
		return nil, err
	}
	recording, err := i.Identifier("recording", &link.Recording)
	if err != nil {
		return nil, err
	}
	record := &PerformerRecordingRecord{
		Performer: performer,
		Recording: recording,
		Role:      link.Role,
	}
	err = i.Update(func(tx *sql.Tx) error {
		source, err := i.createSource(tx, source)
		if err != nil {
			return err
		}
		record.Source = source.ID
		res, err := tx.Exec(
			"INSERT INTO performer_recording (performer_identifier, recording_identifier, role, source) VALUES ($1, $2, $3, $4)",
			performer.ID, recording.ID, record.Role, record.Source,
		)
		if err == nil {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			record.ID = id
			return nil
		} else if isUniqueErr(err) {
			return tx.QueryRow(
				"SELECT id FROM performer_recording WHERE performer_identifier = $1 AND recording_identifier = $2 AND role = $3 AND source = $4",
				performer.ID, recording.ID, record.Role, record.Source,
			).Scan(&record.ID)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Index) CreateComposerWork(link *ComposerWorkLink, source *Source) (*ComposerWorkRecord, error) {
	composer, err := i.Identifier("composer", &link.Composer)
	if err != nil {
		return nil, err
	}
	work, err := i.Identifier("work", &link.Work)
	if err != nil {
		return nil, err
	}
	record := &ComposerWorkRecord{
		Composer: composer,
		Work:     work,
		Role:     link.Role,
	}
	err = i.Update(func(tx *sql.Tx) error {
		source, err := i.createSource(tx, source)
		if err != nil {
			return err
		}
		record.Source = source.ID
		res, err := tx.Exec(
			"INSERT INTO composer_work (composer_identifier, work_identifier, role, source) VALUES ($1, $2, $3, $4)",
			composer.ID, work.ID, record.Role, record.Source,
		)
		if err == nil {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			record.ID = id
			return nil
		} else if isUniqueErr(err) {
			return tx.QueryRow(
				"SELECT id FROM composer_work WHERE composer_identifier = $1 AND work_identifier = $2 AND role = $3 AND source = $4",
				composer.ID, work.ID, record.Role, record.Source,
			).Scan(&record.ID)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Index) CreateRecordLabelSong(link *RecordLabelSongLink, source *Source) (*RecordLabelSongRecord, error) {
	recordLabel, err := i.Identifier("record_label", &link.RecordLabel)
	if err != nil {
		return nil, err
	}
	song, err := i.Identifier("song", &link.Song)
	if err != nil {
		return nil, err
	}
	record := &RecordLabelSongRecord{
		RecordLabel: recordLabel,
		Song:        song,
	}
	err = i.Update(func(tx *sql.Tx) error {
		source, err := i.createSource(tx, source)
		if err != nil {
			return err
		}
		record.Source = source.ID
		res, err := tx.Exec(
			"INSERT INTO record_label_song (record_label_identifier, song_identifier, source) VALUES ($1, $2, $3)",
			recordLabel.ID, song.ID, record.Source,
		)
		if err == nil {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			record.ID = id
			return nil
		} else if isUniqueErr(err) {
			return tx.QueryRow(
				"SELECT id FROM record_label_song WHERE record_label_identifier = $1 AND song_identifier = $2 AND source = $3",
				recordLabel.ID, song.ID, record.Source,
			).Scan(&record.ID)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Index) CreateRecordLabelRelease(link *RecordLabelReleaseLink, source *Source) (*RecordLabelReleaseRecord, error) {
	recordLabel, err := i.Identifier("record_label", &link.RecordLabel)
	if err != nil {
		return nil, err
	}
	release, err := i.Identifier("release", &link.Release)
	if err != nil {
		return nil, err
	}
	record := &RecordLabelReleaseRecord{
		RecordLabel: recordLabel,
		Release:     release,
	}
	err = i.Update(func(tx *sql.Tx) error {
		source, err := i.createSource(tx, source)
		if err != nil {
			return err
		}
		record.Source = source.ID
		res, err := tx.Exec(
			"INSERT INTO record_label_release (record_label_identifier, release_identifier, source) VALUES ($1, $2, $3)",
			recordLabel.ID, release.ID, record.Source,
		)
		if err == nil {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			record.ID = id
			return nil
		} else if isUniqueErr(err) {
			return tx.QueryRow(
				"SELECT id FROM record_label_release WHERE record_label_identifier = $1 AND release_identifier = $2 AND source = $3",
				recordLabel.ID, release.ID, record.Source,
			).Scan(&record.ID)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Index) CreatePublisherWork(link *PublisherWorkLink, source *Source) (*PublisherWorkRecord, error) {
	publisher, err := i.Identifier("publisher", &link.Publisher)
	if err != nil {
		return nil, err
	}
	work, err := i.Identifier("work", &link.Work)
	if err != nil {
		return nil, err
	}
	record := &PublisherWorkRecord{
		Publisher: publisher,
		Work:      work,
	}
	err = i.Update(func(tx *sql.Tx) error {
		source, err := i.createSource(tx, source)
		if err != nil {
			return err
		}
		record.Source = source.ID
		res, err := tx.Exec(
			"INSERT INTO publisher_work (publisher_identifier, work_identifier, source) VALUES ($1, $2, $3)",
			publisher.ID, work.ID, record.Source,
		)
		if err == nil {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			record.ID = id
			return nil
		} else if isUniqueErr(err) {
			return tx.QueryRow(
				"SELECT id FROM publisher_work WHERE publisher_identifier = $1 AND work_identifier = $2 AND source = $3",
				publisher.ID, work.ID, record.Source,
			).Scan(&record.ID)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Index) CreateSongRecording(link *SongRecordingLink, source *Source) (*SongRecordingRecord, error) {
	song, err := i.Identifier("song", &link.Song)
	if err != nil {
		return nil, err
	}
	recording, err := i.Identifier("recording", &link.Recording)
	if err != nil {
		return nil, err
	}
	record := &SongRecordingRecord{
		Song:      song,
		Recording: recording,
	}
	err = i.Update(func(tx *sql.Tx) error {
		source, err := i.createSource(tx, source)
		if err != nil {
			return err
		}
		record.Source = source.ID
		res, err := tx.Exec(
			"INSERT INTO song_recording (song_identifier, recording_identifier, source) VALUES ($1, $2, $3)",
			song.ID, recording.ID, record.Source,
		)
		if err == nil {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			record.ID = id
			return nil
		} else if isUniqueErr(err) {
			return tx.QueryRow(
				"SELECT id FROM song_recording WHERE song_identifier = $1 AND recording_identifier = $2 AND source = $3",
				song.ID, recording.ID, record.Source,
			).Scan(&record.ID)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Index) CreateReleaseRecording(link *ReleaseRecordingLink, source *Source) (*ReleaseRecordingRecord, error) {
	release, err := i.Identifier("release", &link.Release)
	if err != nil {
		return nil, err
	}
	recording, err := i.Identifier("recording", &link.Recording)
	if err != nil {
		return nil, err
	}
	record := &ReleaseRecordingRecord{
		Release:   release,
		Recording: recording,
	}
	err = i.Update(func(tx *sql.Tx) error {
		source, err := i.createSource(tx, source)
		if err != nil {
			return err
		}
		record.Source = source.ID
		res, err := tx.Exec(
			"INSERT INTO release_recording (release_identifier, recording_identifier, source) VALUES ($1, $2, $3)",
			release.ID, recording.ID, record.Source,
		)
		if err == nil {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			record.ID = id
			return nil
		} else if isUniqueErr(err) {
			return tx.QueryRow(
				"SELECT id FROM release_recording WHERE release_identifier = $1 AND recording_identifier = $2 AND source = $3",
				release.ID, recording.ID, record.Source,
			).Scan(&record.ID)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Index) CreateRecordingWork(link *RecordingWorkLink, source *Source) (*RecordingWorkRecord, error) {
	recording, err := i.Identifier("recording", &link.Recording)
	if err != nil {
		return nil, err
	}
	work, err := i.Identifier("work", &link.Work)
	if err != nil {
		return nil, err
	}
	record := &RecordingWorkRecord{
		Recording: recording,
		Work:      work,
	}
	err = i.Update(func(tx *sql.Tx) error {
		source, err := i.createSource(tx, source)
		if err != nil {
			return err
		}
		record.Source = source.ID
		res, err := tx.Exec(
			"INSERT INTO recording_work (recording_identifier, work_identifier, source) VALUES ($1, $2, $3)",
			recording.ID, work.ID, record.Source,
		)
		if err == nil {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			record.ID = id
			return nil
		} else if isUniqueErr(err) {
			return tx.QueryRow(
				"SELECT id FROM recording_work WHERE recording_identifier = $1 AND work_identifier = $2 AND source = $3",
				recording.ID, work.ID, record.Source,
			).Scan(&record.ID)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Index) CreateReleaseSong(link *ReleaseSongLink, source *Source) (*ReleaseSongRecord, error) {
	release, err := i.Identifier("release", &link.Release)
	if err != nil {
		return nil, err
	}
	song, err := i.Identifier("song", &link.Song)
	if err != nil {
		return nil, err
	}
	record := &ReleaseSongRecord{
		Release: release,
		Song:    song,
	}
	err = i.Update(func(tx *sql.Tx) error {
		source, err := i.createSource(tx, source)
		if err != nil {
			return err
		}
		record.Source = source.ID
		res, err := tx.Exec(
			"INSERT INTO release_song (release_identifier, song_identifier, source) VALUES ($1, $2, $3)",
			release.ID, song.ID, record.Source,
		)
		if err == nil {
			id, err := res.LastInsertId()
			if err != nil {
				return err
			}
			record.ID = id
			return nil
		} else if isUniqueErr(err) {
			return tx.QueryRow(
				"SELECT id FROM release_song WHERE release_identifier = $1 AND song_identifier = $2 AND source = $3",
				release.ID, song.ID, record.Source,
			).Scan(&record.ID)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (i *Index) createSource(tx *sql.Tx, source *Source) (*SourceRecord, error) {
	if _, err := tx.Exec("INSERT OR IGNORE INTO source (name) VALUES ($1)", source.Name); err != nil {
		return nil, err
	}
	record := &SourceRecord{Name: source.Name}
	if err := tx.QueryRow("SELECT id FROM source WHERE name = $1", source.Name).Scan(&record.ID); err != nil {
		return nil, err
	}
	return record, nil
}

// isUniqueErr determines whether an error is a SQLite3 uniqueness error.
func isUniqueErr(err error) bool {
	e, ok := err.(sqlite3.Error)
	if !ok {
		return false
	}
	return e.Code == sqlite3.ErrConstraint && e.ExtendedCode == sqlite3.ErrConstraintUnique
}
