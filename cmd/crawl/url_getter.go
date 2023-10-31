package crawl

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type URLGetter interface {
	Get(string) (io.ReadCloser, error)
}

type SingleRequestGetter struct {
	lock       sync.Mutex
	httpClient http.Client
}

func NewHttpGetter(timeout time.Duration) URLGetter {
	return &SingleRequestGetter{
		httpClient: http.Client{Timeout: timeout},
	}
}

// allows a single request at a time to not overwhelm the server
func (h *SingleRequestGetter) Get(url string) (io.ReadCloser, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") || response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("url (%v) is not parseable html or an error page", url)
	}
	return response.Body, nil
}
