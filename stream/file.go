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
	"fmt"
	"os"

	"github.com/hpcloud/tail"
	"github.com/ipfs/go-cid"
)

// NewFileStream returns a Stream which reads and writes CIDs to a single file.
func NewFileStream(path string) Stream {
	return &fileStream{path: path}
}

type fileStream struct {
	path string
}

func (f *fileStream) NewWriter() (StreamWriter, error) {
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return &fileStreamWriter{file}, nil
}

type fileStreamWriter struct {
	file *os.File
}

func (f *fileStreamWriter) Write(ids ...*cid.Cid) error {
	for _, id := range ids {
		if _, err := fmt.Fprintln(f.file, id.String()); err != nil {
			return err
		}
	}
	return nil
}

func (f *fileStreamWriter) Close() error {
	return f.file.Close()
}

func (f *fileStream) NewReader() (StreamReader, error) {
	t, err := tail.TailFile(f.path, tail.Config{Follow: true})
	if err != nil {
		return nil, err
	}
	r := &fileStreamReader{
		ch:   make(chan *cid.Cid),
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go func() {
		defer func() {
			t.Stop()
			close(r.ch)
			close(r.done)
		}()
		for {
			select {
			case line, ok := <-t.Lines:
				if !ok {
					return
				}
				id, err := cid.Parse(line.Text)
				if err != nil {
					r.err = err
					return
				}
				select {
				case r.ch <- id:
				case <-r.stop:
					return
				}
			case <-r.stop:
				return
			}
		}
	}()
	return r, nil
}

type fileStreamReader struct {
	ch   chan *cid.Cid
	stop chan struct{}
	done chan struct{}
	err  error
}

func (f *fileStreamReader) Ch() <-chan *cid.Cid {
	return f.ch
}

func (f *fileStreamReader) Close() error {
	close(f.stop)
	<-f.done
	return nil
}

func (f *fileStreamReader) Err() error {
	return f.err
}
