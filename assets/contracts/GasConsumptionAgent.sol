// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

contract GasConsumptionAgent {
    fallback() external payable {
        int64 value = 0;
        while (value < 100) {
            value++;
        }
    }

    receive() external payable {
        int64 value = 0;
        while (value < 100) {
            value++;
        }
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

    function TransferToContract(address payable to) public payable {
        to.transfer(msg.value);
    }

    function SendToContract(address payable to) public payable {
        require(to.send(msg.value), "Send failed");
    }
}
