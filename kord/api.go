// This file is part of the go-kord library.
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

package kord

import (
	"github.com/cayleygraph/cayley/graph"
	"github.com/ethereum/go-ethereum/common"
)

type PublicAPI struct {
	kord *Kord
}

func NewPublicAPI(kord *Kord) *PublicAPI {
	return &PublicAPI{kord}
}

func (api *PublicAPI) CreateGraph(name string) (common.Hash, error) {
	return api.kord.driver.Create(name)
}

func (api *PublicAPI) CommitGraph(name string) (common.Hash, error) {
	return api.kord.driver.Commit(name)
}

func (api *PublicAPI) SetGraph(hash common.Hash, sig []byte) error {
	return api.kord.registry.SetGraph(hash, sig)
}

func (api *PublicAPI) SetRootDapp(dappURI string) error {
	return api.kord.setRootDapp(dappURI)
}

func (api *PublicAPI) HttpAddr() string {
	return api.kord.srv.Addr
}

func (api *PublicAPI) ApplyDeltas(name string, in []graph.Delta, opts graph.IgnoreOpts) (common.Hash, error) {
	qs, err := api.kord.driver.Get(name)
	if err != nil {
		return common.Hash{}, err
	}
	if err := qs.ApplyDeltas(in, opts); err != nil {
		return common.Hash{}, err
	}
	return api.kord.driver.Commit(name)
}
