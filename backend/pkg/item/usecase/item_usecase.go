package usecase

import (
	"backend/pkg/item"
)

type ItemRepo interface {
	// Items will return batch of items from repo
	ItemsWithOffcetLimit(offset, limit int) (map[string]item.Item, error)
	Items() ([]item.Item, error)

	TextSearch(string) ([]interface{}, error)
}

type ItemUsecase struct {
	repo ItemRepo
}

func NewItemUsecase(repo ItemRepo) *ItemUsecase {
	return &ItemUsecase{
		repo: repo,
	}
}

func (iu *ItemUsecase) Search(s string) ([]interface{}, error) {
	return iu.repo.TextSearch(s)
}
