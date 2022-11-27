package generators

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/pkg/common"
)

func CommonContractGen() *common.Contract {
	return &common.Contract{
		Name:          gofakeit.Word(),
		AbiDefinition: gofakeit.Word(),
		CompiledCode:  gofakeit.HexUint256(),
	}
}
