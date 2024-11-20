// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.25;

contract Blobhash {
    function storeBlobHash(uint256 index) external {
        assembly {
            sstore(0, blobhash(index))
        }
    }
}

// build/bin/abigen --abi ./tests/immutable/cancun/blobhash/abi.json --pkg blobhash --type Blobhash --out ./tests/immutable/cancun/blobhash/blobhash.go --bin ./tests/immutable/cancun/blobhash/bytecode.bin
