package app

import "errors"

var NotFoundErr = errors.New("not found")

var ConflictErr = errors.New("conflict")

var InternalErr = errors.New("internal error")
