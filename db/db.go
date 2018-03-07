// This file is part of the go-kord library.
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
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/kord-network/go-kord/registry"
)

// Driver implements the driver.Conn interface by wrapping a SQLite3 driver
// with connections that use databases stored in Swarm.
type Driver struct {
	name     string
	dpa      *storage.DPA
	registry registry.Registry
	dir      string

	sqlite sqlite3.SQLiteDriver

	dbs   map[string]*db
	dbMtx sync.Mutex
}

// NewDriver creates and registers a new database driver.
func NewDriver(name string, dpa *storage.DPA, registry registry.Registry, dir string) *Driver {
	d := &Driver{
		name:     name,
		dpa:      dpa,
		registry: registry,
		dir:      dir,
		dbs:      make(map[string]*db),
	}
	sql.Register(name, d)
	return d
}

// Open opens the SQLite graph database with the given name, wrapping it in a
// connection which is re-opened if the underlying graph is updated.
func (d *Driver) Open(name string) (driver.Conn, error) {
	db, err := d.openDB(name)
	if err != nil {
		return nil, err
	}
	return db.newConn()
}

// Commit commits the SQLite graph database with the given name by storing it
// in Swarm and returning the resulting Swarm hash.
func (d *Driver) Commit(name string) (common.Hash, error) {
	path := filepath.Join(d.dir, name)
	f, err := os.Open(path)
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

func (d *Driver) openDB(name string) (*db, error) {
	d.dbMtx.Lock()
	defer d.dbMtx.Unlock()

	// if the db is already open, return it
	if db, ok := d.dbs[name]; ok {
		return db, nil
	}

	// get the current hash from the registry
	addr := common.HexToAddress(name)
	hash, err := d.registry.Graph(addr)
	if err != nil {
		return nil, err
	}

	// fetch the database
	path := filepath.Join(d.dir, name)
	if err := d.fetchDB(hash, path); err != nil {
		return nil, err
	}

	// subscribe to registry updates
	updates := make(chan common.Hash)
	sub, err := d.registry.SubscribeGraph(addr, updates)
	if err != nil {
		return nil, err
	}

	db := newDB(d, path)
	d.dbs[name] = db
	go func() {
		defer func() {
			sub.Close()
			db.close()
			d.dbMtx.Lock()
			delete(d.dbs, name)
			d.dbMtx.Unlock()
		}()
		for {
			select {
			case hash, ok := <-updates:
				if !ok {
					return
				}
				if err := d.fetchDB(hash, path); err != nil {
					return
				}
				if err := db.reopenConns(); err != nil {
					return
				}
			case <-db.closed:
				return
			}
		}
	}()

	return db, nil
}

func (d *Driver) fetchDB(hash common.Hash, path string) error {
	tmp, err := ioutil.TempFile("", "kord-db")
	if err != nil {
		return err
	}
	defer tmp.Close()
	if err := d.fetchHash(hash, tmp); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return nil
}

func (d *Driver) fetchHash(hash common.Hash, dst io.Writer) error {
	if common.EmptyHash(hash) {
		return nil
	}
	reader := d.dpa.Retrieve(storage.Key(hash[:]))
	size, err := reader.Size(nil)
	if err != nil {
		return err
	}
	n, err := io.Copy(dst, io.LimitReader(reader, size))
	if err != nil {
		return err
	} else if n != size {
		return fmt.Errorf("failed to fetch database, expected %d bytes, copied %d", size, n)
	}
	return nil
}

type db struct {
	driver *Driver

	path string

	conns    map[*Conn]struct{}
	connsMtx sync.RWMutex

	closed chan struct{}
}

func newDB(driver *Driver, path string) *db {
	return &db{
		driver: driver,
		path:   path,
		conns:  make(map[*Conn]struct{}),
		closed: make(chan struct{}),
	}
}

func (db *db) newConn() (driver.Conn, error) {
	sqliteConn, err := db.driver.sqlite.Open(db.path)
	if err != nil {
		return nil, err
	}

	conn := newConn(db, sqliteConn.(*sqlite3.SQLiteConn))
	db.addConn(conn)
	return conn, nil
}

func (db *db) addConn(conn *Conn) {
	db.connsMtx.Lock()
	defer db.connsMtx.Unlock()
	db.conns[conn] = struct{}{}
}

func (db *db) removeConn(conn *Conn) {
	db.connsMtx.Lock()
	defer db.connsMtx.Unlock()
	delete(db.conns, conn)
}

func (db *db) reopenConns() error {
	db.connsMtx.RLock()
	defer db.connsMtx.RUnlock()
	for conn := range db.conns {
		sqliteConn, err := db.driver.sqlite.Open(db.path)
		if err != nil {
			return err
		}
		conn.sqliteConn.Store(sqliteConn)
	}
	return nil
}

func (db *db) close() {
	close(db.closed)
	db.connsMtx.Lock()
	for conn := range db.conns {
		conn.Close()
	}
	db.conns = nil
	db.connsMtx.Unlock()
}

type Conn struct {
	sqliteConn atomic.Value

	db *db

	closeOnce sync.Once
	closed    chan struct{}
}

func newConn(db *db, sqliteConn *sqlite3.SQLiteConn) *Conn {
	conn := &Conn{
		db:     db,
		closed: make(chan struct{}),
	}
	conn.sqliteConn.Store(sqliteConn)
	return conn
}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return c.SQLiteConn().Prepare(query)
}

func (c *Conn) Close() error {
	c.closeOnce.Do(func() { close(c.closed) })
	c.db.removeConn(c)
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
