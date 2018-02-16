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

package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/swarm"
	swarmapi "github.com/ethereum/go-ethereum/swarm/api"
	"github.com/meta-network/go-meta/meta"
	"github.com/naoina/toml"
)

var testnetBootnodes = []string{
	"enode://21cd1409c28106062f79dbae8d9a69d4e1050c6f8a40ab63ec507c03970ed152c6f20708262f23a7334061fde7943b10ead6249bb88b2d7375d36f40ff471e82@35.176.243.138:30303",
}

func init() {
	registerCommand("node", RunNode, `
usage: meta node [--datadir <dir>] [--config <path>] [--dev] [--testnet] [--mine] [--cors-domain <domain>...] [--verbosity <n>]

Run a META node.

options:
	-d, --datadir <dir>         Node data directory
	-c, --config <path>         Path to the TOML config file
	--dev                       Run a dev node
	--testnet                   Connect to the testnet
	--mine                      Mine the Ethereum chain
	--cors-domain <domain>...   The allowed CORS domains
	--verbosity <n>             Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail [default: 3]
`[1:])
}

func RunNode(ctx context.Context, args Args) error {
	cfg := defaultConfig()

	if v := args.String("--verbosity"); v != "" {
		if _, err := setLogVerbosity(v); err != nil {
			return err
		}
	}

	if file := args.String("--config"); file != "" {
		if err := loadConfig(file, &cfg); err != nil {
			return err
		}
	}

	if dir := args.String("--datadir"); dir != "" {
		cfg.Node.DataDir = dir
	}

	if _, ok := args["--cors-domain"]; ok {
		domains := args.List("--cors-domain")
		cfg.Swarm.Cors = strings.Join(domains, ",")
		cfg.Meta.CORSDomains = domains
	}

	if args.Bool("--dev") && args.Bool("--testnet") {
		return errors.New("--dev and --testnet cannot both be set")
	} else if args.Bool("--dev") {
		// --dev mode can't use p2p networking.
		cfg.Node.P2P.MaxPeers = 0
		cfg.Node.P2P.ListenAddr = ":0"
		cfg.Node.P2P.NoDiscovery = true
		cfg.Node.P2P.DiscoveryV5 = false
	} else if args.Bool("--testnet") {
		cfg.Eth.NetworkId = 1035

		if !args.Bool("--mine") {
			cfg.Node.P2P.BootstrapNodes = make([]*discover.Node, 0, len(testnetBootnodes))
			for _, url := range testnetBootnodes {
				node, err := discover.ParseNode(url)
				if err != nil {
					return fmt.Errorf("invalid testnet bootnode: %s: %s", url, err)
				}
				cfg.Node.P2P.BootstrapNodes = append(cfg.Node.P2P.BootstrapNodes, node)
			}
		}
	}

	stack, err := node.New(&cfg.Node)
	if err != nil {
		return err
	}

	if args.Bool("--dev") {
		if err := setupDevAccount(stack, &cfg); err != nil {
			return err
		}
	}

	utils.RegisterEthService(stack, &cfg.Eth)

	if err := registerSwarmService(stack, &cfg.Swarm); err != nil {
		return err
	}

	if err := registerMetaService(stack, &cfg.Meta); err != nil {
		return err
	}

	// start the node
	if err := stack.Start(); err != nil {
		return err
	}

	// start mining if required or in dev mode
	if args.Bool("--mine") || args.Bool("--dev") {
		if err := startMining(stack, &cfg); err != nil {
			stack.Stop()
			return err
		}
	}

	// stop the node if the context is cancelled
	go func() {
		<-ctx.Done()
		stack.Stop()
	}()

	// wait for the node to exit
	stack.Wait()
	return nil
}

func registerSwarmService(stack *node.Node, cfg *swarmapi.Config) error {
	cfg.Path = stack.InstanceDir()

	// load the bzzaccount private key to initialise the config
	//
	// TODO: support getting the password from the user
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	account, err := ks.Find(accounts.Account{Address: common.HexToAddress(cfg.BzzAccount)})
	if err != nil {
		return err
	}
	keyjson, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		return err
	}
	key, err := keystore.DecryptKey(keyjson, "")
	if err != nil {
		return err
	}
	cfg.Init(key.PrivateKey)

	return stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return swarm.NewSwarm(
			ctx,
			nil,
			nil,
			cfg,
			cfg.SwapEnabled,
			cfg.SyncEnabled,
			cfg.Cors,
		)
	})
}

func registerMetaService(stack *node.Node, cfg *meta.Config) error {
	return stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return meta.New(ctx, stack, cfg)
	})
}

type config struct {
	Node  node.Config
	Eth   eth.Config
	Swarm swarmapi.Config
	Meta  meta.Config
}

func loadConfig(file string, cfg *config) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	return err
}

func defaultConfig() config {
	swarmCfg := swarmapi.NewDefaultConfig()
	return config{
		Node:  defaultNodeConfig(),
		Eth:   eth.DefaultConfig,
		Swarm: *swarmCfg,
		Meta:  meta.DefaultConfig,
	}
}

func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = "meta"
	cfg.Version = "0.0.1"
	cfg.HTTPModules = append(cfg.HTTPModules, "eth")
	cfg.WSModules = append(cfg.WSModules, "eth")
	cfg.IPCPath = "meta.ipc"
	return cfg
}

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		link := ""
		if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
			link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
		}
		return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
	},
}

func setupDevAccount(stack *node.Node, cfg *config) error {
	// Create new developer account or reuse existing one
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	var developer accounts.Account
	if accs := ks.Accounts(); len(accs) > 0 {
		developer = accs[0]
	} else {
		var err error
		developer, err = ks.NewAccount("")
		if err != nil {
			return fmt.Errorf("error creating developer account: %s", err)
		}
	}
	log.Info("Using developer account", "address", developer.Address)
	cfg.Swarm.BzzAccount = developer.Address.String()

	cfg.Eth.Genesis = core.DeveloperGenesisBlock(0, developer.Address)
	cfg.Eth.GasPrice = big.NewInt(1)

	return nil
}

func setLogVerbosity(v string) (int, error) {
	lvl, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid --verbosity: %s", err)
	}
	handler := log.StreamHandler(os.Stderr, log.TerminalFormat(true))
	handler = log.LvlFilterHandler(log.Lvl(lvl), handler)
	log.Root().SetHandler(handler)
	return lvl, nil
}

func startMining(stack *node.Node, cfg *config) error {
	var ethereum *eth.Ethereum
	if err := stack.Service(&ethereum); err != nil {
		return fmt.Errorf("error getting Ethereum service: %s", err)
	}
	etherbase, err := ethereum.Etherbase()
	if err != nil {
		return fmt.Errorf("error getting Etherbase: %s", err)
	}
	// TODO: support keys with non-empty passphrase
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	if err := ks.Unlock(accounts.Account{Address: etherbase}, ""); err != nil {
		return fmt.Errorf("error unlocking Etherbase: %s", err)
	}
	ethereum.TxPool().SetGasPrice(cfg.Eth.GasPrice)
	if err := ethereum.StartMining(true); err != nil {
		return fmt.Errorf("error starting Ethereum mining: %s", err)
	}
	return nil
}
