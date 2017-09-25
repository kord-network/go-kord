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
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

func getSlice(s string, from int, to int) string {
	if len(s) < from || len(s) < to {
		return ""
	}
	return s[from:to]
}
func ParseCWRFile(cwrFileReader io.Reader) (records []*Record, err error) {
	scanner := bufio.NewScanner(cwrFileReader)

	for scanner.Scan() {
		var record *Record
		var recordBytes []byte
		if strings.HasPrefix(scanner.Text(), "NWR") ||
			strings.HasPrefix(scanner.Text(), "REV") {
			registeredWork := RegisteredWork{}
			registeredWork.RecordType = getSlice(scanner.Text(), 0, 19)
			registeredWork.Title = strings.TrimSpace(getSlice(scanner.Text(), 19, 79))
			registeredWork.LanguageCode = getSlice(scanner.Text(), 79, 81)
			registeredWork.SubmitteWorkNumber = getSlice(scanner.Text(), 81, 95)
			registeredWork.ISWC = getSlice(scanner.Text(), 95, 106)
			registeredWork.CopyRightDate = getSlice(scanner.Text(), 106, 113)
			registeredWork.DistributionCategory = getSlice(scanner.Text(), 127, 129)
			registeredWork.Duration = getSlice(scanner.Text(), 129, 135)
			registeredWork.RecordedIndicator = getSlice(scanner.Text(), 135, 136)
			registeredWork.TextMusicRelationship = getSlice(scanner.Text(), 136, 139)
			registeredWork.CompositeType = getSlice(scanner.Text(), 140, 142)
			registeredWork.VersionType = getSlice(scanner.Text(), 142, 145)
			registeredWork.PriorityFlag = getSlice(scanner.Text(), 259, 260)
			recordBytes, err = json.Marshal(registeredWork)
			if err != nil {
				return nil, err
			}
		}
		if strings.HasPrefix(scanner.Text(), "SPU") {
			publisherControllBySubmitter := PublisherControllBySubmitter{}
			publisherControllBySubmitter.RecordType = getSlice(scanner.Text(), 0, 19)
			publisherControllBySubmitter.PublisherSequenceNumber = getSlice(scanner.Text(), 19, 21)
			recordBytes, err = json.Marshal(publisherControllBySubmitter)
			if err != nil {
				return nil, err
			}
		}
		if len(recordBytes) > 0 {
			if err := json.NewDecoder(bytes.NewReader(recordBytes)).Decode(&record); err != nil {
				return nil, err
			}
			records = append(records, record)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return
}
