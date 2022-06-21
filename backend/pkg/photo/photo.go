package photo

import "fmt"

const (
	original = iota + 1
	thumbnail
	small
	medium
	large
)

type Photo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	Path string `json:"path"`

	// Фото в разных размерах
	Sizes [5]string `json:"sizes"`

	OriginalPath string `json:"originalPath"`
}

var ErrSizeNotSupported = fmt.Errorf("size not supported")

func (p Photo) GetPhoto(size uint) (string, error) {
	if size > 5 {
		return "", fmt.Errorf("%w: %d", ErrSizeNotSupported, size)
	}

	return p.Sizes[size], nil
}
