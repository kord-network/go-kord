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

package meta

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ipfs/go-cid"
)

func TestObjectJSON(t *testing.T) {
	children := []*cid.Cid{
		MustEncode(Properties{"name": "child0"}).Cid(),
		MustEncode(Properties{"name": "child1"}).Cid(),
	}
	parent := MustEncode(Properties{
		"name":     "parent",
		"children": children,
	})

	data, err := json.MarshalIndent(parent, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte(`
{
  "children": [
    {
      "/": "zdpuAvQ9wYysgruG4D7iqv8Rvm6n3tLAtWkC6MJrGMdddyxuY"
    },
    {
      "/": "zdpuAxdhwSiu1J3ZE5gxKrhpU9QxVcAqEJWrTySCg7K3GyPUC"
    }
  ],
  "name": "parent"
}`[1:])
	if !bytes.Equal(data, expected) {
		t.Fatalf("unexpected JSON:\nexpected: %v\nactual:   %v", data, expected)
	}
}
