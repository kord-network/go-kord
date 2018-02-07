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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/log"
	"github.com/meta-network/go-meta/ens"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
}

func TestLoad(t *testing.T) {
	// generate test config
	tmpDir, err := ioutil.TempDir("", "meta-cli-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	ks := keystore.NewKeyStore(
		filepath.Join(tmpDir, "keystore"),
		keystore.LightScryptN,
		keystore.LightScryptP,
	)
	if _, err := ks.ImportECDSA(ens.DevKey, ""); err != nil {
		t.Fatal(err)
	}

	// start a node
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := Run(ctx, "node", "--datadir", tmpDir, "--dev"); err != nil {
			log.Error("error running node", "err", err)
		}
	}()

	// wait for the node to start
	ipcPath := filepath.Join(tmpDir, "meta.ipc")
	for start := time.Now(); time.Since(start) < 10*time.Second; time.Sleep(50 * time.Millisecond) {
		if _, err := os.Stat(ipcPath); err == nil {
			break
		}
	}

	// deploy ENS
	if err := ens.Deploy(ipcPath, ens.DefaultConfig, log.New()); err != nil {
		t.Fatal(err)
	}

	// create a graph
	if err := Run(context.Background(), "create", "--url", ipcPath, "test.meta"); err != nil {
		t.Fatal(err)
	}

	// load test data
	if err := Run(context.Background(), "load", "--url", ipcPath, "../graph/data/testdata.nq", "test.meta"); err != nil {
		t.Fatal(err)
	}
}
