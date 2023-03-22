package solidity

import (
	"errors"
	"math/big"
	"math/rand"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

var ErrInvalidUnsignedBigInt = errors.New("the provided json does not correspond to a unsigned big int type")

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

func (h *unsignedBigIntHandler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	options := []common.TypeIdentifier{UINT8, UINT16, UINT24, UINT32, UINT40, UINT48, UINT56, UINT64, UINT72, UINT80, UINT88, UINT96, UINT104, UINT112, UINT120, UINT128, UINT136, UINT144, UINT152, UINT160, UINT168, UINT176, UINT184, UINT192, UINT200, UINT208, UINT216, UINT224, UINT232, UINT240, UINT248, UINT256}
	typ := options[h.bitSize/8-1]
	addressSeeds := seeds[typ]
	chosenSeed := common.RandomChoice(addressSeeds)
	return h.Deserialize(chosenSeed)
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
	rnd := rand.New(rand.NewSource(common.Now().Unix()))

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
	base := big.NewInt(2)
	exponent := big.NewInt(int64(h.bitSize))

	max := new(big.Int)
	max.Exp(base, exponent, nil)
	max.Sub(max, big.NewInt(1))

	newHandler := NewUnsignedBigIntHandler(h.bitSize)
	newHandler.Generate()
	value := newHandler.GetValue().(*big.Int)

	if h.value.Cmp(big.NewInt(0).Sub(max, value)) > 0 {
		// if a + b > MAX then a + b - MAX
		// to not overflow, a + b - MAX ~~ a - (MAX - b)
		h.value.Sub(h.value, big.NewInt(0).Sub(max, value))
	} else {
		h.value.Add(h.value, value)
	}
}

func (h *unsignedBigIntHandler) SafeSub() {
	base := big.NewInt(2)
	exponent := big.NewInt(int64(h.bitSize))

	max := new(big.Int)
	max.Exp(base, exponent, nil)
	max.Sub(max, big.NewInt(1))

	newHandler := NewUnsignedBigIntHandler(h.bitSize)
	newHandler.Generate()
	value := newHandler.GetValue().(*big.Int)

	if h.value.Cmp(value) < 0 {
		// if a - b < 0 then 0 - (a - b)
		// to not undeflow, 0 - (a - b) ~~ b - a
		h.value.Sub(value, h.value)
	} else {
		h.value.Sub(h.value, value)
	}
}

func (h *unsignedBigIntHandler) SafeMul() {
	base := big.NewInt(2)
	exponent := big.NewInt(int64(h.bitSize))

	max := new(big.Int)
	max.Exp(base, exponent, nil)
	max.Sub(max, big.NewInt(1))

	newHandler := NewUnsignedBigIntHandler(h.bitSize)
	newHandler.Generate()
	value := newHandler.GetValue().(*big.Int)

	if value.Cmp(big.NewInt(0)) == 0 {
		h.value = big.NewInt(0)
	} else if h.value.Cmp(big.NewInt(0).Div(max, value)) > 0 {
		h.value.Mod(big.NewInt(1).Mul(h.value, value), max)
	} else {
		h.value.Mul(h.value, value)
	}
}

func (h *unsignedBigIntHandler) SafeDiv() {
	base := big.NewInt(2)
	exponent := big.NewInt(int64(h.bitSize))

	max := new(big.Int)
	max.Exp(base, exponent, nil)
	max.Sub(max, big.NewInt(1))

	newHandler := NewUnsignedBigIntHandler(h.bitSize)
	newHandler.Generate()
	value := newHandler.GetValue().(*big.Int)

	if value.Cmp(big.NewInt(0)) == 0 {
		h.value = big.NewInt(0)
	} else {
		h.value.Div(h.value, value)
	}
}
