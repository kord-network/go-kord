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
	"io"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// Converter converts CWR data from a CWR file to META
// objects.
type Converter struct {
	store *meta.Store
}

// NewConverter returns a Converter which reads data from the given CWR io.Reader
// and stores META object in the given META store.
func NewConverter(store *meta.Store) *Converter {

	return &Converter{
		store: store,
	}
}

// ConvertCWR converts the given source CWR file into a META object graph and
// returns the CID of the graph's root META object.
func (c *Converter) ConvertCWR(cwrFileReader io.Reader) (*cid.Cid, error) {
	// get records from the db
	cwr := make(map[string]interface{})
	nwr := make(map[string]interface{})
	var txs map[string][]map[string]interface{}
	var groups []map[string][]map[string]interface{}
	var spus []*cid.Cid

	scanner := bufio.NewScanner(cwrFileReader)

	for scanner.Scan() {
		var record = &Record{}
		switch substring(scanner.Text(), 0, 3) {
		case "HDR":
			record.RecordType = substring(scanner.Text(), 0, 3)
			record.SenderType = substring(scanner.Text(), 3, 6)
			record.SenderID = substring(scanner.Text(), 6, 14)
			record.SenderName = strings.TrimSpace(substring(scanner.Text(), 14, 59))
		case "GRH":
			record.RecordType = substring(scanner.Text(), 0, 3)
			record.TransactionType = substring(scanner.Text(), 3, 6)
			record.GroupID = substring(scanner.Text(), 6, 11)
		case "GRT":
			record.RecordType = substring(scanner.Text(), 0, 3)
			record.GroupID = substring(scanner.Text(), 3, 8)
		case "TRL":
			record.RecordType = substring(scanner.Text(), 0, 3)
		case "NWR", "REV":
			record.RecordType = substring(scanner.Text(), 0, 3)
			record.TransactionSequenceN = substring(scanner.Text(), 3, 12)
			record.RecordSequenceN = substring(scanner.Text(), 12, 19)
			record.Title = strings.TrimSpace(substring(scanner.Text(), 19, 79))
			record.LanguageCode = substring(scanner.Text(), 79, 81)
			record.SubmitteWorkNumber = substring(scanner.Text(), 81, 95)
			record.ISWC = substring(scanner.Text(), 95, 106)
			record.CopyRightDate = substring(scanner.Text(), 106, 113)
			record.DistributionCategory = substring(scanner.Text(), 127, 129)
			record.Duration = substring(scanner.Text(), 129, 135)
			record.RecordedIndicator = substring(scanner.Text(), 135, 136)
			record.TextMusicRelationship = substring(scanner.Text(), 136, 139)
			record.CompositeType = substring(scanner.Text(), 140, 142)
			record.VersionType = substring(scanner.Text(), 142, 145)
			record.PriorityFlag = substring(scanner.Text(), 259, 260)
		case "SPU":
			record.RecordType = substring(scanner.Text(), 0, 3)
			record.TransactionSequenceN = substring(scanner.Text(), 3, 12)
			record.RecordSequenceN = substring(scanner.Text(), 12, 19)
			record.PublisherSequenceNumber = substring(scanner.Text(), 19, 21)
		default:
		}

		if record.RecordType != "" {

			obj, err := meta.Encode(record)
			if err != nil {
				return nil, err
			}
			if err := c.store.Put(obj); err != nil {
				return nil, err
			}
			switch record.RecordType {
			case "HDR", "TRL":
				cwr[record.RecordType] = obj.Cid()
			case "GRH":
				txs = make(map[string][]map[string]interface{})
			case "GRT":
				nwr["SPU"] = spus // accumulate last spus
				spus = []*cid.Cid{}
				txs["NWR"] = append(txs["NWR"], nwr) //accumulate last transaction
				nwr = make(map[string]interface{})
				if len(txs) > 0 { // accumulate txs and continue
					groups = append(groups, txs)
				}
			case "NWR", "REV":
				if nwr[record.RecordType] != nil { //accumulate the current nwr and continue
					nwr["SPU"] = spus // accumulate spus
					spus = []*cid.Cid{}
					txs["NWR"] = append(txs["NWR"], nwr)
					nwr = make(map[string]interface{}) //re initilize nwr

				}
				nwr[record.RecordType] = obj.Cid()
			case "SPU": // NWR/REV transaction records
				spus = append(spus, obj.Cid())
			}
		}
	}
	cwr["GRH"] = groups

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	obj, err := meta.Encode(cwr)
	if err != nil {
		return nil, err
	}
	if err := c.store.Put(obj); err != nil {
		return nil, err
	}

	return obj.Cid(), nil
}

func substring(s string, from int, to int) string {
	if len(s) < from || len(s) < to {
		return ""
	}
	return s[from:to]
}
