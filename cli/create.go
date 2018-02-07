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
	"github.com/ethereum/go-ethereum/log"
	"github.com/meta-network/go-meta/meta"
)

func init() {
	registerCommand("create", RunCreate, `
usage: meta create [options] <name>

Create META graph with name <name>.

options:
        -u, --url <url>   URL of the META node [default: dev/meta.ipc]
	--verbosity <n>   Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail [default: 3]
`[1:])
}

func RunCreate(ctx context.Context, args Args) error {
	if v := args.String("--verbosity"); v != "" {
		if _, err := setLogVerbosity(v); err != nil {
			return err
		}
	}
	log.Info("creating graph", "url", args.String("--url"), "name", args.String("<name>"))

	client, err := meta.NewClient(args.String("--url"))
	if err != nil {
		return err
	}

	if err := client.CreateGraph(ctx, args.String("<name>")); err != nil {
		return err
	}

	log.Info("graph created successfully", "name", args.String("<name>"))
	return nil
}
