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
