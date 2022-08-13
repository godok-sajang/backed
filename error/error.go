package error

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound          = errors.New("Not Found")
	ErrNoRows            = errors.New("No Rows")
	ErrNicknameDuplicate = errors.New("nickname duplicate")
	ErrNilPtrDeref       = errors.New("Invalid memory address or nil pointer dereference")
)
