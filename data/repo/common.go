package repo

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/data"
)

type Env interface {
	DbConnection() data.Connection
}

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)
