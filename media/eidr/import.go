// This file is part of the go-meta library.
//
// Copyright (C) 2018 JAAK MUSIC LTD
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
	"encoding/xml"
	"fmt"
	"io"

	"github.com/meta-network/go-meta/media"
)

type Importer struct {
	client *media.Client
}

func NewImporter(client *media.Client) *Importer {
	return &Importer{client}
}

func (i *Importer) ImportParty(src io.Reader) error {
	var msg Party
	if err := xml.NewDecoder(src).Decode(&msg); err != nil {
		return err
	}
	return i.importParty(&msg)
}

func (i *Importer) ImportFullMetadata(src io.Reader) error {
	var msg FullMetadata
	if err := xml.NewDecoder(src).Decode(&msg); err != nil {
		return err
	}
	switch msg.BaseObjectData.ReferentType {
	case "Series":
		return i.importSeries(&msg)
	case "Season":
		return i.importSeason(&msg)
	case "TV":
		return i.importEpisode(&msg)
	default:
		return fmt.Errorf("unhandled ReferentType: %s", msg.BaseObjectData.ReferentType)
	}
}

func (i *Importer) importParty(msg *Party) error {
	identifier := &media.Identifier{
		Type:  "doid",
		Value: msg.ID,
	}
	org := &media.Organisation{
		Name: msg.PartyName.DisplayName,
	}
	return i.client.CreateOrganisation(org, identifier)
}

func (i *Importer) importSeries(msg *FullMetadata) error {
	identifier := &media.Identifier{
		Type:  "doid",
		Value: msg.BaseObjectData.ID,
	}
	series := &media.Series{
		Name: msg.BaseObjectData.ResourceName.Value,
	}
	if err := i.client.CreateSeries(series, identifier); err != nil {
		return err
	}
	org := media.Identifier{
		Type:  "doid",
		Value: msg.BaseObjectData.Administrators.Registrant,
	}
	return i.client.CreateOrganisationSeriesLink(&media.OrganisationSeriesLink{
		Organisation: org,
		Series:       *identifier,
	})
}

func (i *Importer) importSeason(msg *FullMetadata) error {
	identifier := &media.Identifier{
		Type:  "doid",
		Value: msg.BaseObjectData.ID,
	}
	season := &media.Season{
		Name: msg.BaseObjectData.ResourceName.Value,
	}
	if err := i.client.CreateSeason(season, identifier); err != nil {
		return err
	}
	org := media.Identifier{
		Type:  "doid",
		Value: msg.BaseObjectData.Administrators.Registrant,
	}
	if err := i.client.CreateOrganisationSeasonLink(&media.OrganisationSeasonLink{
		Organisation: org,
		Season:       *identifier,
	}); err != nil {
		return err
	}
	series := media.Identifier{
		Type:  "doid",
		Value: msg.ExtraObjectMetadata.SeasonInfo.Parent,
	}
	return i.client.CreateSeriesSeasonLink(&media.SeriesSeasonLink{
		Series: series,
		Season: *identifier,
	})
}

func (i *Importer) importEpisode(msg *FullMetadata) error {
	identifier := &media.Identifier{
		Type:  "doid",
		Value: msg.BaseObjectData.ID,
	}
	episode := &media.Episode{
		Name: msg.BaseObjectData.ResourceName.Value,
	}
	if err := i.client.CreateEpisode(episode, identifier); err != nil {
		return err
	}
	org := media.Identifier{
		Type:  "doid",
		Value: msg.BaseObjectData.Administrators.Registrant,
	}
	if err := i.client.CreateOrganisationEpisodeLink(&media.OrganisationEpisodeLink{
		Organisation: org,
		Episode:      *identifier,
	}); err != nil {
		return err
	}
	season := media.Identifier{
		Type:  "doid",
		Value: msg.ExtraObjectMetadata.EpisodeInfo.Parent,
	}
	return i.client.CreateSeasonEpisodeLink(&media.SeasonEpisodeLink{
		Season:  season,
		Episode: *identifier,
	})
}
