// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.25;

contract SelfDestructConstructor {
    constructor(address recipient) payable {
        selfdestruct(payable(recipient));
    }
}

// build/bin/abigen --abi ./tests/immutable/cancun/selfdestruct/constructor/abi.json --pkg selfdestructconstructor --type SelfDestructConstructor --out ./tests/immutable/cancun/selfdestruct/constructor/selfdestructconstructor.go --bin ./tests/immutable/cancun/selfdestruct/constructor/bytecode.bin
