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
	"crypto/ecdsa"
	"fmt"
	"os"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/quad/nquads"
	_ "github.com/cayleygraph/cayley/writer"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/meta-network/go-meta/api"
	"github.com/meta-network/go-meta/node"
	"github.com/meta-network/go-meta/store"
)

var defaultPrivKey = mustHexToECDSA("5ea708a1d733fb6f2254a2f092c4e6e748ec9564c84247f293f32199833d62bf")

func init() {
	registerCommand("load", RunLoad, `
usage: meta load [options] <file> <db>

Load quads from <file> into META database <db>.

options:
        -d, --datadir <dir>  META data directory
        -k, --key <key>      Hex encoded private key
        -u, --url <url>      URL of the META node [default: http://127.0.0.1:5000]
`[1:])
}

func RunLoad(ctx context.Context, args Args) error {
	config := node.DefaultConfig
	if dir := args.String("--datadir"); dir != "" {
		config.DataDir = dir
	}
	node := node.New(config)
	if err := node.Start(); err != nil {
		return err
	}
	defer node.Stop()

	f, err := os.Open(args.String("<file>"))
	if err != nil {
		return err
	}
	defer f.Close()

	privKey := defaultPrivKey
	if k := args.String("--key"); k != "" {
		key, err := crypto.HexToECDSA(k)
		if err != nil {
			return fmt.Errorf("error parsing --key: %s", err)
		}
		privKey = key
	}

	address := crypto.PubkeyToAddress(privKey.PublicKey)
	signer := store.NewPrivateKeySigner(privKey)
	url := args.String("--url")
	client := api.NewClient(url)
	store, err := store.NewClientStore(
		address,
		signer,
		client,
		args.String("<db>"),
	)
	if err != nil {
		return err
	}

	qw, err := graph.NewQuadWriter("single", store, nil)
	if err != nil {
		return err
	}
	qr := nquads.NewReader(f, false)
	_, err = quad.CopyBatch(graph.NewWriter(qw), qr, 100)
	return err
}

func mustHexToECDSA(hexkey string) *ecdsa.PrivateKey {
	key, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		panic(err)
	}
	return key
}
