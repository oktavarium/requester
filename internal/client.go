package requester

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

func fetchURL(chURL <-chan string, chOut chan<- string) {
	defer close(chOut)
	client := resty.New()
	client.SetTimeout(1 * time.Second)
	for url := range chURL {
		start := time.Now()
		resp, err := client.R().Get(url)
		duration := time.Since(start)

		if err != nil {
			chOut <- fmt.Sprintf("error on requesting %s. Error: %s", url, err)
			continue
		}

		chOut <- fmt.Sprintf("Requesting %s. Size: %d. Duration: %d", url, resp.Size(), duration)
	}
}
