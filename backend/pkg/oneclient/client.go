package oneclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewOneClient(host, token string) *ClientOnec {
	c := &ClientOnec{host, token, http.Client{
		Timeout: time.Second * 120, // for long requests
	}} // todo: config timeout

	return c
}

type ClientOnec struct {
	Host, Token string
	hc          http.Client
}

//http://10.51.0.10:19201/sm1/hs/ - host

func (c *ClientOnec) Request(method, path string, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, c.Host+path, body)
	if err != nil {
		return nil, fmt.Errorf("cant create new once request %w", err)
	}

	req.Header.Set("Authorization", "Basic "+c.Token)

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cant do request %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		defer resp.Body.Close()

		return nil, fmt.Errorf("error: %w : %s : %d",
			onecErrorFromJSON(resp.Body), resp.Status, resp.StatusCode)
	}

	return resp.Body, nil
}

func (c *ClientOnec) Proxy(w io.Writer, method, path string) error {
	req, err := http.NewRequest(method, c.Host+path, nil)
	if err != nil {
		return fmt.Errorf("cant create new once request %w", err)
	}

	req.Header.Set("Authorization", "Basic "+c.Token)

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

type onecError struct {
	Err              string `json:"error"`
	ModuleNameLine   string `json:"moduleNameLine"`
	Raw, Description string
}

func onecErrorFromJSON(r io.Reader) (err error) {
	var oneErr onecError
	err = json.NewDecoder(r).Decode(&oneErr)

	if err != nil {
		return fmt.Errorf("cannot decode error: %w", err)
	}

	return oneErr
}

func (e onecError) Error() string {
	return fmt.Sprintf("%s - %s - %s - %s",
		e.ModuleNameLine, e.Err, e.Raw, e.Description,
	)
}
