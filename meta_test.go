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
	"crypto/ecdsa"
	"database/sql"
	"testing"
	"time"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	_ "github.com/cayleygraph/cayley/writer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	metasql "github.com/meta-network/go-meta/sql"
	"github.com/meta-network/go-meta/testutil"
)

func TestAddQuad(t *testing.T) {
	// setup Swarm SQLite3 driver
	dpa, err := testutil.NewTestDPA()
	if err != nil {
		t.Fatal(err)
	}
	defer dpa.Cleanup()
	driver := metasql.NewDriver(dpa.DPA, &testutil.ENS{}, dpa.Dir)
	sql.Register("meta", driver)

	// generate account
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	// create state
	state := NewState(driver)

	// create client
	client, err := NewClient(
		crypto.PubkeyToAddress(key.PublicKey),
		"test.meta",
		state,
		&signer{key},
		driver,
	)
	if err != nil {
		t.Fatal(err)
	}

	qw, err := graph.NewQuadWriter("single", client, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := qw.AddQuad(quad.Make("phrase of the day", "is of course", "Hello World!", nil)); err != nil {
		t.Fatal(err)
	}

	path := cayley.StartPath(client, quad.String("phrase of the day")).Out(quad.String("is of course"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	results, err := path.Iterate(ctx).All()
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	v := client.NameOf(results[0])
	s, ok := v.Native().(string)
	if !ok {
		t.Fatalf("expected string, got %T", v.Native())
	}
	if s != "Hello World!" {
		t.Fatalf(`expected "Hello World!", got %q`, s)
	}
}

type signer struct {
	privKey *ecdsa.PrivateKey
}

func (s *signer) SignHash(_ common.Address, hash []byte) ([]byte, error) {
	return crypto.Sign(hash, s.privKey)
}
