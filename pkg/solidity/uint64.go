package solidity

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

const UINT64_BIT_SIZE = 64

var ErrInvalidUint64 = errors.New("the provided string does not correspond to a uint64 type")

type uint64Handler struct {
	value uint64
}

func NewUint64Handler() *uint64Handler {
	return &uint64Handler{}
}

func (h *uint64Handler) GetValue() interface{} {
	return h.value
}

func (h *uint64Handler) SetValue(value interface{}) {
	h.value = value.(uint64)
}

func (h *uint64Handler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	addressSeeds := seeds[UINT64]
	chosenSeed := common.RandomChoice(addressSeeds)
	return h.Deserialize(chosenSeed)
}

func (h *uint64Handler) Serialize() string {
	return strconv.FormatUint(uint64(h.value), 10)
}

func (h *uint64Handler) Deserialize(value string) error {
	val, err := strconv.ParseUint(value, 10, UINT64_BIT_SIZE)
	if err != nil {
		return ErrInvalidUint64
	}
	h.value = uint64(val)
	return nil
}

func (h *uint64Handler) Generate() {
	rand.Seed(common.Now().Unix())
	h.value = rand.Uint64()
}

func (h *uint64Handler) GetMutators() []func() {
	return []func(){
		h.SafeAdd,
		h.SafeSub,
		h.SafeMul,
		h.SafeDiv,
	}
}

func (h *uint64Handler) SafeAdd() {
	newHandler := NewUint64Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint64)
	if h.value > math.MaxUint64-value {
		// if a + b > MAX then a + b - MAX
		// to not overflow, a + b - MAX ~~ a - (MAX - b)
		h.value = h.value - (math.MaxUint64 - value)
	} else {
		h.value = h.value + value
	}
}

func (h *uint64Handler) SafeSub() {
	newHandler := NewUint64Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint64)
	if h.value < value {
		// if a - b < 0 then 0 - (a - b)
		// to not undeflow, 0 - (a - b) ~~ b - a
		h.value = value - h.value
	} else {
		h.value = h.value - value
	}
}

func (h *uint64Handler) SafeMul() {
	newHandler := NewUint64Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint64)
	if value == 0 {
		h.value = 0
	} else if h.value > math.MaxUint64/value {
		a := new(big.Int).SetUint64(h.value)
		b := new(big.Int).SetUint64(value)
		max := new(big.Int).SetUint64(math.MaxUint64)
		c := new(big.Int).Mod(new(big.Int).Mul(a, b), max)
		h.value = c.Uint64()
	} else {
		h.value = h.value * value
	}
}

func (h *uint64Handler) SafeDiv() {
	newHandler := NewUint64Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint64)
	if value == 0 {
		h.value = 0
	} else {
		h.value = h.value / value
	}
}
