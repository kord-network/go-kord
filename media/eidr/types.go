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

type FullMetadata struct {
	BaseObjectData      *BaseObjectData      `xml:"BaseObjectData"`
	ExtraObjectMetadata *ExtraObjectMetadata `xml:"ExtraObjectMetadata"`
}

type BaseObjectData struct {
	ID                string            `xml:"ID"`
	StructuralType    string            `xml:"StructuralType"`
	Mode              string            `xml:"Mode"`
	ReferentType      string            `xml:"ReferentType"`
	ResourceName      *ResourceName     `xml:"ResourceName"`
	OriginalLanguage  *OriginalLanguage `xml:"OriginalLanguage"`
	AssociatedOrg     *AssociatedOrg    `xml:"AssociatedOrg"`
	ReleaseDate       string            `xml:"ReleaseDate"`
	CountryOfOrigin   string            `xml:"CountryOfOrigin"`
	Status            string            `xml:"Status"`
	ApproximateLength string            `xml:"ApproximateLength"`
	Administrators    *Administrators   `xml:"Administrators"`
	Credits           *Credits          `xml:"Credits"`
}

type ResourceName struct {
	Value      string `xml:",chardata"`
	Lang       string `xml:"lang,attr"`
	TitleClass string `xml:"titleClass,attr"`
}

type OriginalLanguage struct {
	Mode  string `xml:"mode,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

type AssociatedOrg struct {
	Role        string `xml:"role,attr"`
	DisplayName string `xml:"DisplayName"`
}

type Administrators struct {
	Registrant string `xml:"Registrant"`
}

type Credits struct {
	Actor []*Actor `xml:"Actor"`
}

type Actor struct {
	DisplayName string `xml:"DisplayName"`
}

type ExtraObjectMetadata struct {
	SeriesInfo  *SeriesInfo  `xml:"SeriesInfo"`
	SeasonInfo  *SeasonInfo  `xml:"SeasonInfo"`
	EpisodeInfo *EpisodeInfo `xml:"EpisodeInfo"`
}

type SeriesInfo struct {
	SeriesClass           string `xml:"SeriesClass"`
	NumberRequired        string `xml:"NumberRequired"`
	DateRequired          string `xml:"DateRequired"`
	OriginalTitleRequired string `xml:"OriginalTitleRequired"`
}

type SeasonInfo struct {
	Parent                string `xml:"Parent"`
	SeasonClass           string `xml:"SeasonClass"`
	NumberRequired        string `xml:"NumberRequired"`
	DateRequired          string `xml:"DateRequired"`
	OriginalTitleRequired string `xml:"OriginalTitleRequired"`
	SequenceNumber        string `xml:"SequenceNumber"`
}

type EpisodeInfo struct {
	Parent       string        `xml:"Parent"`
	SequenceInfo *SequenceInfo `xml:"SequenceInfo"`
}

type SequenceInfo struct {
	DistributionNumber string `xml:"DistributionNumber"`
}
