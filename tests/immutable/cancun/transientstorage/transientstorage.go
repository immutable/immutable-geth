// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package transientstorage

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

// TransientStorageMetaData contains all meta data concerning the TransientStorage contract.
var TransientStorageMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"tStoreLoad\",\"inputs\":[{\"name\":\"key\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Value\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x608060405234801561001057600080fd5b5060f48061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806366a82e4a14602d575b600080fd5b603c6038366004607c565b603e565b005b80825d50604051815c808252907f2a27502c345a4cd966daa061d5537f54cd60d2d20b73680b3bf195c91e806a4b9060200160405180910390a15050565b60008060408385031215608e57600080fd5b5050803592602090910135915056fea2646970667358221220632312d849916834c328572910b3d44422cfb2ab00f8a234e5eaab69d71d0b6a64736f6c637823302e382e32322d63692e323032332e392e32312b636f6d6d69742e33633536396462390054",
}

// TransientStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use TransientStorageMetaData.ABI instead.
var TransientStorageABI = TransientStorageMetaData.ABI

// TransientStorageBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TransientStorageMetaData.Bin instead.
var TransientStorageBin = TransientStorageMetaData.Bin

// DeployTransientStorage deploys a new Ethereum contract, binding an instance of TransientStorage to it.
func DeployTransientStorage(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TransientStorage, error) {
	parsed, err := TransientStorageMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TransientStorageBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TransientStorage{TransientStorageCaller: TransientStorageCaller{contract: contract}, TransientStorageTransactor: TransientStorageTransactor{contract: contract}, TransientStorageFilterer: TransientStorageFilterer{contract: contract}}, nil
}

// TransientStorage is an auto generated Go binding around an Ethereum contract.
type TransientStorage struct {
	TransientStorageCaller     // Read-only binding to the contract
	TransientStorageTransactor // Write-only binding to the contract
	TransientStorageFilterer   // Log filterer for contract events
}

// TransientStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type TransientStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransientStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TransientStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransientStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransientStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransientStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransientStorageSession struct {
	Contract     *TransientStorage // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TransientStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransientStorageCallerSession struct {
	Contract *TransientStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// TransientStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransientStorageTransactorSession struct {
	Contract     *TransientStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// TransientStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type TransientStorageRaw struct {
	Contract *TransientStorage // Generic contract binding to access the raw methods on
}

// TransientStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransientStorageCallerRaw struct {
	Contract *TransientStorageCaller // Generic read-only contract binding to access the raw methods on
}

// TransientStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransientStorageTransactorRaw struct {
	Contract *TransientStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransientStorage creates a new instance of TransientStorage, bound to a specific deployed contract.
func NewTransientStorage(address common.Address, backend bind.ContractBackend) (*TransientStorage, error) {
	contract, err := bindTransientStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransientStorage{TransientStorageCaller: TransientStorageCaller{contract: contract}, TransientStorageTransactor: TransientStorageTransactor{contract: contract}, TransientStorageFilterer: TransientStorageFilterer{contract: contract}}, nil
}

// NewTransientStorageCaller creates a new read-only instance of TransientStorage, bound to a specific deployed contract.
func NewTransientStorageCaller(address common.Address, caller bind.ContractCaller) (*TransientStorageCaller, error) {
	contract, err := bindTransientStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransientStorageCaller{contract: contract}, nil
}

// NewTransientStorageTransactor creates a new write-only instance of TransientStorage, bound to a specific deployed contract.
func NewTransientStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*TransientStorageTransactor, error) {
	contract, err := bindTransientStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransientStorageTransactor{contract: contract}, nil
}

// NewTransientStorageFilterer creates a new log filterer instance of TransientStorage, bound to a specific deployed contract.
func NewTransientStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*TransientStorageFilterer, error) {
	contract, err := bindTransientStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransientStorageFilterer{contract: contract}, nil
}

// bindTransientStorage binds a generic wrapper to an already deployed contract.
func bindTransientStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransientStorageMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransientStorage *TransientStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransientStorage.Contract.TransientStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransientStorage *TransientStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransientStorage.Contract.TransientStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransientStorage *TransientStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransientStorage.Contract.TransientStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransientStorage *TransientStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransientStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransientStorage *TransientStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransientStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransientStorage *TransientStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransientStorage.Contract.contract.Transact(opts, method, params...)
}

// TStoreLoad is a paid mutator transaction binding the contract method 0x66a82e4a.
//
// Solidity: function tStoreLoad(uint256 key, uint256 value) returns()
func (_TransientStorage *TransientStorageTransactor) TStoreLoad(opts *bind.TransactOpts, key *big.Int, value *big.Int) (*types.Transaction, error) {
	return _TransientStorage.contract.Transact(opts, "tStoreLoad", key, value)
}

// TStoreLoad is a paid mutator transaction binding the contract method 0x66a82e4a.
//
// Solidity: function tStoreLoad(uint256 key, uint256 value) returns()
func (_TransientStorage *TransientStorageSession) TStoreLoad(key *big.Int, value *big.Int) (*types.Transaction, error) {
	return _TransientStorage.Contract.TStoreLoad(&_TransientStorage.TransactOpts, key, value)
}

// TStoreLoad is a paid mutator transaction binding the contract method 0x66a82e4a.
//
// Solidity: function tStoreLoad(uint256 key, uint256 value) returns()
func (_TransientStorage *TransientStorageTransactorSession) TStoreLoad(key *big.Int, value *big.Int) (*types.Transaction, error) {
	return _TransientStorage.Contract.TStoreLoad(&_TransientStorage.TransactOpts, key, value)
}

// TransientStorageValueIterator is returned from FilterValue and is used to iterate over the raw logs and unpacked data for Value events raised by the TransientStorage contract.
type TransientStorageValueIterator struct {
	Event *TransientStorageValue // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TransientStorageValueIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TransientStorageValue)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TransientStorageValue)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TransientStorageValueIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TransientStorageValueIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TransientStorageValue represents a Value event raised by the TransientStorage contract.
type TransientStorageValue struct {
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterValue is a free log retrieval operation binding the contract event 0x2a27502c345a4cd966daa061d5537f54cd60d2d20b73680b3bf195c91e806a4b.
//
// Solidity: event Value(uint256 value)
func (_TransientStorage *TransientStorageFilterer) FilterValue(opts *bind.FilterOpts) (*TransientStorageValueIterator, error) {

	logs, sub, err := _TransientStorage.contract.FilterLogs(opts, "Value")
	if err != nil {
		return nil, err
	}
	return &TransientStorageValueIterator{contract: _TransientStorage.contract, event: "Value", logs: logs, sub: sub}, nil
}

// WatchValue is a free log subscription operation binding the contract event 0x2a27502c345a4cd966daa061d5537f54cd60d2d20b73680b3bf195c91e806a4b.
//
// Solidity: event Value(uint256 value)
func (_TransientStorage *TransientStorageFilterer) WatchValue(opts *bind.WatchOpts, sink chan<- *TransientStorageValue) (event.Subscription, error) {

	logs, sub, err := _TransientStorage.contract.WatchLogs(opts, "Value")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TransientStorageValue)
				if err := _TransientStorage.contract.UnpackLog(event, "Value", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValue is a log parse operation binding the contract event 0x2a27502c345a4cd966daa061d5537f54cd60d2d20b73680b3bf195c91e806a4b.
//
// Solidity: event Value(uint256 value)
func (_TransientStorage *TransientStorageFilterer) ParseValue(log types.Log) (*TransientStorageValue, error) {
	event := new(TransientStorageValue)
	if err := _TransientStorage.contract.UnpackLog(event, "Value", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
