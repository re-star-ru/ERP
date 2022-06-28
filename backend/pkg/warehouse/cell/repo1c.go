package cell

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Requester interface {
	Request(method, path string, body io.Reader) (io.ReadCloser, error)
}

func NewRepoOnec(r Requester) *RepoOnec {
	c := &RepoOnec{r}

	return c
}

type RepoOnec struct {
	r Requester
}

func (repo *RepoOnec) GetCellByID(cellID string) (*Cell, error) {
	body, err := repo.r.Request(http.MethodGet, "warehouse/cell/"+cellID, nil)
	if err != nil {
		return nil, fmt.Errorf("cant get cell by cellID %w", err)
	}

	defer body.Close()

	cell := new(Cell)

	if err = json.NewDecoder(body).Decode(cell); err != nil {
		return nil, fmt.Errorf("cant decode cell %w", err)
	}

	cell.ID = cellID

	return cell, nil
}
