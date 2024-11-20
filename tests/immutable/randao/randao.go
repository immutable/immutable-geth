// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package randao

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

// RandaoMetaData contains all meta data concerning the Randao contract.
var RandaoMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"difficulty\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rand\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x6080604052348015600f57600080fd5b5060808061001e6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c806319cae4621460375780633b3dca76146037575b600080fd5b4460405190815260200160405180910390f3fea264697066735822122073d353e8eb69c22f774d562d71a74ea7cb9835935aa5f222c0c41eec2c15a7c364736f6c63430008180033",
}

// RandaoABI is the input ABI used to generate the binding from.
// Deprecated: Use RandaoMetaData.ABI instead.
var RandaoABI = RandaoMetaData.ABI

// RandaoBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RandaoMetaData.Bin instead.
var RandaoBin = RandaoMetaData.Bin

// DeployRandao deploys a new Ethereum contract, binding an instance of Randao to it.
func DeployRandao(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Randao, error) {
	parsed, err := RandaoMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RandaoBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Randao{RandaoCaller: RandaoCaller{contract: contract}, RandaoTransactor: RandaoTransactor{contract: contract}, RandaoFilterer: RandaoFilterer{contract: contract}}, nil
}

// Randao is an auto generated Go binding around an Ethereum contract.
type Randao struct {
	RandaoCaller     // Read-only binding to the contract
	RandaoTransactor // Write-only binding to the contract
	RandaoFilterer   // Log filterer for contract events
}

// RandaoCaller is an auto generated read-only Go binding around an Ethereum contract.
type RandaoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RandaoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RandaoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RandaoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RandaoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RandaoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RandaoSession struct {
	Contract     *Randao           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RandaoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RandaoCallerSession struct {
	Contract *RandaoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RandaoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RandaoTransactorSession struct {
	Contract     *RandaoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RandaoRaw is an auto generated low-level Go binding around an Ethereum contract.
type RandaoRaw struct {
	Contract *Randao // Generic contract binding to access the raw methods on
}

// RandaoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RandaoCallerRaw struct {
	Contract *RandaoCaller // Generic read-only contract binding to access the raw methods on
}

// RandaoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RandaoTransactorRaw struct {
	Contract *RandaoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRandao creates a new instance of Randao, bound to a specific deployed contract.
func NewRandao(address common.Address, backend bind.ContractBackend) (*Randao, error) {
	contract, err := bindRandao(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Randao{RandaoCaller: RandaoCaller{contract: contract}, RandaoTransactor: RandaoTransactor{contract: contract}, RandaoFilterer: RandaoFilterer{contract: contract}}, nil
}

// NewRandaoCaller creates a new read-only instance of Randao, bound to a specific deployed contract.
func NewRandaoCaller(address common.Address, caller bind.ContractCaller) (*RandaoCaller, error) {
	contract, err := bindRandao(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RandaoCaller{contract: contract}, nil
}

// NewRandaoTransactor creates a new write-only instance of Randao, bound to a specific deployed contract.
func NewRandaoTransactor(address common.Address, transactor bind.ContractTransactor) (*RandaoTransactor, error) {
	contract, err := bindRandao(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RandaoTransactor{contract: contract}, nil
}

// NewRandaoFilterer creates a new log filterer instance of Randao, bound to a specific deployed contract.
func NewRandaoFilterer(address common.Address, filterer bind.ContractFilterer) (*RandaoFilterer, error) {
	contract, err := bindRandao(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RandaoFilterer{contract: contract}, nil
}

// bindRandao binds a generic wrapper to an already deployed contract.
func bindRandao(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RandaoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Randao *RandaoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Randao.Contract.RandaoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Randao *RandaoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Randao.Contract.RandaoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Randao *RandaoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Randao.Contract.RandaoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Randao *RandaoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Randao.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Randao *RandaoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Randao.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Randao *RandaoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Randao.Contract.contract.Transact(opts, method, params...)
}

// Difficulty is a free data retrieval call binding the contract method 0x19cae462.
//
// Solidity: function difficulty() view returns(uint256)
func (_Randao *RandaoCaller) Difficulty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Randao.contract.Call(opts, &out, "difficulty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Difficulty is a free data retrieval call binding the contract method 0x19cae462.
//
// Solidity: function difficulty() view returns(uint256)
func (_Randao *RandaoSession) Difficulty() (*big.Int, error) {
	return _Randao.Contract.Difficulty(&_Randao.CallOpts)
}

// Difficulty is a free data retrieval call binding the contract method 0x19cae462.
//
// Solidity: function difficulty() view returns(uint256)
func (_Randao *RandaoCallerSession) Difficulty() (*big.Int, error) {
	return _Randao.Contract.Difficulty(&_Randao.CallOpts)
}

// Rand is a free data retrieval call binding the contract method 0x3b3dca76.
//
// Solidity: function rand() view returns(uint256)
func (_Randao *RandaoCaller) Rand(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Randao.contract.Call(opts, &out, "rand")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Rand is a free data retrieval call binding the contract method 0x3b3dca76.
//
// Solidity: function rand() view returns(uint256)
func (_Randao *RandaoSession) Rand() (*big.Int, error) {
	return _Randao.Contract.Rand(&_Randao.CallOpts)
}

// Rand is a free data retrieval call binding the contract method 0x3b3dca76.
//
// Solidity: function rand() view returns(uint256)
func (_Randao *RandaoCallerSession) Rand() (*big.Int, error) {
	return _Randao.Contract.Rand(&_Randao.CallOpts)
}
