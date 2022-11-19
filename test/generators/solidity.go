package generators

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

func SmartContractGen() string {
	const ADDRESS_LENGTH = 20
	rand.Seed(time.Now().UnixNano())

	parts := make([]string, ADDRESS_LENGTH)
	for idx := 0; idx < len(parts); idx++ {
		parts[idx] = fmt.Sprintf("%x", rand.Intn(256))
	}
	return strings.Join(parts, "")
}

func OverflowedNumberAsStringGen(bitSize int) string {
	base := big.NewInt(2)
	exponent := big.NewInt(int64(bitSize))
	number := base.Exp(base, exponent, nil)
	return number.String()
}

func UnderflowedNumberAsStringGen(bitSize int) string {
	base := big.NewInt(2)
	exponent := big.NewInt(int64(bitSize))
	number := new(big.Int)
	number.Exp(base, exponent, nil)
	number.Add(number, big.NewInt(1))
	number.Neg(number)
	return number.String()
}

func UnsignedBigIntGen(bitSize int) *big.Int {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	base := big.NewInt(2)
	exponent := big.NewInt(int64(bitSize))
	max := new(big.Int)
	max.Exp(base, exponent, nil)
	max.Sub(max, big.NewInt(1))

	value := new(big.Int)
	value.Rand(rnd, max)
	return value
}

func SignedBigIntGen(bitSize int) *big.Int {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	base := big.NewInt(2)
	exponent := big.NewInt(int64(bitSize / 2))

	max := new(big.Int)
	max.Exp(base, exponent, nil)
	max.Sub(max, big.NewInt(1))

	min := new(big.Int)
	min.Exp(base, exponent, nil)
	min.Neg(min)

	randomRange := new(big.Int)
	randomRange.Add(randomRange, max)
	randomRange.Sub(randomRange, min)

	value := new(big.Int)
	value.Rand(rnd, randomRange)
	value.Add(value, min)

	return value
}
