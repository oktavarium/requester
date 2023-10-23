package requester

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func getUrls(urlCh chan<- string, errCh chan<- error, filePath string) {
	defer close(urlCh)
	defer close(errCh)
	f, err := os.Open(filePath)
	if err != nil {
		errCh <- fmt.Errorf("error occured on opening file: %w", err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		urlCh <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		errCh <- fmt.Errorf("error occured on reading file: %w", err)
		return
	}
}
