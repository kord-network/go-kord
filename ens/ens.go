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
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"
	"sync"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/ens"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	DevRegistryAddr = common.HexToAddress("0x241be96854Fc2f0172dAA660EE7A14410957C15d")
	DevResolverAddr = common.HexToAddress("0xD277b08f085121d287878A991e0C496488AAaEc6")
	DevKey          = mustKey("476e921a198fd2744f270da0bb80dce2dab24e9105473d9bb19e540fcbd04bb0")
	DevAddress      = common.HexToAddress("0xEe078019375fBFEef8f6278754d54Cf415e83329")
)

type ENS interface {
	Register(name string) error
	Content(name string) (common.Hash, error)
	SetContent(name string, hash common.Hash) error
	SubscribeContent(name string, updates chan common.Hash) (Subscription, error)
}

type Subscription interface {
	Close() error
	Err() error
}

type Config struct {
	URL          string
	Key          *ecdsa.PrivateKey
	RegistryAddr common.Address
	ResolverAddr common.Address
}

var DefaultConfig = Config{
	URL:          "./dev/geth.ipc",
	Key:          DevKey,
	RegistryAddr: DevRegistryAddr,
	ResolverAddr: DevResolverAddr,
}

type Client struct {
	ens          *ens.ENS
	client       *rpc.Client
	ethClient    *ethclient.Client
	blocks       event.Feed
	transactOpts *bind.TransactOpts
	resolverAddr common.Address
	closed       chan struct{}
	closeOnce    sync.Once
}

func NewClient() (*Client, error) {
	return NewClientWithConfig(DefaultConfig)
}

func NewClientWithConfig(config Config) (*Client, error) {
	if config.Key == nil {
		key, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}
		config.Key = key
	}

	client, err := rpc.Dial(config.URL)
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(client)

	transactOpts := bind.NewKeyedTransactor(config.Key)
	transactOpts.GasLimit = params.GenesisGasLimit

	ens, err := ens.NewENS(
		transactOpts,
		config.RegistryAddr,
		ethClient,
	)
	if err != nil {
		client.Close()
		return nil, err
	}

	c := &Client{
		ens:          ens,
		client:       client,
		ethClient:    ethClient,
		transactOpts: transactOpts,
		resolverAddr: DevResolverAddr,
		closed:       make(chan struct{}),
	}
	if err := c.subscribeBlocks(); err != nil {
		client.Close()
		return nil, err
	}
	return c, nil
}

func (c *Client) Register(name string) error {
	if err := c.register(name); err != nil {
		return err
	}
	return c.setResolver(name)
}

func (c *Client) Content(name string) (common.Hash, error) {
	return c.ens.Resolve(name)
}

func (c *Client) SetContent(name string, hash common.Hash) error {
	return c.setContent(name, hash)
}

func (c *Client) SubscribeContent(name string, updates chan common.Hash) (Subscription, error) {
	// TODO
	return nil, nil
}

func (c *Client) Close() {
	c.closeOnce.Do(func() { close(c.closed) })
	c.client.Close()
}

func (c *Client) register(name string) error {
	_, err := c.do(func() (*types.Transaction, error) {
		return c.ens.Register(name)
	})
	return err
}

func (c *Client) setResolver(name string) error {
	_, err := c.do(func() (*types.Transaction, error) {
		return c.ens.SetResolver(Namehash(name), c.resolverAddr)
	})
	return err
}

func (c *Client) setContent(name string, hash common.Hash) error {
	_, err := c.do(func() (*types.Transaction, error) {
		return c.ens.SetContentHash(name, hash)
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
			receipt, err := c.ethClient.TransactionReceipt(context.Background(), tx.Hash())
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
	sub, err := c.ethClient.SubscribeNewHead(context.Background(), heads)
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

func Sha3(s string) common.Hash {
	return crypto.Keccak256Hash([]byte(s))
}

func Namehash(name string) (node common.Hash) {
	if name != "" {
		parts := strings.Split(name, ".")
		for i := len(parts) - 1; i >= 0; i-- {
			label := Sha3(parts[i])
			node = crypto.Keccak256Hash(append(node[:], label[:]...))
		}
	}
	return
}
