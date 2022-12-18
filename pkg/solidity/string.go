package solidity

import (
	"errors"
	"math/rand"
	"time"
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

func (h *stringHandler) Serialize() string {
	return h.value
}

func (h *stringHandler) Deserialize(value string) error {
	h.value = value
	return nil
}

func (h *stringHandler) Generate() {
	rand.Seed(time.Now().Unix())
	var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	length := rand.Intn(256)
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

}

func (h *stringHandler) AddCharacter() {

}

func (h *stringHandler) RemoveCharacter() {

}
