// This file is part of the go-meta library.
//
// Copyright (C) 2018 JAAK MUSIC LTD
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

package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cayleygraph/cayley/clog"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/log"
	cayleysql "github.com/cayleygraph/cayley/graph/sql"
	"github.com/cayleygraph/cayley/quad"
	"github.com/lib/pq"
	sqlite3 "github.com/mattn/go-sqlite3"
)

func init() {
	// TODO: Update OpIsTrue to handle the fact that SQLite does not have
	//       a built-in 'true' literal
	//
	// TODO: Add a SQLite regex extension so that the REGEXP operator works
	//
	// TODO: Update cayley so that it supports creating indexes in the
	//       CREATE TABLE statement (currently we are just not creating
	//       the indexes by setting NoForeignKeys)
	cayleysql.Register("meta", cayleysql.Registration{
		Driver:               "meta",
		HashType:             `BLOB`,
		BytesType:            `BLOB`,
		TimeType:             `TIMESTAMP`,
		HorizonType:          `BIGINT`,
		ConditionalIndexes:   true,
		NoForeignKeys:        true,
		QueryDialect:         QueryDialect,
		NoOffsetWithoutLimit: true,
		Error: func(err error) error {
			return err
		},
		RunTx: runTx,
	})
}

var QueryDialect = cayleysql.QueryDialect{
	RegexpOp:   "REGEXP", // TODO: add regexp extension
	FieldQuote: pq.QuoteIdentifier,
	Placeholder: func(n int) string {
		return fmt.Sprintf("$%d", n)
	},
}

func runTx(tx *sql.Tx, nodes []graphlog.NodeUpdate, quads []graphlog.QuadUpdate, opts graph.IgnoreOpts) error {
	// update node ref counts and insert nodes
	var (
		// prepared statements for each value type
		insertValue = make(map[cayleysql.ValueType]*sql.Stmt)
		updateValue *sql.Stmt
	)

	for _, n := range nodes {
		if n.RefInc < 0 {
			return fmt.Errorf("invalid node.RefInc: %d", n.RefInc)
		}
		nodeKey, values, err := cayleysql.NodeValues(cayleysql.NodeHash{n.Hash}, n.Val)
		if err != nil {
			return err
		}
		values = append([]interface{}{n.RefInc}, values...)
		stmt, ok := insertValue[nodeKey]
		if !ok {
			var ph = make([]string, len(values))
			for i := range ph {
				ph[i] = QueryDialect.Placeholder(i + 1)
			}
			stmt, err = tx.Prepare(fmt.Sprintf(
				`INSERT INTO nodes(refs, hash, %s) VALUES (%s)`,
				strings.Join(nodeKey.Columns(), ", "),
				strings.Join(ph, ", "),
			))
			if err != nil {
				return err
			}
			defer stmt.Close()
			insertValue[nodeKey] = stmt
		}
		_, err = stmt.Exec(values...)
		if isUniqueErr(err) {
			if updateValue == nil {
				updateValue, err = tx.Prepare(`UPDATE nodes SET refs = $1 WHERE hash = $2`)
				if err != nil {
					return err
				}
				defer updateValue.Close()
			}
			if _, err := updateValue.Exec(n.RefInc, cayleysql.NodeHash{n.Hash}.SQLValue()); err != nil {
				return err
			}
		} else if err != nil {
			clog.Errorf("couldn't exec INSERT statement: %#v", err)
			return err
		}
	}

	// now we can deal with quads
	var insertQuad *sql.Stmt
	for _, d := range quads {
		if d.Del {
			return fmt.Errorf("unexpected quad delete: %v", d)
		}
		dirs := make([]interface{}, 0, len(quad.Directions))
		for _, h := range d.Quad.Dirs() {
			dirs = append(dirs, cayleysql.NodeHash{h}.SQLValue())
		}
		if insertQuad == nil {
			var err error
			insertQuad, err = tx.Prepare(
				`INSERT INTO quads(subject_hash, predicate_hash, object_hash, label_hash, ts) VALUES ($1, $2, $3, $4, datetime('now'))`,
			)
			if err != nil {
				return err
			}
			defer insertQuad.Close()
		}
		_, err := insertQuad.Exec(dirs...)
		if isUniqueErr(err) {
			if !opts.IgnoreDup {
				return &graph.DeltaError{Err: graph.ErrQuadExists}
			}
		} else if err != nil {
			clog.Errorf("couldn't exec INSERT statement: %v", err)
			return err
		}
	}
	return nil
}

// isUniqueErr determines whether an error is a SQLite3 uniqueness error.
func isUniqueErr(err error) bool {
	e, ok := err.(sqlite3.Error)
	if !ok {
		return false
	}
	return e.Code == sqlite3.ErrConstraint &&
		(e.ExtendedCode == sqlite3.ErrConstraintUnique || e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey)
}
