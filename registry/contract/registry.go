// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// METARegistryABI is the input ABI used to generate the binding from.
const METARegistryABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"setGraph\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"metaID\",\"type\":\"address\"}],\"name\":\"graph\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// METARegistryBin is the compiled bytecode used for deploying new contracts.
const METARegistryBin = `0x6060604052341561000f57600080fd5b6102388061001e6000396000f30060606040526004361061004b5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166375dfb0e08114610050578063ab2c7776146100a8575b600080fd5b341561005b57600080fd5b6100a6600480359060446024803590810190830135806020601f820181900481020160405190810160405281815292919060208401838380828437509496506100e695505050505050565b005b34156100b357600080fd5b6100d473ffffffffffffffffffffffffffffffffffffffff600435166101e4565b60405190815260200160405180910390f35b60008060008084516041146100fa57600080fd5b6020850151925060408501519150606085015160001a9350601b8460ff16101561012557601b840193505b8360ff16601b1415801561013d57508360ff16601c14155b1561014757600080fd5b6001868585856040516000815260200160405260006040516020015260405193845260ff90921660208085019190915260408085019290925260608401929092526080909201915160208103908084039060008661646e5a03f115156101ac57600080fd5b50506020604051035173ffffffffffffffffffffffffffffffffffffffff166000908152602081905260409020959095555050505050565b73ffffffffffffffffffffffffffffffffffffffff16600090815260208190526040902054905600a165627a7a72305820bc3cd2f1cc7e0da514b92c0bf67d62293da625bb79164b7a48a3b7c7c58cd52c0029`

// DeployMETARegistry deploys a new Ethereum contract, binding an instance of METARegistry to it.
func DeployMETARegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *METARegistry, error) {
	parsed, err := abi.JSON(strings.NewReader(METARegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(METARegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &METARegistry{METARegistryCaller: METARegistryCaller{contract: contract}, METARegistryTransactor: METARegistryTransactor{contract: contract}, METARegistryFilterer: METARegistryFilterer{contract: contract}}, nil
}

// METARegistry is an auto generated Go binding around an Ethereum contract.
type METARegistry struct {
	METARegistryCaller     // Read-only binding to the contract
	METARegistryTransactor // Write-only binding to the contract
	METARegistryFilterer   // Log filterer for contract events
}

// METARegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type METARegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// METARegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type METARegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// METARegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type METARegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// METARegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type METARegistrySession struct {
	Contract     *METARegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// METARegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type METARegistryCallerSession struct {
	Contract *METARegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// METARegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type METARegistryTransactorSession struct {
	Contract     *METARegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// METARegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type METARegistryRaw struct {
	Contract *METARegistry // Generic contract binding to access the raw methods on
}

// METARegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type METARegistryCallerRaw struct {
	Contract *METARegistryCaller // Generic read-only contract binding to access the raw methods on
}

// METARegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type METARegistryTransactorRaw struct {
	Contract *METARegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMETARegistry creates a new instance of METARegistry, bound to a specific deployed contract.
func NewMETARegistry(address common.Address, backend bind.ContractBackend) (*METARegistry, error) {
	contract, err := bindMETARegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &METARegistry{METARegistryCaller: METARegistryCaller{contract: contract}, METARegistryTransactor: METARegistryTransactor{contract: contract}, METARegistryFilterer: METARegistryFilterer{contract: contract}}, nil
}

// NewMETARegistryCaller creates a new read-only instance of METARegistry, bound to a specific deployed contract.
func NewMETARegistryCaller(address common.Address, caller bind.ContractCaller) (*METARegistryCaller, error) {
	contract, err := bindMETARegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &METARegistryCaller{contract: contract}, nil
}

// NewMETARegistryTransactor creates a new write-only instance of METARegistry, bound to a specific deployed contract.
func NewMETARegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*METARegistryTransactor, error) {
	contract, err := bindMETARegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &METARegistryTransactor{contract: contract}, nil
}

// NewMETARegistryFilterer creates a new log filterer instance of METARegistry, bound to a specific deployed contract.
func NewMETARegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*METARegistryFilterer, error) {
	contract, err := bindMETARegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &METARegistryFilterer{contract: contract}, nil
}

// bindMETARegistry binds a generic wrapper to an already deployed contract.
func bindMETARegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(METARegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_METARegistry *METARegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _METARegistry.Contract.METARegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_METARegistry *METARegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _METARegistry.Contract.METARegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_METARegistry *METARegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _METARegistry.Contract.METARegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_METARegistry *METARegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _METARegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_METARegistry *METARegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _METARegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_METARegistry *METARegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _METARegistry.Contract.contract.Transact(opts, method, params...)
}

// Graph is a free data retrieval call binding the contract method 0xab2c7776.
//
// Solidity: function graph(metaID address) constant returns(bytes32)
func (_METARegistry *METARegistryCaller) Graph(opts *bind.CallOpts, metaID common.Address) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _METARegistry.contract.Call(opts, out, "graph", metaID)
	return *ret0, err
}

// Graph is a free data retrieval call binding the contract method 0xab2c7776.
//
// Solidity: function graph(metaID address) constant returns(bytes32)
func (_METARegistry *METARegistrySession) Graph(metaID common.Address) ([32]byte, error) {
	return _METARegistry.Contract.Graph(&_METARegistry.CallOpts, metaID)
}

// Graph is a free data retrieval call binding the contract method 0xab2c7776.
//
// Solidity: function graph(metaID address) constant returns(bytes32)
func (_METARegistry *METARegistryCallerSession) Graph(metaID common.Address) ([32]byte, error) {
	return _METARegistry.Contract.Graph(&_METARegistry.CallOpts, metaID)
}

// SetGraph is a paid mutator transaction binding the contract method 0x75dfb0e0.
//
// Solidity: function setGraph(hash bytes32, sig bytes) returns()
func (_METARegistry *METARegistryTransactor) SetGraph(opts *bind.TransactOpts, hash [32]byte, sig []byte) (*types.Transaction, error) {
	return _METARegistry.contract.Transact(opts, "setGraph", hash, sig)
}

// SetGraph is a paid mutator transaction binding the contract method 0x75dfb0e0.
//
// Solidity: function setGraph(hash bytes32, sig bytes) returns()
func (_METARegistry *METARegistrySession) SetGraph(hash [32]byte, sig []byte) (*types.Transaction, error) {
	return _METARegistry.Contract.SetGraph(&_METARegistry.TransactOpts, hash, sig)
}

// SetGraph is a paid mutator transaction binding the contract method 0x75dfb0e0.
//
// Solidity: function setGraph(hash bytes32, sig bytes) returns()
func (_METARegistry *METARegistryTransactorSession) SetGraph(hash [32]byte, sig []byte) (*types.Transaction, error) {
	return _METARegistry.Contract.SetGraph(&_METARegistry.TransactOpts, hash, sig)
}
