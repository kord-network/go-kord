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

package node

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"github.com/meta-network/go-meta/api"
	"github.com/meta-network/go-meta/db"
	"github.com/meta-network/go-meta/ens"
)

type Config struct {
	DataDir string

	API api.Config
	ENS ens.Config
}

var DefaultConfig = Config{
	DataDir: ".meta",
	API:     api.DefaultConfig,
	ENS:     ens.DefaultConfig,
}

var DefaultLogger = log.New()

func init() {
	DefaultLogger.SetHandler(
		log.LvlFilterHandler(
			log.LvlInfo,
			log.StreamHandler(
				os.Stderr,
				log.TerminalFormat(true),
			),
		),
	)
}

type Node struct {
	config Config

	dpa *storage.DPA
	srv *http.Server
	ens *ens.Client

	log log.Logger

	err      error
	done     chan struct{}
	doneOnce sync.Once
}

func New(config Config) *Node {
	return &Node{
		config: config,
		done:   make(chan struct{}),
		log:    DefaultLogger,
	}
}

func (n *Node) Start() error {
	n.log.Info("starting swarm store", "dir", n.config.DataDir)
	localStore, err := storage.NewLocalStore(
		storage.MakeHashFunc("SHA3"),
		&storage.StoreParams{
			ChunkDbPath:   n.config.DataDir,
			DbCapacity:    5000000,
			CacheCapacity: 5000,
			Radius:        0,
		},
	)
	if err != nil {
		n.log.Error("error opening local swarm store", "dir", n.config.DataDir, "err", err)
		return err
	}
	chunker := storage.NewTreeChunker(storage.NewChunkerParams())
	n.dpa = &storage.DPA{
		Chunker:    chunker,
		ChunkStore: localStore,
	}
	n.dpa.Start()

	n.log.Info("connecting ENS client", "url", n.config.ENS.URL)
	ens, err := ens.NewClientWithConfig(n.config.ENS)
	if err != nil {
		n.dpa.Stop()
		return err
	}
	n.ens = ens

	n.log.Info("registering the META storage")
	db.Init(n.dpa, ens, n.config.DataDir)

	addr := fmt.Sprintf("%s:%d", n.config.API.HTTPAddr, n.config.API.HTTPPort)
	n.log.Info("starting HTTP server", "addr", addr)
	n.srv = &http.Server{
		Addr:    addr,
		Handler: api.NewServer(),
	}
	go func() {
		if err := n.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			n.log.Error("error starting HTTP server", "err", err)
			n.err = err
			n.Stop()
		}
	}()
	return nil

}

func (n *Node) Stop() error {
	n.log.Info("stopping META node")
	if n.ens != nil {
		n.log.Info("closing ENS client")
		n.ens.Close()
	}
	if n.srv != nil {
		n.log.Info("stopping HTTP server", "addr", n.srv.Addr)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		n.srv.Shutdown(ctx)
	}
	if n.dpa != nil {
		n.log.Info("stopping swarm store", "dir", n.config.DataDir)
		n.dpa.Stop()
	}
	n.doneOnce.Do(func() { close(n.done) })
	return nil
}

func (n *Node) Wait() error {
	<-n.done
	n.log.Info("META node exited", "err", n.err)
	return n.err
}
