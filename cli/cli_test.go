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
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/meta-network/go-meta/ens"
	metanode "github.com/meta-network/go-meta/node"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
}

func TestLoad(t *testing.T) {
	// start eth node
	ensConfig := ens.DefaultConfig
	ethURL, cleanup := startEthNode(t, ensConfig.Key)
	defer cleanup()

	// deploy ENS
	ensConfig.URL = ethURL
	if err := ens.Deploy(ensConfig, log.New()); err != nil {
		t.Fatal(err)
	}

	// start meta node
	node, cleanup := startMetaNode(t, &ensConfig)
	defer cleanup()

	// create test.meta
	if err := RunCreate(context.Background(), Args(map[string]interface{}{
		"meta":   true,
		"create": true,
		"--url":  fmt.Sprintf("http://%s", node.HTTPAddr()),
		"<db>":   "test.meta",
	})); err != nil {
		t.Fatal(err)
	}

	// load test data
	if err := RunLoad(context.Background(), Args(map[string]interface{}{
		"meta":   true,
		"load":   true,
		"--url":  fmt.Sprintf("http://%s", node.HTTPAddr()),
		"<file>": "../db/data/testdata.nq",
		"<db>":   "test.meta",
	})); err != nil {
		t.Fatal(err)
	}
}

func startEthNode(t *testing.T, key *ecdsa.PrivateKey) (string, func()) {
	tmpDir, err := ioutil.TempDir("", "meta-cli-test")
	if err != nil {
		t.Fatal(err)
	}
	config := node.DefaultConfig
	config.Name = "geth"
	config.IPCPath = "geth.ipc"
	config.DataDir = tmpDir
	config.P2P.MaxPeers = 0
	config.P2P.ListenAddr = ":0"
	config.P2P.NoDiscovery = true
	config.P2P.DiscoveryV5 = false
	config.P2P.NAT = nil
	config.NoUSB = true
	config.UseLightweightKDF = true

	ethConfig := eth.DefaultConfig

	stack, err := node.New(&config)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	acct, err := ks.ImportECDSA(key, "")
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}
	if err := ks.Unlock(acct, ""); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	ethConfig.Genesis = core.DeveloperGenesisBlock(0, acct.Address)
	ethConfig.GasPrice = big.NewInt(1)

	var ethereum *eth.Ethereum
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		var err error
		ethereum, err = eth.New(ctx, &ethConfig)
		return ethereum, err
	}); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	if err := stack.Start(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	if err := ethereum.StartMining(true); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	return filepath.Join(tmpDir, "geth.ipc"), func() {
		stack.Stop()
		os.RemoveAll(tmpDir)
	}
}

func startMetaNode(t *testing.T, ensConfig *ens.Config) (*metanode.Node, func()) {
	tmpDir, err := ioutil.TempDir("", "meta-cli-test")
	if err != nil {
		t.Fatal(err)
	}
	config := metanode.DefaultConfig
	config.DataDir = tmpDir
	config.API.HTTPPort = 0
	config.ENS = *ensConfig
	node := metanode.New(config)
	if err := node.Start(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}
	return node, func() {
		node.Stop()
		os.RemoveAll(tmpDir)
	}
}
