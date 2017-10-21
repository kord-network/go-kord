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

package stream

import (
	"sync"

	"github.com/ipfs/go-cid"
)

// NewMemoryStream returns a Stream which reads and writes CIDs in memory, with
// writes being synchronously broadcast to all readers.
func NewMemoryStream() Stream {
	return &memoryStream{
		readers: make(map[*memoryStreamReader]struct{}),
	}
}

type memoryStream struct {
	mtx     sync.RWMutex
	readers map[*memoryStreamReader]struct{}
}

func (m *memoryStream) NewWriter() (StreamWriter, error) {
	return &memoryStreamWriter{m: m}, nil
}

type memoryStreamWriter struct {
	m *memoryStream
}

func (w *memoryStreamWriter) Write(ids ...*cid.Cid) error {
	w.m.mtx.RLock()
	defer w.m.mtx.RUnlock()
	for _, id := range ids {
		for r := range w.m.readers {
			select {
			case r.ch <- id:
			case <-r.done:
			}
		}
	}
	return nil
}

func (m *memoryStreamWriter) Close() error {
	return nil
}

func (m *memoryStream) NewReader() (StreamReader, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	r := &memoryStreamReader{
		m:    m,
		ch:   make(chan *cid.Cid),
		done: make(chan struct{}),
	}
	m.readers[r] = struct{}{}
	return r, nil
}

type memoryStreamReader struct {
	m    *memoryStream
	ch   chan *cid.Cid
	done chan struct{}
}

func (m *memoryStreamReader) Ch() <-chan *cid.Cid {
	return m.ch
}

func (r *memoryStreamReader) Close() error {
	r.m.mtx.Lock()
	delete(r.m.readers, r)
	r.m.mtx.Unlock()
	close(r.done)
	return nil
}

func (r *memoryStreamReader) Err() error {
	return nil
}
