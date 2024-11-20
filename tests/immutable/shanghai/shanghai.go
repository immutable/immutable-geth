// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package shanghai

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

// ShanghaiMetaData contains all meta data concerning the Shanghai contract.
var ShanghaiMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"SetCoinbase\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b50606f80601a5f395ff3fe6080604052348015600e575f80fd5b50600436106026575f3560e01c8063de4436a614602a575b5f80fd5b60306032565b005b5f4190505056fea2646970667358221220ac8ec15e4867e5d8e3f28efa2be3d944029f496c1c0465c0f2fd611ac8d7d3e664736f6c63430008140033",
}

// ShanghaiABI is the input ABI used to generate the binding from.
// Deprecated: Use ShanghaiMetaData.ABI instead.
var ShanghaiABI = ShanghaiMetaData.ABI

// ShanghaiBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ShanghaiMetaData.Bin instead.
var ShanghaiBin = ShanghaiMetaData.Bin

// DeployShanghai deploys a new Ethereum contract, binding an instance of Shanghai to it.
func DeployShanghai(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Shanghai, error) {
	parsed, err := ShanghaiMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ShanghaiBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Shanghai{ShanghaiCaller: ShanghaiCaller{contract: contract}, ShanghaiTransactor: ShanghaiTransactor{contract: contract}, ShanghaiFilterer: ShanghaiFilterer{contract: contract}}, nil
}

// Shanghai is an auto generated Go binding around an Ethereum contract.
type Shanghai struct {
	ShanghaiCaller     // Read-only binding to the contract
	ShanghaiTransactor // Write-only binding to the contract
	ShanghaiFilterer   // Log filterer for contract events
}

// ShanghaiCaller is an auto generated read-only Go binding around an Ethereum contract.
type ShanghaiCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShanghaiTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ShanghaiTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShanghaiFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ShanghaiFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShanghaiSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ShanghaiSession struct {
	Contract     *Shanghai         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ShanghaiCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ShanghaiCallerSession struct {
	Contract *ShanghaiCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ShanghaiTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ShanghaiTransactorSession struct {
	Contract     *ShanghaiTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ShanghaiRaw is an auto generated low-level Go binding around an Ethereum contract.
type ShanghaiRaw struct {
	Contract *Shanghai // Generic contract binding to access the raw methods on
}

// ShanghaiCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ShanghaiCallerRaw struct {
	Contract *ShanghaiCaller // Generic read-only contract binding to access the raw methods on
}

// ShanghaiTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ShanghaiTransactorRaw struct {
	Contract *ShanghaiTransactor // Generic write-only contract binding to access the raw methods on
}

// NewShanghai creates a new instance of Shanghai, bound to a specific deployed contract.
func NewShanghai(address common.Address, backend bind.ContractBackend) (*Shanghai, error) {
	contract, err := bindShanghai(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Shanghai{ShanghaiCaller: ShanghaiCaller{contract: contract}, ShanghaiTransactor: ShanghaiTransactor{contract: contract}, ShanghaiFilterer: ShanghaiFilterer{contract: contract}}, nil
}

// NewShanghaiCaller creates a new read-only instance of Shanghai, bound to a specific deployed contract.
func NewShanghaiCaller(address common.Address, caller bind.ContractCaller) (*ShanghaiCaller, error) {
	contract, err := bindShanghai(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ShanghaiCaller{contract: contract}, nil
}

// NewShanghaiTransactor creates a new write-only instance of Shanghai, bound to a specific deployed contract.
func NewShanghaiTransactor(address common.Address, transactor bind.ContractTransactor) (*ShanghaiTransactor, error) {
	contract, err := bindShanghai(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ShanghaiTransactor{contract: contract}, nil
}

// NewShanghaiFilterer creates a new log filterer instance of Shanghai, bound to a specific deployed contract.
func NewShanghaiFilterer(address common.Address, filterer bind.ContractFilterer) (*ShanghaiFilterer, error) {
	contract, err := bindShanghai(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ShanghaiFilterer{contract: contract}, nil
}

// bindShanghai binds a generic wrapper to an already deployed contract.
func bindShanghai(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ShanghaiMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Shanghai *ShanghaiRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Shanghai.Contract.ShanghaiCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Shanghai *ShanghaiRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Shanghai.Contract.ShanghaiTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Shanghai *ShanghaiRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Shanghai.Contract.ShanghaiTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Shanghai *ShanghaiCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Shanghai.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Shanghai *ShanghaiTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Shanghai.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Shanghai *ShanghaiTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Shanghai.Contract.contract.Transact(opts, method, params...)
}

// SetCoinbase is a paid mutator transaction binding the contract method 0xde4436a6.
//
// Solidity: function SetCoinbase() returns()
func (_Shanghai *ShanghaiTransactor) SetCoinbase(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Shanghai.contract.Transact(opts, "SetCoinbase")
}

// SetCoinbase is a paid mutator transaction binding the contract method 0xde4436a6.
//
// Solidity: function SetCoinbase() returns()
func (_Shanghai *ShanghaiSession) SetCoinbase() (*types.Transaction, error) {
	return _Shanghai.Contract.SetCoinbase(&_Shanghai.TransactOpts)
}

// SetCoinbase is a paid mutator transaction binding the contract method 0xde4436a6.
//
// Solidity: function SetCoinbase() returns()
func (_Shanghai *ShanghaiTransactorSession) SetCoinbase() (*types.Transaction, error) {
	return _Shanghai.Contract.SetCoinbase(&_Shanghai.TransactOpts)
}
