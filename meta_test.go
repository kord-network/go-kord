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
	"crypto/sha256"
	"encoding/json"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/swarm/testutil"
	"github.com/ipfs/go-cid"
)

func TestObjectJSON(t *testing.T) {
	type Person struct {
		Name     string     `json:"name"`
		Children []*cid.Cid `json:"children,omitempty"`
	}
	obj := MustEncode(&Person{
		Name: "parent",
		Children: []*cid.Cid{
			MustEncode(&Person{Name: "child0"}).Cid(),
			MustEncode(&Person{Name: "child1"}).Cid(),
		},
	})

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte(`
{
  "children": [
    {
      "/": "zdpuB3PyNStxWJ9P7Qho6MYMPinMiKLfCERUdrVuRvdExT6nn"
    },
    {
      "/": "zdpuAxwfhw8WVtTSdumNdrVgWori8BuPHieBq1QkQi4FDtpWZ"
    }
  ],
  "name": "parent"
}`[1:])
	if !bytes.Equal(data, expected) {
		t.Fatalf("unexpected JSON:\nexpected: %s\nactual:   %s", expected, data)
	}
}

func TestEncodeDecode(t *testing.T) {
	type test struct {
		Null   interface{}       `json:"null"`
		Bool   bool              `json:"bool"`
		Int    int64             `json:"int"`
		Float  float64           `json:"float"`
		String string            `json:"string"`
		Array  []string          `json:"array"`
		Map    map[string]string `json:"map"`
	}
	v := &test{
		Null:   nil,
		Bool:   true,
		Int:    42,
		Float:  42.24,
		String: "42",
		Array:  []string{"42", "42", "42"},
		Map:    map[string]string{"foo": "42", "bar": "42"},
	}
	obj, err := Encode(v)
	if err != nil {
		t.Fatal(err)
	}
	w := &test{}
	if err := obj.Decode(w); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v, w) {
		t.Fatalf("decoded object not equal:\nexpected %#v\nactual   %#v", v, w)
	}
}

type testMeta struct {
	store *Store
	srv   *testutil.TestSwarmServer
}

func (tm *testMeta) openSwarmStore(t *testing.T) *Store {
	tm.srv = testutil.NewTestSwarmServer(t)
	return NewSwarmDatastore(tm.srv.URL)
}

func TestSwarmDatastore(t *testing.T) {
	type Person struct {
		Name     string     `json:"name"`
		Children []*cid.Cid `json:"children,omitempty"`
	}

	x := &testMeta{}

	store := x.openSwarmStore(t)

	defer x.srv.Close()

	storeObj := store.MustPut(&Person{
		Name: "parent",
		Children: []*cid.Cid{
			store.MustPut(&Person{Name: "child0"}).Cid(),
			store.MustPut(&Person{Name: "child1"}).Cid(),
		},
	})

	retriveObj, err := store.Get(storeObj.Cid())
	if err != nil {
		t.Fatal(err)
	}

	if retriveObj.Cid() != storeObj.Cid() {
		t.Fatalf("retrived object's cid %s should be equal to the stored object cid %s ", retriveObj.Cid(), storeObj.Cid())
	}

	obj := retriveObj

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte(`
{
  "children": [
    {
      "/": "zdqaWD2eehRTP7Abff7EVHc6Hyhpm379guLTgf1acpW4scgMm"
    },
    {
      "/": "zdqaWT5zs1T1NVFvF49bpAQYqtzUshDmWf9uMN42e1uW4dp88"
    }
  ],
  "name": "parent"
}`[1:])
	if !bytes.Equal(data, expected) {
		t.Fatalf("unexpected JSON:\nexpected: %s\nactual:   %s", expected, data)
	}
}

func benchmarkEncode(n int, t *testing.B) {
	type test struct {
		String string `json:"string"`
	}
	v := &test{}
	v.String = strings.Repeat("t", n)

	t.ReportAllocs()
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		_, err := Encode(v)
		if err != nil {
			t.Fatal(err)
		}
	}

	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
}

func benchmarkJSONEncode(n int, t *testing.B) {
	type test struct {
		String string `json:"string"`
	}
	v := &test{}
	v.String = strings.Repeat("t", n)

	t.ReportAllocs()
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		obj, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		_ = sha256.Sum256(obj)
	}

	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
}

func BenchmarkEncode1(b *testing.B)     { benchmarkEncode(1, b) }
func BenchmarkEncode10(b *testing.B)    { benchmarkEncode(10, b) }
func BenchmarkEncode100(b *testing.B)   { benchmarkEncode(100, b) }
func BenchmarkEncode1000(b *testing.B)  { benchmarkEncode(1000, b) }
func BenchmarkEncode10000(b *testing.B) { benchmarkEncode(10000, b) }

func BenchmarkJSONEncode1(b *testing.B)     { benchmarkJSONEncode(1, b) }
func BenchmarkJSONEncode10(b *testing.B)    { benchmarkJSONEncode(10, b) }
func BenchmarkJSONEncode100(b *testing.B)   { benchmarkJSONEncode(100, b) }
func BenchmarkJSONEncode1000(b *testing.B)  { benchmarkJSONEncode(1000, b) }
func BenchmarkJSONEncode10000(b *testing.B) { benchmarkJSONEncode(10000, b) }
