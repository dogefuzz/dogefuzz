package solidity

import (
	"errors"
	"math/big"
	"math/rand"
	"time"
)

var ErrInvalidUnsignedBigInt = errors.New("the provided json does not correspond to a boolean type")

type unsignedBigIntHandler struct {
	bitSize int
	value   *big.Int
}

func NewUnsignedBigIntHandler(bitSize int) *unsignedBigIntHandler {
	return &unsignedBigIntHandler{bitSize: bitSize}
}

func (h *unsignedBigIntHandler) GetValue() interface{} {
	return h.value
}

func (h *unsignedBigIntHandler) SetValue(value interface{}) {
	h.value = value.(*big.Int)
}

func (h *unsignedBigIntHandler) Serialize() string {
	return h.value.String()
}

func (h *unsignedBigIntHandler) Deserialize(value string) error {
	number := new(big.Int)
	number, ok := number.SetString(value, 10)
	if !ok || number.Sign() == -1 {
		return ErrInvalidUnsignedBigInt
	}

	base := big.NewInt(2)
	exponent := big.NewInt(int64(h.bitSize))
	max := new(big.Int)
	max.Exp(base, exponent, nil)
	max.Sub(max, big.NewInt(1))
	if number.Cmp(max) > 0 {
		return ErrInvalidUnsignedBigInt
	}
	h.value = number
	return nil
}

func (h *unsignedBigIntHandler) Generate() {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	base := big.NewInt(2)
	exponent := big.NewInt(int64(h.bitSize))
	max := new(big.Int)
	max.Exp(base, exponent, nil)

	value := new(big.Int)
	value.Rand(rnd, max)
	h.value = value
}

func (h *unsignedBigIntHandler) GetMutators() []func() {
	return []func(){
		h.SafeAdd,
		h.SafeSub,
		h.SafeDiv,
		h.SafeMul,
	}
}

func (h *unsignedBigIntHandler) SafeAdd() {
	// value := s.Generate()
	// if input > math.MaxUint8-value {
	// 	return input - (math.MaxUint8 - value) - 1
	// }
	// return input + value
}

func (h *unsignedBigIntHandler) SafeSub() {
	// value := s.Generate()
	// if input < value {
	// 	return value - input - 1
	// }
	// return input - value
}

func (h *unsignedBigIntHandler) SafeMul() {
	// value := s.Generate()
	// if value.Cmp(big.NewInt(0)) {
	// 	return big.NewInt(0)
	// }
	// if input > math.MaxUint8/value {
	// 	return input * (math.MaxUint8 / value)
	// }
	// return input * value
}

func (h *unsignedBigIntHandler) SafeDiv() {
	// value := s.Generate()
	// if value == 0 {
	// 	return 0
	// }
	// return input / value
}
