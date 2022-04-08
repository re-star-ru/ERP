package repo

import (
	"backend/pkg/item"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ErrorOnec struct {
	Err                                       string `json:"error"`
	ModuleName, Raw, Line, Description, Cause string
}

func (e *ErrorOnec) fromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(e)
}

func (e ErrorOnec) Error() string {
	return fmt.Sprintf("%s:%s - %s - %s - %s - %s",
		e.Line, e.Cause, e.Err, e.ModuleName, e.Raw, e.Description,
	)
}

func NewRepoOnec(host, token string) *ClientOnec {
	c := &ClientOnec{host, token, http.Client{Timeout: time.Second * 60}} // todo: config timeout
	return c
}

type ClientOnec struct {
	Host, Token string
	hc          http.Client
}

func (c *ClientOnec) ItemsWithOffcetLimit(offset, limit int) (map[string]item.Item, error) {
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

func (c *ClientOnec) Items() ([]item.Item, error) {
	r, err := http.NewRequest("GET", c.Host+"products/batch", nil)
	if err != nil {
		return nil, fmt.Errorf("cant create new products barch request %w", err)
	}
	r.Header.Set("Authorization", "Basic "+c.Token)

	resp, err := c.hc.Do(r)
	if err != nil {
		return nil, fmt.Errorf("cant do request %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("cant read body: %w", err)
		}

		return nil, fmt.Errorf("error: %s : %s : %d", body, resp.Status, resp.StatusCode)
	}

	m := []item.Item{}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("cannot decode body to products: %w", err)
	}

	return m, nil
}

func (c *ClientOnec) TextSearch(s string) ([]interface{}, error) {
	w, err := c.newRequest("GET", c.Host+"products/text-search/"+s)
	if err != nil {
		return nil, err
	}

	r, err := c.hc.Do(w)
	if err != nil {
		return nil, fmt.Errorf("cant do request %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode >= 400 {
		oe := ErrorOnec{}
		if err = oe.fromJSON(r.Body); err != nil {
			return nil, fmt.Errorf("cant read onecerror: %w", err)
		}

		return nil, fmt.Errorf("error - %s - %d - %w", r.Status, r.StatusCode, oe)
	}

	var m []interface{}
	if err = json.NewDecoder(r.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("cannot decode body to products %w", err)
	}

	return m, nil
}

func (c *ClientOnec) newRequest(method, path string) (*http.Request, error) {
	w, err := http.NewRequest(method, path, nil)
	if err != nil {
		return nil, fmt.Errorf("cant create new products barch request %w", err)
	}
	w.Header.Set("Authorization", "Basic "+c.Token)

	return w, err
}

func (c *ClientOnec) Proxy(w io.Writer, method, path string) error {
	req, err := c.newRequest(method, path)
	if err != nil {
		return err
	}

	r, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("cant do request %w", err)
	}
	defer r.Body.Close()

	if _, err = io.Copy(w, r.Body); err != nil {
		return err
	}

	return nil
}
