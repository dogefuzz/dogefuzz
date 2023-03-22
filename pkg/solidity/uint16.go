package solidity

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

const UINT16_BIT_SIZE = 16

var ErrInvalidUint16 = errors.New("the provided string does not correspond to a uint16 type")

type uint16Handler struct {
	value uint16
}

func NewUint16Handler() *uint16Handler {
	return &uint16Handler{}
}

func (h *uint16Handler) GetValue() interface{} {
	return h.value
}

func (h *uint16Handler) SetValue(value interface{}) {
	h.value = value.(uint16)
}

func (h *uint16Handler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	addressSeeds := seeds[UINT16]
	chosenSeed := common.RandomChoice(addressSeeds)
	return h.Deserialize(chosenSeed)
}

func (h *uint16Handler) Serialize() string {
	return strconv.FormatUint(uint64(h.value), 10)
}

func (h *uint16Handler) Deserialize(value string) error {
	val, err := strconv.ParseUint(value, 10, UINT16_BIT_SIZE)
	if err != nil {
		return ErrInvalidUint16
	}
	h.value = uint16(val)
	return nil
}

func (h *uint16Handler) Generate() {
	rand.Seed(common.Now().Unix())
	h.value = uint16(rand.Intn(2 << (UINT16_BIT_SIZE - 1)))
}

func (h *uint16Handler) GetMutators() []func() {
	return []func(){
		h.SafeAdd,
		h.SafeSub,
		h.SafeMul,
		h.SafeDiv,
	}
}

func (h *uint16Handler) SafeAdd() {
	newHandler := NewUint16Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint16)
	if h.value > math.MaxUint16-value {
		// if a + b > MAX then a + b - MAX
		// to not overflow, a + b - MAX ~~ a - (MAX - b)
		h.value = h.value - (math.MaxUint16 - value)
	} else {
		h.value = h.value + value
	}
}

func (h *uint16Handler) SafeSub() {
	newHandler := NewUint16Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint16)
	if h.value < value {
		// if a - b < 0 then 0 - (a - b)
		// to not undeflow, 0 - (a - b) ~~ b - a
		h.value = value - h.value
	} else {
		h.value = h.value - value
	}
}

func (h *uint16Handler) SafeMul() {
	newHandler := NewUint16Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint16)
	if value == 0 {
		h.value = 0
	} else if h.value > math.MaxUint16/value {
		a := new(big.Int).SetUint64(uint64(h.value))
		b := new(big.Int).SetUint64(uint64(value))
		max := new(big.Int).SetUint64(uint64(math.MaxUint16))
		c := new(big.Int).Mod(new(big.Int).Mul(a, b), max)
		h.value = uint16(c.Uint64())
	} else {
		h.value = h.value * value
	}
}

func (h *uint16Handler) SafeDiv() {
	newHandler := NewUint16Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint16)
	if value == 0 {
		h.value = 0
	} else {
		h.value = h.value / value
	}
}
