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

package sql

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type ENS interface {
	Resolve(name string) (common.Hash, error)
}

type Driver struct {
	sqlite3.SQLiteDriver

	ens     ENS
	dpa     *storage.DPA
	dataDir string

	conns    map[string][]*Conn
	connsMtx sync.RWMutex
}

func NewDriver(dpa *storage.DPA, ens ENS, dataDir string) *Driver {
	return &Driver{
		ens:     ens,
		dpa:     dpa,
		dataDir: dataDir,
		conns:   make(map[string][]*Conn),
	}
}

func (d *Driver) Open(name string) (driver.Conn, error) {
	hash, err := d.ens.Resolve(name)
	if err != nil {
		return nil, err
	}

	sqliteConn, err := d.open(name, hash)
	if err != nil {
		return nil, err
	}

	conn := &Conn{}
	conn.sqliteConn.Store(sqliteConn)

	d.connsMtx.Lock()
	d.conns[name] = append(d.conns[name], conn)
	d.connsMtx.Unlock()

	return conn, nil
}

func (d *Driver) Update(name string, hash common.Hash) error {
	d.connsMtx.RLock()
	defer d.connsMtx.RUnlock()
	conns, ok := d.conns[name]
	if !ok {
		return nil
	}
	for _, conn := range conns {
		sqliteConn, err := d.open(name, hash)
		if err != nil {
			return err
		}
		conn.sqliteConn.Store(sqliteConn)
	}
	return nil
}

func (d *Driver) Save(name string) (common.Hash, error) {
	dbPath := filepath.Join(d.dataDir, name)
	f, err := os.Open(dbPath)
	if err != nil {
		return common.Hash{}, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return common.Hash{}, err
	}
	key, err := d.dpa.Store(f, info.Size(), &sync.WaitGroup{}, &sync.WaitGroup{})
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(key[:]), nil
}

func (d *Driver) open(name string, hash common.Hash) (driver.Conn, error) {
	dbPath := filepath.Join(d.dataDir, name)
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		if hash != (common.Hash{}) {
			reader := d.dpa.Retrieve(storage.Key(hash[:]))
			size, err := reader.Size(nil)
			if err != nil {
				return nil, err
			}
			dst, err := os.Create(dbPath)
			if err != nil {
				return nil, err
			}
			n, err := io.Copy(dst, io.LimitReader(reader, size))
			dst.Close()
			if err != nil {
				return nil, err
			} else if n != size {
				return nil, fmt.Errorf("failed to copy database, expected %d bytes, copied %d", size, n)
			}
		}
	} else if err != nil {
		return nil, err
	}
	return d.SQLiteDriver.Open(dbPath)
}

type Conn struct {
	sqliteConn atomic.Value
}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return c.SQLiteConn().Prepare(query)
}

func (c *Conn) Close() error {
	return c.SQLiteConn().Close()
}

func (c *Conn) Begin() (driver.Tx, error) {
	return c.SQLiteConn().Begin()
}

func (c *Conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.SQLiteConn().BeginTx(ctx, opts)
}

func (c *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return c.SQLiteConn().PrepareContext(ctx, query)
}

func (c *Conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	return c.SQLiteConn().Exec(query, args)
}

func (c *Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return c.SQLiteConn().ExecContext(ctx, query, args)
}

func (c *Conn) Ping(ctx context.Context) error {
	return c.SQLiteConn().Ping(ctx)
}

func (c *Conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	return c.SQLiteConn().Query(query, args)
}

func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return c.SQLiteConn().QueryContext(ctx, query, args)
}

func (c *Conn) SQLiteConn() *sqlite3.SQLiteConn {
	return c.sqliteConn.Load().(*sqlite3.SQLiteConn)
}
