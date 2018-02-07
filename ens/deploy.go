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

package ens

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/ens/contract"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

// Deploy deploys an ENS registry, resolver and FIFS registrar.
func Deploy(url string, config Config, log log.Logger) error {
	c, err := rpc.Dial(url)
	if err != nil {
		return err
	}
	defer c.Close()

	client, err := NewClientWithConfig(c, config)
	if err != nil {
		return err
	}

	d := &deployer{client}

	log.Info("deploying ENS registry contract")
	registryAddr, err := d.deployRegistry(config.RegistryAddr)
	if err != nil {
		log.Error("error deploying ENS registry contract", "err", err)
		return err
	}
	log.Info("deployed ENS registry contract", "addr", registryAddr.Hex())

	log.Info("deploying ENS resolver contract")
	resolverAddr, err := d.deployResolver(config.ResolverAddr, registryAddr)
	if err != nil {
		log.Error("error deploying ENS resolver contract", "err", err)
		return err
	}
	log.Info("deployed ENS resolver contract", "addr", resolverAddr.Hex())

	log.Info("deploying ENS .meta registrar contract")
	registrarAddr, err := d.deployMetaRegistrar(registryAddr)
	if err != nil {
		log.Error("error deploying ENS .meta registrar contract", "err", err)
		return err
	}
	log.Info("deployed ENS .meta registrar contract", "address", registrarAddr.Hex())

	return nil
}

type deployer struct {
	*Client
}

func (d *deployer) deployRegistry(addr common.Address) (common.Address, error) {
	if data, err := d.CodeAt(context.Background(), addr, nil); err == nil && len(data) > 0 {
		return addr, nil
	}
	receipt, err := d.do(func() (tx *types.Transaction, err error) {
		_, tx, _, err = contract.DeployENS(d.transactOpts, d)
		return
	})
	if err != nil {
		return common.Address{}, err
	}
	return receipt.ContractAddress, nil
}

func (d *deployer) deployResolver(addr common.Address, registryAddr common.Address) (common.Address, error) {
	if data, err := d.CodeAt(context.Background(), addr, nil); err == nil && len(data) > 0 {
		return addr, nil
	}
	receipt, err := d.do(func() (tx *types.Transaction, err error) {
		_, tx, _, err = contract.DeployPublicResolver(d.transactOpts, d, registryAddr)
		return
	})
	if err != nil {
		return common.Address{}, err
	}
	return receipt.ContractAddress, nil
}

func (d *deployer) deployMetaRegistrar(registryAddr common.Address) (common.Address, error) {
	if addr, err := d.ens.Owner(Namehash("meta")); err == nil && addr != (common.Address{}) {
		return addr, nil
	}
	receipt, err := d.do(func() (tx *types.Transaction, err error) {
		_, tx, _, err = contract.DeployFIFSRegistrar(d.transactOpts, d, registryAddr, Namehash("meta"))
		return
	})
	if err != nil {
		return common.Address{}, err
	}
	registrarAddr := receipt.ContractAddress
	_, err = d.do(func() (*types.Transaction, error) {
		return d.ens.SetSubnodeOwner(common.Hash{}, Sha3("meta"), registrarAddr)
	})
	if err != nil {
		return common.Address{}, err
	}
	return registrarAddr, nil
}
