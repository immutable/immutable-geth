// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blobbasefee

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

// BlobBaseFeeMetaData contains all meta data concerning the BlobBaseFee contract.
var BlobBaseFeeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"blobBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b5060ae80601a5f395ff3fe6080604052348015600e575f80fd5b50600436106026575f3560e01c8063f820614014602a575b5f80fd5b60306044565b604051603b91906061565b60405180910390f35b5f4a905090565b5f819050919050565b605b81604b565b82525050565b5f60208201905060725f8301846054565b9291505056fea2646970667358221220f36b1e170b1a918aa653ca6f6ee810faa51fe43e25828bb9b8caedb51031bee264736f6c63430008190033",
}

// BlobBaseFeeABI is the input ABI used to generate the binding from.
// Deprecated: Use BlobBaseFeeMetaData.ABI instead.
var BlobBaseFeeABI = BlobBaseFeeMetaData.ABI

// BlobBaseFeeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BlobBaseFeeMetaData.Bin instead.
var BlobBaseFeeBin = BlobBaseFeeMetaData.Bin

// DeployBlobBaseFee deploys a new Ethereum contract, binding an instance of BlobBaseFee to it.
func DeployBlobBaseFee(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BlobBaseFee, error) {
	parsed, err := BlobBaseFeeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BlobBaseFeeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BlobBaseFee{BlobBaseFeeCaller: BlobBaseFeeCaller{contract: contract}, BlobBaseFeeTransactor: BlobBaseFeeTransactor{contract: contract}, BlobBaseFeeFilterer: BlobBaseFeeFilterer{contract: contract}}, nil
}

// BlobBaseFee is an auto generated Go binding around an Ethereum contract.
type BlobBaseFee struct {
	BlobBaseFeeCaller     // Read-only binding to the contract
	BlobBaseFeeTransactor // Write-only binding to the contract
	BlobBaseFeeFilterer   // Log filterer for contract events
}

// BlobBaseFeeCaller is an auto generated read-only Go binding around an Ethereum contract.
type BlobBaseFeeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlobBaseFeeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BlobBaseFeeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlobBaseFeeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BlobBaseFeeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlobBaseFeeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BlobBaseFeeSession struct {
	Contract     *BlobBaseFee      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BlobBaseFeeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BlobBaseFeeCallerSession struct {
	Contract *BlobBaseFeeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BlobBaseFeeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BlobBaseFeeTransactorSession struct {
	Contract     *BlobBaseFeeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BlobBaseFeeRaw is an auto generated low-level Go binding around an Ethereum contract.
type BlobBaseFeeRaw struct {
	Contract *BlobBaseFee // Generic contract binding to access the raw methods on
}

// BlobBaseFeeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BlobBaseFeeCallerRaw struct {
	Contract *BlobBaseFeeCaller // Generic read-only contract binding to access the raw methods on
}

// BlobBaseFeeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BlobBaseFeeTransactorRaw struct {
	Contract *BlobBaseFeeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBlobBaseFee creates a new instance of BlobBaseFee, bound to a specific deployed contract.
func NewBlobBaseFee(address common.Address, backend bind.ContractBackend) (*BlobBaseFee, error) {
	contract, err := bindBlobBaseFee(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BlobBaseFee{BlobBaseFeeCaller: BlobBaseFeeCaller{contract: contract}, BlobBaseFeeTransactor: BlobBaseFeeTransactor{contract: contract}, BlobBaseFeeFilterer: BlobBaseFeeFilterer{contract: contract}}, nil
}

// NewBlobBaseFeeCaller creates a new read-only instance of BlobBaseFee, bound to a specific deployed contract.
func NewBlobBaseFeeCaller(address common.Address, caller bind.ContractCaller) (*BlobBaseFeeCaller, error) {
	contract, err := bindBlobBaseFee(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BlobBaseFeeCaller{contract: contract}, nil
}

// NewBlobBaseFeeTransactor creates a new write-only instance of BlobBaseFee, bound to a specific deployed contract.
func NewBlobBaseFeeTransactor(address common.Address, transactor bind.ContractTransactor) (*BlobBaseFeeTransactor, error) {
	contract, err := bindBlobBaseFee(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BlobBaseFeeTransactor{contract: contract}, nil
}

// NewBlobBaseFeeFilterer creates a new log filterer instance of BlobBaseFee, bound to a specific deployed contract.
func NewBlobBaseFeeFilterer(address common.Address, filterer bind.ContractFilterer) (*BlobBaseFeeFilterer, error) {
	contract, err := bindBlobBaseFee(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BlobBaseFeeFilterer{contract: contract}, nil
}

// bindBlobBaseFee binds a generic wrapper to an already deployed contract.
func bindBlobBaseFee(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BlobBaseFeeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BlobBaseFee *BlobBaseFeeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BlobBaseFee.Contract.BlobBaseFeeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BlobBaseFee *BlobBaseFeeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlobBaseFee.Contract.BlobBaseFeeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BlobBaseFee *BlobBaseFeeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlobBaseFee.Contract.BlobBaseFeeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BlobBaseFee *BlobBaseFeeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BlobBaseFee.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BlobBaseFee *BlobBaseFeeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BlobBaseFee.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BlobBaseFee *BlobBaseFeeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BlobBaseFee.Contract.contract.Transact(opts, method, params...)
}

// BlobBaseFee is a free data retrieval call binding the contract method 0xf8206140.
//
// Solidity: function blobBaseFee() view returns(uint256)
func (_BlobBaseFee *BlobBaseFeeCaller) BlobBaseFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BlobBaseFee.contract.Call(opts, &out, "blobBaseFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BlobBaseFee is a free data retrieval call binding the contract method 0xf8206140.
//
// Solidity: function blobBaseFee() view returns(uint256)
func (_BlobBaseFee *BlobBaseFeeSession) BlobBaseFee() (*big.Int, error) {
	return _BlobBaseFee.Contract.BlobBaseFee(&_BlobBaseFee.CallOpts)
}

// BlobBaseFee is a free data retrieval call binding the contract method 0xf8206140.
//
// Solidity: function blobBaseFee() view returns(uint256)
func (_BlobBaseFee *BlobBaseFeeCallerSession) BlobBaseFee() (*big.Int, error) {
	return _BlobBaseFee.Contract.BlobBaseFee(&_BlobBaseFee.CallOpts)
}
