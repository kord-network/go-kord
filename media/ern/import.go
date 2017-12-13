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

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/meta-network/go-meta/media"
)

type Importer struct {
	client *media.Client
}

func NewImporter(client *media.Client) *Importer {
	return &Importer{client}
}

type importContext struct {
	soundRecordings map[string]*media.Identifier
	recordLabel     *media.Identifier
	mainRelease     *media.Identifier
	songs           []*media.Identifier
}

func newImportContext() *importContext {
	return &importContext{
		soundRecordings: make(map[string]*media.Identifier),
	}
}

func (i *Importer) ImportERN(src io.Reader) error {
	var msg NewReleaseMessage
	if err := xml.NewDecoder(src).Decode(&msg); err != nil {
		return err
	}

	if msg.MessageHeader == nil {
		return errors.New("missing MessageHeader")
	}
	if msg.MessageHeader.MessageSender == nil {
		return errors.New("missing MessageSender")
	}

	ctx := newImportContext()

	recordLabel, err := i.importRecordLabel(ctx, msg.MessageHeader.MessageSender)
	if err != nil {
		return err
	}
	ctx.recordLabel = recordLabel

	if msg.WorkList != nil {
		for _, work := range msg.WorkList.MusicalWork {
			if err := i.importMusicalWork(ctx, work); err != nil {
				return err
			}
		}
	}

	if msg.ResourceList != nil {
		for _, soundRecording := range msg.ResourceList.SoundRecording {
			if err := i.importSoundRecording(ctx, soundRecording); err != nil {
				return err
			}
		}
	}

	if msg.ReleaseList != nil {
		for _, release := range msg.ReleaseList.Release {
			if err := i.importRelease(ctx, release); err != nil {
				return err
			}
		}
	}

	if ctx.mainRelease != nil {
		if err := i.client.CreateRecordLabelReleaseLink(&media.RecordLabelReleaseLink{
			RecordLabel: *recordLabel,
			Release:     *ctx.mainRelease,
		}); err != nil {
			return err
		}
		for _, song := range ctx.songs {
			if err := i.client.CreateReleaseSongLink(&media.ReleaseSongLink{
				Release: *ctx.mainRelease,
				Song:    *song,
			}); err != nil {
				return err
			}
		}
	}

	for _, song := range ctx.songs {
		if err := i.client.CreateRecordLabelSongLink(&media.RecordLabelSongLink{
			RecordLabel: *recordLabel,
			Song:        *song,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (i *Importer) importRecordLabel(ctx *importContext, sender *MessagingParty) (*media.Identifier, error) {
	var recordLabel media.RecordLabel
	if v := sender.PartyName; v != nil {
		recordLabel.Name = v.FullName.Value
	}
	identifier, err := i.partyIdentifier(sender.PartyId, sender)
	if err != nil {
		return nil, err
	}
	if err := i.client.CreateRecordLabel(&recordLabel, identifier); err != nil {
		return nil, err
	}
	return identifier, nil
}

func (i *Importer) importMusicalWork(ctx *importContext, musicalWork *MusicalWork) error {
	var work media.Work
	if v := musicalWork.ReferenceTitle; v != nil && v.TitleText != nil {
		work.Title = v.TitleText.Value
	}
	var identifier media.Identifier
	if id := musicalWork.MusicalWorkId; id != nil {
		identifier.Type = "iswc"
		identifier.Value = id.ISWC
	}
	if err := i.client.CreateWork(&work, &identifier); err != nil {
		return err
	}
	return nil
}

func (i *Importer) importSoundRecording(ctx *importContext, soundRecording *SoundRecording) error {
	recording := &media.Recording{
		Duration: soundRecording.Duration,
	}
	if v := soundRecording.ReferenceTitle; v != nil && v.TitleText != nil {
		recording.Title = v.TitleText.Value
	}
	var identifier media.Identifier
	if id := soundRecording.SoundRecordingId; id != nil {
		identifier.Type = "isrc"
		identifier.Value = id.ISRC
	}
	if err := i.client.CreateRecording(recording, &identifier); err != nil {
		return err
	}
	ctx.soundRecordings[soundRecording.ResourceReference] = &identifier
	for _, details := range soundRecording.SoundRecordingDetailsByTerritory {
		for _, artist := range details.DisplayArtist {
			if err := i.importArtist(ctx, artist, &identifier); err != nil {
				return err
			}
		}
		for _, contributor := range details.ResourceContributor {
			if err := i.importContributor(ctx, contributor, &identifier); err != nil {
				return err
			}
		}
		for _, contributor := range details.IndirectResourceContributor {
			if err := i.importIndirectContributor(ctx, contributor, &identifier); err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *Importer) importRelease(ctx *importContext, release *Release) error {
	// only index releases which have sound recordings
	if len(ctx.soundRecordings) == 0 {
		return nil
	}

	var identifier media.Identifier
	if id := release.ReleaseId; id != nil {
		switch {
		case id.GRid != "":
			identifier.Type = "grid"
			identifier.Value = id.GRid
		case id.ISRC != "":
			identifier.Type = "isrc"
			identifier.Value = id.ISRC
		case id.ICPN != nil:
			identifier.Type = "icpn"
			identifier.Value = id.ICPN.Value
		default:
			return fmt.Errorf("invalid release identifier: %v", id)
		}
	}

	var title string
	if v := release.ReferenceTitle; v != nil && v.TitleText != nil {
		title = v.TitleText.Value
	}

	var releaseType string
	if v := release.ReleaseType; v != nil {
		releaseType = v.Value
	}

	if release.IsMainRelease {
		mediaRelease := &media.Release{
			Type:  releaseType,
			Title: title,
		}
		for _, details := range release.ReleaseDetailsByTerritory {
			if v := details.TerritoryCode; v != nil && v.Value == "Worldwide" && details.ReleaseDate != nil {
				mediaRelease.Date = details.ReleaseDate.Value
			}
		}
		if err := i.client.CreateRelease(mediaRelease, &identifier); err != nil {
			return err
		}
		ctx.mainRelease = &identifier
		if list := release.ReleaseResourceReferenceList; list != nil {
			for _, ref := range list.ReleaseResourceReference {
				recording, ok := ctx.soundRecordings[ref.Value]
				if !ok {
					continue
				}
				link := &media.ReleaseRecordingLink{
					Release:   identifier,
					Recording: *recording,
				}
				if err := i.client.CreateReleaseRecordingLink(link); err != nil {
					return err
				}
			}
		}
	} else if releaseType == "TrackRelease" {
		song := &media.Song{
			Title:    title,
			Duration: release.Duration,
		}
		if err := i.client.CreateSong(song, &identifier); err != nil {
			return err
		}
		ctx.songs = append(ctx.songs, &identifier)
		if list := release.ReleaseResourceReferenceList; list != nil {
			for _, ref := range list.ReleaseResourceReference {
				recording, ok := ctx.soundRecordings[ref.Value]
				if !ok {
					continue
				}
				link := &media.SongRecordingLink{
					Song:      identifier,
					Recording: *recording,
				}
				if err := i.client.CreateSongRecordingLink(link); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (i *Importer) importArtist(ctx *importContext, artist *Artist, recording *media.Identifier) error {
	var performer media.Performer
	if v := artist.PartyName; v != nil {
		performer.Name = v.FullName.Value
	}

	identifier, err := i.partyIdentifier(artist.PartyId, artist)
	if err != nil {
		return err
	}

	if err := i.client.CreatePerformer(&performer, identifier); err != nil {
		return err
	}

	link := &media.PerformerRecordingLink{
		Performer: *identifier,
		Recording: *recording,
	}
	if v := artist.ArtistRole; v != nil {
		if v.Value == "UserDefined" {
			link.Role = i.joinNamespace(v.Namespace, v.UserDefinedValue)
		} else {
			link.Role = i.joinNamespace(v.Namespace, v.Value)
		}
	}
	return i.client.CreatePerformerRecordingLink(link)
}

func (i *Importer) importContributor(ctx *importContext, contributor *DetailedResourceContributor, recording *media.Identifier) error {
	var performer media.Performer
	if v := contributor.PartyName; v != nil {
		performer.Name = v.FullName.Value
	}

	identifier, err := i.partyIdentifier(contributor.PartyId, contributor)
	if err != nil {
		return err
	}

	if err := i.client.CreatePerformer(&performer, identifier); err != nil {
		return err
	}

	for _, v := range contributor.ResourceContributorRole {
		link := &media.PerformerRecordingLink{
			Performer: *identifier,
			Recording: *recording,
		}
		if v.Value == "UserDefined" {
			link.Role = i.joinNamespace(v.Namespace, v.UserDefinedValue)
		} else {
			link.Role = i.joinNamespace(v.Namespace, v.Value)
		}
		if err := i.client.CreatePerformerRecordingLink(link); err != nil {
			return err
		}
	}
	return nil
}

func (i *Importer) importIndirectContributor(ctx *importContext, contributor *IndirectResourceContributor, recording *media.Identifier) error {
	var performer media.Performer
	if v := contributor.PartyName; v != nil {
		performer.Name = v.FullName.Value
	}

	identifier, err := i.partyIdentifier(contributor.PartyId, contributor)
	if err != nil {
		return err
	}

	if err := i.client.CreatePerformer(&performer, identifier); err != nil {
		return err
	}

	for _, v := range contributor.IndirectResourceContributorRole {
		link := &media.PerformerRecordingLink{
			Performer: *identifier,
			Recording: *recording,
		}
		if v.Value == "UserDefined" {
			link.Role = i.joinNamespace(v.Namespace, v.UserDefinedValue)
		} else {
			link.Role = i.joinNamespace(v.Namespace, v.Value)
		}
		if err := i.client.CreatePerformerRecordingLink(link); err != nil {
			return err
		}
	}
	return nil
}

func (i *Importer) joinNamespace(namespace, value string) string {
	if namespace == "" {
		return value
	}
	return namespace + "|" + value
}

func (i *Importer) partyIdentifier(partyID *PartyId, v interface{}) (*media.Identifier, error) {
	if partyID != nil {
		var identifier media.Identifier
		if partyID.IsISNI {
			identifier.Type = "isni"
		} else {
			identifier.Type = "dpid"
		}
		identifier.Value = i.joinNamespace(partyID.Namespace, partyID.Value)
		return &identifier, nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return &media.Identifier{
		Type:  "sha3",
		Value: hexutil.Encode(crypto.Keccak256(data)),
	}, nil
}
