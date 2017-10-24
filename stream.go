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

package meta

import (
	"fmt"
	"os"

	"github.com/ipfs/go-cid"
	"github.com/lmars/tail"
)

// NewStreamWriter returns a StreamWriter which writes CIDs to a META stream.
func NewStreamWriter(path string) (*StreamWriter, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &StreamWriter{file}, nil
}

// StreamWriter writes CIDs to a META stream.
type StreamWriter struct {
	file *os.File
}

// Write writes the given CIDs to the META stream.
func (s *StreamWriter) Write(ids ...*cid.Cid) error {
	for _, id := range ids {
		if _, err := fmt.Fprintln(s.file, id.String()); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the META stream.
func (s *StreamWriter) Close() error {
	return s.file.Close()
}

// StreamOpts are options for opening a StreamReader.
type StreamOpts func(*StreamReader)

// StreamLimit limits the amount of CIDs read from a META stream.
func StreamLimit(limit int) StreamOpts {
	return func(reader *StreamReader) {
		reader.limit = &limit
	}
}

// NewStreamReader returns a StreamReader which reads CIDs from a META stream.
func NewStreamReader(path string, opts ...StreamOpts) (*StreamReader, error) {
	t, err := tail.TailFile(path, tail.Config{Follow: true})
	if err != nil {
		return nil, err
	}
	reader := &StreamReader{
		ch:   make(chan *cid.Cid),
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(reader)
	}
	go func() {
		defer func() {
			t.Stop()
			close(reader.ch)
			close(reader.done)
		}()
		for read := 0; reader.limit == nil || read < *reader.limit; read++ {
			select {
			case line, ok := <-t.Lines:
				if !ok {
					return
				}
				id, err := cid.Parse(line.Text)
				if err != nil {
					reader.err = err
					return
				}
				select {
				case reader.ch <- id:
				case <-reader.stop:
					return
				}
			case <-reader.stop:
				return
			}
		}
	}()
	return reader, nil
}

// StreamReader reads CIDs from a META stream.
type StreamReader struct {
	ch    chan *cid.Cid
	limit *int
	stop  chan struct{}
	done  chan struct{}
	err   error
}

// Ch returns a channel which receives CIDs read from the META stream.
func (r *StreamReader) Ch() <-chan *cid.Cid {
	return r.ch
}

// Close closes the StreamReader.
func (r *StreamReader) Close() error {
	close(r.stop)
	<-r.done
	return nil
}

// Err returns any error encountered whilst reading the META stream.
func (r *StreamReader) Err() error {
	return r.err
}
