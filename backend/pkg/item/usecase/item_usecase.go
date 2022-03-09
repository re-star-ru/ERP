package usecase

import (
	"backend/pkg/item"
	"io"
)

type IItemUsecase interface {
	UploadFrom1cToS3drom() error
}

type Renderer interface {
	Render([]item.Item) io.Reader
}

type ItemUsecase struct {
	renders []Renderer
}

func NewItemUsecase() *ItemUsecase {
	return &ItemUsecase{
		renders: []Renderer{&DromRender{}},
	}
}
