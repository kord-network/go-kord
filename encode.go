package meta

import (
	"bytes"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-multihash"
	"github.com/whyrusleeping/cbor/go"
)

func Encode(properties Properties) (*Object, error) {
	p := make(map[string]interface{}, len(properties))
	for key, val := range properties {
		switch val.(type) {
		case string, []byte, *cid.Cid, map[string]string, map[string]*cid.Cid, []*cid.Cid, *Object:
			p[key] = val
		default:
			return nil, fmt.Errorf("meta: unsupported property value: %T", val)
		}
	}

	var buf bytes.Buffer
	enc := cbor.NewEncoder(&buf)
	enc.SetFilter(cbornode.EncoderFilter)
	if err := enc.Encode(p); err != nil {
		return nil, err
	}
	data := buf.Bytes()

	cid, err := cid.Prefix{
		Version:  1,
		Codec:    cid.DagCBOR,
		MhType:   multihash.SHA2_256,
		MhLength: -1,
	}.Sum(data)
	if err != nil {
		return nil, err
	}

	return NewObject(cid, data)
}

func MustEncode(properties Properties) *Object {
	obj, err := Encode(properties)
	if err != nil {
		panic(err)
	}
	return obj
}
