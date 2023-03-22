pragma solidity 0.4.26;

contract GasConsumptionFallback {
    function() external {
        int64 value = 0;
        while (true) {
            value++;
        }
    }
}
