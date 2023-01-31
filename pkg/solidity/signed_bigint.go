package solidity

import (
	"errors"
	"math/big"
	"math/rand"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

var ErrInvalidSignedBigInt = errors.New("the provided json does not correspond to a Big.Int type")

type signedBigIntHandler struct {
	bitSize int
	value   *big.Int
}

func NewSignedBigIntHandler(bitSize int) *signedBigIntHandler {
	return &signedBigIntHandler{bitSize: bitSize}
}

func (h *signedBigIntHandler) GetValue() interface{} {
	return h.value
}

func (h *signedBigIntHandler) SetValue(value interface{}) {
	h.value = value.(*big.Int)
}

func (h *signedBigIntHandler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	options := []common.TypeIdentifier{INT8, INT16, INT24, INT32, INT40, INT48, INT56, INT64, INT72, INT80, INT88, INT96, INT104, INT112, INT120, INT128, INT136, INT144, INT152, INT160, INT168, INT176, INT184, INT192, INT200, INT208, INT216, INT224, INT232, INT240, INT248, INT256}
	typ := options[h.bitSize/8-1]
	addressSeeds := seeds[typ]
	chosenSeed := common.RandomChoice(addressSeeds)
	return h.Deserialize(chosenSeed)
}

func (h *signedBigIntHandler) Serialize() string {
	return h.value.String()
}

func (h *signedBigIntHandler) Deserialize(value string) error {
	number := new(big.Int)
	number, ok := number.SetString(value, 10)
	if !ok {
		return ErrInvalidSignedBigInt
	}

	base := big.NewInt(2)
	exponent := big.NewInt(int64(h.bitSize) - 1)

	min := new(big.Int)
	min.Exp(base, exponent, nil)
	min.Neg(min)
	if number.Cmp(min) < 0 {
		return ErrInvalidSignedBigInt
	}

	max := new(big.Int)
	max.Exp(base, exponent, nil)
	max.Sub(max, big.NewInt(1))
	if number.Cmp(max) > 0 {
		return ErrInvalidSignedBigInt
	}
	h.value = number
	return nil
}

func (h *signedBigIntHandler) Generate() {
	rnd := rand.New(rand.NewSource(common.Now().Unix()))

	base := big.NewInt(2)
	exponent := big.NewInt(int64(h.bitSize / 2))

	max := new(big.Int)
	max.Exp(base, exponent, nil)
	max.Sub(max, big.NewInt(1))

	min := new(big.Int)
	min.Exp(base, exponent, nil)
	min.Neg(min)

	randomRange := new(big.Int)
	randomRange.Add(randomRange, max)
	randomRange.Sub(randomRange, min)

	value := new(big.Int)
	value.Rand(rnd, randomRange)
	value.Add(value, min)
	h.value = value
}

func (h *signedBigIntHandler) GetMutators() []func() {
	return []func(){
		h.SafeAdd,
		h.SafeSub,
		h.SafeDiv,
		h.SafeMul,
	}
}

func (h *signedBigIntHandler) SafeAdd() {
	// value := s.Generate()
	// if input > math.MaxUint8-value {
	// 	return input - (math.MaxUint8 - value) - 1
	// }
	// return input + value
}

func (h *signedBigIntHandler) SafeSub() {
	// value := s.Generate()
	// if input < value {
	// 	return value - input - 1
	// }
	// return input - value
}

func (h *signedBigIntHandler) SafeMul() {
	// value := s.Generate()
	// if value.Cmp(big.NewInt(0)) {
	// 	return big.NewInt(0)
	// }
	// if input > math.MaxUint8/value {
	// 	return input * (math.MaxUint8 / value)
	// }
	// return input * value
}

func (h *signedBigIntHandler) SafeDiv() {
	// value := s.Generate()
	// if value == 0 {
	// 	return 0
	// }
	// return input / value
}
