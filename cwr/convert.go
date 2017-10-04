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
	"sync"

	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// Converter converts CWR data from a CWR file to META
// objects.
type Converter struct {
	store *meta.Store
}

type recordJob struct {
	record *Record
	index  int
}

type objectResult struct {
	obj   *meta.Object
	index int
	err   error
}

var concurrentWorkNum = 16

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

	jobs := make(chan recordJob)
	results := make(chan objectResult)
	recordObjs := make(map[int]*meta.Object)
	cwr := make(map[string]interface{})
	nwr := make(map[string]interface{})
	var txs map[string][]map[string]interface{}
	var groups []map[string][]map[string]interface{}
	var spus []*cid.Cid

	wg := new(sync.WaitGroup)

	for w := 1; w <= concurrentWorkNum; w++ {
		wg.Add(1)
		go c.worker(jobs, results, wg)
	}

	go func() {
		scanner := bufio.NewScanner(cwrFileReader)
		index := 0
		for scanner.Scan() {
			if record := newRecord(scanner.Text()); record != nil {
				jobs <- recordJob{record, index}
				index++
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for v := range results {
		if v.err != nil {
			return nil, v.err
		}
		recordObjs[v.index] = v.obj
	}

	for i := 0; i < len(recordObjs); i++ {
		obj := recordObjs[i]
		recordType, err := obj.GetString("record_type")
		if err != nil {
			return nil, err
		}
		switch recordType {
		case "HDR", "TRL":
			cwr[recordType] = obj.Cid()
		case "GRH":
			txs = make(map[string][]map[string]interface{})
		case "GRT":
			nwr["SPU"] = spus // accumulate last spus
			spus = []*cid.Cid{}
			txs["NWR"] = append(txs["NWR"], nwr) // accumulate last transaction
			nwr = make(map[string]interface{})
			if len(txs) > 0 { //accumulate txs and continue
				groups = append(groups, txs)
			}
		case "NWR", "REV":
			if nwr[recordType] != nil { // accumulate the current nwr and continue
				nwr["SPU"] = spus // accumulate spus
				spus = []*cid.Cid{}
				txs["NWR"] = append(txs["NWR"], nwr)
				nwr = make(map[string]interface{}) //re initilize nwr
			}
			nwr[recordType] = obj.Cid()
		case "SPU": // NWR/REV transaction records
			spus = append(spus, obj.Cid())
		}
	}
	cwr["GRH"] = groups
	obj, err := meta.Encode(cwr)
	if err != nil {
		return nil, err
	}

	if err := c.store.Put(obj); err != nil {
		return nil, err
	}
	return obj.Cid(), nil
}

func newRecord(line string) *Record {
	record := &Record{}
	switch substring(line, 0, 3) {
	case "HDR":
		record.RecordType = substring(line, 0, 3)
		record.SenderType = substring(line, 3, 6)
		record.SenderID = substring(line, 6, 14)
		record.SenderName = strings.TrimSpace(substring(line, 14, 59))
	case "GRH":
		record.RecordType = substring(line, 0, 3)
		record.TransactionType = substring(line, 3, 6)
		record.GroupID = substring(line, 6, 11)
	case "GRT":
		record.RecordType = substring(line, 0, 3)
		record.GroupID = substring(line, 3, 8)
	case "TRL":
		record.RecordType = substring(line, 0, 3)
	case "NWR", "REV":
		record.RecordType = substring(line, 0, 3)
		record.TransactionSequenceN = substring(line, 3, 12)
		record.RecordSequenceN = substring(line, 12, 19)
		record.Title = strings.TrimSpace(substring(line, 19, 79))
		record.LanguageCode = substring(line, 79, 81)
		record.SubmitteWorkNumber = substring(line, 81, 95)
		record.ISWC = substring(line, 95, 106)
		record.CopyRightDate = substring(line, 106, 113)
		record.DistributionCategory = substring(line, 127, 129)
		record.Duration = substring(line, 129, 135)
		record.RecordedIndicator = substring(line, 135, 136)
		record.TextMusicRelationship = substring(line, 136, 139)
		record.CompositeType = substring(line, 140, 142)
		record.VersionType = substring(line, 142, 145)
		record.PriorityFlag = substring(line, 259, 260)
	case "SPU":
		record.RecordType = substring(line, 0, 3)
		record.TransactionSequenceN = substring(line, 3, 12)
		record.RecordSequenceN = substring(line, 12, 19)
		record.PublisherSequenceNumber = substring(line, 19, 21)
	default:
		return nil
	}
	return record
}

func (c *Converter) worker(jobs <-chan recordJob, results chan<- objectResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		if job.record.RecordType != "" {
			obj, err := meta.Encode(job.record)
			if err != nil {
				results <- objectResult{nil, 0, err}
			}
			if err := c.store.Put(obj); err != nil {
				results <- objectResult{nil, 0, err}
			}
			results <- objectResult{obj, job.index, nil}
		}
	}
}

func substring(s string, from int, to int) string {
	if len(s) < from || len(s) < to {
		return ""
	}
	return s[from:to]
}
