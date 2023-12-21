package solidity

import (
	"errors"
	"math/rand"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

var ErrInvalidString = errors.New("the provided json does not correspond to a string type")

type stringHandler struct {
	value string
}

func NewStringHandler() *stringHandler {
	return &stringHandler{}
}

func (h *stringHandler) GetValue() interface{} {
	return h.value
}

func (h *stringHandler) SetValue(value interface{}) {
	h.value = value.(string)
}

func (h *stringHandler) LoadSeedsAndChooseOneRandomly(seeds common.Seeds) error {
	addressSeeds := seeds[STRING]
	chosenSeed := common.RandomChoice(addressSeeds)
	return h.Deserialize(chosenSeed)
}

func (h *stringHandler) Serialize() string {
	return h.value
}

func (h *stringHandler) Deserialize(value string) error {
	h.value = value
	return nil
}

func (h *stringHandler) Generate() {
	rand.Seed(common.Now().UnixNano())
	var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := rand.Intn(256)
	if length == 0 {
		length = 1
	}
	wordSlice := make([]rune, length)
	for i := range wordSlice {
		wordSlice[i] = alphabet[rand.Intn(len(alphabet))]
	}
	h.value = string(wordSlice)
}

func (h *stringHandler) GetMutators() []func() {
	return []func(){
		h.ChangeCharacter,
		h.AddCharacter,
		h.RemoveCharacter,
	}
}

func (h *stringHandler) ChangeCharacter() {
	rand.Seed(common.Now().UnixNano())
	idx := rand.Intn(len(h.value))
	var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	letter := alphabet[rand.Intn(len(alphabet))]
	h.value = h.value[:idx] + string([]rune{letter}) + h.value[idx+1:]
}

func (h *stringHandler) AddCharacter() {
	rand.Seed(common.Now().UnixNano())
	idx := rand.Intn(len(h.value))
	var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	letter := alphabet[rand.Intn(len(alphabet))]
	h.value = h.value[:idx] + string([]rune{letter}) + h.value[idx:]
}

func (h *stringHandler) RemoveCharacter() {
	rand.Seed(common.Now().UnixNano())
	idx := rand.Intn(len(h.value))
	h.value = h.value[:idx] + h.value[idx+1:]
}
