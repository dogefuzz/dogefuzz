pragma solidity 0.4.26;

contract ExceptionFallback {
    function() external {
        revert("Exception thrown");
    }
}
