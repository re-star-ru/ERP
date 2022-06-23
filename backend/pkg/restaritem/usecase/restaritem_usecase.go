package usecase

import (
	"backend/pkg/photo"
	"backend/pkg/restaritem"
	"context"
)

type IRestaritemRepo interface {
	Create(ctx context.Context, item restaritem.RestarItem) (*restaritem.RestarItem, error)
	List(context.Context) ([]*restaritem.RestarItem, error)
	Get(ctx context.Context, id int) (*restaritem.RestarItem, error)
	AddPhoto(ctx context.Context, id int, photo photo.Photo) error
}

func NewRestaritemUsecase(repo IRestaritemRepo) *RestarItemUsecase {
	return &RestarItemUsecase{
		repo: repo,
	}
}

type RestarItemUsecase struct {
	repo IRestaritemRepo
}

func (r RestarItemUsecase) AddPhoto(ctx context.Context, id int, photo photo.Photo) (*restaritem.RestarItem, error) {
	if err := r.repo.AddPhoto(ctx, id, photo); err != nil {
		return nil, err
	}

	return r.repo.Get(ctx, id)
}

func (r RestarItemUsecase) Create(ctx context.Context, restaritem restaritem.RestarItem) (*restaritem.RestarItem, error) {
	return r.repo.Create(ctx, restaritem)
}

func (r RestarItemUsecase) GetAll(ctx context.Context) ([]*restaritem.RestarItem, error) {
	return r.repo.List(ctx)
}

func (r RestarItemUsecase) GetByID(ctx context.Context, id int) (*restaritem.RestarItem, error) {
	return r.repo.Get(ctx, id)
}
