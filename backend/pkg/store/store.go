package store

import "io"

type Storer interface {
	Store(string, io.Reader) (string, error)
}
