// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Shanghai {
    function SetCoinbase() external {
        address coinbase = block.coinbase;
    }
}

// build/bin/abigen --abi ./tests/immutable/shanghai/abi.json --pkg shanghai --type Shanghai --out ./tests/immutable/shanghai/shanghai.go --bin ./tests/immutable/shanghai/bytecode.bin