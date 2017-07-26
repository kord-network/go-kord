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

/*
Package meta encodes sets of properties as META objects, stores them in a
key-value store, and provides a mechanism for traversing through graphs of
those objects.

The Object type is an immutable representation of a META object which has a
Content IDentifier (a.k.a. CID, see https://github.com/ipld/cid) and a raw
byte representation which is the IPLD Canonical CBOR format
(see https://github.com/ipld/specs/tree/master/ipld#canonical-format).

To encode a set of properties:

	obj, err := Encode(Properties{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
	})

The object's CID can be retrieved by calling `obj.Cid()` and the raw CBOR
representation with `obj.RawData()`.

Objects can be linked by assigning an object's CID as the value of a property:

	jane := MustEncode(Properties{"name": "Jane"})
	john := MustEncode(Properties{"name": "John"})
	jack := MustEncode(Properties{"name": "Jack"})

	me := MustEncode(Properties{
		"sister": jane.Cid(),
		"children": []*cid.Cid{
			john.Cid(),
			jack.Cid(),
		},
	})

Encoding the object as JSON is valid IPLD:

	json.MarshalIndent(me, "", "  ")

	{
          "children": [
            {
              "/": "zdpuAoqDTaSJuFifCN1EDYexQhdz3b4WchYNfvcDWHoEMcXii"
            },
            {
              "/": "zdpuAqWtnAKfAG7RPihvYE4c9iR2nK28hHcJYBFqGnMEpSaeS"
            }
          ],
          "sister": {
            "/": "zdpuArJB9DZzwceaB91z5RE9v6ALhn5tkQxhWw9zTEVBqnkd4"
          }
        }

Objects can be stored and retrieved using a Store object:

	store := NewStore(datastore.NewMapDatastore())

	err := store.Put(obj)

	obj, err := store.Get(cid)

An object graph can be traversed using a Graph object:

	graph := NewGraph(store, root)

	v, err := graph.Get("some", "path", "through", "the", "graph")

*/
package meta
