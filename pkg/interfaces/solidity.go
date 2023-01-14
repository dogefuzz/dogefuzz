package interfaces

import "github.com/dogefuzz/dogefuzz/pkg/common"

type TypeHandler interface {
	GetValue() interface{}
	SetValue(value interface{})
	GetType() common.TypeIdentifier
	Serialize() string
	Deserialize(value string) error
	Generate() // Add Random provider to be mocked in tests
	GetMutators() []func()
}
