package generators

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dogefuzz/dogefuzz/domain"
)

func ContractGen() *domain.Contract {
	return &domain.Contract{
		Id:      gofakeit.LetterN(255),
		Address: SmartContractGen(),
		Source:  gofakeit.LetterN(255),
		Name:    gofakeit.Name(),
	}
}
