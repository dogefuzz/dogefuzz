package solidity

import (
	"fmt"
	"reflect"
)

type Empty struct{}

var ErrNotImplementedType = func(typ reflect.Type) error {
	return fmt.Errorf("the desired type's mapping was not implemented yet: %v", typ)
}
