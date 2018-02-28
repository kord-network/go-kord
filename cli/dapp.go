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
	"os"
	"path/filepath"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/schema"
	"github.com/ethereum/go-ethereum/log"
	swarm "github.com/ethereum/go-ethereum/swarm/api/client"
	"github.com/meta-network/go-meta/dapp"
)

func init() {
	registerCommand("dapp", RunDapp, `
usage: meta dapp deploy [options] <dir> <uri>
       meta dapp set-root [options] <uri>

Deploy a META Dapp.

options:
        -u, --url <url>        URL of the META node
        -s, --swarm-api <url>  URL of the Swarm API [default: http://localhost:5000]
        -k, --keystore <dir>   Keystore directory

example:
        meta dapp deploy path/to/dapp meta://xyz123/cool-dapp

        meta dapp set-root meta://xyz123/cool-dapp
`[1:])
}

func RunDapp(ctx *Context) error {
	switch {
	case ctx.Args.Bool("deploy"):
		return RunDappDeploy(ctx)
	case ctx.Args.Bool("set-root"):
		return RunDappSetRoot(ctx)
	default:
		return errors.New("unknown dapp command")
	}
}

func RunDappDeploy(ctx *Context) error {
	u, err := ctx.URI()
	if err != nil {
		return err
	}
	id := u.ID

	client, err := ctx.Client()
	if err != nil {
		return err
	}

	swarm := swarm.NewClient(ctx.Args.String("--swarm-api"))
	dir := ctx.Args.String("<dir>")
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

	log.Info("committing graph")
	hash, err := client.CommitGraph(ctx, id.Hex())
	if err != nil {
		return err
	}

	if err := setGraph(ctx, client, id, hash); err != nil {
		return err
	}

	log.Info("dapp deployed", "uri", d.ID, "hash", hash)
	return nil
}

func RunDappSetRoot(ctx *Context) error {
	client, err := ctx.Client()
	if err != nil {
		return err
	}
	if err := client.SetRootDapp(ctx, ctx.Args.String("<uri>")); err != nil {
		return err
	}
	log.Info("root dapp set", "uri", ctx.Args.String("<uri>"))
	return nil
}
