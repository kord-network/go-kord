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

// xmlschema is a collection of pre-generated Content IDentifiers which
// can be used as the context for META objects which originate from
// XML documents
package xmlschema

import "github.com/ipfs/go-cid"

// Generated with:
//
// $ meta import xsd xs \
//     https://www.w3.org/2009/XMLSchema \
//     <(curl -fSL https://www.w3.org/2009/XMLSchema/XMLSchema.xsd)
//
var XML_Schema = Schema{
	URI: "http://www.w3.org/2001/XMLSchema",
	Cid: mustCid("zdpuAniQsjiwxZyLdBBsqhTbrf2YTuoZNnXQctZrYytAFUU7D"),
}

// Generated with:
//
// $ meta import xsd ds \
//     http://www.w3.org/2000/09/xmldsig# \
//     <(curl -fSL https://www.w3.org/TR/2002/REC-xmldsig-core-20020212/xmldsig-core-schema.xsd)
//
var XML_Dsig = Schema{
	URI: "http://www.w3.org/2000/09/xmldsig#",
	Cid: mustCid("zdpuAz5xgov4sKBWFXvXz5h9kLAX2Tqnt2yBD6Ea79Wq3exfu"),
}

// Generated with:
//
// $ meta import xsd avs http://ddex.net/xml/avs/avs
//
var DDEX_Avs = Schema{
	URI: "http://ddex.net/xml/avs/avs",
	Cid: mustCid("zdpuAyfSqxq1pSSumk2QKFuz28HCkTaRmngvV52ku6Rsdq4Nh"),
}

// Generated with:
//
// $ meta import xsd ern \
//     http://ddex.net/xml/ern/382 \
//     <(curl -fSL http://service.ddex.net/xml/ern/382/release-notification.xsd)
//
var DDEX_Ern382 = Schema{
	URI: "http://ddex.net/xml/ern/382",
	Cid: mustCid("zdpuAo4f7WzDbiHVfigrWoywgzZB3xegd4Vc9BczgjJBUgaK7"),
}

// Generated with:
//
// $ meta import xsd eidr \
//	http://www.eidr.org/schema \
//	<(curl http://www.eidr.org/schema/common.xsd)
var EIDR_common = Schema{
	URI: "http://www.eidr.org/schema",
	Cid: mustCid("zdpuAxFejco7Jim9RaKCE5UyBDdJYVKp4az5UVB8ZcYjrNWpJ"),
}

// Generated with:
// $ meta import xsd md \
//	http://www.movielabs.com/schema/md/v2.1/md \
// 	<(curl http://www.eidr.org/schema/md-v21-eidr.xsd)
//
var EIDR_md = Schema{
	URI: "http://www.movielabs.com/schema/md/v2.1/md",
	Cid: mustCid("zdpuAqqfUHqQDEzCj8nDLxdmXiqpKXJ9vCVrZRrK8e8MewRRj"),
}

type Schema struct {
	URI string
	Cid *cid.Cid
}

func mustCid(v string) *cid.Cid {
	cid, err := cid.Decode(v)
	if err != nil {
		panic(err)
	}
	return cid
}
