// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

contract ReentrancyAgent {
    address globalCalledAddress;
    bytes globalData;
    uint balance = 0;
    mapping(address => bool) sent;

    fallback() external payable {
        balance += msg.value;
        if (!sent[globalCalledAddress]) {
            sent[globalCalledAddress] = true;
            // globalCalledAddress.call(globalData);
            (bool success, ) = globalCalledAddress.call(globalData);
            assert(success);
        } else {
            sent[globalCalledAddress] = false;
        }
    }

    receive() external payable {
        balance += msg.value;
        if (!sent[globalCalledAddress]) {
            sent[globalCalledAddress] = true;
            // globalCalledAddress.call(globalData);
            (bool success, ) = globalCalledAddress.call(globalData);
            assert(success);
        } else {
            sent[globalCalledAddress] = false;
        }
    }

    function CallContract(
        address contractAddress,
        bytes calldata data
    ) public payable {
        globalCalledAddress = contractAddress;
        globalData = data;
        bool success;
        if (msg.value > 0) {
            (success, ) = contractAddress.call{value: msg.value}(data);
        } else {
            (success, ) = contractAddress.call(data);
        }
        assert(success);
    }
}
