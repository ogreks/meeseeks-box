package service

import "errors"

var (
	ErrDataNotFound       = errors.New("data does not exist")
	ErrMissRequiredParams = errors.New("missing required parameters")
)
