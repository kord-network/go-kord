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
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/meta-network/go-meta"
)

// EncodeXMLSchema encodes an XML Schema document as a META object graph.
func EncodeXMLSchema(src io.Reader, namespace, uri string) (*meta.Object, error) {
	dec := xml.NewDecoder(src)

	context := map[string]string{
		namespace: uri,
	}

	// walk the XML document, adding any namespaces or element types to
	// the context
	for {
		token, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		el, ok := token.(xml.StartElement)
		if !ok {
			continue
		}

		// add XML namespaces
		for _, attr := range el.Attr {
			if attr.Name.Space == "xmlns" {
				context[attr.Name.Local] = attr.Value
			}
		}

		// add element, simple and complex types
		switch el.Name.Local {
		case "element", "simpleType", "complexType":
			for _, attr := range el.Attr {
				if attr.Name.Local == "name" {
					name := attr.Value
					context[name] = fmt.Sprintf("%s:%s", namespace, name)
				}
			}
		}
	}

	return meta.Encode(map[string]interface{}{"@context": context})
}

// EncodeXML encodes an XML document as a META object graph.
func EncodeXML(src io.Reader, context []*cid.Cid, callback func(*meta.Object) error) (*meta.Object, error) {
	dec := xml.NewDecoder(src)

	// read tokens until we find the root element (i.e. the first
	// xml.StartElement)
	var root xml.StartElement
	for {
		token, err := dec.Token()
		if err != nil {
			return nil, err
		}
		if t, ok := token.(xml.StartElement); ok {
			root = t
			break
		}
	}

	// convert the root element
	obj, err := encodeXML(dec, &root, context, callback)
	if err != nil {
		return nil, err
	}

	// wrap it in an XML object
	properties := map[string]interface{}{
		"@type":         "meta:xml",
		root.Name.Local: obj.Cid(),
	}
	if len(context) > 0 {
		properties["@context"] = context
	}
	xml, err := meta.Encode(properties)
	if err != nil {
		return nil, err
	}
	if callback != nil {
		if err := callback(xml); err != nil {
			return nil, err
		}
	}
	return xml, nil
}

func encodeXML(dec *xml.Decoder, el *xml.StartElement, context []*cid.Cid, callback func(*meta.Object) error) (*meta.Object, error) {
	// create a new node with the type as the name of the element
	node := map[string]interface{}{"@type": el.Name.Local}

	// add the context
	if len(context) > 0 {
		node["@context"] = context
	}

	// add the attributes
	for _, attr := range el.Attr {
		key := attr.Name.Local
		if attr.Name.Space != "" {
			key = attr.Name.Space + ":" + key
		}
		node[key] = attr.Value
	}

	// keep decoding until we see the end of the current element
	for {
		token, err := dec.Token()
		if err != nil {
			return nil, err
		}

		switch token := token.(type) {

		// xml.StartElement is the start of a child element so convert
		// it and add it as a property
		case xml.StartElement:
			child, err := encodeXML(dec, &token, context, callback)
			if err != nil {
				return nil, err
			}

			switch v := node[token.Name.Local].(type) {
			case nil:
				node[token.Name.Local] = child.Cid()
			case *cid.Cid:
				node[token.Name.Local] = []*cid.Cid{v, child.Cid()}
			case []*cid.Cid:
				node[token.Name.Local] = append(v, child.Cid())
			}

		// xml.CharData is text data inside the element so treat it
		// like a value object
		case xml.CharData:
			// ignore pure whitespace
			if strings.TrimSpace(string(token)) == "" {
				continue
			}
			if v, ok := node["@value"]; ok {
				node["@value"] = v.(string) + string(token)
			} else {
				node["@value"] = string(token)
			}

		// xml.EndElement marks the end of the current element,
		// return it as a META object
		case xml.EndElement:
			obj, err := meta.Encode(node)
			if err != nil {
				return nil, err
			}
			if callback != nil {
				if err := callback(obj); err != nil {
					return nil, err
				}
			}
			return obj, nil
		}
	}
}
