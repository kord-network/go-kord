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
)

func substring(s string, from int, to int) string {
	if len(s) < from || len(s) < to {
		return ""
	}
	return s[from:to]
}
func ParseCWRFile(cwrFileReader io.Reader) (records []*Record, err error) {
	scanner := bufio.NewScanner(cwrFileReader)

	for scanner.Scan() {
		var record *Record = &Record{}

		if strings.HasPrefix(scanner.Text(), "NWR") ||
			strings.HasPrefix(scanner.Text(), "REV") {
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
		}
		if strings.HasPrefix(scanner.Text(), "SPU") {
			record.RecordType = substring(scanner.Text(), 0, 3)
			record.TransactionSequenceN = substring(scanner.Text(), 3, 12)
			record.RecordSequenceN = substring(scanner.Text(), 12, 19)
			record.PublisherSequenceNumber = substring(scanner.Text(), 19, 21)
		}

		if record.RecordType != "" {
			records = append(records, record)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return
}
