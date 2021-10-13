package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/context/ctxhttp"
)

type Fetcher struct {
	client *http.Client
}

func NewFetcher(c *http.Client) *Fetcher {
	return &Fetcher{
		client: c,
	}
}

func (f *Fetcher) GetDOMContent(ctx context.Context, url string) ([]byte, error) {
	resp, err := ctxhttp.Get(ctx, f.client, url)
	if err != nil {
		return nil, fmt.Errorf("error while getting response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code is not 200 OK")
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %w", err)
	}

	return content, nil
}
