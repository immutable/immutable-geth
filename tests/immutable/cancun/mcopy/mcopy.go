// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mcopy

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

// McopyMetaData contains all meta data concerning the Mcopy contract.
var McopyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"memoryCopy\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"x\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f80fd5b5060b980601a5f395ff3fe6080604052348015600e575f80fd5b50600436106026575f3560e01c80632dbaeee914602a575b5f80fd5b60306044565b604051603b9190606c565b60405180910390f35b5f60506020526020805f5e5f51905090565b5f819050919050565b6066816056565b82525050565b5f602082019050607d5f830184605f565b9291505056fea2646970667358221220710c9da5e12c71adc7661a10c85eaef4d3563725e2fa50835f0b6875b0c5d3cd64736f6c63430008190033",
}

// McopyABI is the input ABI used to generate the binding from.
// Deprecated: Use McopyMetaData.ABI instead.
var McopyABI = McopyMetaData.ABI

// McopyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use McopyMetaData.Bin instead.
var McopyBin = McopyMetaData.Bin

// DeployMcopy deploys a new Ethereum contract, binding an instance of Mcopy to it.
func DeployMcopy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Mcopy, error) {
	parsed, err := McopyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(McopyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Mcopy{McopyCaller: McopyCaller{contract: contract}, McopyTransactor: McopyTransactor{contract: contract}, McopyFilterer: McopyFilterer{contract: contract}}, nil
}

// Mcopy is an auto generated Go binding around an Ethereum contract.
type Mcopy struct {
	McopyCaller     // Read-only binding to the contract
	McopyTransactor // Write-only binding to the contract
	McopyFilterer   // Log filterer for contract events
}

// McopyCaller is an auto generated read-only Go binding around an Ethereum contract.
type McopyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// McopyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type McopyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// McopyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type McopyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// McopySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type McopySession struct {
	Contract     *Mcopy            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// McopyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type McopyCallerSession struct {
	Contract *McopyCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// McopyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type McopyTransactorSession struct {
	Contract     *McopyTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// McopyRaw is an auto generated low-level Go binding around an Ethereum contract.
type McopyRaw struct {
	Contract *Mcopy // Generic contract binding to access the raw methods on
}

// McopyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type McopyCallerRaw struct {
	Contract *McopyCaller // Generic read-only contract binding to access the raw methods on
}

// McopyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type McopyTransactorRaw struct {
	Contract *McopyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMcopy creates a new instance of Mcopy, bound to a specific deployed contract.
func NewMcopy(address common.Address, backend bind.ContractBackend) (*Mcopy, error) {
	contract, err := bindMcopy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Mcopy{McopyCaller: McopyCaller{contract: contract}, McopyTransactor: McopyTransactor{contract: contract}, McopyFilterer: McopyFilterer{contract: contract}}, nil
}

// NewMcopyCaller creates a new read-only instance of Mcopy, bound to a specific deployed contract.
func NewMcopyCaller(address common.Address, caller bind.ContractCaller) (*McopyCaller, error) {
	contract, err := bindMcopy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &McopyCaller{contract: contract}, nil
}

// NewMcopyTransactor creates a new write-only instance of Mcopy, bound to a specific deployed contract.
func NewMcopyTransactor(address common.Address, transactor bind.ContractTransactor) (*McopyTransactor, error) {
	contract, err := bindMcopy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &McopyTransactor{contract: contract}, nil
}

// NewMcopyFilterer creates a new log filterer instance of Mcopy, bound to a specific deployed contract.
func NewMcopyFilterer(address common.Address, filterer bind.ContractFilterer) (*McopyFilterer, error) {
	contract, err := bindMcopy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &McopyFilterer{contract: contract}, nil
}

// bindMcopy binds a generic wrapper to an already deployed contract.
func bindMcopy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := McopyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mcopy *McopyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mcopy.Contract.McopyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mcopy *McopyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mcopy.Contract.McopyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mcopy *McopyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mcopy.Contract.McopyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mcopy *McopyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mcopy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mcopy *McopyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mcopy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mcopy *McopyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mcopy.Contract.contract.Transact(opts, method, params...)
}

// MemoryCopy is a free data retrieval call binding the contract method 0x2dbaeee9.
//
// Solidity: function memoryCopy() pure returns(bytes32 x)
func (_Mcopy *McopyCaller) MemoryCopy(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Mcopy.contract.Call(opts, &out, "memoryCopy")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MemoryCopy is a free data retrieval call binding the contract method 0x2dbaeee9.
//
// Solidity: function memoryCopy() pure returns(bytes32 x)
func (_Mcopy *McopySession) MemoryCopy() ([32]byte, error) {
	return _Mcopy.Contract.MemoryCopy(&_Mcopy.CallOpts)
}

// MemoryCopy is a free data retrieval call binding the contract method 0x2dbaeee9.
//
// Solidity: function memoryCopy() pure returns(bytes32 x)
func (_Mcopy *McopyCallerSession) MemoryCopy() ([32]byte, error) {
	return _Mcopy.Contract.MemoryCopy(&_Mcopy.CallOpts)
}
