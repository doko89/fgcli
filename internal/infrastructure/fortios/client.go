package fortios

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	base   string
	apiKey string
	http   *http.Client
}

func newHTTP(insecure bool) *http.Client {
	tr, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return &http.Client{Timeout: 15 * time.Second}
	}
	tr = tr.Clone()
	if insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // opt-in via FG_INSECURE
	}
	return &http.Client{Timeout: 15 * time.Second, Transport: tr}
}

func NewClient(host, apiKey string, insecure bool) *Client {
	return &Client{base: strings.TrimSuffix(host, "/"), apiKey: apiKey, http: newHTTP(insecure)}
}

type apiResp struct {
	Results any `json:"results"`
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any) ([]byte, error) {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = bytes.NewReader(b)
	}
	if query == nil {
		query = url.Values{}
	}
	u := c.base + path + "?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, method, u, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt*attempt) * 500 * time.Millisecond):
			}
		}
		resp, err := c.http.Do(req.Clone(ctx))
		if err != nil {
			lastErr = err
			continue
		}
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
		resp.Body.Close()
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("fortios: %s %s -> http %d", method, path, resp.StatusCode)
			continue
		}
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("fortios: %s %s -> http %d: %s", method, path, resp.StatusCode, strings.TrimSpace(string(b)))
		}
		return b, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("fortios: request failed")
	}
	return nil, lastErr
}

func decodeResults(b []byte, out any) error {
	var raw struct {
		Results json.RawMessage `json:"results"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if len(raw.Results) == 0 {
		return fmt.Errorf("fortios: empty results")
	}
	return json.Unmarshal(raw.Results, out)
}

func decodeResultsAllowEmpty(b []byte, out any) error {
	if len(bytes.TrimSpace(b)) == 0 {
		return nil
	}
	var raw struct {
		Results json.RawMessage `json:"results"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if len(raw.Results) == 0 || string(raw.Results) == "null" {
		return nil
	}
	return json.Unmarshal(raw.Results, out)
}
