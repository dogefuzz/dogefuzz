package solidity

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/dogefuzz/dogefuzz/pkg/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

var ErrInvalidAddress = errors.New("the provided json does not correspond to a address type")

type addressHandler struct {
	value gethcommon.Address
}

func NewAddressHandler() *addressHandler {
	return &addressHandler{}
}

func (h *addressHandler) GetValue() interface{} {
	return h.value
}

func (h *addressHandler) SetValue(value interface{}) {
	h.value = value.(gethcommon.Address)
}

func (h *addressHandler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	addressSeeds := seeds[ADDRESS]
	chosenSeed := common.RandomChoice(addressSeeds)
	return h.Deserialize(chosenSeed)
}

func (h *addressHandler) Serialize() string {
	return h.value.Hex()
}

func (h *addressHandler) Deserialize(value string) error {
	val := gethcommon.HexToAddress(value)
	if (val == gethcommon.Address{}) {
		return ErrInvalidAddress
	}
	h.value = val
	return nil
}

func (h *addressHandler) Generate() {
	const ADDRESS_LENGTH = 20
	rand.Seed(common.Now().UnixNano())

	parts := make([]string, ADDRESS_LENGTH)
	for idx := 0; idx < len(parts); idx++ {
		parts[idx] = fmt.Sprintf("%x", rand.Intn(256))
	}
	hexValue := strings.Join(parts, "")
	h.value = gethcommon.HexToAddress(hexValue)
}

func (h *addressHandler) GetMutators() []func() {
	return []func(){h.GenerateAgainOp}
}

func (h *addressHandler) GenerateAgainOp() {
	h.Generate()
}
