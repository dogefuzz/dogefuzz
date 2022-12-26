package solidity

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const INT8_BIT_SIZE = 8

var ErrInvalidInt8 = errors.New("the provided string does not correspond to a int8 type")

type int8Handler struct {
	value int8
}

func NewInt8Handler() *int8Handler {
	return &int8Handler{}
}

func (h *int8Handler) GetValue() interface{} {
	return h.value
}

func (h *int8Handler) SetValue(value interface{}) {
	h.value = value.(int8)
}

func (h *int8Handler) GetType() TypeIdentifier {
	return INT8
}

func (h *int8Handler) Serialize() string {
	return strconv.FormatInt(int64(h.value), 10)
}

func (h *int8Handler) Deserialize(value string) error {
	val, err := strconv.ParseInt(value, 10, INT8_BIT_SIZE)
	if err != nil {
		return ErrInvalidInt8
	}
	h.value = int8(val)
	return nil
}

func (h *int8Handler) Generate() {
	rand.Seed(time.Now().Unix())
	h.value = int8(rand.Intn(2<<(INT8_BIT_SIZE*2)) + math.MinInt8)
}

func (h *int8Handler) GetMutators() []func() {
	return []func(){
		h.SafeAdd,
		h.SafeSub,
		h.SafeMul,
		h.SafeDiv,
	}
}

func (h *int8Handler) SafeAdd() {
	newHandler := NewInt8Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int8)
	if h.value > math.MaxInt8-value {
		h.value = h.value - (math.MaxInt8 - value) - 1
	}
	h.value = h.value + value
}

func (h *int8Handler) SafeSub() {
	newHandler := NewInt8Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int8)
	if h.value < value {
		h.value = value - h.value - 1
	}
	h.value = h.value - value
}

func (h *int8Handler) SafeMul() {
	newHandler := NewInt8Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int8)
	if value == 0 {
		h.value = 0
	}
	if h.value > math.MaxInt8/value {
		h.value = h.value * (math.MaxInt8 / value)
	}
	h.value = h.value * value
}

func (h *int8Handler) SafeDiv() {
	newHandler := NewInt8Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(int8)
	if value == 0 {
		h.value = 0
	}
	h.value = h.value / value
}
