package renderer

import (
	"backend/pkg/item"
	"bytes"
	"io"
)

func DromRender([]item.Item) io.Reader {
	return bytes.NewBuffer([]byte{})
}
