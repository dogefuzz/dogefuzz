package solidity

import (
	"errors"
	"math"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

const INT16_BIT_SIZE = 16

var ErrInvalidInt16 = errors.New("the provided string does not correspond to a int16 type")

type int16Handler struct {
	value int16
}

func NewInt16Handler() *int16Handler {
	return &int16Handler{}
}

func (h *int16Handler) GetValue() interface{} {
	return h.value
}

func (h *int16Handler) SetValue(value interface{}) {
	h.value = value.(int16)
}

func (h *int16Handler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	addressSeeds := seeds[INT16]
	chosenSeed := common.RandomChoice(addressSeeds)
	return h.Deserialize(chosenSeed)
}

func (h *int16Handler) Serialize() string {
	return strconv.FormatInt(int64(h.value), 10)
}

func (h *int16Handler) Deserialize(value string) error {
	val, err := strconv.ParseInt(value, 10, INT16_BIT_SIZE)
	if err != nil {
		return ErrInvalidInt16
	}
	h.value = int16(val)
	return nil
}

func (h *int16Handler) Generate() {
	rand.Seed(common.Now().Unix())
	h.value = int16(rand.Intn(2<<(INT16_BIT_SIZE*2)) + math.MinInt16)
}

func (h *int16Handler) GetMutators() []func() {
	return []func(){
		h.SafeAdd,
		h.SafeSub,
		h.SafeMul,
		h.SafeDiv,
	}
}

func (h *int16Handler) SafeAdd() {
	newHandler := NewInt16Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int16)
	if h.value > math.MaxInt16-value {
		// if a + b > MAX then a + b - MAX
		// to not overflow, a + b - MAX ~~ a - (MAX - b)
		h.value = h.value - (math.MaxInt16 - value)
	} else {
		h.value = h.value + value
	}
}

func (h *int16Handler) SafeSub() {
	newHandler := NewInt16Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int16)
	if h.value < math.MinInt16+value {
		// if a - b < MIN then MIN - (a - b)
		// to not undeflow, MIN - (a - b) ~~ (b + MIN) - a
		h.value = (value + math.MinInt16) - h.value
	} else {
		h.value = h.value - value
	}
}

func (h *int16Handler) SafeMul() {
	newHandler := NewInt16Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int16)
	if value == 0 {
		h.value = 0
	} else if h.value > math.MaxInt16/value {
		a := new(big.Int).SetInt64(int64(h.value))
		b := new(big.Int).SetInt64(int64(value))
		max := new(big.Int).SetInt64(int64(math.MaxInt16))
		c := new(big.Int).Mod(new(big.Int).Mul(a, b), max)
		h.value = int16(c.Int64())
	} else {
		h.value = h.value * value
	}
}

func (h *int16Handler) SafeDiv() {
	newHandler := NewInt16Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int16)
	if value == 0 {
		h.value = 0
	} else {
		h.value = h.value / value
	}
}
