package solidity

import (
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func GetTypeHandler(typ abi.Type) (interfaces.TypeHandler, error) {
	if typ.Elem != nil {
		switch typ.GetType() {
		case common.SliceT(typ.Elem.GetType()):
			return NewSliceHandler(typ)
		case common.ArrayT(typ.Size, typ.Elem.GetType()):
			return NewArrayHandler(typ.Size, typ)
		default:
			return nil, ErrNotImplementedType(typ.GetType())
		}
	}

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
	case common.AddressT:
		return NewAddressHandler(), nil
	case common.FixedBytesT(typ.Size):
		return NewFixedBytesHandler(typ)
	case common.BytesT:
		return NewBytesHandler(typ)
	default:
		return nil, ErrNotImplementedType(typ.GetType())
	}
}
