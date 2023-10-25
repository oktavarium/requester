package requester

import (
	"context"
	"fmt"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func Run() error {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return fmt.Errorf("error occured on getting syscall limits: %w", err)
	}

	flagsConfig, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error occured on parsing config: %w", err)
	}

	eg, ctx := errgroup.WithContext(context.Background())

	urlCh := getURLs(ctx, eg, limit.Cur, flagsConfig.FilePath)
	resultCh := fetchURLs(ctx, eg, limit.Cur, urlCh)
	writer(ctx, eg, resultCh)

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
