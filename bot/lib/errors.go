package lib

import "errors"

var (
	ErrNilHandler = errors.New("nil handler")
	ErrNilLogger  = errors.New("nil logger")
	ErrNilService = errors.New("nil service")
)
