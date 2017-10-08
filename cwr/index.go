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
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// Indexer is a META indexer which indexes a stream of META objects
// representing records (represented in cwr files) into a SQLite3 database, getting the
// associated META objects from a META store.
type Indexer struct {
	indexDB *sql.DB
	sqlTx   *sql.Tx
	store   *meta.Store
}

type jobIn struct {
	cwrID   *cid.Cid
	tx      map[string]interface{}
	indexFn func(cwrID *cid.Cid, tx map[string]interface{}) error
}

// NewIndexer returns an Indexer which updates the indexes in the given SQLite3
// database connection, getting META objects from the given META store.
func NewIndexer(indexDB *sql.DB, store *meta.Store) (*Indexer, error) {
	// migrate the db to ensure it has an up-to-date schema
	if err := migrations.Run(indexDB); err != nil {
		return nil, err
	}

	return &Indexer{
		indexDB: indexDB,
		store:   store,
	}, nil
}

// Index indexes a stream of META object links which are expected to
// point at CWRs.
func (i *Indexer) Index(ctx context.Context, stream chan *cid.Cid) error {
	for {
		select {
		case cid, ok := <-stream:
			if !ok {
				return nil
			}
			obj, err := i.store.Get(cid)
			if err != nil {
				return err
			}
			i.sqlTx, err = i.indexDB.Begin()
			if err != nil {
				return err
			}
			if err := i.index(obj); err != nil {
				i.sqlTx.Rollback()
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// index indexes a CWR based on its NWR,REV
func (i *Indexer) index(cwr *meta.Object) (err error) {
	graph := meta.NewGraph(i.store, cwr)
	jobs := make(chan jobIn)
	results := make(chan error)
	wg := new(sync.WaitGroup)

	for field, indexFn := range map[string]func(*cid.Cid, *meta.Object) error{
		"HDR": i.indexTransmissionHeader,
		"TRL": i.indexTransmissionTrailer,
	} {
		v, err := graph.Get("Records", field)
		if meta.IsPathNotFound(err) {
			continue
		} else if err != nil {
			return err
		}
		id, ok := v.(*cid.Cid)
		if !ok {
			return fmt.Errorf("unexpected field type for %q, expected *cid.Cid, got %T", field, v)
		}
		if err := i.indexRecord(cwr.Cid(), id, indexFn); err != nil {
			return err
		}
	}
	v, err := graph.Get("Groups")
	if err != nil {
		return err
	}

	numberOfGroups := len(v.([]interface{}))

	for w := 1; w <= concurrentWorkNum; w++ {
		wg.Add(1)
		go i.worker(jobs, results, wg)
	}
	go func() {
		for field, indexFn := range map[string]func(cwrID *cid.Cid, tx map[string]interface{}) error{
			"NWR": i.indexNWR,
			"REV": i.indexNWR,
			"ISW": i.indexISW,
			"EXC": i.indexEXC,
			"AGR": i.indexAGR,
			"ACK": i.indexACK,
		} {

			for k := 0; k < numberOfGroups; k++ {

				v, err := graph.Get("Groups", strconv.Itoa(k), "GRH")
				if meta.IsPathNotFound(err) {
					continue
				} else if err != nil {
					results <- err
					return
				}

				id, ok := v.(*cid.Cid)
				if !ok {
					results <- fmt.Errorf("unexpected field type for %q, expected *cid.Cid, got %T", field, v)
				}
				if err := i.indexRecord(cwr.Cid(), id, i.indexGroupHeader); err != nil {
					results <- err
					return
				}

				v, err = graph.Get("Groups", strconv.Itoa(k), "Transactions", field)
				if meta.IsPathNotFound(err) {
					continue
				} else if err != nil {
					results <- err
					return
				}
				numberOfTx := len(v.([]interface{}))

				for j := 0; j < numberOfTx; j++ {
					v, err := graph.Get("Groups", strconv.Itoa(k), "Transactions", field, strconv.Itoa(j))
					if meta.IsPathNotFound(err) {
						continue
					} else if err != nil {
						results <- err
						return
					}
					tx, ok := v.(map[string]interface{})
					if !ok {
						results <- err
						return
					}
					jobs <- jobIn{cwr.Cid(), tx, indexFn}
				}
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result != nil {
			err = result
			continue //continue to drain results channel.
		}
	}
	if err != nil {
		return err
	}
	return i.sqlTx.Commit()
}

func (i *Indexer) worker(jobs <-chan jobIn, results chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		results <- job.indexFn(job.cwrID, job.tx)
	}
}

// indexRecord indexes a particular CWR record using the provided index
// function.
func (i *Indexer) indexRecord(cwrID, cid *cid.Cid, indexFn func(*cid.Cid, *meta.Object) error) error {
	obj, err := i.store.Get(cid)
	if err != nil {
		return err
	}
	return indexFn(cwrID, obj)
}

// indexRegisteredWork indexes the given registeredWork record on its title,iswc,CompositeType and record_type
// properties.
func (i *Indexer) indexRegisteredWork(cwrID *cid.Cid, obj *meta.Object) error {
	registeredWork := &RegisteredWork{}

	if err := obj.Decode(registeredWork); err != nil {
		return err
	}

	log.Info("indexing nwr (registered work)", "object_id", obj.Cid().String(), "Title", registeredWork.Title, "ISWC", registeredWork.ISWC, "Composite Type", registeredWork.CompositeType, "Record Type", registeredWork.RecordType)
	stmt, err := i.sqlTx.Prepare(`INSERT INTO registered_work (cwr_id,object_id, title, iswc, composite_type,record_type) VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cwrID.String(), obj.Cid().String(), registeredWork.Title, registeredWork.ISWC, registeredWork.CompositeType, registeredWork.RecordType)
	return err
}

// indexTransmissionHeader indexes the given transmission header (HDR) record on its sender_type,sender_id,sender_name and record_type
// properties.
func (i *Indexer) indexTransmissionHeader(cwrID *cid.Cid, hdr *meta.Object) error {
	transmissionHeader := &TransmissionHeader{}

	if err := hdr.Decode(transmissionHeader); err != nil {
		return err
	}
	log.Info("indexing cwr transmission  header", "Sender  Type", transmissionHeader.SenderType, "Sender Id", transmissionHeader.SenderID, "Record Type", transmissionHeader.RecordType)

	stmt, err := i.sqlTx.Prepare(`INSERT INTO transmission_header (cwr_id,object_id,sender_type,sender_id,sender_name,record_type) VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cwrID.String(), hdr.Cid().String(), transmissionHeader.SenderType, transmissionHeader.SenderID, transmissionHeader.SenderName, transmissionHeader.RecordType)

	return err
}

// indexTransmissionTrailer indexes the given TRL record
func (i *Indexer) indexTransmissionTrailer(cwrID *cid.Cid, obj *meta.Object) error {
	// TODO: index TRL
	return nil
}

// indexGroupHeader indexes the given GRH record
func (i *Indexer) indexGroupHeader(cwrID *cid.Cid, grhr *meta.Object) error {

	// TODO: index GRH
	return nil
}

// indexPublisherControlledBySubmiter indexes the given SPU record on its publisher_sequence_n and record_type
// properties.
func (i *Indexer) indexPublisherControlledBySubmiter(cwrID *cid.Cid, txCid *cid.Cid, obj *meta.Object) error {
	publisherControlledBySubmitter := &PublisherControllBySubmitter{}

	if err := obj.Decode(publisherControlledBySubmitter); err != nil {
		return err
	}
	log.Info("indexing publisherControlledBySubmitter ", "object_id", obj.Cid().String(), "publisher_sequence_n", publisherControlledBySubmitter.PublisherSequenceNumber, "Record Type", publisherControlledBySubmitter.RecordType)

	stmt, err := i.sqlTx.Prepare(`INSERT INTO publisher_control (cwr_id,tx_id,object_id, publisher_sequence_n, record_type) VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cwrID.String(), txCid.String(), obj.Cid().String(), publisherControlledBySubmitter.PublisherSequenceNumber, publisherControlledBySubmitter.RecordType)
	return err
}

// indexNWR indexes the given cwr transaction by indexing each transacation's record and link it to its
// transaction.
func (i *Indexer) indexNWR(cwrID *cid.Cid, tx map[string]interface{}) error {
	mainRecordTx, ok := tx["MainRecord"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("error indexing CWR: expected MainRecord property to be map[string]interface{}, got %T", tx["MainRecord"])
	}
	nwrCid, ok := mainRecordTx["NWR"].(*cid.Cid)
	if !ok {
		nwrCid, ok = mainRecordTx["REV"].(*cid.Cid)
		if !ok {
			return fmt.Errorf("error indexing CWR: expected REV property to be *cid.Cid, got %T", mainRecordTx["REV"])
		}
	}
	obj, err := i.store.Get(nwrCid)
	if err != nil {
		return err
	}

	if err := i.indexRegisteredWork(cwrID, obj); err != nil {
		return err
	}

	for _, spuCid := range tx["DetailRecords"].(map[string]interface{})["SPU"].([]interface{}) {
		obj, err := i.store.Get(spuCid.(*cid.Cid))
		if err != nil {
			return err
		}
		if err := i.indexPublisherControlledBySubmiter(cwrID, nwrCid, obj); err != nil {
			return err
		}
	}
	return nil
}

func (i *Indexer) indexISW(cwrID *cid.Cid, tx map[string]interface{}) error {
	// TODO: index ISW
	return nil
}

func (i *Indexer) indexACK(cwrID *cid.Cid, tx map[string]interface{}) error {
	// TODO: index ACK
	return nil
}

func (i *Indexer) indexEXC(cwrID *cid.Cid, tx map[string]interface{}) error {
	// TODO: index EXC
	return nil
}

func (i *Indexer) indexAGR(cwrID *cid.Cid, tx map[string]interface{}) error {
	// TODO: index AGR
	return nil
}
