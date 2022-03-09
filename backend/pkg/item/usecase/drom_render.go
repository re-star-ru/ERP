package usecase

import (
	"backend/pkg/item"
	"io"
	"strings"
)

type DromRender struct {
	path string // name
}

func (r *DromRender) Render(price []item.Item) io.Reader {
	return strings.NewReader("string")
}
