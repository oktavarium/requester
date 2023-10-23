package requester

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func getResponseSize(r io.Reader) (int64, error) {
	b := new(bytes.Buffer)
	size, err := b.ReadFrom(r)
	if err != nil {
		return 0, fmt.Errorf("error occured on readig response: %w", err)
	}

	return size, nil
}

func fetchURLs(chURL <-chan string, chOut chan<- string) {
	defer close(chOut)
	t := &http.Transport{}
	t.MaxIdleConns = 1
	t.MaxConnsPerHost = 1
	t.MaxIdleConnsPerHost = 1
	t.IdleConnTimeout = 1 * time.Millisecond

	client := &http.Client{
		Timeout:   1 * time.Second,
		Transport: t,
	}
	var wg sync.WaitGroup
	for URL := range chURL {
		wg.Add(1)
		go func(URL string) {
			defer wg.Done()

			if _, err := url.ParseRequestURI(URL); err != nil {
				chOut <- fmt.Sprintf("error on parsing %s. Error: %s\n", URL, err)
				return
			}

			start := time.Now()
			resp, err := client.Get(URL)
			duration := time.Since(start)

			if err != nil {
				chOut <- fmt.Sprintf("error on requesting %s. Error: %s\n", URL, err)
				return
			}

			defer resp.Body.Close()
			size, err := getResponseSize(resp.Body)
			if err != nil {
				chOut <- fmt.Sprintf("error on requesting %s. Error: %s\n", URL, err)
				return
			}

			chOut <- fmt.Sprintf("Requesting %s. Size: %d. Duration: %d ms\n", URL, size, duration.Milliseconds())
		}(URL)

	}
	wg.Wait()
}
