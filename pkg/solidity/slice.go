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

var ErrInvalidSlice = errors.New("the provided json does not correspond to a slice type")
var ErrInvalidSliceSize = errors.New("the provided json does not correspond to a slice with size required by this handler")

type sliceHandler struct {
	typ     abi.Type
	value   interface{}
	handler interfaces.TypeHandler
}

func NewSliceHandler(typ abi.Type) (*sliceHandler, error) {
	handler, err := GetTypeHandler(*typ.Elem)
	if err != nil {
		return nil, err
	}
	return &sliceHandler{typ: typ, value: make([]any, 0), handler: handler}, nil
}

func (h *sliceHandler) GetValue() interface{} {
	return h.value
}

func (h *sliceHandler) SetValue(value interface{}) {
	h.value = value
}

func (h *sliceHandler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	rand.Seed(time.Now().UnixNano())

	handler, err := GetTypeHandler(h.typ)
	if err != nil {
		return err
	}

	randomSize := rand.Intn(16)
	sliceValue := reflect.MakeSlice(h.typ.GetType(), randomSize, randomSize)
	for idx := 0; idx < randomSize; idx++ {
		err = handler.LoadSeedsAndChooseOneRandomly(seeds)
		if err != nil {
			return err
		}

		valueAsInterface := handler.GetValue()
		sliceValue.Index(idx).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))
	}

	h.value = sliceValue.Interface()
	return nil
}

func (h *sliceHandler) Serialize() string {
	js, _ := json.Marshal(h.value)
	return string(js)
}

func (h *sliceHandler) Deserialize(value string) error {
	var val []interface{}
	err := json.Unmarshal([]byte(value), &val)
	if err != nil {
		return ErrInvalidSlice
	}
	h.value = val
	return nil
}

func (h *sliceHandler) Generate() {
	rand.Seed(time.Now().UnixNano())
	size := rand.Intn(16)
	val := reflect.MakeSlice(h.typ.GetType(), size, size)
	for idx := 0; idx < size; idx++ {
		h.handler.Generate()
		valueAsInterface := h.handler.GetValue()
		val.Index(idx).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))
	}
	h.value = val.Interface()
}

func (h *sliceHandler) GetMutators() []func() {
	return []func(){
		h.MutateElementOp,
		h.AddElementOp,
		h.RemoveElementOp,
	}
}

func (h *sliceHandler) MutateElementOp() {
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

func (h *sliceHandler) AddElementOp() {
	arrayValue := reflect.ValueOf(h.value)
	appendedArray := reflect.MakeSlice(arrayValue.Type(), arrayValue.Len()+1, arrayValue.Len()+1)

	rand.Seed(time.Now().UnixNano())
	insertionIdx := rand.Intn(arrayValue.Len())
	for idx := 0; idx < insertionIdx; idx++ {
		value := arrayValue.Index(idx)
		appendedArray.Index(idx).Set(value)
	}
	h.handler.Generate()
	valueAsInterface := h.handler.GetValue()
	appendedArray.Index(insertionIdx).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))

	for idx := insertionIdx + 1; idx < arrayValue.Len(); idx++ {
		value := arrayValue.Index(idx - 1)
		appendedArray.Index(idx).Set(value)
	}

	h.value = appendedArray.Interface()
}

func (h *sliceHandler) RemoveElementOp() {
	arrayValue := reflect.ValueOf(h.value)
	resultArray := reflect.MakeSlice(arrayValue.Type(), arrayValue.Len()-1, arrayValue.Len()-1)

	rand.Seed(time.Now().UnixNano())
	removeIdx := rand.Intn(arrayValue.Len())
	for idx := 0; idx < removeIdx; idx++ {
		value := arrayValue.Index(idx)
		resultArray.Index(idx).Set(value)
	}

	for idx := removeIdx + 1; idx < arrayValue.Len(); idx++ {
		value := arrayValue.Index(idx - 1)
		resultArray.Index(idx).Set(value)
	}

	h.value = resultArray.Interface()
}
