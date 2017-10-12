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
	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta/doi"
)

type AssociatedOrg struct {
	ID          *string
	IDType      *string
	DisplayName string
	Role        string `json:"role"`
}

type AlternateID struct {
	ID       string `xml:"xs:string"`
	Type     string
	Domain   *string
	Relation *string
}

type Sequence struct {
	Value  string
	Domain string
}

type SequenceInfo struct {
	DistributionNumber *Sequence
	HouseSequence      *Sequence
	AlternateNumber    []*Sequence
}
type BaseObjectData struct {
	ID             doi.ID
	StructuralType string
	ReferentType   string
	ResourceName   string
	AlternateID    []AlternateID
	Status         string
	ExtraMetaData  cid.Cid
}

type Series struct {
	SeriesClass           string
	EndDate               string
	NumberRequired        bool
	DateRequired          bool
	OriginalTitleRequired bool
}

type Season struct {
	Parent                *cid.Cid
	SeasonClass           []string
	EndDate               string
	NumberRequired        bool
	DateRequired          bool
	OriginalTitleRequired bool
	SequenceNumber        int
}

type episode struct {
	Parent       *cid.Cid
	EpisodeClass string
	SequenceInfo *SequenceInfo
	TimeSlot     string
}
