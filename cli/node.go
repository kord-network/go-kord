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
	"fmt"
	"strconv"

	"github.com/meta-network/go-meta/node"
)

func init() {
	registerCommand("node", RunNode, `
usage: meta node [options]

Run a META node.

options:
        -p, --port <port>   HTTP port
`[1:])
}

func RunNode(ctx context.Context, args Args) error {
	config := node.DefaultConfig
	if p := args.String("--port"); p != "" {
		port, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("error parsing --port: %s", err)
		}
		config.API.HTTPPort = port
	}
	node := node.New(config)

	// start the node
	if err := node.Start(); err != nil {
		return err
	}

	// stop the node if the context is cancelled
	go func() {
		<-ctx.Done()
		node.Stop()
	}()

	// wait for the node to exit
	return node.Wait()
}
