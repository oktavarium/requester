package requester

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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

func fetchURL(chURL <-chan string, chOut chan<- string) {
	defer close(chOut)
	t := &http.Transport{}
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
	var wg sync.WaitGroup
	for url := range chURL {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			start := time.Now()
			resp, err := client.Get(url)
			duration := time.Since(start)

			if err != nil {
				chOut <- fmt.Sprintf("error on requesting %s. Error: %s\n", url, err)
				return
			}

			defer resp.Body.Close()
			size, err := getResponseSize(resp.Body)
			if err != nil {
				chOut <- fmt.Sprintf("error on requesting %s. Error: %s\n", url, err)
				return
			}

			chOut <- fmt.Sprintf("Requesting %s. Size: %d. Duration: %d ms\n", url, size, duration.Milliseconds())
		}(url)

	}
	wg.Wait()
}
