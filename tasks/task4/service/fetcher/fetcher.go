package fetcher

import (
	"fmt"
	"io"
	"net/http"
)

type Fetcher struct {
	client http.Client
}

func NewFetcher(c http.Client) *Fetcher {
	return &Fetcher{
		client: c,
	}
}

func (f *Fetcher) GetDOMContent(url string) ([]byte, error) {
	resp, err := f.client.Get(url)
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
