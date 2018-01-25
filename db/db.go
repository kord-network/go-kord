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
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
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
	meta "github.com/meta-network/go-meta"
)

var driver *Driver

type Driver struct {
	dpa *storage.DPA
	ens meta.ENS
	dir string

	sqlite sqlite3.SQLiteDriver

	connsMtx sync.RWMutex
	conns    map[string]map[*Conn]struct{}
}

func Init(dpa *storage.DPA, ens meta.ENS, dir string) {
	if driver != nil {
		panic("db: driver already initialised")
	}

	driver = &Driver{
		dpa:   dpa,
		ens:   ens,
		dir:   dir,
		conns: make(map[string]map[*Conn]struct{}),
	}

	sql.Register("meta", driver)
}

func (d *Driver) Open(name string) (sqldriver.Conn, error) {
	// fetch the database if it doesn't exist locally
	_, err := os.Stat(d.path(name))
	if os.IsNotExist(err) {
		err = d.fetchDB(name)
	}
	if err != nil {
		return nil, err
	}

	// open a SQLite connection to the db
	sqliteConn, err := d.sqlite.Open(d.path(name))
	if err != nil {
		return nil, err
	}

	// return the wrapped connection
	conn := newConn(name, d, sqliteConn.(*sqlite3.SQLiteConn))
	d.addConn(conn)
	return conn, nil
}

func Update(name string, hash common.Hash) error {
	if driver == nil {
		panic("db: uninitialised driver")
	}
	return driver.Update(name, hash)
}

// Update fetches the db stored at the given hash, stores it at the db's path
// and then re-opens connections to the database.
func (d *Driver) Update(name string, hash common.Hash) error {
	tmp, err := ioutil.TempFile("", "meta-db")
	if err != nil {
		return err
	}
	defer tmp.Close()
	if err := d.fetchHash(hash, tmp); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err := os.Rename(tmp.Name(), d.path(name)); err != nil {
		os.Remove(tmp.Name())
		return err
	}

	d.connsMtx.RLock()
	defer d.connsMtx.RUnlock()
	conns, ok := d.conns[name]
	if !ok {
		return nil
	}
	for conn := range conns {
		sqliteConn, err := d.sqlite.Open(d.path(name))
		if err != nil {
			return err
		}
		conn.sqliteConn.Store(sqliteConn)
	}
	return nil
}

func Commit(name string) (common.Hash, error) {
	if driver == nil {
		panic("db: uninitialised driver")
	}
	return driver.Commit(name)
}

func (d *Driver) Commit(name string) (common.Hash, error) {
	f, err := os.Open(d.path(name))
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

func (d *Driver) addConn(conn *Conn) {
	d.connsMtx.Lock()
	defer d.connsMtx.Unlock()
	if _, ok := d.conns[conn.name]; !ok {
		d.conns[conn.name] = make(map[*Conn]struct{})
	}
	d.conns[conn.name][conn] = struct{}{}
}

func (d *Driver) removeConn(conn *Conn) {
	d.connsMtx.Lock()
	defer d.connsMtx.Unlock()
	conns, ok := d.conns[conn.name]
	if !ok {
		return
	}
	delete(conns, conn)
}

func (d *Driver) path(name string) string {
	return filepath.Join(d.dir, name)
}

func (d *Driver) fetchDB(name string) error {
	hash, err := d.ens.Content(name)
	if err != nil {
		return err
	}
	if common.EmptyHash(hash) {
		return nil
	}
	dst, err := os.Create(d.path(name))
	if err != nil {
		return err
	}
	defer dst.Close()
	return d.fetchHash(hash, dst)
}

func (d *Driver) fetchHash(hash common.Hash, dst io.Writer) error {
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

type Conn struct {
	sqliteConn atomic.Value

	name   string
	driver *Driver

	closeOnce sync.Once
	closed    chan struct{}
}

func newConn(name string, driver *Driver, sqliteConn *sqlite3.SQLiteConn) *Conn {
	conn := &Conn{
		name:   name,
		driver: driver,
		closed: make(chan struct{}),
	}
	conn.sqliteConn.Store(sqliteConn)
	return conn
}

func (c *Conn) Prepare(query string) (sqldriver.Stmt, error) {
	return c.SQLiteConn().Prepare(query)
}

func (c *Conn) Close() error {
	c.closeOnce.Do(func() { close(c.closed) })
	c.driver.removeConn(c)
	return c.SQLiteConn().Close()
}

func (c *Conn) Begin() (sqldriver.Tx, error) {
	return c.SQLiteConn().Begin()
}

func (c *Conn) BeginTx(ctx context.Context, opts sqldriver.TxOptions) (sqldriver.Tx, error) {
	return c.SQLiteConn().BeginTx(ctx, opts)
}

func (c *Conn) PrepareContext(ctx context.Context, query string) (sqldriver.Stmt, error) {
	return c.SQLiteConn().PrepareContext(ctx, query)
}

func (c *Conn) Exec(query string, args []sqldriver.Value) (sqldriver.Result, error) {
	return c.SQLiteConn().Exec(query, args)
}

func (c *Conn) ExecContext(ctx context.Context, query string, args []sqldriver.NamedValue) (sqldriver.Result, error) {
	return c.SQLiteConn().ExecContext(ctx, query, args)
}

func (c *Conn) Ping(ctx context.Context) error {
	return c.SQLiteConn().Ping(ctx)
}

func (c *Conn) Query(query string, args []sqldriver.Value) (sqldriver.Rows, error) {
	return c.SQLiteConn().Query(query, args)
}

func (c *Conn) QueryContext(ctx context.Context, query string, args []sqldriver.NamedValue) (sqldriver.Rows, error) {
	return c.SQLiteConn().QueryContext(ctx, query, args)
}

func (c *Conn) SQLiteConn() *sqlite3.SQLiteConn {
	return c.sqliteConn.Load().(*sqlite3.SQLiteConn)
}
