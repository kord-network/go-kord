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

import "github.com/meta-network/go-meta/media"

type Importer struct {
	client *media.Client
}

func NewImporter(client *media.Client) *Importer {
	return &Importer{client}
}

func (i *Importer) ImportArtist(artist *Artist) error {
	identifiers := []*media.Identifier{
		{
			Type:  "mbid",
			Value: artist.MBID,
		},
	}
	for _, ipi := range artist.IPI {
		identifiers = append(identifiers, &media.Identifier{
			Type:  "ipi",
			Value: ipi,
		})
	}
	for _, isni := range artist.ISNI {
		identifiers = append(identifiers, &media.Identifier{
			Type:  "isni",
			Value: isni,
		})
	}
	for _, identifier := range identifiers {
		if err := i.client.CreatePerformer(
			&media.Performer{Name: artist.Name},
			identifier,
		); err != nil {
			return err
		}
	}
	return nil
}

func (i *Importer) ImportRecordingWorkLink(link *RecordingWorkLink) error {
	recording := media.Identifier{
		Type:  "isrc",
		Value: link.ISRC,
	}
	if err := i.client.CreateRecording(
		&media.Recording{Title: link.RecordingTitle},
		&recording,
	); err != nil {
		return err
	}
	work := media.Identifier{
		Type:  "iswc",
		Value: link.ISWC,
	}
	if err := i.client.CreateWork(
		&media.Work{Title: link.WorkTitle},
		&work,
	); err != nil {
		return err
	}
	return i.client.CreateRecordingWorkLink(
		&media.RecordingWorkLink{
			Recording: recording,
			Work:      work,
		},
	)
}
