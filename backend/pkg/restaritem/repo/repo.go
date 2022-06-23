package repo

import (
	"backend/ent"
	entrestaritem "backend/ent/restaritem"
	"backend/pkg/photo"
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

func (r *RestaritemRepo) AddPhoto(ctx context.Context, id int, photo photo.Photo) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("create tx error: %w", err)
	}

	u, err := tx.Restaritem.Query().Where(entrestaritem.IDEQ(id)).Only(ctx)
	if err != nil {
		return fmt.Errorf("get item error: %w", err)
	}

	if err = tx.Restaritem.
		Update().
		SetPhotos(append(u.Photos, photo)).
		Where(entrestaritem.IDEQ(id)).
		Exec(ctx); err != nil {
		return rollback(tx, fmt.Errorf("update item error: %w", err))
	}

	return tx.Commit()
}

// rollback calls to tx.Rollback and wraps the given error
// with the rollback error if occurred.
func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
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
