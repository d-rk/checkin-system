package app

import "errors"

var ErrNotFound = errors.New("not found")

var ErrConflict = errors.New("conflict")

var ErrInternal = errors.New("internal error")
