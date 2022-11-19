package solidity

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"

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
	rand.Seed(time.Now().Unix())
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
		h.value = h.value - (math.MaxInt64 - value) - 1
	}
	h.value = h.value + value
}

func (h *int64Handler) SafeSub() {
	newHandler := NewInt64Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int64)
	if h.value < value {
		h.value = value - h.value - 1
	}
	h.value = h.value - value
}

func (h *int64Handler) SafeMul() {
	newHandler := NewInt64Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int64)
	if value == 0 {
		h.value = 0
	}
	if h.value > math.MaxInt64/value {
		h.value = h.value * (math.MaxInt64 / value)
	}
	h.value = h.value * value
}

func (h *int64Handler) SafeDiv() {
	newHandler := NewInt64Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int64)
	if value == 0 {
		h.value = 0
	}
	h.value = h.value / value
}
