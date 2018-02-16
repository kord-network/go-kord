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

package meta

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/swarm"
	"github.com/julienschmidt/httprouter"
	"github.com/meta-network/go-meta/ens"
	"github.com/meta-network/go-meta/graph"
	"github.com/meta-network/go-meta/identity"
	"github.com/rs/cors"
)

type Config struct {
	HTTPAddr    string
	HTTPPort    int
	CORSDomains []string
}

var DefaultConfig = Config{
	HTTPAddr: "localhost",
	HTTPPort: 5000,
}

type Meta struct {
	driver *graph.Driver
	config *Config
	srv    *http.Server
}

func New(ctx *node.ServiceContext, stack *node.Node, cfg *Config) (*Meta, error) {
	var swarm *swarm.Swarm
	if err := ctx.Service(&swarm); err != nil {
		return nil, fmt.Errorf("error getting Swarm service: %s", err)
	}
	dir := ctx.ResolvePath("db")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	driver := graph.NewDriver("meta", swarm.DPA(), &lazyENS{stack: stack}, dir)
	return &Meta{driver: driver, config: cfg}, nil
}

func (m *Meta) Protocols() []p2p.Protocol {
	return nil
}

func (m *Meta) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "meta",
			Version:   "0.1",
			Service:   NewPublicAPI(m),
			Public:    true,
		},
	}
}

func (m *Meta) Start(_ *p2p.Server) error {
	router := httprouter.New()

	identityAPI, err := identity.NewAPI(m.driver)
	if err != nil {
		return err
	}
	router.Handler("GET", "/meta-id/*path", http.StripPrefix("/meta-id", identityAPI))
	router.Handler("POST", "/meta-id/*path", http.StripPrefix("/meta-id", identityAPI))

	// handle OPTIONS requests
	router.OPTIONS("/", func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		w.Header().Set("Allow", "OPTIONS, GET, HEAD, POST")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	})

	addr := fmt.Sprintf("%s:%d", m.config.HTTPAddr, m.config.HTTPPort)
	log.Info("starting META HTTP server", "addr", addr)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	m.srv = &http.Server{
		Addr:    ln.Addr().String(),
		Handler: router,
	}

	if len(m.config.CORSDomains) > 0 {
		log.Info("configuring Cross-Origin Resource Sharing", "domains", m.config.CORSDomains)
		m.srv.Handler = cors.New(cors.Options{
			AllowedOrigins: m.config.CORSDomains,
			AllowedMethods: []string{"POST", "GET", "DELETE", "PATCH", "PUT"},
			MaxAge:         600,
			AllowedHeaders: []string{"*"},
		}).Handler(m.srv.Handler)
	}

	go func() {
		if err := m.srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Error("error starting HTTP server", "err", err)
		}
	}()

	return nil
}

func (m *Meta) Stop() error {
	if m.srv != nil {
		log.Info("stopping META HTTP server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return m.srv.Shutdown(ctx)
	}
	return nil
}

type lazyENS struct {
	ens.ENS

	mtx   sync.Mutex
	stack *node.Node
}

func (e *lazyENS) ens() (ens.ENS, error) {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	if e.ENS != nil {
		return e.ENS, nil
	}
	client, err := e.stack.Attach()
	if err != nil {
		return nil, err
	}
	ens, err := ens.NewClient(client)
	if err != nil {
		return nil, err
	}
	e.ENS = ens
	return ens, nil
}

func (e *lazyENS) Register(name string) error {
	ens, err := e.ens()
	if err != nil {
		return err
	}
	return ens.Register(name)
}

func (e *lazyENS) Content(name string) (common.Hash, error) {
	ens, err := e.ens()
	if err != nil {
		return common.Hash{}, err
	}
	return ens.Content(name)
}

func (e *lazyENS) SetContent(name string, hash common.Hash) error {
	ens, err := e.ens()
	if err != nil {
		return err
	}
	return ens.SetContent(name, hash)
}

func (e *lazyENS) SubscribeContent(name string, updates chan common.Hash) (ens.Subscription, error) {
	ens, err := e.ens()
	if err != nil {
		return nil, err
	}
	return ens.SubscribeContent(name, updates)
}
