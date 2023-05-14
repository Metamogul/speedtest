package speedtester

import (
	"bufio"
	"strings"
	"time"

	"github.com/Metamogul/speedtest-series/resultfile"
)

func addTimestampToResult(result string) string {
	var newResult string

	lineScanner := bufio.NewScanner(strings.NewReader(result))
	for lineScanner.Scan() {
		currentLine := lineScanner.Text()

		if resultfile.ContainsHeader(currentLine) {
			newResult += addTimestampToHeaderLine(currentLine) + "\n"
			continue
		}

		newResult += addTimestampToRecordLine(currentLine) + "\n"
	}

	return newResult
}

func addTimestampToHeaderLine(header string) string {
	return `"timestamp",` + header
}

func addTimestampToRecordLine(record string) string {

	return `"` + time.Now().String() + `",` + record
}
