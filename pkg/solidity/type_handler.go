package solidity

import (
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type TypeHandler interface {
	GetValue() interface{}
	Serialize() string
	Deserialize(value string) error
	Generate() // Add Random provider to be mocked in tests
	GetMutators() []func()
}

func GetTypeHandler[T any](typ abi.Type) (TypeHandler, error) {
	switch typ.GetType() {
	case common.BoolT:
		return NewBoolHandler(), nil
	case common.Uint8T:
		return NewUint8Handler(), nil
	case common.Uint16T:
		return NewUint16Handler(), nil
	case common.Uint32T:
		return NewUint32Handler(), nil
	case common.Uint64T:
		return NewUint64Handler(), nil
	case common.Int8T:
		return NewInt8Handler(), nil
	case common.Int16T:
		return NewInt16Handler(), nil
	case common.Int32T:
		return NewInt32Handler(), nil
	case common.Int64T:
		return NewInt64Handler(), nil
	case common.BigIntT:
		if typ.T == abi.UintTy {
			return NewUnsignedBigIntHandler(typ.Size), nil
		}
		return NewSignedBigIntHandler(typ.Size), nil
	case common.StringT:
		return NewStringHandler(), nil
	}
	return nil, ErrNotImplementedType(typ.GetType())
}
