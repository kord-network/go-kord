// This file is part of the go-meta library.
//
// Copyright (C) 2017 JAAK MUSIC LTD
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

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ipfs/go-datastore/fs"
	"github.com/meta-network/go-meta"
	"github.com/meta-network/go-meta/cli"
)

func main() {
	log.Root().SetHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(true)))

	store, err := openStore()
	if err != nil {
		log.Crit("error opening meta store", "err", err)
	}

	// shutdown gracefully on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Info("received signal, exiting...")
	}()

	if err := cli.New(store, os.Stdin, os.Stdout).Run(ctx, os.Args[1:]...); err != nil {
		log.Crit("error running meta command", "err", err)
	}
}

func openStore() (*meta.Store, error) {
	metaDir := ".meta"
	if err := os.MkdirAll(metaDir, 0755); err != nil {
		return nil, err
	}
	store, err := fs.NewDatastore(metaDir)
	if err != nil {
		return nil, err
	}
	return meta.NewStore(store), nil
}
