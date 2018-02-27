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
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/quad/nquads"
	_ "github.com/cayleygraph/cayley/writer"
	"github.com/cheggaaa/pb"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/meta-network/go-meta/meta"
	"github.com/meta-network/go-meta/registry"
	"github.com/moby/moby/pkg/term"
)

func init() {
	registerCommand("graph", RunGraph, `
usage: meta graph create [options] <id>
       meta graph load [options] <id> <file>

Create, update or query a META graph.

options:
        -u, --url <url>        URL of the META node
	-k, --keystore <dir>   Keystore directory
`[1:])
}

func RunGraph(ctx *Context) error {
	switch {
	case ctx.Args.Bool("create"):
		return RunGraphCreate(ctx)
	case ctx.Args.Bool("load"):
		return RunGraphLoad(ctx)
	default:
		return errors.New("unknown graph command")
	}
}

func RunGraphCreate(ctx *Context) error {
	idArg := ctx.Args.String("<id>")
	if !common.IsHexAddress(idArg) {
		return fmt.Errorf("invalid META ID, must be a hex string: %s", idArg)
	}
	id := common.HexToAddress(idArg)

	client, err := ctx.Client()
	if err != nil {
		return err
	}

	log.Info("creating graph", "id", id)
	hash, err := client.CreateGraph(ctx, id.Hex())
	if err != nil {
		return err
	}

	if err := setGraph(ctx, client, id, hash); err != nil {
		return err
	}

	log.Info("graph created successfully", "id", id, "hash", hash)
	return nil
}

func RunGraphLoad(ctx *Context) error {
	idArg := ctx.Args.String("<id>")
	if !common.IsHexAddress(idArg) {
		return fmt.Errorf("invalid META ID, must be a hex string: %s", idArg)
	}
	id := common.HexToAddress(idArg)

	client, err := ctx.Client()
	if err != nil {
		return err
	}

	file := ctx.Args.String("<file>")
	log.Info("loading quads", "id", id, "file", file)
	count, err := loadQuads(ctx, client, id, file)
	if err != nil {
		return err
	}

	log.Info("committing graph")
	hash, err := client.CommitGraph(ctx, id.Hex())
	if err != nil {
		return err
	}

	if err := setGraph(ctx, client, id, hash); err != nil {
		return err
	}

	log.Info("quads loaded successfully", "id", id, "count", count, "hash", hash)
	return nil
}

func loadQuads(ctx *Context, client *meta.Client, id common.Address, file string) (int, error) {
	var in io.Reader
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	in = f

	if stderr, ok := ctx.Stderr.(*os.File); ok && term.IsTerminal(stderr.Fd()) {
		info, err := f.Stat()
		if err != nil {
			return 0, err
		}
		bar := pb.New(int(info.Size())).SetUnits(pb.U_BYTES)
		bar.Output = stderr
		bar.Start()
		defer bar.Finish()
		in = bar.NewProxyReader(in)
	}

	qs := client.QuadStore(id.Hex())
	qw, err := graph.NewQuadWriter("single", qs, nil)
	if err != nil {
		return 0, err
	}
	qr := nquads.NewReader(in, false)
	return quad.CopyBatch(graph.NewWriter(qw), qr, quad.DefaultBatch)
}

func setGraph(ctx *Context, client *meta.Client, id common.Address, hash common.Hash) error {
	log.Info("signing graph hash", "hash", hash)
	sig, err := signHash(ctx, id, hash)
	if err != nil {
		return err
	}

	log.Info("updating registry")
	return client.SetGraph(ctx, hash, sig)
}

func signHash(ctx *Context, id common.Address, hash common.Hash) ([]byte, error) {
	if id == registry.DevAddr {
		return crypto.Sign(hash[:], registry.DevKey)
	}
	ks := keystore.NewKeyStore(
		ctx.Args.String("--keystore"),
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)
	account, err := ks.Find(accounts.Account{Address: id})
	if err != nil {
		return nil, err
	}
	passphrase, err := getPassphrase(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("error reading passphrase: %s", err)
	}
	return ks.SignHashWithPassphrase(account, string(passphrase), hash[:])
}
