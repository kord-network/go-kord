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

package cli

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/meta-network/go-meta/meta"
	"github.com/meta-network/go-meta/pkg/uri"
)

type Context struct {
	context.Context

	Args Args

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func NewContext(ctx context.Context) *Context {
	return &Context{Context: ctx}
}

func (c *Context) NodeURL() string {
	if url := c.Args.String("--url"); url != "" {
		return url
	}
	return filepath.Join(os.TempDir(), "meta.ipc")
}

func (c *Context) URI() (*uri.URI, error) {
	return uri.Parse(c.Args.String("<uri>"))
}

func (c *Context) Client() (*meta.Client, error) {
	return meta.NewClient(c.NodeURL())
}
