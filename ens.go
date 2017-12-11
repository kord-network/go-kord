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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
)

// ErrNameNotExist is returned when a name cannot be resolved because it does
// not exist.
var ErrNameNotExist = errors.New("meta: ENS name does not exist")

// ENS supports resolving and updating ENS names.
type ENS interface {
	Resolve(name string) (common.Hash, error)
	SetContentHash(name string, hash common.Hash) error
}

// LocalENS returns a local ENS object which stores name to hash mappings
// on the local filesystem.
func LocalENS(dir string) ENS {
	return &localENS{dir}
}

type localENS struct {
	dir string
}

func (l *localENS) Resolve(name string) (common.Hash, error) {
	hash, err := ioutil.ReadFile(l.path(name))
	if os.IsNotExist(err) {
		return common.Hash{}, ErrNameNotExist
	} else if err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(string(hash)), nil
}

func (l *localENS) SetContentHash(name string, hash common.Hash) error {
	return ioutil.WriteFile(l.path(name), []byte(hash.String()), 0644)
}

func (l *localENS) path(name string) string {
	return filepath.Join(l.dir, name)
}
