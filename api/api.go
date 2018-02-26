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

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/meta-network/go-meta/graph"
	graphql "github.com/neelance/graphql-go"
)

type API struct {
	schema *graphql.Schema
}

func NewAPI(driver *graph.Driver) (*API, error) {
	resolver := NewResolver(driver)
	schema, err := graphql.ParseSchema(GraphQLSchema, resolver)
	if err != nil {
		return nil, err
	}
	return &API{
		schema: schema,
	}, nil
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, fmt.Sprintf("error decoding request: %s", err), http.StatusBadRequest)
		return
	}

	swarmHash := common.Hash{}
	ctx := context.WithValue(r.Context(), "swarmHash", &swarmHash)
	response := a.schema.Exec(ctx, params.Query, params.OperationName, params.Variables)

	if response.Extensions == nil {
		response.Extensions = make(map[string]interface{})
	}
	response.Extensions["meta"] = map[string]interface{}{"swarmHash": swarmHash}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
