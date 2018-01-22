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
	"sync"
	"sync/atomic"

	sqlite3 "github.com/mattn/go-sqlite3"
)

type Driver struct {
	sqlite3.SQLiteDriver

	storage Storage
}

func NewDriver(storage Storage) *Driver {
	return &Driver{storage: storage}
}

func (d *Driver) Open(name string) (driver.Conn, error) {
	notify := make(chan struct{})
	path, notifier, err := d.storage.GetDB(name, notify)
	if err != nil {
		return nil, err
	}
	sqliteConn, err := d.SQLiteDriver.Open(path)
	if err != nil {
		notifier.Close()
		return nil, err
	}
	conn := newConn(sqliteConn.(*sqlite3.SQLiteConn))
	go func() {
		defer notifier.Close()
		for {
			select {
			case _, ok := <-notify:
				if !ok {
					conn.Close()
					return
				}
				sqliteConn, err := d.SQLiteDriver.Open(path)
				if err != nil {
					conn.Close()
					return
				}
				conn.sqliteConn.Store(sqliteConn)
			case <-conn.closed:
				return
			}
		}
	}()
	return conn, nil
}

type Conn struct {
	sqliteConn atomic.Value

	closeOnce sync.Once
	closed    chan struct{}
}

func newConn(sqliteConn *sqlite3.SQLiteConn) *Conn {
	conn := &Conn{
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
