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

package metaxml

import (
	"bytes"
	"os"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

func TestConvertXML(t *testing.T) {
	// create a context
	context := []*cid.Cid{
		meta.MustEncode(map[string]string{"foo": "bar"}).Cid(),
	}

	// convert the XML with the context, storing objects in memory
	store := meta.NewMapDatastore()
	converter := NewConverter(store)
	source := "test"
	xml, err := converter.ConvertXML(bytes.NewReader(testXML), context, source)
	if err != nil {
		t.Fatal(err)
	}

	// check it got stored
	get := func(cid *cid.Cid) *meta.Object {
		obj, err := store.Get(cid)
		if err != nil {
			t.Fatal(err)
		}
		return obj
	}
	xml = get(xml.Cid())

	assertString := func(obj *meta.Object, key, expected string) {
		actual, err := obj.GetString(key)
		if err != nil {
			t.Fatal(err)
		}
		if actual != expected {
			t.Fatalf("expected obj[%q] to be %q, got %q", key, expected, actual)
		}
	}
	assertType := func(obj *meta.Object, typ string) {
		assertString(obj, "@type", typ)
	}
	assertSource := func(obj *meta.Object) {
		assertString(obj, "@source", source)
	}
	assertLink := func(obj *meta.Object, key string) *meta.Object {
		link, err := obj.GetLink(key)
		if err != nil {
			t.Fatal(err)
		}
		v := get(link.Cid)
		assertType(v, key)
		return v
	}
	assertContext := func(obj *meta.Object) {
		list, err := obj.GetList("@context")
		if err != nil {
			t.Fatal(err)
		}
		if len(list) != 1 {
			t.Fatalf("expected @context to have 1 entry, got %d", len(list))
		}
		cid, ok := list[0].(*cid.Cid)
		if !ok {
			t.Fatalf("expected @context[0] to have type *cid.Cid, got %T", list[0])
		}
		if !cid.Equals(context[0]) {
			t.Fatalf("expected @context[0] to be %s, got %s", context[0], cid)
		}
	}

	// check the object is META XML
	assertType(xml, "meta:xml")

	// check it has the context
	assertContext(xml)

	// check it has the source
	assertSource(xml)

	// check the XML has a single catalog
	catalog := assertLink(xml, "catalog")
	assertContext(catalog)

	// check the catalog has a single product
	product := assertLink(catalog, "product")
	assertContext(product)

	// check the product's attributes
	assertSource(product)
	assertString(product, "description", "Cardigan Sweater")
	assertString(product, "product_image", "cardigan.jpg")

	// check the product has two items
	itemList, err := product.GetList("catalog_item")
	if err != nil {
		t.Fatal(err)
	}
	if len(itemList) != 2 {
		t.Fatalf("expected product to have 2 items, got %d", len(itemList))
	}
	itemLinks := make([]*cid.Cid, 2)
	for i, v := range itemList {
		link, ok := v.(*cid.Cid)
		if !ok {
			t.Fatalf("expected *cid.Cid, got %T", v)
		}
		itemLinks[i] = link
	}

	// check the two items have the correct properties
	item := get(itemLinks[0])
	assertContext(item)
	assertSource(item)
	assertType(item, "catalog_item")
	assertString(item, "gender", "Men's")
	id := assertLink(item, "item_number")
	assertType(id, "item_number")
	assertSource(id)
	assertString(id, "@value", "QWZ5671")

	item = get(itemLinks[1])
	assertContext(item)
	assertSource(item)
	assertType(item, "catalog_item")
	assertString(item, "gender", "Women's")
	id = assertLink(item, "item_number")
	assertSource(id)
	assertType(id, "item_number")
	assertString(id, "@value", "RRX9856")
}

func TestConvertXMLSchema(t *testing.T) {
	f, err := os.Open("testdata/xmldsig-core-schema.xsd")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	store := meta.NewMapDatastore()
	converter := NewConverter(store)
	obj, err := converter.ConvertXMLSchema(f, "ds", "http://www.w3.org/2000/09/xmldsig#", "test")
	if err != nil {
		t.Fatal(err)
	}

	graph := meta.NewGraph(store, obj)

	// check some expected terms appear in the context
	expectedTerms := []string{
		"CryptoBinary",
		"Signature",
		"SignatureType",
		"SignatureValue",
		"Reference",
		"ReferenceType",
	}
	for _, typ := range expectedTerms {
		v, err := graph.Get("@context", typ)
		if err != nil {
			t.Fatal(err)
		}
		id, ok := v.(string)
		if !ok {
			t.Fatalf("expected context value to be a string, got %T", v)
		}
		expected := "ds:" + typ
		if id != expected {
			t.Fatalf("expected @context[%q] to be %q, got %q", typ, expected, id)
		}
	}
}

// testXML is used to test encoding XML, adapted from
// http://www.service-architecture.com/articles/object-oriented-databases/xml_file_for_complex_data.html
var testXML = []byte(`
<?xml version="1.0" encoding="utf-8" ?>
<catalog>
   <product description="Cardigan Sweater" product_image="cardigan.jpg">
      <catalog_item gender="Men's">
         <item_number>QWZ5671</item_number>
      </catalog_item>
      <catalog_item gender="Women's">
         <item_number>RRX9856</item_number>
      </catalog_item>
   </product>
</catalog>
`[1:])
