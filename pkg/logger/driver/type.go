package driver

import "io"

type Driver interface {
	io.Writer
}

type Option func(Driver)
