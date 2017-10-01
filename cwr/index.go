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

	"github.com/ethereum/go-ethereum/log"
	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// Indexer is a META indexer which indexes a stream of META objects
// representing records (represented in cwr files) into a SQLite3 database, getting the
// associated META objects from a META store.
type Indexer struct {
	indexDB *sql.DB
	store   *meta.Store
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
			if err := i.index(obj); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// index indexes a CWR based on its NWR,REV
func (i *Indexer) index(cwr *meta.Object) error {
	graph := meta.NewGraph(i.store, cwr)

	for field, indexFn := range map[string]func(*cid.Cid, *meta.Object) error{
		"HDR": i.indexTransmissionHeader,
		"TRL": i.indexTransmissionTrailer,
	} {
		v, err := graph.Get(field)
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
	v, err := graph.Get("GRH")
	if err != nil {
		return err
	}

	numberOfGroups := len(v.([]interface{}))

	for field, indexFn := range map[string]func(cwrID *cid.Cid, tx map[string]interface{}) error{
		"NWR": i.indexNWR,
		"REV": i.indexNWR,
		"ISW": i.indexISW,
		"EXC": i.indexEXC,
		"AGR": i.indexAGR,
		"ACK": i.indexACK,
	} {

		for k := 0; k < numberOfGroups; k++ {
			v, err := graph.Get("GRH", strconv.Itoa(k), field)
			if meta.IsPathNotFound(err) {
				continue
			} else if err != nil {
				return err
			}
			numberOfTx := len(v.([]interface{}))

			for j := 0; j < numberOfTx; j++ {
				v, err := graph.Get("GRH", strconv.Itoa(k), field, strconv.Itoa(j))
				if meta.IsPathNotFound(err) {
					continue
				} else if err != nil {
					return err
				}
				tx, ok := v.(map[string]interface{})
				if !ok {
					return fmt.Errorf("unexpected field type for %q, expected *cid.Cid, got %T", field, v)
				}
				if err := indexFn(cwr.Cid(), tx); err != nil {
					return err
				}
			}
		}
	}
	return nil
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

	_, err := i.indexDB.Exec(
		`INSERT INTO registered_work (cwr_id,object_id, title, iswc, composite_type,record_type) VALUES ($1, $2, $3, $4, $5, $6)`,
		cwrID.String(), obj.Cid().String(), registeredWork.Title, registeredWork.ISWC, registeredWork.CompositeType, registeredWork.RecordType,
	)
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
	_, err := i.indexDB.Exec(
		`INSERT INTO transmission_header (cwr_id,object_id,sender_type,sender_id,sender_name,record_type) VALUES ($1, $2, $3, $4, $5, $6)`,
		cwrID.String(), hdr.Cid().String(), transmissionHeader.SenderType, transmissionHeader.SenderID, transmissionHeader.SenderName, transmissionHeader.RecordType,
	)
	return err
}

// indexTransmissionTrailer indexes the given TRL record
func (i *Indexer) indexTransmissionTrailer(cwrID *cid.Cid, obj *meta.Object) error {
	// TODO: index TRL
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

	_, err := i.indexDB.Exec(
		`INSERT INTO publisher_control (cwr_id,tx_id,object_id, publisher_sequence_n, record_type) VALUES ($1, $2, $3, $4, $5)`,
		cwrID.String(), txCid.String(), obj.Cid().String(), publisherControlledBySubmitter.PublisherSequenceNumber, publisherControlledBySubmitter.RecordType,
	)
	return err
}

// indexNWR indexes the given cwr transaction by indexing each transacation's record and link it to its
// transaction.
func (i *Indexer) indexNWR(cwrID *cid.Cid, tx map[string]interface{}) error {
	nwrCid, ok := tx["NWR"].(*cid.Cid)
	if !ok {
		nwrCid, ok = tx["REV"].(*cid.Cid)
		if !ok {
			return fmt.Errorf("indexNWR :cannot cast tx ")
		}
	}
	obj, err := i.store.Get(nwrCid)
	if err != nil {
		return err
	}

	if err := i.indexRegisteredWork(cwrID, obj); err != nil {
		return err
	}

	for _, spuCid := range tx["SPU"].([]interface{}) {
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
