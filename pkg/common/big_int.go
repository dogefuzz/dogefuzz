package common

import (
	"errors"
	"math/big"
	"regexp"
)

var ErrInvalidHexadecimal = errors.New("the provided string is not a valid hexadecimal")

func ConvertHexadecimalToInt(hexadecimal string) (*big.Int, error) {
	pattern := regexp.MustCompile("0x.*")
	if pattern.MatchString(hexadecimal) {
		hexadecimal = hexadecimal[2:]
	}
	val := new(big.Int)
	if _, ok := val.SetString(hexadecimal, 16); !ok {
		return nil, ErrInvalidHexadecimal
	}
	return val, nil
}

func MustConvertHexadecimalToInt(hexadecimal string) *big.Int {
	val, err := ConvertHexadecimalToInt(hexadecimal)
	if err != nil {
		panic(err)
	}
	return val
}
