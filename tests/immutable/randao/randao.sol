// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract RandDao {
    function rand() external view returns (uint) {
        return block.prevrandao;
    }
    function difficulty() external view returns (uint) {
        return block.difficulty;
    }
}

// build/bin/abigen --abi ./tests/immutable/randao/abi.json --pkg randao --type Randao --out ./tests/immutable/randao/randao.go --bin ./tests/immutable/randao/bytecode.bin
