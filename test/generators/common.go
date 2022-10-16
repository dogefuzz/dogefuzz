package generators

import (
	"fmt"
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
