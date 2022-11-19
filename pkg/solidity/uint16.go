package solidity

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"
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
	rand.Seed(time.Now().Unix())
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
		h.value = h.value - (math.MaxUint16 - value) - 1
	}
	h.value = h.value + value
}

func (h *uint16Handler) SafeSub() {
	newHandler := NewUint16Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint16)
	if h.value < value {
		h.value = value - h.value - 1
	}
	h.value = h.value - value
}

func (h *uint16Handler) SafeMul() {
	newHandler := NewUint16Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint16)
	if value == 0 {
		h.value = 0
	}
	if h.value > math.MaxUint16/value {
		h.value = h.value * (math.MaxUint16 / value)
	}
	h.value = h.value * value
}

func (h *uint16Handler) SafeDiv() {
	newHandler := NewUint16Handler()
	newHandler.Generate()
	value := newHandler.GetValue().(uint16)
	if value == 0 {
		h.value = 0
	}
	h.value = h.value / value
}
