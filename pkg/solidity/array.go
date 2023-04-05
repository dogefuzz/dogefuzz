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
var ErrInvalidElement = errors.New("the provided json does not correspond to a valid element for this array")

type arrayHandler struct {
	size              int
	typ               abi.Type
	value             interface{}
	handler           interfaces.TypeHandler
	blockchainContext *BlockchainContext
}

func NewArrayHandler(size int, typ abi.Type, blockchainContext *BlockchainContext) (*arrayHandler, error) {
	handler, err := GetTypeHandler(*typ.Elem, blockchainContext)
	if err != nil {
		return nil, err
	}
	instance := &arrayHandler{
		size:              size,
		typ:               typ,
		value:             make([]any, size),
		handler:           handler,
		blockchainContext: blockchainContext,
	}

	return instance, nil
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
	handler, err := GetTypeHandler(h.typ, h.blockchainContext)
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
	var values []string
	err := json.Unmarshal([]byte(value), &values)
	if err != nil {
		return ErrInvalidArray
	}
	if len(values) != h.size {
		return ErrInvalidArraySize
	}

	sliceType := reflect.SliceOf(h.typ.GetType().Elem())
	arrayValue := reflect.MakeSlice(sliceType, h.size, h.size)
	for idx, value := range values {
		err := h.handler.Deserialize(value)
		if err != nil {
			return ErrInvalidElement
		}
		valueAsInterface := h.handler.GetValue()
		arrayValue.Index(idx).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))
	}
	h.value = arrayValue.Interface()
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

	h.handler.SetValue(temp.Interface())
	mutator := common.RandomChoice(h.handler.GetMutators())
	mutator()

	valueAsInterface := h.handler.GetValue()
	arrayValue.Index(idx).Set(reflect.ValueOf(valueAsInterface).Convert(h.typ.Elem.GetType()))

	h.value = arrayValue.Interface()
}
