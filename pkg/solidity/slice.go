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
	typ               abi.Type
	value             interface{}
	handler           interfaces.TypeHandler
	blockchainContext *BlockchainContext
}

func NewSliceHandler(typ abi.Type, blockchainContext *BlockchainContext) (*sliceHandler, error) {
	handler, err := GetTypeHandler(*typ.Elem, blockchainContext)
	if err != nil {
		return nil, err
	}
	instance := &sliceHandler{
		typ:               typ,
		value:             make([]any, 0),
		handler:           handler,
		blockchainContext: blockchainContext,
	}
	return instance, nil
}

func (h *sliceHandler) GetValue() interface{} {
	return h.value
}

func (h *sliceHandler) SetValue(value interface{}) {
	h.value = value
}

func (h *sliceHandler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	rand.Seed(time.Now().UnixNano())

	randomSize := rand.Intn(16)
	sliceValue := reflect.MakeSlice(h.typ.GetType(), randomSize, randomSize)
	for idx := 0; idx < randomSize; idx++ {
		err := h.handler.LoadSeedsAndChooseOneRandomly(seeds)
		if err != nil {
			return err
		}

		valueAsInterface := h.handler.GetValue()
		sliceValue.Index(idx).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))
	}

	h.value = sliceValue.Interface()
	return nil
}

func (h *sliceHandler) Serialize() string {
	arrayValue := reflect.ValueOf(h.value)
	values := make([]string, arrayValue.Len())

	for idx := 0; idx < arrayValue.Len(); idx++ {
		h.handler.SetValue(arrayValue.Index(idx).Interface())
		values[idx] = h.handler.Serialize()
	}

	js, _ := json.Marshal(values)
	return string(js)
}

func (h *sliceHandler) Deserialize(value string) error {
	var values []string
	err := json.Unmarshal([]byte(value), &values)
	if err != nil {
		return ErrInvalidSlice
	}

	slice := reflect.MakeSlice(h.typ.GetType(), len(values), len(values))
	for idx, value := range values {
		err := h.handler.Deserialize(value)
		if err != nil {
			return err
		}
		valueAsInterface := h.handler.GetValue()
		slice.Index(idx).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))
	}
	h.value = slice.Interface()
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
	if arrayValue.Len() == 0 {
		return
	}

	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(arrayValue.Len())
	temp := arrayValue.Index(idx)

	h.handler.SetValue(temp.Interface())
	mutator := common.RandomChoice(h.handler.GetMutators())
	mutator()

	valueAsInterface := h.handler.GetValue()
	arrayValue.Index(idx).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))

	h.value = arrayValue.Interface()

}

func (h *sliceHandler) AddElementOp() {
	arrayValue := reflect.ValueOf(h.value)
	if arrayValue.Len() == 0 {
		array := reflect.MakeSlice(arrayValue.Type(), 1, 1)
		h.handler.Generate()
		valueAsInterface := h.handler.GetValue()
		array.Index(0).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))
		h.value = array.Interface()
		return
	}
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

	for idx := insertionIdx; idx < arrayValue.Len(); idx++ {
		value := arrayValue.Index(idx)
		appendedArray.Index(idx + 1).Set(value)
	}

	h.value = appendedArray.Interface()
}

func (h *sliceHandler) RemoveElementOp() {
	arrayValue := reflect.ValueOf(h.value)
	if arrayValue.Len() == 0 {
		return
	}

	resultArray := reflect.MakeSlice(arrayValue.Type(), arrayValue.Len()-1, arrayValue.Len()-1)

	rand.Seed(time.Now().UnixNano())
	removeIdx := rand.Intn(arrayValue.Len())
	for idx := 0; idx < removeIdx; idx++ {
		value := arrayValue.Index(idx)
		resultArray.Index(idx).Set(value)
	}

	for idx := removeIdx + 1; idx < arrayValue.Len(); idx++ {
		value := arrayValue.Index(idx)
		resultArray.Index(idx - 1).Set(value)
	}

	h.value = resultArray.Interface()
}
