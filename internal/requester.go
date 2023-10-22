package requester

import "fmt"

func Run() error {
	flagsConfig, err := loadConfig()
	if err != nil {
		return fmt.Errorf("error occured on parsing config: %w", err)
	}

	if !fileExists(flagsConfig.FilePath) {
		return fmt.Errorf("file does not exist")
	}

}
