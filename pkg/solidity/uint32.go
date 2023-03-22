package solidity

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

const UINT32_BIT_SIZE = 32

var ErrInvalidUint32 = errors.New("the provided string does not correspond to a uint32 type")

type uint32Handler struct {
	value uint32
}

func NewUint32Handler() *uint32Handler {
	return &uint32Handler{}
}

func (h *uint32Handler) GetValue() interface{} {
	return h.value
}

func (h *uint32Handler) SetValue(value interface{}) {
	h.value = value.(uint32)
}

func (h *uint32Handler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	addressSeeds := seeds[UINT32]
	chosenSeed := common.RandomChoice(addressSeeds)
	return h.Deserialize(chosenSeed)
}

func (h *uint32Handler) Serialize() string {
	return strconv.FormatUint(uint64(h.value), 10)
}

func (h *uint32Handler) Deserialize(value string) error {
	val, err := strconv.ParseUint(value, 10, UINT32_BIT_SIZE)
	if err != nil {
		return ErrInvalidUint32
	}
	h.value = uint32(val)
	return nil
}

func (h *uint32Handler) Generate() {
	rand.Seed(common.Now().Unix())
	h.value = rand.Uint32()
}

func (h *uint32Handler) GetMutators() []func() {
	return []func(){
		h.SafeAdd,
		h.SafeSub,
		h.SafeMul,
		h.SafeDiv,
	}
}

func (h *uint32Handler) SafeAdd() {
	newHandler := NewUint32Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint32)
	if h.value > math.MaxUint32-value {
		// if a + b > MAX then a + b - MAX
		// to not overflow, a + b - MAX ~~ a - (MAX - b)
		h.value = h.value - (math.MaxUint32 - value)
	} else {
		h.value = h.value + value
	}
}

func (h *uint32Handler) SafeSub() {
	newHandler := NewUint32Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint32)
	if h.value < value {
		// if a - b < 0 then 0 - (a - b)
		// to not undeflow, 0 - (a - b) ~~ b - a
		h.value = value - h.value
	} else {
		h.value = h.value - value
	}
}

func (h *uint32Handler) SafeMul() {
	newHandler := NewUint32Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint32)
	if value == 0 {
		h.value = 0
	} else if h.value > math.MaxUint32/value {
		a := new(big.Int).SetUint64(uint64(h.value))
		b := new(big.Int).SetUint64(uint64(value))
		max := new(big.Int).SetUint64(uint64(math.MaxUint32))
		c := new(big.Int).Mod(new(big.Int).Mul(a, b), max)
		h.value = uint32(c.Uint64())
	} else {
		h.value = h.value * value
	}
}

func (h *uint32Handler) SafeDiv() {
	newHandler := NewUint32Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint32)
	if value == 0 {
		h.value = 0
	} else {
		h.value = h.value / value
	}
}
