package generators

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/entities"
)

func ContractGen() *entities.Contract {
	return &entities.Contract{
		Id:      gofakeit.LetterN(255),
		Address: SmartContractGen(),
		Source:  gofakeit.LetterN(255),
		Name:    gofakeit.Name(),
	}
}
