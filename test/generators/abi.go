package generators

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func BoolTypeGen() abi.Type {
	boolTyp, _ := abi.NewType("bool", "", nil)
	return boolTyp
}

func UintTypeGen(bitSize int) abi.Type {
	uintTyp, _ := abi.NewType(fmt.Sprintf("uint%d", bitSize), "", nil)
	return uintTyp
}

func IntTypeGen(bitSize int) abi.Type {
	uintTyp, _ := abi.NewType(fmt.Sprintf("int%d", bitSize), "", nil)
	return uintTyp
}
