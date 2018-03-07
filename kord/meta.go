// This file is part of the go-kord library.
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

package kord

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/schema"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/swarm"
	"github.com/kord-network/go-kord/api"
	"github.com/kord-network/go-kord/dapp"
	"github.com/kord-network/go-kord/graph"
	"github.com/kord-network/go-kord/pkg/uri"
	"github.com/kord-network/go-kord/registry"
	"github.com/rs/cors"
)

type Config struct {
	HTTPAddr    string
	HTTPPort    int
	RootDapp    string
	CORSDomains []string
}

var DefaultConfig = Config{
	HTTPAddr: "localhost",
	HTTPPort: 5000,
}

type Kord struct {
	driver   *graph.Driver
	registry registry.Registry
	config   *Config
	srv      *http.Server
	kordSrv  *Server
}

func New(ctx *node.ServiceContext, stack *node.Node, cfg *Config) (*Kord, error) {
	var swarm *swarm.Swarm
	if err := ctx.Service(&swarm); err != nil {
		return nil, fmt.Errorf("error getting Swarm service: %s", err)
	}
	dir := ctx.ResolvePath("db")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	registry := &lazyRegistry{stack: stack}
	driver := graph.NewDriver("kord", swarm.DPA(), registry, dir)
	api, err := api.NewAPI(driver)
	if err != nil {
		return nil, err
	}
	return &Kord{
		driver:   driver,
		registry: registry,
		config:   cfg,
		kordSrv:  NewServer(api, swarm.Api()),
	}, nil
}

func (m *Kord) Protocols() []p2p.Protocol {
	return nil
}

func (m *Kord) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "kord",
			Version:   "0.1",
			Service:   NewPublicAPI(m),
			Public:    true,
		},
	}
}

func (m *Kord) Start(_ *p2p.Server) error {
	if m.config.RootDapp != "" {
		if err := m.setRootDapp(m.config.RootDapp); err != nil {
			return err
		}
	}

	addr := fmt.Sprintf("%s:%d", m.config.HTTPAddr, m.config.HTTPPort)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Info("starting KORD HTTP server", "addr", ln.Addr().String())
	m.srv = &http.Server{
		Addr:    ln.Addr().String(),
		Handler: m.kordSrv,
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

func (m *Kord) Stop() error {
	if m.srv != nil {
		log.Info("stopping KORD HTTP server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return m.srv.Shutdown(ctx)
	}
	return nil
}

func (m *Kord) setRootDapp(dappURI string) error {
	u, err := uri.Parse(dappURI)
	if err != nil {
		return err
	}
	qs, err := m.driver.Get(u.ID.Hex())
	if err != nil {
		return err
	}
	var dapp dapp.Dapp
	path := path.StartPathNodes(qs, qs.ValueOf(quad.IRI(dappURI)))
	if err := schema.LoadPathTo(context.Background(), qs, &dapp, path); err != nil {
		return err
	}
	m.kordSrv.setDapp(&dapp)
	return nil
}

type lazyRegistry struct {
	registry.Registry

	mtx   sync.Mutex
	stack *node.Node
}

func (r *lazyRegistry) registry() (registry.Registry, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.Registry != nil {
		return r.Registry, nil
	}
	client, err := r.stack.Attach()
	if err != nil {
		return nil, err
	}
	registry, err := registry.NewClient(client, registry.DefaultConfig)
	if err != nil {
		return nil, err
	}
	r.Registry = registry
	return registry, nil
}

func (r *lazyRegistry) Graph(kordID common.Address) (common.Hash, error) {
	registry, err := r.registry()
	if err != nil {
		return common.Hash{}, err
	}
	return registry.Graph(kordID)
}

func (r *lazyRegistry) SetGraph(graph common.Hash, sig []byte) error {
	registry, err := r.registry()
	if err != nil {
		return err
	}
	return registry.SetGraph(graph, sig)
}

func (r *lazyRegistry) SubscribeGraph(kordID common.Address, updates chan common.Hash) (registry.Subscription, error) {
	registry, err := r.registry()
	if err != nil {
		return nil, err
	}
	return registry.SubscribeGraph(kordID, updates)
}
