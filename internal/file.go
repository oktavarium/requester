package requester

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

func getURLs(ctx context.Context,
	eg *errgroup.Group,
	bufferSize uint64,
	filePath string) <-chan string {

	outCh := make(chan string, bufferSize)

	eg.Go(func() error {
		defer close(outCh)
		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error occured on opening file: %w", err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			if ctx.Err() != nil {
				return nil
			}
			outCh <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error occured on reading file: %w", err)
		}

		return nil
	})

	return outCh
}
