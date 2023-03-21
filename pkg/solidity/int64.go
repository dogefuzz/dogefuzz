package solidity

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

const INT64_BIT_SIZE = 64

var ErrInvalidInt64 = errors.New("the provided string does not correspond to a int64 type")

type int64Handler struct {
	value int64
}

func NewInt64Handler() *int64Handler {
	return &int64Handler{}
}

func (h *int64Handler) GetValue() interface{} {
	return h.value
}

func (h *int64Handler) SetValue(value interface{}) {
	h.value = value.(int64)
}

func (h *int64Handler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	addressSeeds := seeds[INT64]
	chosenSeed := common.RandomChoice(addressSeeds)
	return h.Deserialize(chosenSeed)
}

func (h *int64Handler) Serialize() string {
	return strconv.FormatInt(int64(h.value), 10)
}

func (h *int64Handler) Deserialize(value string) error {
	val, err := strconv.ParseInt(value, 10, INT64_BIT_SIZE)
	if err != nil {
		return ErrInvalidInt64
	}
	h.value = int64(val)
	return nil
}

func (h *int64Handler) Generate() {
	rand.Seed(common.Now().Unix())
	h.value = common.RandomChoice([]int64{1, -1}) * rand.Int63()
}

func (h *int64Handler) GetMutators() []func() {
	return []func(){
		h.SafeAdd,
		h.SafeSub,
		h.SafeMul,
		h.SafeDiv,
	}
}

func (h *int64Handler) SafeAdd() {
	newHandler := NewInt64Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int64)
	if h.value > math.MaxInt64-value {
		// if a + b > MAX then a + b - MAX
		// to not overflow, a + b - MAX ~~ a - (MAX - b)
		h.value = h.value - (math.MaxInt64 - value)
	} else {
		h.value = h.value + value
	}
}

func (h *int64Handler) SafeSub() {
	newHandler := NewInt64Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int64)
	if h.value < math.MinInt64+value {
		// if a - b < MIN then MIN - (a - b)
		// to not undeflow, MIN - (a - b) ~~ (b + MIN) - a
		h.value = (value + math.MinInt64) - h.value
	} else {
		h.value = h.value - value
	}
}

func (h *int64Handler) SafeMul() {
	newHandler := NewInt64Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int64)
	if value == 0 {
		h.value = 0
	} else if h.value > math.MaxInt64/value {
		a := new(big.Int).SetInt64(h.value)
		b := new(big.Int).SetInt64(value)
		max := new(big.Int).SetInt64(math.MaxInt64)
		c := new(big.Int).Mod(new(big.Int).Mul(a, b), max)
		h.value = c.Int64()
	} else {
		h.value = h.value * value
	}
}

func (h *int64Handler) SafeDiv() {
	newHandler := NewInt64Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int64)
	if value == 0 {
		h.value = 0
	} else {
		h.value = h.value / value
	}
}
