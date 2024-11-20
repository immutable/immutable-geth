// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blobhash

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

// BlobhashMetaData contains all meta data concerning the Blobhash contract.
var BlobhashMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"blobHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b506101198061001c5f395ff3fe6080604052348015600e575f80fd5b50600436106026575f3560e01c80630ba54e3214602a575b5f80fd5b60406004803603810190603c91906090565b6054565b604051604b919060cc565b60405180910390f35b5f81499050919050565b5f80fd5b5f819050919050565b6072816062565b8114607b575f80fd5b50565b5f81359050608a81606b565b92915050565b5f6020828403121560a25760a1605e565b5b5f60ad84828501607e565b91505092915050565b5f819050919050565b60c68160b6565b82525050565b5f60208201905060dd5f83018460bf565b9291505056fea26469706673582212208da4d91290c36a37fd192cd72cc7b29ba6328b2de08606cf869614c5ddd9b37364736f6c63430008190033",
}

// BlobhashABI is the input ABI used to generate the binding from.
// Deprecated: Use BlobhashMetaData.ABI instead.
var BlobhashABI = BlobhashMetaData.ABI

// BlobhashBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BlobhashMetaData.Bin instead.
var BlobhashBin = BlobhashMetaData.Bin

// DeployBlobhash deploys a new Ethereum contract, binding an instance of Blobhash to it.
func DeployBlobhash(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Blobhash, error) {
	parsed, err := BlobhashMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BlobhashBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Blobhash{BlobhashCaller: BlobhashCaller{contract: contract}, BlobhashTransactor: BlobhashTransactor{contract: contract}, BlobhashFilterer: BlobhashFilterer{contract: contract}}, nil
}

// Blobhash is an auto generated Go binding around an Ethereum contract.
type Blobhash struct {
	BlobhashCaller     // Read-only binding to the contract
	BlobhashTransactor // Write-only binding to the contract
	BlobhashFilterer   // Log filterer for contract events
}

// BlobhashCaller is an auto generated read-only Go binding around an Ethereum contract.
type BlobhashCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlobhashTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BlobhashTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlobhashFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BlobhashFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlobhashSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BlobhashSession struct {
	Contract     *Blobhash         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BlobhashCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BlobhashCallerSession struct {
	Contract *BlobhashCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// BlobhashTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BlobhashTransactorSession struct {
	Contract     *BlobhashTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// BlobhashRaw is an auto generated low-level Go binding around an Ethereum contract.
type BlobhashRaw struct {
	Contract *Blobhash // Generic contract binding to access the raw methods on
}

// BlobhashCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BlobhashCallerRaw struct {
	Contract *BlobhashCaller // Generic read-only contract binding to access the raw methods on
}

// BlobhashTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BlobhashTransactorRaw struct {
	Contract *BlobhashTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBlobhash creates a new instance of Blobhash, bound to a specific deployed contract.
func NewBlobhash(address common.Address, backend bind.ContractBackend) (*Blobhash, error) {
	contract, err := bindBlobhash(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Blobhash{BlobhashCaller: BlobhashCaller{contract: contract}, BlobhashTransactor: BlobhashTransactor{contract: contract}, BlobhashFilterer: BlobhashFilterer{contract: contract}}, nil
}

// NewBlobhashCaller creates a new read-only instance of Blobhash, bound to a specific deployed contract.
func NewBlobhashCaller(address common.Address, caller bind.ContractCaller) (*BlobhashCaller, error) {
	contract, err := bindBlobhash(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BlobhashCaller{contract: contract}, nil
}

// NewBlobhashTransactor creates a new write-only instance of Blobhash, bound to a specific deployed contract.
func NewBlobhashTransactor(address common.Address, transactor bind.ContractTransactor) (*BlobhashTransactor, error) {
	contract, err := bindBlobhash(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BlobhashTransactor{contract: contract}, nil
}

// NewBlobhashFilterer creates a new log filterer instance of Blobhash, bound to a specific deployed contract.
func NewBlobhashFilterer(address common.Address, filterer bind.ContractFilterer) (*BlobhashFilterer, error) {
	contract, err := bindBlobhash(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BlobhashFilterer{contract: contract}, nil
}

// bindBlobhash binds a generic wrapper to an already deployed contract.
func bindBlobhash(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BlobhashMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Blobhash *BlobhashRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Blobhash.Contract.BlobhashCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Blobhash *BlobhashRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blobhash.Contract.BlobhashTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Blobhash *BlobhashRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Blobhash.Contract.BlobhashTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Blobhash *BlobhashCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Blobhash.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Blobhash *BlobhashTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blobhash.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Blobhash *BlobhashTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Blobhash.Contract.contract.Transact(opts, method, params...)
}

// BlobHash is a free data retrieval call binding the contract method 0x0ba54e32.
//
// Solidity: function blobHash(uint256 index) view returns(bytes32)
func (_Blobhash *BlobhashCaller) BlobHash(opts *bind.CallOpts, index *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Blobhash.contract.Call(opts, &out, "blobHash", index)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BlobHash is a free data retrieval call binding the contract method 0x0ba54e32.
//
// Solidity: function blobHash(uint256 index) view returns(bytes32)
func (_Blobhash *BlobhashSession) BlobHash(index *big.Int) ([32]byte, error) {
	return _Blobhash.Contract.BlobHash(&_Blobhash.CallOpts, index)
}

// BlobHash is a free data retrieval call binding the contract method 0x0ba54e32.
//
// Solidity: function blobHash(uint256 index) view returns(bytes32)
func (_Blobhash *BlobhashCallerSession) BlobHash(index *big.Int) ([32]byte, error) {
	return _Blobhash.Contract.BlobHash(&_Blobhash.CallOpts, index)
}
