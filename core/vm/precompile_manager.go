// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

// precompileManager is used as a default PrecompileManager for the EVM.
type precompileManager struct {
	rules       params.Rules
	precompiles map[common.Address]PrecompiledContract
}

// NewPrecompileManager returns a new PrecompileManager for the current chain rules.
func NewPrecompileManager(rules params.Rules) PrecompileManager {
	return &precompileManager{
		rules: rules,
	}
}

// Has returns whether a precompiled contract is deployed at the given address.
func (pm *precompileManager) Has(addr common.Address) bool {
	if pm.precompiles == nil {
		pm.precompiles = pm.activePrecompiles()
	}
	_, found := pm.precompiles[addr]
	return found
}

// Get returns the precompiled contract deployed at the given address.
func (pm *precompileManager) Get(addr common.Address) PrecompiledContract {
	if pm.precompiles == nil {
		pm.precompiles = pm.activePrecompiles()
	}
	return pm.precompiles[addr]
}

// Run runs the given precompiled contract with the given input data and returns the remaining gas.
func (pm *precompileManager) Run(
	_ StateDB, p PrecompiledContract, input []byte,
	caller common.Address, value *big.Int, suppliedGas uint64, readonly bool,
) (ret []byte, remainingGas uint64, err error) {
	gasCost := p.RequiredGas(input)
	if gasCost > suppliedGas {
		return nil, 0, ErrOutOfGas
	}

	suppliedGas -= gasCost
	output, err := p.Run(context.Background(), input, caller, value, readonly)

	return output, suppliedGas, err
}

// activePrecompiles returns the precompiled contracts for the current chain rules.
func (pm *precompileManager) activePrecompiles() map[common.Address]PrecompiledContract {
	switch {
	case pm.rules.IsBerlin:
		return PrecompiledContractsBerlin
	case pm.rules.IsIstanbul:
		return PrecompiledContractsIstanbul
	case pm.rules.IsByzantium:
		return PrecompiledContractsByzantium
	default:
		return PrecompiledContractsHomestead
	}
}