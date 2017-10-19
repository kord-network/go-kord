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
	index *meta.Index
	store *meta.Store
}

type jobIn struct {
	tx      *sql.Tx
	cwrID   *cid.Cid
	txs     map[string]interface{}
	indexFn func(*sql.Tx, *cid.Cid, map[string]interface{}) error
}

// NewIndexer returns an Indexer which updates the indexes in the given SQLite3
// database connection, getting META objects from the given META store.
func NewIndexer(index *meta.Index, store *meta.Store) (*Indexer, error) {
	// migrate the db to ensure it has an up-to-date schema
	if err := migrations.Run(index.DB); err != nil {
		return nil, err
	}

	return &Indexer{
		index: index,
		store: store,
	}, nil
}

// Index indexes a stream of META object links which are expected to
// point at CWRs.
func (i *Indexer) Index(ctx context.Context, stream *meta.StreamReader) error {
	return i.index.Update(func(tx *sql.Tx) error {
		for {
			select {
			case cid, ok := <-stream.Ch():
				if !ok {
					return stream.Err()
				}
				obj, err := i.store.Get(cid)
				if err != nil {
					return err
				}
				if err := i.indexCWR(tx, obj); err != nil {
					return err
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
}

// index indexes a CWR based on its NWR,REV
func (i *Indexer) indexCWR(tx *sql.Tx, cwr *meta.Object) (err error) {
	graph := meta.NewGraph(i.store, cwr)
	jobs := make(chan jobIn)
	results := make(chan error)

	for field, indexFn := range map[string]func(*sql.Tx, *cid.Cid, *meta.Object) error{
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
		if err := i.indexRecord(tx, cwr.Cid(), id, indexFn); err != nil {
			return err
		}
	}
	v, err := graph.Get("Groups")
	if err != nil {
		return err
	}

	numberOfGroups := len(v.([]interface{}))

	var wg sync.WaitGroup
	wg.Add(concurrentWorkNum)
	for w := 1; w <= concurrentWorkNum; w++ {
		go func() {
			defer wg.Done()
			i.worker(jobs, results)
		}()
	}
	go func() (err error) {
		defer func() {
			results <- err
		}()
		for field, indexFn := range map[string]func(*sql.Tx, *cid.Cid, map[string]interface{}) error{
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
					return err
				}

				id, ok := v.(*cid.Cid)
				if !ok {
					return fmt.Errorf("unexpected field type for %q, expected *cid.Cid, got %T", field, v)
				}
				if err := i.indexRecord(tx, cwr.Cid(), id, i.indexGroupHeader); err != nil {
					return err
				}
				v, err = graph.Get("Groups", strconv.Itoa(k), "Transactions", field)
				if meta.IsPathNotFound(err) {
					continue
				} else if err != nil {
					return err
				}
				numberOfTx := len(v.([]interface{}))

				for j := 0; j < numberOfTx; j++ {
					v, err := graph.Get("Groups", strconv.Itoa(k), "Transactions", field, strconv.Itoa(j))
					if meta.IsPathNotFound(err) {
						continue
					} else if err != nil {
						return err
					}
					txs, ok := v.(map[string]interface{})
					if !ok {
						return fmt.Errorf("unexpected field type .Expected map[string]interface{}, got %T", v)
					}
					jobs <- jobIn{tx, cwr.Cid(), txs, indexFn}
				}
			}
		}
		close(jobs)
		return nil
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
	return err
}

func (i *Indexer) worker(jobs <-chan jobIn, results chan<- error) {
	for job := range jobs {
		results <- job.indexFn(job.tx, job.cwrID, job.txs)
	}
}

// indexRecord indexes a particular CWR record using the provided index
// function.
func (i *Indexer) indexRecord(tx *sql.Tx, cwrID, cid *cid.Cid, indexFn func(*sql.Tx, *cid.Cid, *meta.Object) error) error {
	obj, err := i.store.Get(cid)
	if err != nil {
		return err
	}
	return indexFn(tx, cwrID, obj)
}

// indexRegisteredWork indexes the given registeredWork record on its title,iswc,CompositeType and record_type
// properties.
func (i *Indexer) indexRegisteredWork(tx *sql.Tx, cwrID *cid.Cid, obj *meta.Object) error {
	registeredWork := &RegisteredWork{}

	if err := obj.Decode(registeredWork); err != nil {
		return err
	}

	log.Info("indexing nwr (registered work)", "object_id", obj.Cid().String(), "Title", registeredWork.Title, "ISWC", registeredWork.ISWC, "Composite Type", registeredWork.CompositeType, "Record Type", registeredWork.RecordType)

	_, err := tx.Exec(`INSERT INTO registered_work (cwr_id,object_id, title, iswc, composite_type,record_type) VALUES ($1, $2, $3, $4, $5, $6)`,
		cwrID.String(), obj.Cid().String(), registeredWork.Title, registeredWork.ISWC, registeredWork.CompositeType, registeredWork.RecordType)
	return err
}

// indexTransmissionHeader indexes the given transmission header (HDR) record on its sender_type,sender_id,sender_name and record_type
// properties.
func (i *Indexer) indexTransmissionHeader(tx *sql.Tx, cwrID *cid.Cid, hdr *meta.Object) error {
	transmissionHeader := &TransmissionHeader{}

	if err := hdr.Decode(transmissionHeader); err != nil {
		return err
	}
	log.Info("indexing cwr transmission  header", "Sender  Type", transmissionHeader.SenderType, "Sender Id", transmissionHeader.SenderID, "Record Type", transmissionHeader.RecordType)

	_, err := tx.Exec(`INSERT INTO transmission_header (cwr_id,object_id,sender_type,sender_id,sender_name,record_type) VALUES ($1, $2, $3, $4, $5, $6)`,
		cwrID.String(), hdr.Cid().String(), transmissionHeader.SenderType, transmissionHeader.SenderID, transmissionHeader.SenderName, transmissionHeader.RecordType)
	return err
}

// indexTransmissionTrailer indexes the given TRL record
func (i *Indexer) indexTransmissionTrailer(tx *sql.Tx, cwrID *cid.Cid, obj *meta.Object) error {
	// TODO: index TRL
	return nil
}

// indexGroupHeader indexes the given GRH record
func (i *Indexer) indexGroupHeader(tx *sql.Tx, cwrID *cid.Cid, grhr *meta.Object) error {

	// TODO: index GRH
	return nil
}

// indexPublisherControlledBySubmiter indexes the given SPU record on its properties.
func (i *Indexer) indexPublisherControlledBySubmiter(tx *sql.Tx, cwrID *cid.Cid, txCid *cid.Cid, obj *meta.Object) error {
	publisherControlledBySubmitter := &PublisherControllBySubmitter{}

	if err := obj.Decode(publisherControlledBySubmitter); err != nil {
		return err
	}
	log.Info("indexing publisherControlledBySubmitter ", "cwr_id", cwrID.String(), "tx_id", txCid.String(), "object_id", obj.Cid().String())
	_, err := tx.Exec(`INSERT INTO publisher_control
		(cwr_id,
		 tx_id,
		 object_id,
		 publisher_sequence_n,
		 transaction_sequence_n,
		 record_sequence_n,
		 record_type,
		 interested_party_number,
		 publisher_name,
		 publisher_type,
		 publisher_ipi_name_number,
		 pr_ownership_share,
		 mr_society,
		 mr_ownership_share,
		 sr_society,
		 sr_ownership_share,
		 publisher_ipi_base_number
	 )
	   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
		cwrID.String(),
		txCid.String(),
		obj.Cid().String(),
		publisherControlledBySubmitter.PublisherSequenceNumber,
		publisherControlledBySubmitter.TransactionSequenceN,
		publisherControlledBySubmitter.RecordSequenceN,
		publisherControlledBySubmitter.RecordType,
		publisherControlledBySubmitter.InterestedPartyNumber,
		publisherControlledBySubmitter.PublisherName,
		publisherControlledBySubmitter.PublisherType,
		publisherControlledBySubmitter.PublisherIPINameNumber,
		publisherControlledBySubmitter.PROwnershipShare,
		publisherControlledBySubmitter.MRSociety,
		publisherControlledBySubmitter.MROwnershipShare,
		publisherControlledBySubmitter.SRSociety,
		publisherControlledBySubmitter.SROwnershipShare,
		publisherControlledBySubmitter.PublisherIPIBaseNumber,
	)
	return err
}

// indexWriterControlledbySubmitter indexes the given SWR or OWR record on its properties.
func (i *Indexer) indexWriterControlledbySubmitter(tx *sql.Tx, cwrID *cid.Cid, txCid *cid.Cid, obj *meta.Object) error {
	writerControlledbySubmitter := &WriterControlledbySubmitter{}

	if err := obj.Decode(writerControlledbySubmitter); err != nil {
		return err
	}
	log.Info("indexing writerControlledbySubmitter ", "cwr_id", cwrID.String(), "tx_id", txCid.String(), "object_id", obj.Cid().String())
	_, err := tx.Exec(`INSERT INTO writer_control
		( cwr_id,
			tx_id,
			object_id,
			record_type,
			transaction_sequence_n,
			record_sequence_n,
			interested_party_number,
			writer_last_name,
			writer_first_name,
			writer_ipi_name,
			writer_ipi_base_number,
			personal_number)
	   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		cwrID.String(),
		txCid.String(),
		obj.Cid().String(),
		writerControlledbySubmitter.RecordType,
		writerControlledbySubmitter.TransactionSequenceN,
		writerControlledbySubmitter.RecordSequenceN,
		writerControlledbySubmitter.InterestedPartyNumber,
		writerControlledbySubmitter.WriterLastName,
		writerControlledbySubmitter.WriterFirstName,
		writerControlledbySubmitter.WriterIPIName,
		writerControlledbySubmitter.WriterIPIBaseNumber,
		writerControlledbySubmitter.PersonalNumber,
	)
	return err
}

// indexNWR indexes the given cwr transaction by indexing each transacation's record and link it to its
// transaction.
func (i *Indexer) indexNWR(tx *sql.Tx, cwrID *cid.Cid, txs map[string]interface{}) error {
	mainRecordTx, ok := txs["MainRecord"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("error indexing CWR: expected MainRecord property to be map[string]interface{}, got %T", txs["MainRecord"])
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

	if err := i.indexRegisteredWork(tx, cwrID, obj); err != nil {
		return err
	}

	for _, spuCid := range txs["DetailRecords"].(map[string]interface{})["SPU"].([]interface{}) {
		obj, err := i.store.Get(spuCid.(*cid.Cid))
		if err != nil {
			return err
		}
		if err := i.indexPublisherControlledBySubmiter(tx, cwrID, nwrCid, obj); err != nil {
			return err
		}
	}
	for _, swrCid := range txs["DetailRecords"].(map[string]interface{})["SWR"].([]interface{}) {
		obj, err := i.store.Get(swrCid.(*cid.Cid))
		if err != nil {
			return err
		}
		if err := i.indexWriterControlledbySubmitter(tx, cwrID, nwrCid, obj); err != nil {
			return err
		}
	}
	return nil
}

func (i *Indexer) indexISW(tx *sql.Tx, cwrID *cid.Cid, txs map[string]interface{}) error {
	// TODO: index ISW
	return nil
}

func (i *Indexer) indexACK(tx *sql.Tx, cwrID *cid.Cid, txs map[string]interface{}) error {
	// TODO: index ACK
	return nil
}

func (i *Indexer) indexEXC(tx *sql.Tx, cwrID *cid.Cid, txs map[string]interface{}) error {
	// TODO: index EXC
	return nil
}

func (i *Indexer) indexAGR(tx *sql.Tx, cwrID *cid.Cid, txs map[string]interface{}) error {
	// TODO: index AGR
	return nil
}
