package requester

import (
	"bufio"
	"os"
)

func print(results chan string) {
	writer := bufio.NewWriter(os.Stdout)
	for result := range results {
		writer.WriteString(result)
	}
}
