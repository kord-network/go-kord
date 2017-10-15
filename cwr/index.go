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
			if err := i.sqlTx.Commit(); err != nil {
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
					return err
				}

				id, ok := v.(*cid.Cid)
				if !ok {
					return fmt.Errorf("unexpected field type for %q, expected *cid.Cid, got %T", field, v)
				}
				if err := i.indexRecord(cwr.Cid(), id, i.indexGroupHeader); err != nil {
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
					tx, ok := v.(map[string]interface{})
					if !ok {
						return fmt.Errorf("unexpected field type .Expected map[string]interface{}, got %T", v)
					}
					jobs <- jobIn{cwr.Cid(), tx, indexFn}
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

	_, err := i.sqlTx.Exec(`INSERT INTO registered_work (cwr_id,object_id, title, iswc, composite_type,record_type) VALUES ($1, $2, $3, $4, $5, $6)`,
		cwrID.String(), obj.Cid().String(), registeredWork.Title, registeredWork.ISWC, registeredWork.CompositeType, registeredWork.RecordType)
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

	_, err := i.sqlTx.Exec(`INSERT INTO transmission_header (cwr_id,object_id,sender_type,sender_id,sender_name,record_type) VALUES ($1, $2, $3, $4, $5, $6)`,
		cwrID.String(), hdr.Cid().String(), transmissionHeader.SenderType, transmissionHeader.SenderID, transmissionHeader.SenderName, transmissionHeader.RecordType)
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

// indexPublisherControlledBySubmiter indexes the given SPU record on its properties.
func (i *Indexer) indexPublisherControlledBySubmiter(cwrID *cid.Cid, txCid *cid.Cid, obj *meta.Object) error {
	publisherControlledBySubmitter := &PublisherControllBySubmitter{}

	if err := obj.Decode(publisherControlledBySubmitter); err != nil {
		return err
	}
	log.Info("indexing publisherControlledBySubmitter ", "cwr_id", cwrID.String(), "tx_id", txCid.String(), "object_id", obj.Cid().String())
	_, err := i.sqlTx.Exec(`INSERT INTO publisher_control
		(cwr_id,
		 tx_id,
		 object_id,
		 publisher_sequence_n,
		 transaction_sequence_n,
		 record_sequence_n,
		 record_type,
		 interested_party_number,
		 publisher_name,
		 publisher_unknown_indicator,
		 publisher_type,
		 tax_id_number,
		 publisher_ipi_name_number,
		 pr_affiliation_society_number,
		 submitter_agreement_number,
		 pr_ownership_share,
		 mr_society,
		 mr_ownership_share,
		 sr_society,
		 sr_ownership_share,
		 special_agreements_indicator,
		 first_recording_refusal_ind,
		 publisher_ipi_base_number,
		 inter_standard_agreement_code,
		 society_assigned_agreement_number,
		 agreement_type,
		 usa_license_ind)
	   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)`,
		cwrID.String(),
		txCid.String(),
		obj.Cid().String(),
		publisherControlledBySubmitter.PublisherSequenceNumber,
		publisherControlledBySubmitter.TransactionSequenceN,
		publisherControlledBySubmitter.RecordSequenceN,
		publisherControlledBySubmitter.RecordType,
		publisherControlledBySubmitter.InterestedPartyNumber,
		publisherControlledBySubmitter.PublisherName,
		publisherControlledBySubmitter.PublisherUnknownIndicator,
		publisherControlledBySubmitter.PublisherType,
		publisherControlledBySubmitter.TaxIDNumber,
		publisherControlledBySubmitter.PublisherIPINameNumber,
		publisherControlledBySubmitter.PRAffiliationSocietyNumber,
		publisherControlledBySubmitter.SubmitterAgreementNumber,
		publisherControlledBySubmitter.PROwnershipShare,
		publisherControlledBySubmitter.MRSociety,
		publisherControlledBySubmitter.MROwnershipShare,
		publisherControlledBySubmitter.SRSociety,
		publisherControlledBySubmitter.SROwnershipShare,
		publisherControlledBySubmitter.SpecialAgreementsIndicator,
		publisherControlledBySubmitter.FirstRecordingRefusalInd,
		publisherControlledBySubmitter.PublisherIPIBaseNumber,
		publisherControlledBySubmitter.InterStandardAgreementCode,
		publisherControlledBySubmitter.SocietyAssignedAgreementNumber,
		publisherControlledBySubmitter.AgreementType,
		publisherControlledBySubmitter.USALicenseInd,
	)
	return err
}

// indexWriterControlledbySubmitter indexes the given SWR or OWR record on its properties.
func (i *Indexer) indexWriterControlledbySubmitter(cwrID *cid.Cid, txCid *cid.Cid, obj *meta.Object) error {
	writerControlledbySubmitter := &WriterControlledbySubmitter{}

	if err := obj.Decode(writerControlledbySubmitter); err != nil {
		return err
	}
	log.Info("indexing writerControlledbySubmitter ", "cwr_id", cwrID.String(), "tx_id", txCid.String(), "object_id", obj.Cid().String())
	_, err := i.sqlTx.Exec(`INSERT INTO writer_control
		( cwr_id,
			tx_id,
			object_id,
			record_type,
			transaction_sequence_n,
			record_sequence_n,
			interested_party_number,
			writer_last_name,
			writer_first_name,
			writer_unknown_indicator,
			writer_designation_code,
			tax_id_number,
			writer_ipi_name,
			writer_ipi_base_number,
			personal_number)
	   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		cwrID.String(),
		txCid.String(),
		obj.Cid().String(),
		writerControlledbySubmitter.RecordType,
		writerControlledbySubmitter.TransactionSequenceN,
		writerControlledbySubmitter.RecordSequenceN,
		writerControlledbySubmitter.InterestedPartyNumber,
		writerControlledbySubmitter.WriterLastName,
		writerControlledbySubmitter.WriterFirstName,
		writerControlledbySubmitter.WriterUnknownIndicator,
		writerControlledbySubmitter.WriterDesignationCode,
		writerControlledbySubmitter.TaxIDNumber,
		writerControlledbySubmitter.WriterIPIName,
		writerControlledbySubmitter.WriterIPIBaseNumber,
		writerControlledbySubmitter.PersonalNumber,
	)
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
	for _, swrCid := range tx["DetailRecords"].(map[string]interface{})["SWR"].([]interface{}) {
		obj, err := i.store.Get(swrCid.(*cid.Cid))
		if err != nil {
			return err
		}
		if err := i.indexWriterControlledbySubmitter(cwrID, nwrCid, obj); err != nil {
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
