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
	"crypto/ecdsa"
	"errors"
	"fmt"
	"sync"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/kord-network/go-kord/registry/contract"
)

//go:generate abigen --sol ../contracts/KORDRegistry.sol --pkg contract --out contract/registry.go

var (
	DevKey          = mustKey("476e921a198fd2744f270da0bb80dce2dab24e9105473d9bb19e540fcbd04bb0")
	DevAddr         = crypto.PubkeyToAddress(DevKey.PublicKey)
	DevContractAddr = common.HexToAddress("0x241be96854Fc2f0172dAA660EE7A14410957C15d")
)

type Registry interface {
	Graph(kordID common.Address) (common.Hash, error)
	SetGraph(graph common.Hash, sig []byte) error
	SubscribeGraph(kordID common.Address, updates chan common.Hash) (Subscription, error)
}

type Subscription interface {
	Close() error
	Err() error
}

type Config struct {
	Key          *ecdsa.PrivateKey
	ContractAddr common.Address
}

var DefaultConfig = Config{
	Key:          DevKey,
	ContractAddr: DevContractAddr,
}

type Client struct {
	*ethclient.Client

	registry     *contract.KORDRegistrySession
	blocks       event.Feed
	transactOpts *bind.TransactOpts
	closed       chan struct{}
	closeOnce    sync.Once
}

func NewClient(client *rpc.Client, config Config) (*Client, error) {
	if config.Key == nil {
		key, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}
		config.Key = key
	}

	ethClient := ethclient.NewClient(client)

	transactOpts := bind.NewKeyedTransactor(config.Key)
	transactOpts.GasLimit = params.GenesisGasLimit

	registry, err := contract.NewKORDRegistry(config.ContractAddr, ethClient)
	if err != nil {
		return nil, err
	}
	session := &contract.KORDRegistrySession{
		Contract:     registry,
		TransactOpts: *transactOpts,
	}

	c := &Client{
		Client:       ethClient,
		registry:     session,
		transactOpts: transactOpts,
		closed:       make(chan struct{}),
	}
	if err := c.subscribeBlocks(); err != nil {
		client.Close()
		return nil, err
	}
	return c, nil
}

func (c *Client) Graph(kordID common.Address) (common.Hash, error) {
	return c.registry.Graph(kordID)
}

func (c *Client) SetGraph(graph common.Hash, sig []byte) error {
	return c.setGraph(graph, sig)
}

func (c *Client) SubscribeGraph(kordID common.Address, updates chan common.Hash) (Subscription, error) {
	// TODO
	return nil, nil
}

func (c *Client) Close() {
	c.closeOnce.Do(func() { close(c.closed) })
}

func (c *Client) setGraph(hash common.Hash, sig []byte) error {
	_, err := c.do(func() (*types.Transaction, error) {
		return c.registry.SetGraph(hash, sig)
	})
	return err
}

const blockTimeout = 50

func (c *Client) do(f func() (*types.Transaction, error)) (*types.Receipt, error) {
	heads := make(chan *types.Header)
	sub := c.blocks.Subscribe(heads)
	defer sub.Unsubscribe()

	tx, err := f()
	if err != nil {
		return nil, err
	}
	var count int
	for {
		select {
		case <-heads:
			count++
			if count > blockTimeout {
				return nil, fmt.Errorf("failed to get transaction receipt after %d blocks", blockTimeout)
			}
			receipt, err := c.TransactionReceipt(context.Background(), tx.Hash())
			if err == nil {
				if receipt.Status == types.ReceiptStatusFailed {
					return nil, errors.New("transaction failed")
				}
				return receipt, nil
			} else if err != ethereum.NotFound {
				return nil, err
			}
		case err := <-sub.Err():
			return nil, err
		case <-c.closed:
			return nil, errors.New("client closed")
		}
	}
}

func (c *Client) subscribeBlocks() error {
	heads := make(chan *types.Header)
	sub, err := c.SubscribeNewHead(context.Background(), heads)
	if err != nil {
		return err
	}
	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case head, ok := <-heads:
				if !ok {
					return
				}
				c.blocks.Send(head)
			case <-c.closed:
				return
			}
		}
	}()
	return nil
}

func mustKey(hex string) *ecdsa.PrivateKey {
	key, err := crypto.HexToECDSA(hex)
	if err != nil {
		panic(err)
	}
	return key
}
