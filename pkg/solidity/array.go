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

var ErrInvalidArray = errors.New("the provided json does not correspond to a array type")
var ErrInvalidArraySize = errors.New("the provided json does not correspond to a array with size required by this handler")

type arrayHandler struct {
	size    int
	typ     abi.Type
	value   interface{}
	handler interfaces.TypeHandler
}

func NewArrayHandler(size int, typ abi.Type) (*arrayHandler, error) {
	handler, err := GetTypeHandler(*typ.Elem)
	if err != nil {
		return nil, err
	}
	return &arrayHandler{size: size, typ: typ, value: make([]any, size), handler: handler}, nil
}

func (h *arrayHandler) GetValue() interface{} {
	arrayTyp := reflect.ArrayOf(h.size, h.typ.GetType().Elem())
	array := reflect.New(arrayTyp).Elem()
	slice := reflect.ValueOf(h.value)
	for idx := 0; idx < h.size; idx++ {
		array.Index(idx).Set(slice.Index(idx))
	}
	return array.Interface()
}

func (h *arrayHandler) SetValue(value interface{}) {
	h.value = value
}

func (h *arrayHandler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	handler, err := GetTypeHandler(h.typ)
	if err != nil {
		return err
	}

	sliceType := reflect.SliceOf(h.typ.GetType().Elem())
	arrayValue := reflect.MakeSlice(sliceType, h.size, h.size)
	for idx := 0; idx < h.size; idx++ {
		err = handler.LoadSeedsAndChooseOneRandomly(seeds)
		if err != nil {
			return err
		}

		valueAsInterface := handler.GetValue()
		arrayValue.Index(idx).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))
	}

	h.value = arrayValue.Interface()
	return nil
}

func (h *arrayHandler) Serialize() string {
	js, _ := json.Marshal(h.value)
	return string(js)
}

func (h *arrayHandler) Deserialize(value string) error {
	var val []interface{}
	err := json.Unmarshal([]byte(value), &val)
	if err != nil {
		return ErrInvalidArray
	}
	if len(val) != h.size {
		return ErrInvalidArraySize
	}
	h.value = val
	return nil
}

func (h *arrayHandler) Generate() {
	sliceType := reflect.SliceOf(h.typ.GetType().Elem())
	arrayValue := reflect.MakeSlice(sliceType, h.size, h.size)
	for idx := 0; idx < h.size; idx++ {
		h.handler.Generate()
		valueAsInterface := h.handler.GetValue()
		arrayValue.Index(idx).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))
	}
	h.value = arrayValue.Interface()
}

func (h *arrayHandler) GetMutators() []func() {
	return []func(){h.MutateElementOp}
}

func (h *arrayHandler) MutateElementOp() {
	arrayValue := reflect.ValueOf(h.value)

	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(arrayValue.Len())
	temp := arrayValue.Index(idx)

	h.handler.SetValue(temp)
	mutator := common.RandomChoice(h.handler.GetMutators())
	mutator()

	valueAsInterface := h.handler.GetValue()
	arrayValue.Index(idx).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))

	h.value = arrayValue.Interface()
}