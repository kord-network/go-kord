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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/meta-network/go-meta/registry"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
}

var n *testNode

func TestMain(m *testing.M) {
	os.Exit(func() int {
		// start the test node
		var err error
		n, err = startTestNode()
		if err != nil {
			log.Error("error starting test node", "err", err)
			return 1
		}
		defer n.stop()

		return m.Run()
	}())
}

func TestLoad(t *testing.T) {
	// create an ID
	cliCtx := NewContext(context.Background())
	cliCtx.Stdin = bytes.NewReader([]byte{'\n', '\n'})
	var stdout bytes.Buffer
	cliCtx.Stdout = &stdout
	if err := Run(
		cliCtx,
		"id",
		"new",
		"--keystore", n.keystore,
	); err != nil {
		t.Fatal(err)
	}
	out := strings.TrimSpace(stdout.String())
	if !common.IsHexAddress(out) {
		t.Fatalf("unexpected ID output: %s", out)
	}
	id := common.HexToAddress(out)

	// create a graph
	cliCtx = NewContext(context.Background())
	cliCtx.Stdin = bytes.NewReader([]byte{'\n'})
	if err := Run(
		cliCtx,
		"graph",
		"create",
		"--url", n.ipcPath,
		"--keystore", n.keystore,
		id.Hex(),
	); err != nil {
		t.Fatal(err)
	}

	// load test data
	cliCtx = NewContext(context.Background())
	cliCtx.Stdin = bytes.NewReader([]byte{'\n'})
	if err := Run(
		cliCtx,
		"graph",
		"load",
		"--url", n.ipcPath,
		"--keystore", n.keystore,
		id.Hex(),
		"../graph/data/testdata.nq",
	); err != nil {
		t.Fatal(err)
	}
}

func TestDapp(t *testing.T) {
	// create an ID
	cliCtx := NewContext(context.Background())
	cliCtx.Stdin = bytes.NewReader([]byte{'\n', '\n'})
	var stdout bytes.Buffer
	cliCtx.Stdout = &stdout
	if err := Run(
		cliCtx,
		"id",
		"new",
		"--keystore", n.keystore,
	); err != nil {
		t.Fatal(err)
	}
	out := strings.TrimSpace(stdout.String())
	if !common.IsHexAddress(out) {
		t.Fatalf("unexpected ID output: %s", out)
	}
	id := common.HexToAddress(out)

	// create a graph
	cliCtx = NewContext(context.Background())
	cliCtx.Stdin = bytes.NewReader([]byte{'\n'})
	if err := Run(
		cliCtx,
		"graph",
		"create",
		"--url", n.ipcPath,
		"--keystore", n.keystore,
		id.Hex(),
	); err != nil {
		t.Fatal(err)
	}

	// deploy a dapp
	dappDir, err := ioutil.TempDir("", "meta-cli-test")
	if err != nil {
		t.Fatal(err)
	}
	dappHTML := []byte(`<html><head><title>Test Dapp</title><body><h1>Test Dapp</h1></body></html>`)
	if err := ioutil.WriteFile(filepath.Join(dappDir, "index.html"), dappHTML, 0644); err != nil {
		t.Fatal(err)
	}
	dappURI := fmt.Sprintf("meta://%s/test-dapp", id.Hex())
	cliCtx = NewContext(context.Background())
	cliCtx.Stdin = bytes.NewReader([]byte{'\n'})
	if err := Run(
		cliCtx,
		"dapp",
		"deploy",
		"--url", n.ipcPath,
		"--keystore", n.keystore,
		dappDir,
		dappURI,
	); err != nil {
		t.Fatal(err)
	}

	// set the dapp as root
	cliCtx = NewContext(context.Background())
	if err := Run(
		cliCtx,
		"dapp",
		"set-root",
		"--url", n.ipcPath,
		dappURI,
	); err != nil {
		t.Fatal(err)
	}

	// check the dapp is available
	res, err := http.Get("http://localhost:5000/")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected HTTP status: %s", res.Status)
	}
	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(html, dappHTML) {
		t.Fatalf(`unexpected HTML:\nexpected: %s\nactual:   %s`, dappHTML, html)
	}
}

type testNode struct {
	keystore string
	ipcPath  string
	stop     func()
}

func startTestNode() (*testNode, error) {
	// generate test config
	tmpDir, err := ioutil.TempDir("", "meta-cli-test")
	if err != nil {
		return nil, err
	}
	ks := keystore.NewKeyStore(
		filepath.Join(tmpDir, "keystore"),
		keystore.LightScryptN,
		keystore.LightScryptP,
	)
	if _, err := ks.ImportECDSA(registry.DevKey, ""); err != nil {
		os.RemoveAll(tmpDir)
		return nil, err
	}

	// start a dev node
	ctx, stopNode := context.WithCancel(context.Background())
	go func() {
		if err := Run(NewContext(ctx), "node", "--dev", "--datadir", tmpDir); err != nil {
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

	// deploy the META registry
	if _, err := registry.Deploy(ipcPath, registry.DefaultConfig); err != nil {
		stopNode()
		os.RemoveAll(tmpDir)
		return nil, err
	}

	return &testNode{
		keystore: filepath.Join(tmpDir, "keystore"),
		ipcPath:  ipcPath,
		stop: func() {
			stopNode()
			os.RemoveAll(tmpDir)
		},
	}, nil
}
