package usecase

import (
	"backend/pkg/restaritem"
	"context"
)

type IRestaritemRepo interface {
	Create(ctx context.Context, item restaritem.RestarItem) (restaritem.RestarItem, error)
	List() ([]restaritem.RestarItem, error)
	GetByID(int) (restaritem.RestarItem, error)
}

func NewRestaritemUsecase(repo IRestaritemRepo) *RestarItemUsecase {
	return &RestarItemUsecase{
		repo: repo,
	}
}

type RestarItemUsecase struct {
	repo IRestaritemRepo
}

func (r RestarItemUsecase) Create(ctx context.Context, restaritem restaritem.RestarItem) (restaritem.RestarItem, error) {
	return r.repo.Create(ctx, restaritem)
}

func (r RestarItemUsecase) GetAll(ctx context.Context) ([]restaritem.RestarItem, error) {
	return r.repo.List()
}

func (r RestarItemUsecase) GetByID(ctx context.Context, id int) (restaritem.RestarItem, error) {
	return r.repo.GetByID(id)
}
