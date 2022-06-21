package usecase

import (
	"backend/pkg/restaritem"
	pgrepo "backend/pkg/restaritem/repo"
	"context"
	"database/sql"
)

func NewRestaritemUsecase(db *sql.DB) *RestarItemUsecase {
	return &RestarItemUsecase{
		repo: pgrepo.New(db),
	}
}

type RestarItemUsecase struct {
	repo *pgrepo.Queries
}

func (r RestarItemUsecase) Create(ctx context.Context, restaritem *restaritem.RestarItem) (restaritem.RestarItem, error) {
	return r.repo.CreateRestaritem(ctx, pgrepo.CreateRestaritemParams{
		Name:     restaritem.Name,
		Onceguid: restaritem.OnecGUID,
	})
}

func (r RestarItemUsecase) GetAll() ([]restaritem.RestarItem, error) {
	return r.repo.GetAll()
}

func (r RestarItemUsecase) GetByID(id int) (*restaritem.RestarItem, error) {
	return r.repo.GetByID(id)
}
