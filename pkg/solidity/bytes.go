package solidity

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var ErrInvalidBytes = errors.New("the provided json does not correspond to a bytes type")
var ErrInvalidBytesSize = errors.New("the provided json does not correspond to a bytes with size required by this handler")

type bytesHandler struct {
	value   []byte
	handler interfaces.TypeHandler
}

func NewBytesHandler(typ abi.Type, blockchainContext *BlockchainContext) (*bytesHandler, error) {
	val := make([]byte, 0)

	uint8Typ, err := abi.NewType("uint8", "", nil)
	if err != nil {
		return nil, err
	}

	handler, err := GetTypeHandler(uint8Typ, blockchainContext)
	if err != nil {
		return nil, err
	}

	return &bytesHandler{value: val, handler: handler}, nil
}

func (h *bytesHandler) GetValue() interface{} {
	return h.value
}

func (h *bytesHandler) SetValue(value interface{}) {
	h.value = value.([]byte)
}

func (h *bytesHandler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	rand.Seed(time.Now().UnixNano())

	randomSize := rand.Intn(16)
	val := make([]byte, randomSize)
	for idx := 0; idx < randomSize; idx++ {
		err := h.handler.LoadSeedsAndChooseOneRandomly(seeds)
		if err != nil {
			return err
		}

		val[idx] = h.handler.GetValue().(byte)
	}

	h.value = val
	return nil
}

func (h *bytesHandler) Serialize() string {
	js, _ := json.Marshal(h.value)
	return string(js)
}

func (h *bytesHandler) Deserialize(value string) error {
	var val []byte
	err := json.Unmarshal([]byte(value), &val)
	if err != nil {
		return ErrInvalidBytes
	}
	h.value = val
	return nil
}

func (h *bytesHandler) Generate() {
	rand.Seed(time.Now().UnixNano())
	size := rand.Intn(16)
	val := make([]byte, size)
	for idx := 0; idx < size; idx++ {
		h.handler.Generate()
		val[idx] = h.handler.GetValue().(byte)
	}
	h.value = val
}

func (h *bytesHandler) GetMutators() []func() {
	return []func(){
		h.MutateElementOp,
		h.AddElementOp,
		h.RemoveElementOp,
	}
}

func (h *bytesHandler) MutateElementOp() {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(h.value))
	temp := h.value[idx]

	h.handler.SetValue(temp)
	mutator := common.RandomChoice(h.handler.GetMutators())
	mutator()

	h.value[idx] = h.handler.GetValue().(byte)
}

func (h *bytesHandler) AddElementOp() {
	rand.Seed(time.Now().UnixNano())
	insertionIdx := rand.Intn(len(h.value))

	h.handler.Generate()
	newValue := h.handler.GetValue().(byte)

	h.value = append(h.value[:insertionIdx], append([]byte{newValue}, h.value[insertionIdx:]...)...)
}

func (h *bytesHandler) RemoveElementOp() {
	rand.Seed(time.Now().UnixNano())
	removeIdx := rand.Intn(len(h.value))
	h.value = append(h.value[:removeIdx], h.value[removeIdx+1:]...)
}
