package repo

import (
	"backend/pkg/ent"
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

func (r RestaritemRepo) Create(ctx context.Context, item restaritem.RestarItem) (restaritem.RestarItem, error) {
	rit, err := r.client.Restaritem.Create().Save(ctx)

	if ent.IsValidationError(err) {
		return restaritem.RestarItem{}, fmt.Errorf("%w: %v", restaritem.ErrValidation, err)
	}

	if err != nil {
		return restaritem.RestarItem{}, fmt.Errorf("create item error: %w", err)
	}

	return toRestaritem(rit), nil
}

func toRestaritem(rit *ent.Restaritem) restaritem.RestarItem {
	ritem := restaritem.RestarItem{
		ID:          rit.ID,
		OnecGUID:    rit.OnecGUID,
		Name:        rit.Name,
		SKU:         rit.Sku,
		ItemGUID:    rit.ItemGUID,
		CharGUID:    rit.CharGUID,
		Description: rit.Description,
		Inspector:   rit.Inspector,
		Inspection:  rit.Inspection,
		Photos:      rit.Photos,
	}

	for _, w := range rit.Works {
		ritem.Works = append(ritem.Works, restaritem.Work{
			WorkGUID:     w.WorkGUID,
			EmployeeGUID: w.EmployeeGUID,
			Price:        w.Price,
		})
	}

	return ritem
}

func (r RestaritemRepo) List() ([]restaritem.RestarItem, error) {
	//TODO implement me
	panic("implement me")
}

func (r RestaritemRepo) GetByID(i int) (restaritem.RestarItem, error) {
	//TODO implement me
	panic("implement me")
}
