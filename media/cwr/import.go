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

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/meta-network/go-meta/media"
)

type Importer struct {
	client *media.Client
}

func NewImporter(client *media.Client) *Importer {
	return &Importer{client}
}

func newInput(src io.Reader) *input {
	return &input{r: bufio.NewReader(src)}
}

type input struct {
	r *bufio.Reader
}

func (i *input) Next() (interface{}, error) {
	line, err := i.r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if len(line) <= 3 {
		return nil, fmt.Errorf("cwr: line too short: %s", line)
	}
	recordType, recordData := line[0:3], line[3:]
	typ, ok := cwrTypes[recordType]
	if !ok {
		return nil, fmt.Errorf("cwr: unknown record type: %q", recordType)
	}
	val := reflect.New(typ)
	start := 0
	for i, length := range cwrLengths[recordType] {
		if start >= len(recordData) {
			break
		}
		end := start + length
		if end > len(recordData) {
			end = len(recordData)
		}
		val.Elem().Field(i).SetString(strings.TrimSpace(recordData[start:end]))
		start = end
	}
	return val.Interface(), nil
}

func (i *Importer) ImportCWR(src io.Reader) error {
	input := newInput(src)

	var work *media.Identifier
	for {
		record, err := input.Next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		switch v := record.(type) {

		case *NWR:
			var identifier media.Identifier
			if v.ISWC != "" {
				identifier.Type = "iswc"
				identifier.Value = v.ISWC
			} else {
				identifier.Type = "cwr.work_id"
				identifier.Value = v.SubmitterWorkNumber
			}
			if err := i.client.CreateWork(
				&media.Work{Title: v.WorkTitle},
				&identifier,
			); err != nil {
				return err
			}
			work = &identifier

		case *REV:
			var identifier media.Identifier
			if v.ISWC != "" {
				identifier.Type = "iswc"
				identifier.Value = v.ISWC
			} else {
				identifier.Type = "cwr.work_id"
				identifier.Value = v.SubmitterWorkNumber
			}
			if err := i.client.CreateWork(
				&media.Work{Title: v.WorkTitle},
				&identifier,
			); err != nil {
				return err
			}
			work = &identifier

		case *SPU:
			var identifier media.Identifier
			if v.PublisherIPINameNumber != "" {
				identifier.Type = "ipi"
				identifier.Value = v.PublisherIPINameNumber
			} else {
				identifier.Type = "cwr.publisher_id"
				identifier.Value = v.InterestedPartyNumber
			}
			if err := i.client.CreatePublisher(
				&media.Publisher{Name: v.PublisherName},
				&identifier,
			); err != nil {
				return err
			}

			link := &media.PublisherWorkLink{
				Publisher: identifier,
				Work:      *work,
				Role:      publisherRole(v.PublisherType),
				PerformanceRightsShare:     v.PRShare,
				MechanicalRightsShare:      v.MRShare,
				SynchronizationRightsShare: v.SRShare,
			}
			if err := i.client.CreatePublisherWorkLink(link); err != nil {
				return err
			}

		case *OPU:
			var identifier media.Identifier
			if v.PublisherIPINameNumber != "" {
				identifier.Type = "ipi"
				identifier.Value = v.PublisherIPINameNumber
			} else {
				identifier.Type = "cwr.publisher_id"
				identifier.Value = v.InterestedPartyNumber
			}
			if err := i.client.CreatePublisher(
				&media.Publisher{Name: v.PublisherName},
				&identifier,
			); err != nil {
				return err
			}

			link := &media.PublisherWorkLink{
				Publisher: identifier,
				Work:      *work,
				Role:      publisherRole(v.PublisherType),
				PerformanceRightsShare:     v.PRShare,
				MechanicalRightsShare:      v.MRShare,
				SynchronizationRightsShare: v.SRShare,
			}
			if err := i.client.CreatePublisherWorkLink(link); err != nil {
				return err
			}

		case *SWR:
			var identifier media.Identifier
			if v.WriterIPINameNumber != "" {
				identifier.Type = "ipi"
				identifier.Value = v.WriterIPINameNumber
			} else {
				identifier.Type = "cwr.writer_id"
				identifier.Value = v.InterestedPartyNumber
			}
			if err := i.client.CreateComposer(
				&media.Composer{
					FirstName: v.WriterFirstName,
					LastName:  v.WriterLastName,
				},
				&identifier,
			); err != nil {
				return err
			}

			link := &media.ComposerWorkLink{
				Composer: identifier,
				Work:     *work,
				Role:     writerRole(v.WriterDesignationCode),
				PerformanceRightsShare:     v.PRShare,
				MechanicalRightsShare:      v.MRShare,
				SynchronizationRightsShare: v.SRShare,
			}
			if err := i.client.CreateComposerWorkLink(link); err != nil {
				return err
			}

		case *OWR:
			var identifier media.Identifier
			if v.WriterIPINameNumber != "" {
				identifier.Type = "ipi"
				identifier.Value = v.WriterIPINameNumber
			} else {
				identifier.Type = "cwr.writer_id"
				identifier.Value = v.InterestedPartyNumber
			}
			if err := i.client.CreateComposer(
				&media.Composer{
					FirstName: v.WriterFirstName,
					LastName:  v.WriterLastName,
				},
				&identifier,
			); err != nil {
				return err
			}

			link := &media.ComposerWorkLink{
				Composer: identifier,
				Work:     *work,
				Role:     writerRole(v.WriterDesignationCode),
				PerformanceRightsShare:     v.PRShare,
				MechanicalRightsShare:      v.MRShare,
				SynchronizationRightsShare: v.SRShare,
			}
			if err := i.client.CreateComposerWorkLink(link); err != nil {
				return err
			}

		}
	}

	return nil
}

func publisherRole(publisherType string) string {
	switch publisherType {
	case "AQ":
		return "Acquirer"
	case "AM":
		return "Administrator"
	case "PA":
		return "Income Participant"
	case "E":
		return "Original Publisher"
	case "ES":
		return "Substituted Publisher"
	case "SE":
		return "Sub Publisher"
	default:
		return ""
	}
}

func writerRole(writerDesignationCode string) string {
	switch writerDesignationCode {
	case "AD":
		return "Adaptor"
	case "AR":
		return "Arranger"
	case "A":
		return "Author"
	case "C", "CA":
		return "Composer"
	case "SR":
		return "Sub Arranger"
	case "SA":
		return "Sub Author"
	case "TR":
		return "Translator"
	case "PA":
		return "Income Participant"
	default:
		return ""
	}
}
