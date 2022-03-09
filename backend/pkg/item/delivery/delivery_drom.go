package delivery

import (
	"backend/pkg/item"
	"io"
	"strings"
)

func DromRender(price []item.Item) io.Reader {
	return strings.NewReader("string")
}
