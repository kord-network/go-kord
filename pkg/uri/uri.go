// This file is part of the go-meta library.
//
// Copyright (C) 2018 JAAK MUSIC LTD
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

package uri

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
)

type URI struct {
	ID   common.Address
	Path string
}

func Parse(s string) (*URI, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "meta" {
		return nil, fmt.Errorf("invalid META URI scheme: %s", u.Scheme)
	}
	if !common.IsHexAddress(u.Host) {
		return nil, fmt.Errorf("invalid META ID in uri: %s", u.Host)
	}
	return &URI{
		ID:   common.HexToAddress(u.Host),
		Path: u.Path,
	}, nil
}
