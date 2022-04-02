package repo

import (
	"backend/pkg/item"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewRepoOnec(host, token string) *ClientOnec {
	c := &ClientOnec{host, token, http.Client{}}
	return c
}

type ClientOnec struct {
	Host, Token string
	hc          http.Client
}

func (c *ClientOnec) Items(offset, limit int) (map[string]item.Item, error) {
	r, err := http.NewRequest("GET", c.Host+"products/batch", nil)
	if err != nil {
		return nil, fmt.Errorf("cant create new products barch request %w", err)
	}
	r.Header.Set("Authorization", "Basic "+c.Token)

	resp, err := c.hc.Do(r)
	if err != nil {
		return nil, fmt.Errorf("cant do request %w", err)
	}
	defer r.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("cant read body: %w", err)
		}

		return nil, fmt.Errorf("error: %s : %s : %d", body, resp.Status, resp.StatusCode)
	}

	m := map[string]item.Item{}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("cannot decode body to products %w", err)
	}

	return m, nil
}
