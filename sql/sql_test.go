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

package sql

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/sql/sqltest"
	"github.com/meta-network/go-meta/testutil"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		dpa, err := testutil.NewTestDPA()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error creating test storage: %s", err)
			return 1
		}
		defer dpa.Cleanup()
		sql.Register("meta", NewDriver(dpa.DPA, &testutil.ENS{}, dpa.Dir))
		return m.Run()
	}())
}

func TestSQL(t *testing.T) {
	sqltest.TestAll(t, "meta", newDB, conf)
}

func newDB(t testing.TB) (string, graph.Options, func()) {
	return fmt.Sprintf("%d.test.meta", rand.Int()), nil, func() {}
}

// TODO: remove the need to set TimeInNs by updating Cayley to support
//       nanosecond precision in the tests here:
//       https://github.com/cayleygraph/cayley/blob/master/graph/graphtest/graphtest.go#L613-L636
var conf = &sqltest.Config{
	TimeInNs: true,
}
