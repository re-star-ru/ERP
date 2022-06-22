package repo

import (
	"backend/ent"
	entrestaritem "backend/ent/restaritem"
	"backend/pkg/restaritem"
	"context"
	"fmt"
)

func NewRestaritemRepo(client *ent.Client) *RestaritemRepo {
	return &RestaritemRepo{client: client}
}

type RestaritemRepo struct {
	client *ent.Client
}

func (r RestaritemRepo) Create(ctx context.Context, item restaritem.RestarItem) (*restaritem.RestarItem, error) {
	rit, err := r.client.Restaritem.Create().
		SetName(item.Name).
		SetOnecGUID(item.OnecGUID).
		SetSku(item.Sku).
		SetItemGUID(item.ItemGUID).
		SetCharGUID(item.CharGUID).
		Save(ctx)

	if ent.IsValidationError(err) {
		return nil, fmt.Errorf("%w: %v", restaritem.ErrValidation, err)
	}

	if err != nil {
		return nil, fmt.Errorf("create item error: %w", err)
	}

	return rit, nil
}

func (r *RestaritemRepo) List(ctx context.Context) ([]*restaritem.RestarItem, error) {
	us, err := r.client.Restaritem.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	return us, nil
}

func (r *RestaritemRepo) Get(ctx context.Context, id int) (*restaritem.RestarItem, error) {

	u, err := r.client.Restaritem.Query().Where(entrestaritem.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, err
	}

	return u, nil
}
