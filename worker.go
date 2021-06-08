package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

func (w *Work) Query(ctx context.Context) result {
	// prepare URL
	endpoint := "query_range"
	values := url.Values{}
	for k, v := range w.query.params {
		if k == "endpoint" {
			endpoint = v
		}
		values.Add(k, v)
	}
	url := fmt.Sprintf("%s/api/v1/%s?%s", w.serverURL, endpoint, values.Encode())
	return callPromScale(ctx, url, w.clientTimeout)
}

func (w *Work) Health(ctx context.Context) result {
	url := fmt.Sprintf("%s/%s", w.serverURL, "healthz")
	return callPromScale(ctx, url, w.clientTimeout)
}

func callPromScale(ctx context.Context, url string, timeoutMillis int) result {
	// log.Printf("req url: %v", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return result{
			isError: true,
			err:     errors.WithStack(err),
		}
	}

	c := http.DefaultClient
	c.Timeout = time.Duration(timeoutMillis * int(time.Millisecond))

	r := result{st: time.Now()}
	resp, err := c.Do(req)
	if err != nil {
		r.isError = true
		r.err = errors.WithStack(err)
		return r
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r.isError = true
		r.err = errors.New("non ok response from Promscale")
		r.statusCode = resp.StatusCode
		return r
	}
	r.et = time.Now()

	return r
}
