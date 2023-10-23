package requester

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
)

func Run() error {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		fmt.Println("Getrlimit:" + err.Error())
	}
	fmt.Printf("%v file descriptors out of a maximum of %v available\n", limit.Cur, limit.Max)

	flagsConfig, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error occured on parsing config: %w", err)
	}

	if !fileExists(flagsConfig.FilePath) {
		return fmt.Errorf("file does not exist")
	}

	writer := bufio.NewWriter(os.Stdout)
	urlCh := make(chan string, 1000)
	errCh := make(chan error)
	resultCh := make(chan string, 1000)

	go getUrls(urlCh, errCh, flagsConfig.FilePath)
	go fetchURLs(urlCh, resultCh)

	for result := range resultCh {
		writer.WriteString(result)
	}
	writer.Flush()

	if err := <-errCh; err != nil {
		return err
	}

	return nil
}
