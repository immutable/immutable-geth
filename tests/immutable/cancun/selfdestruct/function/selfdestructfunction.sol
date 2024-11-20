// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.25;

contract SelfDestructFunction {
    function selfDestruct(address recipient) external {
        selfdestruct(payable(recipient));
    }

    receive() external payable {}
}

// build/bin/abigen --abi ./tests/immutable/cancun/selfdestruct/function/abi.json --pkg selfdestructfunction --type SelfDestructFunction --out ./tests/immutable/cancun/selfdestruct/function/selfdestructfunction.go --bin ./tests/immutable/cancun/selfdestruct/function/bytecode.bin
