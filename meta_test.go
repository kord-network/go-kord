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
	"reflect"
	"runtime"
	"testing"

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

func benchmarkEncode(n int, t *testing.B) {
	type test struct {
		Array []string `json:"array"`
	}
	v := &test{}
	for j := 0; j < n; j++ {
		v.Array = append(v.Array, "t")
	}
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {

		_, err := Encode(v)
		if err != nil {
			t.Fatal(err)
		}
	}
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
}

func benchmarkJsonEncode(n int, t *testing.B) {
	type test struct {
		Array []string `json:"array"`
	}
	v := &test{}
	for j := 0; j < n; j++ {
		v.Array = append(v.Array, "t")
	}
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {

		_, err := json.Marshal(v)

		if err != nil {
			t.Fatal(err)
		}
	}
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
}

func BenchmarkEncode1(b *testing.B)    { benchmarkEncode(1, b) }
func BenchmarkEncode10(b *testing.B)   { benchmarkEncode(10, b) }
func BenchmarkEncode100(b *testing.B)  { benchmarkEncode(100, b) }
func BenchmarkEncode1000(b *testing.B) { benchmarkEncode(1000, b) }

func BenchmarkJsonEncode1(b *testing.B)    { benchmarkJsonEncode(1, b) }
func BenchmarkJsonEncode10(b *testing.B)   { benchmarkJsonEncode(10, b) }
func BenchmarkJsonEncode100(b *testing.B)  { benchmarkJsonEncode(100, b) }
func BenchmarkJsonEncode1000(b *testing.B) { benchmarkJsonEncode(1000, b) }
