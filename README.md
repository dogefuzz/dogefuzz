# Dogefuzz: A flexible fuzzer to detect common vulnerabilities in Smart Contracts

## Setup services using docker

First, run the docker-compose file to start the local EVM (geth) node and Vandal API server.

For Docker version >= v2.20.2:

```
docker compose -f ./infra/docker-compose.yml up -d
```

For previous Docker versions:

```
docker-compose -f ./infra/docker-compose.yml up -d
```

## Setup geth service without docker (for debugging)

First, you need to run vandal using docker: 

```
docker compose  -f "infra/docker-compose.yml" up -d --build vandal
```

Second, to run geth outside of docker, follow the steps below:

1. Clone the repo https://github.com/PAMunb/dogefuzz-enhanced-go-ethereum
2. Run the following commands inside the cloned dir
```
go run build/ci.go install -static ./cmd/geth
ln -s build/bin/geth geth
source ./vars_for_debuging_geth.sh
./entrypoint.sh
```

## Run and execute dogefuzz fuzzer

To build and run the fuzzer server (dogefuzz), run the following command:

```
go build ./cmd/dogefuzz && go run ./cmd/dogefuzz 
```

The server will start listening to port 3456 (default).
To configure the fuzzer behavior, look into the [config.json](config.json) file.

To execute a fuzzing process, here is an example of a request:

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

It will execute the contract `AddressLotteryV2` contract per 15 minutes using the `directed_greybox` fuzzing strategy. And, it will try to detect the following weaknesses:

- delegate
- exception-disorder
- gasless-send
- number-dependency
- reentrancy
- timestamp-dependency

Available options for `fuzzingType` are:

- blackbox
- greybox
- directed_greybox

When no `arguments` are passed, the fuzzer will generate the constructor arguments.

The final report will be generated at the end of the fuzzing campaign and named `result.json`.

Note: The `contractSource` must be a JSON string, you can use sites and tools to do a JSON Stringify operation, as we need to convert CR/LF characters to their textual representation.

## Code Structure

1. ./assets/contracts - Directory where the agent's source code is located.
2. ./infra - Contains Docker files and Docker compose files.
