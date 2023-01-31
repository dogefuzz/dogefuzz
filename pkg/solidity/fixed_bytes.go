package solidity

import (
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var ErrInvalidFixedBytes = errors.New("the provided json does not correspond to a fixed bytes type")
var ErrInvalidFixedBytesSize = errors.New("the provided json does not correspond to a fixed bytes with size required by this handler")

type fixedBytesHandler struct {
	size    int
	value   []byte
	handler interfaces.TypeHandler
}

func NewFixedBytesHandler(typ abi.Type) (*fixedBytesHandler, error) {
	val := make([]byte, typ.Size)
	uint8Typ, err := abi.NewType("uint8", "", nil)
	if err != nil {
		return nil, err
	}

	handler, err := GetTypeHandler(uint8Typ)
	if err != nil {
		return nil, err
	}

	return &fixedBytesHandler{size: typ.Size, value: val, handler: handler}, nil
}

func (h *fixedBytesHandler) GetValue() interface{} {
	arrayTyp := reflect.ArrayOf(h.size, reflect.TypeOf(byte(0)))
	array := reflect.New(arrayTyp).Elem()
	slice := reflect.ValueOf(h.value)
	for idx := 0; idx < h.size; idx++ {
		array.Index(idx).Set(slice.Index(idx))
	}
	return array.Interface()
}

func (h *fixedBytesHandler) SetValue(value interface{}) {
	h.value = value.([]byte)
}

func (h *fixedBytesHandler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	val := make([]byte, h.size)
	for idx := 0; idx < h.size; idx++ {
		err := h.handler.LoadSeedsAndChooseOneRandomly(seeds)
		if err != nil {
			return err
		}

		val[idx] = h.handler.GetValue().(byte)
	}

	h.value = val
	return nil
}

func (h *fixedBytesHandler) Serialize() string {
	js, _ := json.Marshal(h.value)
	return string(js)
}

func (h *fixedBytesHandler) Deserialize(value string) error {
	var val []byte
	err := json.Unmarshal([]byte(value), &val)
	if err != nil {
		return ErrInvalidFixedBytes
	}
	if len(val) != h.size {
		return ErrInvalidFixedBytesSize
	}
	h.value = val
	return nil
}

func (h *fixedBytesHandler) Generate() {
	val := make([]byte, h.size)
	for idx := 0; idx < h.size; idx++ {
		h.handler.Generate()
		val[idx] = h.handler.GetValue().(byte)
	}
	h.value = val
}

func (h *fixedBytesHandler) GetMutators() []func() {
	return []func(){h.MutateElementOp}
}

func (h *fixedBytesHandler) MutateElementOp() {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(h.size)
	temp := h.value[idx]

	h.handler.SetValue(temp)
	mutator := common.RandomChoice(h.handler.GetMutators())
	mutator()

	h.value[idx] = h.handler.GetValue().(byte)
}
