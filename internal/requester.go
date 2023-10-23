package requester

import (
	"fmt"
)

func Run() error {
	flagsConfig, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error occured on parsing config: %w", err)
	}

	if !fileExists(flagsConfig.FilePath) {
		return fmt.Errorf("file does not exist")
	}

	urlCh := make(chan string, 1000)
	errCh := make(chan error)
	resultCh := make(chan string, 1000)

	go getUrls(urlCh, errCh, flagsConfig.FilePath)
	go fetchURL(urlCh, resultCh)
	go print(resultCh)

	if err := <-errCh; err != nil {
		fmt.Println(err)
	}

	return nil
}
