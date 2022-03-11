package repo

import (
	"backend/pkg/item"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

func NewClient1c(host, token string) *Client1c {
	c := &Client1c{host, token, http.Client{}}
	return c
}

type Client1c struct {
	Host, Token string
	hc          http.Client
}

type IClient1c interface {
	Products(limit int) (map[string]item.Item, error)
}

func (c *Client1c) Products(limit int) (map[string]item.Item, error) {
	r, err := http.NewRequest("GET", c.Host+"products", nil)
	if err != nil {
		log.Fatal().Err(err).Msg("create new request")
	}
	q := r.URL.Query()
	q.Add("limit", strconv.Itoa(limit))

	r.URL.RawQuery = q.Encode()

	r.Header.Set("Authorization", "Basic "+c.Token)

	resp, err := c.hc.Do(r)
	if err != nil {
		return nil, fmt.Errorf("do request %w", err)
	}

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
