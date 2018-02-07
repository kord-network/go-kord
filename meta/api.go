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

package meta

import "github.com/cayleygraph/cayley/graph"

type PublicAPI struct {
	meta *Meta
}

func NewPublicAPI(meta *Meta) *PublicAPI {
	return &PublicAPI{meta}
}

func (api *PublicAPI) CreateGraph(name string) error {
	return api.meta.driver.Create(name)
}

func (api *PublicAPI) ApplyDeltas(name string, in []graph.Delta, opts graph.IgnoreOpts) error {
	qs, err := api.meta.driver.Get(name)
	if err != nil {
		return err
	}
	if err := qs.ApplyDeltas(in, opts); err != nil {
		return err
	}
	return api.meta.driver.Commit(name)
}
