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

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/quad/nquads"
	_ "github.com/cayleygraph/cayley/writer"
	"github.com/cheggaaa/pb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/meta-network/go-meta/meta"
)

func init() {
	registerCommand("load", RunLoad, `
usage: meta load [options] <file> <name>

Load quads from <file> into META graph with name <name>.

options:
        -u, --url <url>   URL of the META node [default: dev/meta.ipc]
	--verbosity <n>   Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail [default: 3]
`[1:])
}

func RunLoad(ctx context.Context, args Args) error {
	var logLvl int
	if v := args.String("--verbosity"); v != "" {
		var err error
		logLvl, err = setLogVerbosity(v)
		if err != nil {
			return err
		}
	}
	log.Info("loading quads", "url", args.String("--url"), "file", args.String("<file>"))

	client, err := meta.NewClient(args.String("--url"))
	if err != nil {
		return err
	}

	var in io.Reader
	f, err := os.Open(args.String("<file>"))
	if err != nil {
		return err
	}
	defer f.Close()
	in = f

	if logLvl >= 3 {
		info, err := f.Stat()
		if err != nil {
			return err
		}
		bar := pb.New(int(info.Size())).SetUnits(pb.U_BYTES)
		bar.Start()
		defer bar.Finish()
		in = bar.NewProxyReader(in)
	}

	qs := client.QuadStore(args.String("<name>"))
	qw, err := graph.NewQuadWriter("single", qs, nil)
	if err != nil {
		return err
	}
	qr := nquads.NewReader(in, false)
	count, err := quad.CopyBatch(graph.NewWriter(qw), qr, quad.DefaultBatch)
	if err != nil {
		return err
	}
	log.Info("quads loaded successfully", "count", count)
	return nil
}
