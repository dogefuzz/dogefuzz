package solidity

import (
	"errors"
	"math/big"
	"math/rand"
	"time"
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
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

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
