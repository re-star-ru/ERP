package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/rs/zerolog/log"
)

type Client1c struct {
	Host, Token string
	hc          http.Client
}

func NewClient1c(host, token string) *Client1c {
	c := &Client1c{host, token, http.Client{}}
	log.Debug().Msgf("%+v", c)

	return c
}

type Product struct {
	ID   string `json:"id"`
	SKU  string `json:"sku"`
	Name string `json:"name"`
	Char string `json:"char"`

	Images []string `json:"images"`

	Amount        int `json:"amount"`
	Price         int `json:"price"`
	DiscountPrice int `json:"discountPrice"`
	Discount      int `json:"discount"`
}

func (c *Client1c) Products() ([]Product, error) {
	r, err := http.NewRequest("GET", c.Host+"products/", nil)
	if err != nil {
		log.Fatal().Err(err).Msg("create new request")
	}

	r.Header.Set("Authorization", "Basic "+c.Token)

	resp, err := c.hc.Do(r)
	if err != nil {
		return nil, fmt.Errorf("do request %w", err)
	}

	if !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("cant read body: %w", err)
		}

		return nil, fmt.Errorf("error: %s : %s : %d", body, resp.Status, resp.StatusCode)
	}

	log.Debug().Msgf("size: %d", resp.ContentLength)

	m := map[string]Product{}

	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("cannot decode body to products %w", err)
	}

	ps := make([]Product, len(m))

	i := 0
	for _, v := range m {
		ps[i] = v
		i++
	}

	log.Debug().Msg("Status: " + resp.Status)

	return ps, nil
}
