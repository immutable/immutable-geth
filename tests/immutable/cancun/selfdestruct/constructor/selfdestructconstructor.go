// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package selfdestructconstructor

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

// SelfDestructConstructorMetaData contains all meta data concerning the SelfDestructConstructor contract.
var SelfDestructConstructorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"payable\"}]",
	Bin: "0x608060405260405160593803806059833981016040819052601e91602a565b806001600160a01b0316ff5b600060208284031215603b57600080fd5b81516001600160a01b0381168114605157600080fd5b939250505056fe",
}

// SelfDestructConstructorABI is the input ABI used to generate the binding from.
// Deprecated: Use SelfDestructConstructorMetaData.ABI instead.
var SelfDestructConstructorABI = SelfDestructConstructorMetaData.ABI

// SelfDestructConstructorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SelfDestructConstructorMetaData.Bin instead.
var SelfDestructConstructorBin = SelfDestructConstructorMetaData.Bin

// DeploySelfDestructConstructor deploys a new Ethereum contract, binding an instance of SelfDestructConstructor to it.
func DeploySelfDestructConstructor(auth *bind.TransactOpts, backend bind.ContractBackend, recipient common.Address) (common.Address, *types.Transaction, *SelfDestructConstructor, error) {
	parsed, err := SelfDestructConstructorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SelfDestructConstructorBin), backend, recipient)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SelfDestructConstructor{SelfDestructConstructorCaller: SelfDestructConstructorCaller{contract: contract}, SelfDestructConstructorTransactor: SelfDestructConstructorTransactor{contract: contract}, SelfDestructConstructorFilterer: SelfDestructConstructorFilterer{contract: contract}}, nil
}

// SelfDestructConstructor is an auto generated Go binding around an Ethereum contract.
type SelfDestructConstructor struct {
	SelfDestructConstructorCaller     // Read-only binding to the contract
	SelfDestructConstructorTransactor // Write-only binding to the contract
	SelfDestructConstructorFilterer   // Log filterer for contract events
}

// SelfDestructConstructorCaller is an auto generated read-only Go binding around an Ethereum contract.
type SelfDestructConstructorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfDestructConstructorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SelfDestructConstructorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfDestructConstructorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SelfDestructConstructorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SelfDestructConstructorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SelfDestructConstructorSession struct {
	Contract     *SelfDestructConstructor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SelfDestructConstructorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SelfDestructConstructorCallerSession struct {
	Contract *SelfDestructConstructorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// SelfDestructConstructorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SelfDestructConstructorTransactorSession struct {
	Contract     *SelfDestructConstructorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// SelfDestructConstructorRaw is an auto generated low-level Go binding around an Ethereum contract.
type SelfDestructConstructorRaw struct {
	Contract *SelfDestructConstructor // Generic contract binding to access the raw methods on
}

// SelfDestructConstructorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SelfDestructConstructorCallerRaw struct {
	Contract *SelfDestructConstructorCaller // Generic read-only contract binding to access the raw methods on
}

// SelfDestructConstructorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SelfDestructConstructorTransactorRaw struct {
	Contract *SelfDestructConstructorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSelfDestructConstructor creates a new instance of SelfDestructConstructor, bound to a specific deployed contract.
func NewSelfDestructConstructor(address common.Address, backend bind.ContractBackend) (*SelfDestructConstructor, error) {
	contract, err := bindSelfDestructConstructor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SelfDestructConstructor{SelfDestructConstructorCaller: SelfDestructConstructorCaller{contract: contract}, SelfDestructConstructorTransactor: SelfDestructConstructorTransactor{contract: contract}, SelfDestructConstructorFilterer: SelfDestructConstructorFilterer{contract: contract}}, nil
}

// NewSelfDestructConstructorCaller creates a new read-only instance of SelfDestructConstructor, bound to a specific deployed contract.
func NewSelfDestructConstructorCaller(address common.Address, caller bind.ContractCaller) (*SelfDestructConstructorCaller, error) {
	contract, err := bindSelfDestructConstructor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SelfDestructConstructorCaller{contract: contract}, nil
}

// NewSelfDestructConstructorTransactor creates a new write-only instance of SelfDestructConstructor, bound to a specific deployed contract.
func NewSelfDestructConstructorTransactor(address common.Address, transactor bind.ContractTransactor) (*SelfDestructConstructorTransactor, error) {
	contract, err := bindSelfDestructConstructor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SelfDestructConstructorTransactor{contract: contract}, nil
}

// NewSelfDestructConstructorFilterer creates a new log filterer instance of SelfDestructConstructor, bound to a specific deployed contract.
func NewSelfDestructConstructorFilterer(address common.Address, filterer bind.ContractFilterer) (*SelfDestructConstructorFilterer, error) {
	contract, err := bindSelfDestructConstructor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SelfDestructConstructorFilterer{contract: contract}, nil
}

// bindSelfDestructConstructor binds a generic wrapper to an already deployed contract.
func bindSelfDestructConstructor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SelfDestructConstructorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SelfDestructConstructor *SelfDestructConstructorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SelfDestructConstructor.Contract.SelfDestructConstructorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SelfDestructConstructor *SelfDestructConstructorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfDestructConstructor.Contract.SelfDestructConstructorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SelfDestructConstructor *SelfDestructConstructorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SelfDestructConstructor.Contract.SelfDestructConstructorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SelfDestructConstructor *SelfDestructConstructorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SelfDestructConstructor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SelfDestructConstructor *SelfDestructConstructorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SelfDestructConstructor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SelfDestructConstructor *SelfDestructConstructorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SelfDestructConstructor.Contract.contract.Transact(opts, method, params...)
}
