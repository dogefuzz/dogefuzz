# Dogefuzz: A flexible fuzzer to detect common vulnerabilities in Smart Contracts

## Setup
First, run the docker compose file to start the local EVM node and Vandal API server.
```
docker compose -f ./infra/docker-compose.yml up -d
```

To configure the fuzzer behavior, look into the [config.json](config.json) file.

## Run and execute
To run the fuzzer server, run the following command:
```
go run ./cmd/dogefuzz
```
And the server will start listen the port 3456 (default).

To execute a fuzzing process, here an example of request:
```
curl -X POST \
     http://localhost:3456/tasks \
     -H 'Accept: application/json' \
     -H 'Content-Type: application/json' \
     -d '{
    "contractSource": "pragma solidity ^0.4.0;\n/*\n * This is a distributed lottery that chooses random addresses as lucky addresses. If these\n * participate, they get the jackpot: the whole balance of the contract, including the ticket\n * price. Of course one address can only win once. The owner regularly reseeds the secret\n * seed of the contract (based on which the lucky addresses are chosen), so if you did not win,\n * just wait for a reseed and try again! Contract addresses cannot play for obvious reasons.\n *\n * Jackpot chance:   1 in 8\n*/\ncontract AddressLotteryV2{\n    struct SeedComponents{\n        uint component1;\n        uint component2;\n        uint component3;\n        uint component4;\n    }\n    \n    address owner;\n    uint private secretSeed;\n    uint private lastReseed;\n    \n    uint winnerLuckyNumber = 7;\n    \n    uint public ticketPrice = 0.1 ether;\n        \n    mapping (address => bool) participated;\n\n    modifier onlyOwner() {\n        require(msg.sender == owner);\n        _;\n    }\n  \n    modifier onlyHuman() {\n        require(msg.sender == tx.origin);\n        _;\n    }\n    \n    function AddressLotteryV2() {\n        owner = msg.sender;\n        reseed(SeedComponents(12345678, 0x12345678, 0xabbaeddaacdc, 0x22222222));\n    }\n    \n    function setTicketPrice(uint newPrice) onlyOwner {\n        ticketPrice = newPrice;\n    }\n    \n    function participate() payable onlyHuman { \n        require(msg.value == ticketPrice);\n        \n        // every address can only win once, obviously\n        require(!participated[msg.sender]);\n        \n        if ( luckyNumberOfAddress(msg.sender) == winnerLuckyNumber)\n        {\n            participated[msg.sender] = true;\n            require(msg.sender.call.value(this.balance)());\n        }\n    }\n    \n    function luckyNumberOfAddress(address addr) constant returns(uint n){\n        // 1 in 8 chance\n        n = uint(keccak256(uint(addr), secretSeed)[0]) % 8;\n    }\n    \n    function reseed(SeedComponents components) internal{\n        secretSeed = uint256(keccak256(\n            components.component1,\n            components.component2,\n            components.component3,\n            components.component4\n        ));\n        lastReseed = block.number;\n    }\n    \n    function kill() onlyOwner {\n        suicide(owner);\n    }\n    \n    function forceReseed() onlyOwner{\n        SeedComponents s;\n        s.component1 = uint(msg.sender);\n        s.component2 = uint256(block.blockhash(block.number - 1));\n        s.component3 = block.number * 1337;\n        s.component4 = tx.gasprice * 7;\n        reseed(s);\n    }\n    \n    function () payable {}\n    \n    // DEBUG, DELETE BEFORE DEPLOYMENT!!\n    function _myLuckyNumber() constant returns(uint n){\n        n = luckyNumberOfAddress(msg.sender);\n    }\n}",
    "contractName": "AddressLotteryV2",
    "arguments": [],
    "duration": "15m",
    "detectors": [
        "delegate",
        "exception-disorder",
        "gasless-send",
        "number-dependency",
        "reentrancy",
        "timestamp-dependency"
    ],
    "fuzzingType": "directed_greybox"
}'
```

It will execute the contract `AddressLotteryV2` contract per 15 minutes using the `directed_greybox` fuzzing strategy. And, it will detect the following weaknesses:

- delegate
- exception-disorder
- gasless-send
- number-dependency
- reentrancy
- timestamp-dependency

Available options for `fuzzingType` are: 

- 	blackbox
-	greybox
-   directed_greybox

As no `arguments` were passed, the fuzzer will generate the contructor arguments.


## Code Structure

1. ./assets/contracts - Directory where attack contracts are located.
2. ./infra - Contains docker files and docker compose files.
