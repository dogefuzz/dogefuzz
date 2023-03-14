pragma solidity 0.4.26;

contract ReentrancyAgent {
    address globalCalledAddress;
    bytes globalData;
    bool sent = false;

    function() external {
        if (!sent) {
            sent = true;
            globalCalledAddress.call(globalData);
        } else {
            sent = false;
        }
        
    }

    function CallContract(address contractAddress, bytes data) public payable {
        globalCalledAddress = contractAddress;
        globalData = data;
        if (msg.value > 0) {
            contractAddress.call.value(msg.value)(data);
        } else {
            contractAddress.call(data);
        }
    }
}
