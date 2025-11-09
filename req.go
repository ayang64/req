package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"
)

func Decode[T any](r io.Reader) (T, error) {
	var v, zero T
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		return zero, err
	}
	return v, nil
}

func Every[T interface{ Page() (int, int) }](ctx context.Context, m string, u string, h http.Header, v url.Values, b io.Reader, s map[int]struct{}) ([]T, error) {
	var a []T
	for v, err := range All[T](ctx, m, u, h, v, b, s) {
		if err != nil {
			return nil, err
		}
		a = append(a, v)
	}
	return a, nil
}

func All[T interface{ Page() (int, int) }](ctx context.Context, m string, u string, h http.Header, v url.Values, b io.Reader, s map[int]struct{}) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for cur, limit := 0, 0; ; {
			pl, err := Send[T](ctx, m, u, h, v, b, s)
			if !yield(pl, err) {
				break
			}
			if err != nil {
				break
			}
			cur, limit = pl.Page()
			if cur >= limit {
				break
			}
		}
	}
}

func Post[T any](ctx context.Context, u string, b io.Reader) (T, error) {
	return Send[T](ctx, http.MethodPost, u, nil, nil, b, nil)
}

func Get[T any](ctx context.Context, u string) (T, error) {
	return Send[T](ctx, http.MethodGet, u, nil, nil, nil, nil)
}

func Send[T any](ctx context.Context, m string, u string, h http.Header, v url.Values, b io.Reader, s map[int]struct{}) (T, error) {
	var zero T
	s = func() map[int]struct{} {
		if s == nil {
			return map[int]struct{}{http.StatusOK: {}}
		}
		return s
	}()

	if h == nil {
		h = http.Header{}
	}
	h["Content-Type"] = []string{"application/json"}

	req, err := http.NewRequestWithContext(ctx, m, u, b)
	if err != nil {
		return zero, err
	}

	req.Header = h.Clone()
	req.URL.RawQuery = v.Encode()
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	if _, ok := s[resp.StatusCode]; !ok {
		return zero, fmt.Errorf("client received %q (%d) response code", http.StatusText(resp.StatusCode), resp.StatusCode)
	}
	return Decode[T](resp.Body)
}
