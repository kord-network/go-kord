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
	"net/url"
	"os"
	"path/filepath"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/schema"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	swarm "github.com/ethereum/go-ethereum/swarm/api/client"
	"github.com/meta-network/go-meta/dapp"
	"github.com/meta-network/go-meta/meta"
)

func init() {
	registerCommand("dapp", RunDapp, `
usage: meta dapp deploy [options] <dir> <uri>
       meta dapp set-root [options] <uri>

Deploy a META Dapp.

options:
        -u, --url <url>        URL of the META node
        -s, --swarm-api <url>  URL of the Swarm API [default: http://localhost:8500]
        -k, --keystore <dir>   Keystore directory [default: dev/keystore]

example:
        meta dapp deploy path/to/dapp meta://xyz123/cool-dapp

        meta dapp set-root meta://xyz123/cool-dapp
`[1:])
}

func RunDapp(ctx *Context, args Args) error {
	switch {
	case args.Bool("deploy"):
		return RunDappDeploy(ctx, args)
	case args.Bool("set-root"):
		return RunDappSetRoot(ctx, args)
	default:
		return errors.New("unknown dapp command")
	}
}

func RunDappDeploy(ctx *Context, args Args) error {
	u, err := url.Parse(args.String("<uri>"))
	if err != nil {
		return err
	}
	if u.Scheme != "meta" {
		return fmt.Errorf("<uri> must have meta scheme, not %s", u.Scheme)
	}
	if !common.IsHexAddress(u.Host) {
		return fmt.Errorf("invalid META ID: %s", u.Host)
	}
	id := common.HexToAddress(u.Host)

	ks := keystore.NewKeyStore(
		args.String("--keystore"),
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)
	account, err := ks.Find(accounts.Account{Address: id})
	if err != nil {
		return err
	}

	client, err := meta.NewClient(args.NodeURL())
	if err != nil {
		return err
	}

	swarm := swarm.NewClient(args.String("--swarm-api"))
	dir := args.String("<dir>")
	var defaultPath string
	if _, err := os.Stat(filepath.Join(dir, "index.html")); err == nil {
		defaultPath = filepath.Join(dir, "index.html")
	}
	manifestHash, err := swarm.UploadDirectory(dir, defaultPath, "")
	if err != nil {
		return err
	}

	qs := client.QuadStore(id.Hex())
	qw, err := graph.NewQuadWriter("single", qs, nil)
	if err != nil {
		return err
	}
	w := graph.NewWriter(qw)
	d := dapp.Dapp{
		ID:           quad.IRI(u.String()),
		ManifestHash: manifestHash,
	}
	if _, err := schema.WriteAsQuads(w, d); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	hash, err := client.CommitGraph(ctx, id.Hex())
	if err != nil {
		return err
	}

	log.Info("signing graph hash", "hash", hash)
	passphrase, err := getPassphrase(ctx, false)
	if err != nil {
		return fmt.Errorf("error reading passphrase: %s", err)
	}
	sig, err := ks.SignHashWithPassphrase(account, string(passphrase), hash[:])
	if err != nil {
		return err
	}

	log.Info("updating registry")
	if err := client.SetGraph(ctx, hash, sig); err != nil {
		return err
	}

	return nil
}

func RunDappSetRoot(ctx *Context, args Args) error {
	client, err := meta.NewClient(args.NodeURL())
	if err != nil {
		return err
	}
	return client.SetRootDapp(ctx, args.String("<uri>"))
}
