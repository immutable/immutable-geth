// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package selfdestructfunction

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// SelfDestructFunctionMetaData contains all meta data concerning the SelfDestructFunction contract.
var SelfDestructFunctionMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"selfDestruct\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x6080604052348015600f57600080fd5b5060b880601d6000396000f3fe60806040526004361060205760003560e01c80633f5a0bdd14602b57600080fd5b36602657005b600080fd5b348015603657600080fd5b50604660423660046054565b6048565b005b806001600160a01b0316ff5b600060208284031215606557600080fd5b81356001600160a01b0381168114607b57600080fd5b939250505056fea2646970667358221220f9c68da8ceb7e06771a69cb7acfbe7b52c3c45d91e44a58a1243d16d54606eae64736f6c63430008190033",
}

// SelfDestructFunctionABI is the input ABI used to generate the binding from.
// Deprecated: Use SelfDestructFunctionMetaData.ABI instead.
var SelfDestructFunctionABI = SelfDestructFunctionMetaData.ABI

// SelfDestructFunctionBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SelfDestructFunctionMetaData.Bin instead.
var SelfDestructFunctionBin = SelfDestructFunctionMetaData.Bin

// DeploySelfDestructFunction deploys a new Ethereum contract, binding an instance of SelfDestructFunction to it.
func DeploySelfDestructFunction(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SelfDestructFunction, error) {
	parsed, err := SelfDestructFunctionMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SelfDestructFunctionBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SelfDestructFunction{SelfDestructFunctionCaller: SelfDestructFunctionCaller{contract: contract}, SelfDestructFunctionTransactor: SelfDestructFunctionTransactor{contract: contract}, SelfDestructFunctionFilterer: SelfDestructFunctionFilterer{contract: contract}}, nil
}

// SelfDestructFunction is an auto generated Go binding around an Ethereum contract.
type SelfDestructFunction struct {
	SelfDestructFunctionCaller     // Read-only binding to the contract
	SelfDestructFunctionTransactor // Write-only binding to the contract
	SelfDestructFunctionFilterer   // Log filterer for contract events
}

// SelfDestructFunctionCaller is an auto generated read-only Go binding around an Ethereum contract.
type SelfDestructFunctionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfDestructFunctionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SelfDestructFunctionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfDestructFunctionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SelfDestructFunctionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfDestructFunctionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SelfDestructFunctionSession struct {
	Contract     *SelfDestructFunction // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// SelfDestructFunctionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SelfDestructFunctionCallerSession struct {
	Contract *SelfDestructFunctionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// SelfDestructFunctionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SelfDestructFunctionTransactorSession struct {
	Contract     *SelfDestructFunctionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// SelfDestructFunctionRaw is an auto generated low-level Go binding around an Ethereum contract.
type SelfDestructFunctionRaw struct {
	Contract *SelfDestructFunction // Generic contract binding to access the raw methods on
}

// SelfDestructFunctionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SelfDestructFunctionCallerRaw struct {
	Contract *SelfDestructFunctionCaller // Generic read-only contract binding to access the raw methods on
}

// SelfDestructFunctionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SelfDestructFunctionTransactorRaw struct {
	Contract *SelfDestructFunctionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSelfDestructFunction creates a new instance of SelfDestructFunction, bound to a specific deployed contract.
func NewSelfDestructFunction(address common.Address, backend bind.ContractBackend) (*SelfDestructFunction, error) {
	contract, err := bindSelfDestructFunction(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SelfDestructFunction{SelfDestructFunctionCaller: SelfDestructFunctionCaller{contract: contract}, SelfDestructFunctionTransactor: SelfDestructFunctionTransactor{contract: contract}, SelfDestructFunctionFilterer: SelfDestructFunctionFilterer{contract: contract}}, nil
}

// NewSelfDestructFunctionCaller creates a new read-only instance of SelfDestructFunction, bound to a specific deployed contract.
func NewSelfDestructFunctionCaller(address common.Address, caller bind.ContractCaller) (*SelfDestructFunctionCaller, error) {
	contract, err := bindSelfDestructFunction(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SelfDestructFunctionCaller{contract: contract}, nil
}

// NewSelfDestructFunctionTransactor creates a new write-only instance of SelfDestructFunction, bound to a specific deployed contract.
func NewSelfDestructFunctionTransactor(address common.Address, transactor bind.ContractTransactor) (*SelfDestructFunctionTransactor, error) {
	contract, err := bindSelfDestructFunction(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SelfDestructFunctionTransactor{contract: contract}, nil
}

// NewSelfDestructFunctionFilterer creates a new log filterer instance of SelfDestructFunction, bound to a specific deployed contract.
func NewSelfDestructFunctionFilterer(address common.Address, filterer bind.ContractFilterer) (*SelfDestructFunctionFilterer, error) {
	contract, err := bindSelfDestructFunction(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SelfDestructFunctionFilterer{contract: contract}, nil
}

// bindSelfDestructFunction binds a generic wrapper to an already deployed contract.
func bindSelfDestructFunction(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SelfDestructFunctionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SelfDestructFunction *SelfDestructFunctionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SelfDestructFunction.Contract.SelfDestructFunctionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SelfDestructFunction *SelfDestructFunctionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfDestructFunction.Contract.SelfDestructFunctionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SelfDestructFunction *SelfDestructFunctionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SelfDestructFunction.Contract.SelfDestructFunctionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SelfDestructFunction *SelfDestructFunctionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SelfDestructFunction.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SelfDestructFunction *SelfDestructFunctionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfDestructFunction.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SelfDestructFunction *SelfDestructFunctionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SelfDestructFunction.Contract.contract.Transact(opts, method, params...)
}

// SelfDestruct is a paid mutator transaction binding the contract method 0x3f5a0bdd.
//
// Solidity: function selfDestruct(address recipient) returns()
func (_SelfDestructFunction *SelfDestructFunctionTransactor) SelfDestruct(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _SelfDestructFunction.contract.Transact(opts, "selfDestruct", recipient)
}

// SelfDestruct is a paid mutator transaction binding the contract method 0x3f5a0bdd.
//
// Solidity: function selfDestruct(address recipient) returns()
func (_SelfDestructFunction *SelfDestructFunctionSession) SelfDestruct(recipient common.Address) (*types.Transaction, error) {
	return _SelfDestructFunction.Contract.SelfDestruct(&_SelfDestructFunction.TransactOpts, recipient)
}

// SelfDestruct is a paid mutator transaction binding the contract method 0x3f5a0bdd.
//
// Solidity: function selfDestruct(address recipient) returns()
func (_SelfDestructFunction *SelfDestructFunctionTransactorSession) SelfDestruct(recipient common.Address) (*types.Transaction, error) {
	return _SelfDestructFunction.Contract.SelfDestruct(&_SelfDestructFunction.TransactOpts, recipient)
}
