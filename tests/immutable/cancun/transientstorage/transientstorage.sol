// SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.20;

contract SimpleTStore {
    event Value(uint value);

    function tStoreLoad(uint key, uint value) external returns (uint) {
        assembly {
            tstore(key, value)
            value := tload(key)
        }
    }
}

// build/bin/abigen --abi ./tests/immutable/cancun/transientstorage/abi.json --pkg transientstorage --type TransientStorage --out ./tests/immutable/cancun/transientstorage/transientstorage.go --bin ./tests/immutable/cancun/transientstorage/bytecode.bin
