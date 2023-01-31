package solidity

import (
	"errors"
	"strconv"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

var ErrInvalidBool = errors.New("the provided json does not correspond to a boolean type")

type boolHandler struct {
	value bool
}

func NewBoolHandler() *boolHandler {
	return &boolHandler{}
}

func (h *boolHandler) GetValue() interface{} {
	return h.value
}

func (h *boolHandler) SetValue(value interface{}) {
	h.value = value.(bool)
}

func (h *boolHandler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	addressSeeds := seeds[BOOL]
	chosenSeed := common.RandomChoice(addressSeeds)
	return h.Deserialize(chosenSeed)
}

func (h *boolHandler) Serialize() string {
	return strconv.FormatBool(h.value)
}

func (h *boolHandler) Deserialize(value string) error {
	val, err := strconv.ParseBool(value)
	if err != nil {
		return ErrInvalidBool
	}
	h.value = val
	return nil
}

func (h *boolHandler) Generate() {
	h.value = common.RandomChoice([]bool{true, false})
}

func (h *boolHandler) GetMutators() []func() {
	return []func(){
		h.NotOp,
	}
}

func (h *boolHandler) NotOp() {
	h.value = !h.value
}
