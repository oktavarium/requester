package requester

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

func writer(ctx context.Context, eg *errgroup.Group, ch <-chan string) {
	writer := bufio.NewWriter(os.Stdout)

	eg.Go(func() error {
		defer writer.Flush()
		for v := range ch {
			select {
			default:
				_, err := writer.WriteString(v)
				if err != nil {
					return fmt.Errorf("error occured on writing to stdout: %w", err)
				}
			case <-ctx.Done():
				return nil
			}
		}
		return nil
	})
}
