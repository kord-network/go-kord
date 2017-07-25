package meta

import (
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"
)

type ErrPathNotFound struct {
	Path []string
}

func (e ErrPathNotFound) Error() string {
	return fmt.Sprintf("meta: path not found: %s", strings.Join(e.Path, "/"))
}

type ErrInvalidCidVersion struct {
	Version uint64
}

func (e ErrInvalidCidVersion) Error() string {
	return fmt.Sprintf("meta: invalid CID version: %d", e.Version)
}

type ErrInvalidCodec struct {
	Codec uint64
}

func (e ErrInvalidCodec) Error() string {
	return fmt.Sprintf("meta: invalid CID codec: %x", e.Codec)
}

type ErrInvalidType struct {
	Type interface{}
}

func (e ErrInvalidType) Error() string {
	return fmt.Sprintf("meta: field @type is not a string (has type %T)", e.Type)
}

type ErrCidMismatch struct {
	Expected *cid.Cid
	Actual   *cid.Cid
}

func (e ErrCidMismatch) Error() string {
	return fmt.Sprintf("meta: CID mismatch, expected %q, got %q", e.Expected, e.Actual)
}
