package main

import (
	"io"
)

type Provider interface {
	Store(path string, reader io.Reader) (string, error)
}
