// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.25;

contract A {
    function memoryCopy() external pure returns (bytes32 x) {
        assembly {
            mstore(0x20, 0x50) // Store 0x50 at word 1 in memory
            mcopy(0, 0x20, 0x20) // Copies 0x50 to word 0 in memory
            x := mload(0) // Returns 32 bytes "0x50"
        }
    }
}

// build/bin/abigen --abi ./tests/immutable/cancun/mcopy/abi.json --pkg mcopy --type mcopy --out ./tests/immutable/cancun/mcopy/mcopy.go --bin ./tests/immutable/cancun/mcopy/bytecode.bin
