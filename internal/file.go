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

func getUrls(ch chan<- string, chErr chan<- error, filePath string) error {
	defer close(ch)
	defer close(chErr)
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error occured on opening file: %w", err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ch <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		chErr <- fmt.Errorf("error occured on reading file: %w", err)
	}

	return nil
}
