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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/meta-network/go-meta/node"
)

func init() {
	registerCommand("node", RunNode, `
usage: meta node [options]

Run a META node.

options:
        -d, --datadir <dir>  META data directory
        --ens.url <url>      ENS URL
        --ens.addr <addr>    ENS registry address
        --ens.key <key>      ENS private key
        --http.port <port>   HTTP server port
`[1:])
}

func RunNode(ctx context.Context, args Args) error {
	config := node.DefaultConfig
	if dir := args.String("--datadir"); dir != "" {
		config.DataDir = dir
	}
	if v := args.String("--http.port"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("error parsing --http.port: %s", err)
		}
		config.API.HTTPPort = port
	}
	if v := args.String("--ens.url"); v != "" {
		config.ENS.URL = v
	}
	if v := args.String("--ens.addr"); v != "" {
		if !common.IsHexAddress(v) {
			return fmt.Errorf("invalid --ens.addr: %s", v)
		}
		config.ENS.RegistryAddr = common.HexToAddress(v)
	}
	if v := args.String("--ens.key"); v != "" {
		key, err := crypto.HexToECDSA(v)
		if err != nil {
			return fmt.Errorf("invalid --ens.key: %s", err)
		}
		config.ENS.Key = key
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
