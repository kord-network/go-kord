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

import "github.com/ipfs/go-cid"

// LimitedReader reads a limitied number of CIDs from an underlying META stream
// and then closes the stream channel (useful when wanting to index an exact
// number of CIDs from a stream).
func LimitedReader(r StreamReader, limit int) StreamReader {
	l := newLimitedReader(r)
	go func() {
		defer close(l.ch)
		for read := 0; read < limit; read++ {
			select {
			case id, ok := <-r.Ch():
				if !ok {
					l.err = r.Err()
					return
				}
				select {
				case l.ch <- id:
				case <-l.stop:
					return
				}
			case <-l.stop:
				return
			}
		}
	}()
	return l
}

type limitedReader struct {
	r    StreamReader
	ch   chan *cid.Cid
	stop chan struct{}
	err  error
}

func newLimitedReader(r StreamReader) *limitedReader {
	return &limitedReader{
		r:    r,
		ch:   make(chan *cid.Cid),
		stop: make(chan struct{}),
	}
}

func (l *limitedReader) Ch() <-chan *cid.Cid {
	return l.ch
}

func (l *limitedReader) Close() error {
	close(l.stop)
	return l.r.Close()
}

func (l *limitedReader) Err() error {
	return l.err
}
