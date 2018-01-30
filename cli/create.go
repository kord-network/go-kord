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

	_ "github.com/cayleygraph/cayley/writer"
	"github.com/meta-network/go-meta/api"
)

func init() {
	registerCommand("create", RunCreate, `
usage: meta create [options] <db>

Create META database <db>.

options:
        -u, --url <url>   URL of the META node [default: http://localhost:5000]
`[1:])
}

func RunCreate(ctx context.Context, args Args) error {
	client := api.NewClient(args.String("--url"), args.String("<db>"))
	return client.Create()
}
