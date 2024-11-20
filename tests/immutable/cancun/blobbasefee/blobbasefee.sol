// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.25;

contract BlobBaseFee {
    function blobBaseFee() external view returns (uint) {
        return block.blobbasefee;
    }
}

// build/bin/abigen --abi ./tests/immutable/cancun/blobbasefee/abi.json --pkg blobbasefee --type BlobBaseFee --out ./tests/immutable/cancun/blobbasefee/blobbasefee.go --bin ./tests/immutable/cancun/blobbasefee/bytecode.bin
