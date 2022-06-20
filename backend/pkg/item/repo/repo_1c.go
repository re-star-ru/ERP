package repo

import (
	"backend/pkg/item"
	"encoding/json"
	"fmt"
	"io"
)

type Requester interface {
	Request(method, path string, body io.Reader) (io.ReadCloser, error)
}

func NewRepoOnec(r Requester) *ClientOnec {
	c := &ClientOnec{r}

	return c
}

type ClientOnec struct {
	r Requester
}

func (c *ClientOnec) ItemsWithOffcetLimit(offset, limit int) (map[string]item.Item, error) {
	body, err := c.r.Request("GET", "products/batch", nil)
	if err != nil {
		return nil, fmt.Errorf("cant get products %w", err)
	}

	defer body.Close()

	m := map[string]item.Item{}
	if err = json.NewDecoder(body).Decode(&m); err != nil {
		return nil, fmt.Errorf("cannot decode body to products %w", err)
	}

	return m, nil
}

func (c *ClientOnec) Items() ([]item.Item, error) {
	body, err := c.r.Request("GET", "products/batch", nil)
	if err != nil {
		return nil, fmt.Errorf("cant get products %w", err)
	}

	defer body.Close()

	var m []item.Item
	if err = json.NewDecoder(body).Decode(&m); err != nil {
		return nil, fmt.Errorf("cannot decode body to products: %w", err)
	}

	return m, nil
}

func (c *ClientOnec) TextSearch(s string) ([]interface{}, error) {
	body, err := c.r.Request("GET", "products/text-search/"+s, nil)
	if err != nil {
		return nil, fmt.Errorf("cant get products %w", err)
	}

	defer body.Close()

	var m []interface{}
	if err = json.NewDecoder(body).Decode(&m); err != nil {
		return nil, fmt.Errorf("cannot decode body to products %w", err)
	}

	return m, nil
}
