package requester

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/oktavarium/requester/internal/semaphore"
	"golang.org/x/sync/errgroup"
)

func getResponseSize(r io.Reader) (int64, error) {
	b := new(bytes.Buffer)
	size, err := b.ReadFrom(r)
	if err != nil {
		return 0, fmt.Errorf("error occured on readig response: %w", err)
	}

	return size, nil
}

func fetchURLs(ctx context.Context,
	eg *errgroup.Group,
	bufferSize uint64,
	urlCh <-chan string) <-chan string {

	outCh := make(chan string, bufferSize)
	smphr := semaphore.NewSemaphore(bufferSize)

	t := &http.Transport{}
	t.MaxIdleConns = 1
	t.MaxConnsPerHost = 1
	t.MaxIdleConnsPerHost = 1
	t.IdleConnTimeout = 1 * time.Second

	client := &http.Client{
		Timeout:   1 * time.Second,
		Transport: t,
	}

	eg.Go(func() error {
		defer close(outCh)
		var wg sync.WaitGroup

		for URL := range urlCh {
			wg.Add(1)
			smphr.Acquire()
			go func(URL string) {
				defer wg.Done()
				defer smphr.Release()

				if _, err := url.ParseRequestURI(URL); err != nil {
					outCh <- fmt.Sprintf("error on parsing %s. Error: %s\n", URL, err)
					return
				}

				start := time.Now()
				resp, err := client.Get(URL)
				duration := time.Since(start)

				if err != nil {
					outCh <- fmt.Sprintf("error on requesting %s. Error: %s\n", URL, err)
					return
				}

				defer resp.Body.Close()
				size, err := getResponseSize(resp.Body)
				if err != nil {
					outCh <- fmt.Sprintf("error on requesting %s. Error: %s\n", URL, err)
					return
				}
				outCh <- fmt.Sprintf("Requesting %s. Size: %d. Duration: %s\n", URL, size, duration.String())
			}(URL)

			if ctx.Err() != nil {
				return nil
			}
		}
		wg.Wait()
		return nil
	})

	return outCh
}
