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

package graph

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/sql/sqltest"
	"github.com/kord-network/go-kord/testutil"
)

var testDriver *Driver

// TestMain runs the Cayley test suite against the Swarm backed SQLite database
// driver.
//
// TODO: support nanosecond precision in the tests here:
//       https://github.com/cayleygraph/cayley/blob/v0.7.1/graph/graphtest/graphtest.go#L655-L673
func TestMain(m *testing.M) {
	os.Exit(func() int {
		dpa, err := testutil.NewTestDPA()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error creating test storage:", err)
			return 1
		}
		defer dpa.Cleanup()
		testDriver = NewDriver("kord-test", dpa.DPA, testutil.NewTestRegistry(), dpa.Dir)
		return m.Run()
	}())
}

func TestSQLBackend(t *testing.T) {
	sqltest.TestAll(t, testDriver.name, newTestDB, nil)
}

func newTestDB(t testing.TB) (string, graph.Options, func()) {
	return fmt.Sprintf("%d.test.kord", rand.Int()), nil, func() {}
}
