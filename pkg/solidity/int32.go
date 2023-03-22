package solidity

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

const INT32_BIT_SIZE = 32

var ErrInvalidInt32 = errors.New("the provided string does not correspond to a int32 type")

type int32Handler struct {
	value int32
}

func NewInt32Handler() *int32Handler {
	return &int32Handler{}
}

func (h *int32Handler) GetValue() interface{} {
	return h.value
}

func (h *int32Handler) SetValue(value interface{}) {
	h.value = value.(int32)
}

func (h *int32Handler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	addressSeeds := seeds[INT32]
	chosenSeed := common.RandomChoice(addressSeeds)
	return h.Deserialize(chosenSeed)
}

func (h *int32Handler) Serialize() string {
	return strconv.FormatInt(int64(h.value), 10)
}

func (h *int32Handler) Deserialize(value string) error {
	val, err := strconv.ParseInt(value, 10, INT32_BIT_SIZE)
	if err != nil {
		return ErrInvalidInt32
	}
	h.value = int32(val)
	return nil
}

func (h *int32Handler) Generate() {
	rand.Seed(common.Now().Unix())
	h.value = common.RandomChoice([]int32{1, -1}) * rand.Int31()
}

func (h *int32Handler) GetMutators() []func() {
	return []func(){
		h.SafeAdd,
		h.SafeSub,
		h.SafeMul,
		h.SafeDiv,
	}
}

func (h *int32Handler) SafeAdd() {
	newHandler := NewInt32Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int32)
	if h.value > math.MaxInt32-value {
		// if a + b > MAX then a + b - MAX
		// to not overflow, a + b - MAX ~~ a - (MAX - b)
		h.value = h.value - (math.MaxInt32 - value)
	} else {
		h.value = h.value + value
	}
}

func (h *int32Handler) SafeSub() {
	newHandler := NewInt32Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int32)
	if h.value < math.MinInt32+value {
		// if a - b < MIN then MIN - (a - b)
		// to not undeflow, MIN - (a - b) ~~ (b + MIN) - a
		h.value = (value + math.MinInt32) - h.value
	} else {
		h.value = h.value - value
	}
}

func (h *int32Handler) SafeMul() {
	newHandler := NewInt32Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int32)
	if value == 0 {
		h.value = 0
	} else if h.value > math.MaxInt32/value {
		a := new(big.Int).SetInt64(int64(h.value))
		b := new(big.Int).SetInt64(int64(value))
		max := new(big.Int).SetInt64(int64(math.MaxInt32))
		c := new(big.Int).Mod(new(big.Int).Mul(a, b), max)
		h.value = int32(c.Int64())
	} else {
		h.value = h.value * value
	}
}

func (h *int32Handler) SafeDiv() {
	newHandler := NewInt32Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int32)
	if value == 0 {
		h.value = 0
	} else {
		h.value = h.value / value
	}
}
