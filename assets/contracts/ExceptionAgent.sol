// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

contract ExceptionAgent {
    fallback() external payable {
        revert("Exception thrown");
    }

    receive() external payable {
        revert("Exception thrown");
    }

    function CallContract(
        address contractAddress,
        bytes calldata data
    ) public payable {
        bool success;
        if (msg.value > 0) {
            (success, ) = contractAddress.call{value: msg.value}(data);
        } else {
            (success, ) = contractAddress.call(data);
        }
        assert(success);
    }
}
