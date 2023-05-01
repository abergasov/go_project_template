// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package swapper

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

// StargateRouterMetaData contains all meta data concerning the StargateRouter contract.
var StargateRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stargateEthVault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stargateRouter\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"_poolId\",\"type\":\"uint16\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"addLiquidityETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"poolId\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stargateEthVault\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stargateRouter\",\"outputs\":[{\"internalType\":\"contractIStargateRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_dstChainId\",\"type\":\"uint16\"},{\"internalType\":\"addresspayable\",\"name\":\"_refundAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_toAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_amountLD\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minAmountLD\",\"type\":\"uint256\"}],\"name\":\"swapETH\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// StargateRouterABI is the input ABI used to generate the binding from.
// Deprecated: Use StargateRouterMetaData.ABI instead.
var StargateRouterABI = StargateRouterMetaData.ABI

// StargateRouter is an auto generated Go binding around an Ethereum contract.
type StargateRouter struct {
	StargateRouterCaller     // Read-only binding to the contract
	StargateRouterTransactor // Write-only binding to the contract
	StargateRouterFilterer   // Log filterer for contract events
}

// StargateRouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type StargateRouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StargateRouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StargateRouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StargateRouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StargateRouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StargateRouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StargateRouterSession struct {
	Contract     *StargateRouter   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StargateRouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StargateRouterCallerSession struct {
	Contract *StargateRouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// StargateRouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StargateRouterTransactorSession struct {
	Contract     *StargateRouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// StargateRouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type StargateRouterRaw struct {
	Contract *StargateRouter // Generic contract binding to access the raw methods on
}

// StargateRouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StargateRouterCallerRaw struct {
	Contract *StargateRouterCaller // Generic read-only contract binding to access the raw methods on
}

// StargateRouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StargateRouterTransactorRaw struct {
	Contract *StargateRouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStargateRouter creates a new instance of StargateRouter, bound to a specific deployed contract.
func NewStargateRouter(address common.Address, backend bind.ContractBackend) (*StargateRouter, error) {
	contract, err := bindStargateRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StargateRouter{StargateRouterCaller: StargateRouterCaller{contract: contract}, StargateRouterTransactor: StargateRouterTransactor{contract: contract}, StargateRouterFilterer: StargateRouterFilterer{contract: contract}}, nil
}

// NewStargateRouterCaller creates a new read-only instance of StargateRouter, bound to a specific deployed contract.
func NewStargateRouterCaller(address common.Address, caller bind.ContractCaller) (*StargateRouterCaller, error) {
	contract, err := bindStargateRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StargateRouterCaller{contract: contract}, nil
}

// NewStargateRouterTransactor creates a new write-only instance of StargateRouter, bound to a specific deployed contract.
func NewStargateRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*StargateRouterTransactor, error) {
	contract, err := bindStargateRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StargateRouterTransactor{contract: contract}, nil
}

// NewStargateRouterFilterer creates a new log filterer instance of StargateRouter, bound to a specific deployed contract.
func NewStargateRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*StargateRouterFilterer, error) {
	contract, err := bindStargateRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StargateRouterFilterer{contract: contract}, nil
}

// bindStargateRouter binds a generic wrapper to an already deployed contract.
func bindStargateRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StargateRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StargateRouter *StargateRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StargateRouter.Contract.StargateRouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StargateRouter *StargateRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StargateRouter.Contract.StargateRouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StargateRouter *StargateRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StargateRouter.Contract.StargateRouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StargateRouter *StargateRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StargateRouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StargateRouter *StargateRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StargateRouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StargateRouter *StargateRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StargateRouter.Contract.contract.Transact(opts, method, params...)
}

// PoolId is a free data retrieval call binding the contract method 0x3e0dc34e.
//
// Solidity: function poolId() view returns(uint16)
func (_StargateRouter *StargateRouterCaller) PoolId(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _StargateRouter.contract.Call(opts, &out, "poolId")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// PoolId is a free data retrieval call binding the contract method 0x3e0dc34e.
//
// Solidity: function poolId() view returns(uint16)
func (_StargateRouter *StargateRouterSession) PoolId() (uint16, error) {
	return _StargateRouter.Contract.PoolId(&_StargateRouter.CallOpts)
}

// PoolId is a free data retrieval call binding the contract method 0x3e0dc34e.
//
// Solidity: function poolId() view returns(uint16)
func (_StargateRouter *StargateRouterCallerSession) PoolId() (uint16, error) {
	return _StargateRouter.Contract.PoolId(&_StargateRouter.CallOpts)
}

// StargateEthVault is a free data retrieval call binding the contract method 0x38e31d39.
//
// Solidity: function stargateEthVault() view returns(address)
func (_StargateRouter *StargateRouterCaller) StargateEthVault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StargateRouter.contract.Call(opts, &out, "stargateEthVault")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StargateEthVault is a free data retrieval call binding the contract method 0x38e31d39.
//
// Solidity: function stargateEthVault() view returns(address)
func (_StargateRouter *StargateRouterSession) StargateEthVault() (common.Address, error) {
	return _StargateRouter.Contract.StargateEthVault(&_StargateRouter.CallOpts)
}

// StargateEthVault is a free data retrieval call binding the contract method 0x38e31d39.
//
// Solidity: function stargateEthVault() view returns(address)
func (_StargateRouter *StargateRouterCallerSession) StargateEthVault() (common.Address, error) {
	return _StargateRouter.Contract.StargateEthVault(&_StargateRouter.CallOpts)
}

// StargateRouter is a free data retrieval call binding the contract method 0xa9e56f3c.
//
// Solidity: function stargateRouter() view returns(address)
func (_StargateRouter *StargateRouterCaller) StargateRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StargateRouter.contract.Call(opts, &out, "stargateRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StargateRouter is a free data retrieval call binding the contract method 0xa9e56f3c.
//
// Solidity: function stargateRouter() view returns(address)
func (_StargateRouter *StargateRouterSession) StargateRouter() (common.Address, error) {
	return _StargateRouter.Contract.StargateRouter(&_StargateRouter.CallOpts)
}

// StargateRouter is a free data retrieval call binding the contract method 0xa9e56f3c.
//
// Solidity: function stargateRouter() view returns(address)
func (_StargateRouter *StargateRouterCallerSession) StargateRouter() (common.Address, error) {
	return _StargateRouter.Contract.StargateRouter(&_StargateRouter.CallOpts)
}

// AddLiquidityETH is a paid mutator transaction binding the contract method 0xed995307.
//
// Solidity: function addLiquidityETH() payable returns()
func (_StargateRouter *StargateRouterTransactor) AddLiquidityETH(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StargateRouter.contract.Transact(opts, "addLiquidityETH")
}

// AddLiquidityETH is a paid mutator transaction binding the contract method 0xed995307.
//
// Solidity: function addLiquidityETH() payable returns()
func (_StargateRouter *StargateRouterSession) AddLiquidityETH() (*types.Transaction, error) {
	return _StargateRouter.Contract.AddLiquidityETH(&_StargateRouter.TransactOpts)
}

// AddLiquidityETH is a paid mutator transaction binding the contract method 0xed995307.
//
// Solidity: function addLiquidityETH() payable returns()
func (_StargateRouter *StargateRouterTransactorSession) AddLiquidityETH() (*types.Transaction, error) {
	return _StargateRouter.Contract.AddLiquidityETH(&_StargateRouter.TransactOpts)
}

// SwapETH is a paid mutator transaction binding the contract method 0x1114cd2a.
//
// Solidity: function swapETH(uint16 _dstChainId, address _refundAddress, bytes _toAddress, uint256 _amountLD, uint256 _minAmountLD) payable returns()
func (_StargateRouter *StargateRouterTransactor) SwapETH(opts *bind.TransactOpts, _dstChainId uint16, _refundAddress common.Address, _toAddress []byte, _amountLD *big.Int, _minAmountLD *big.Int) (*types.Transaction, error) {
	return _StargateRouter.contract.Transact(opts, "swapETH", _dstChainId, _refundAddress, _toAddress, _amountLD, _minAmountLD)
}

// SwapETH is a paid mutator transaction binding the contract method 0x1114cd2a.
//
// Solidity: function swapETH(uint16 _dstChainId, address _refundAddress, bytes _toAddress, uint256 _amountLD, uint256 _minAmountLD) payable returns()
func (_StargateRouter *StargateRouterSession) SwapETH(_dstChainId uint16, _refundAddress common.Address, _toAddress []byte, _amountLD *big.Int, _minAmountLD *big.Int) (*types.Transaction, error) {
	return _StargateRouter.Contract.SwapETH(&_StargateRouter.TransactOpts, _dstChainId, _refundAddress, _toAddress, _amountLD, _minAmountLD)
}

// SwapETH is a paid mutator transaction binding the contract method 0x1114cd2a.
//
// Solidity: function swapETH(uint16 _dstChainId, address _refundAddress, bytes _toAddress, uint256 _amountLD, uint256 _minAmountLD) payable returns()
func (_StargateRouter *StargateRouterTransactorSession) SwapETH(_dstChainId uint16, _refundAddress common.Address, _toAddress []byte, _amountLD *big.Int, _minAmountLD *big.Int) (*types.Transaction, error) {
	return _StargateRouter.Contract.SwapETH(&_StargateRouter.TransactOpts, _dstChainId, _refundAddress, _toAddress, _amountLD, _minAmountLD)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_StargateRouter *StargateRouterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StargateRouter.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_StargateRouter *StargateRouterSession) Receive() (*types.Transaction, error) {
	return _StargateRouter.Contract.Receive(&_StargateRouter.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_StargateRouter *StargateRouterTransactorSession) Receive() (*types.Transaction, error) {
	return _StargateRouter.Contract.Receive(&_StargateRouter.TransactOpts)
}
