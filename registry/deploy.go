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

package registry

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/kord-network/go-kord/registry/contract"
)

func Deploy(url string, config Config) (common.Address, error) {
	c, err := rpc.Dial(url)
	if err != nil {
		return common.Address{}, err
	}
	defer c.Close()

	client, err := NewClient(c, config)
	if err != nil {
		return common.Address{}, err
	}

	return deployRegistry(client, config.ContractAddr)
}

func deployRegistry(client *Client, addr common.Address) (common.Address, error) {
	if data, err := client.CodeAt(context.Background(), addr, nil); err == nil && len(data) > 0 {
		return addr, nil
	}
	receipt, err := client.do(func() (tx *types.Transaction, err error) {
		_, tx, _, err = contract.DeployKORDRegistry(client.transactOpts, client)
		return
	})
	if err != nil {
		return common.Address{}, err
	}
	return receipt.ContractAddress, nil
}
