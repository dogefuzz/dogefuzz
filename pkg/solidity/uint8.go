package solidity

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const UINT8_BIT_SIZE = 8

var ErrInvalidUint8 = errors.New("the provided string does not correspond to a uint8 type")

type uint8Handler struct {
	value uint8
}

func NewUint8Handler() *uint8Handler {
	return &uint8Handler{}
}

func (h *uint8Handler) GetValue() interface{} {
	return h.value
}

func (h *uint8Handler) SetValue(value interface{}) {
	h.value = value.(uint8)
}

func (h *uint8Handler) Serialize() string {
	return strconv.FormatUint(uint64(h.value), 10)
}

func (h *uint8Handler) Deserialize(value string) error {
	val, err := strconv.ParseUint(value, 10, UINT8_BIT_SIZE)
	if err != nil {
		return ErrInvalidUint8
	}
	h.value = uint8(val)
	return nil
}

func (h *uint8Handler) Generate() {
	rand.Seed(time.Now().Unix())
	h.value = uint8(rand.Intn(2 << (UINT8_BIT_SIZE - 1)))
}

func (h *uint8Handler) GetMutators() []func() {
	return []func(){
		h.SafeAdd,
		h.SafeSub,
		h.SafeMul,
		h.SafeDiv,
	}
}

func (h *uint8Handler) SafeAdd() {
	newHandler := NewUint8Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint8)
	if h.value > math.MaxUint8-value {
		h.value = h.value - (math.MaxUint8 - value) - 1
	}
	h.value = h.value + value
}

func (h *uint8Handler) SafeSub() {
	newHandler := NewUint8Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint8)
	if h.value < value {
		h.value = value - h.value - 1
	}
	h.value = h.value - value
}

func (h *uint8Handler) SafeMul() {
	newHandler := NewUint8Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint8)
	if value == 0 {
		h.value = 0
	}
	if h.value > math.MaxUint8/value {
		h.value = h.value * (math.MaxUint8 / value)
	}
	h.value = h.value * value
}

func (h *uint8Handler) SafeDiv() {
	newHandler := NewUint8Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint8)
	if value == 0 {
		h.value = 0
	}
	h.value = h.value / value
}
