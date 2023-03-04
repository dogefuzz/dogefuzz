package solidity

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/pkg/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

var ErrInvalidAddress = errors.New("the provided json does not correspond to a address type")

type addressHandler struct {
	value              gethcommon.Address
	availableAddresses []gethcommon.Address
}

func NewAddressHandler(addresses []string) (*addressHandler, error) {
	availableAddresses := make([]gethcommon.Address, len(addresses))
	for idx, address := range addresses {
		val, err := convertStringToGethAddress(address)
		if err != nil {
			return nil, err
		}
		availableAddresses[idx] = val
	}
	return &addressHandler{availableAddresses: availableAddresses}, nil
}

func (h *addressHandler) GetValue() interface{} {
	return h.value
}

func (h *addressHandler) SetValue(value interface{}) {
	h.value = value.(gethcommon.Address)
}

func (h *addressHandler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	addressSeeds := seeds[ADDRESS]
	if len(addressSeeds) > 0 {
		chosenSeed := common.RandomChoice(addressSeeds)
		return h.Deserialize(chosenSeed)
	}
	return nil
}

func (h *addressHandler) Serialize() string {
	return h.value.Hex()
}

func (h *addressHandler) Deserialize(value string) error {
	val, err := convertStringToGethAddress(value)
	if err != nil {
		return ErrInvalidAddress
	}
	h.value = val
	return nil
}

func (h *addressHandler) Generate() {
	h.value = common.RandomChoice(h.availableAddresses)
}

func (h *addressHandler) GetMutators() []func() {
	return []func(){h.NotOp, h.ChooseAgain}
}

func (h *addressHandler) NotOp() {
}

func (h *addressHandler) ChooseAgain() {
	h.value = common.RandomChoice(h.availableAddresses)
}

func convertStringToGethAddress(address string) (gethcommon.Address, error) {
	val := gethcommon.HexToAddress(address)
	if (val == gethcommon.Address{}) {
		return gethcommon.Address{}, ErrInvalidAddress
	}
	return val, nil
}
