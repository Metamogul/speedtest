package speedtester

import (
	"bufio"
	"strings"
	"time"
)

func addTimestampToResult(result string) string {
	var newResult string

	lineScanner := bufio.NewScanner(strings.NewReader(result))
	for lineScanner.Scan() {
		currentLine := lineScanner.Text()

		if containsHeader(currentLine) {
			newResult += addTimestampToHeaderLine(currentLine) + "\n"
			continue
		}

		newResult += addTimestampToRecordLine(currentLine) + "\n"
	}

	return newResult
}

func containsHeader(header string) bool {
	return (strings.Contains(header, `"server name"`) &&
		strings.Contains(header, `"server id"`) &&
		strings.Contains(header, `"idle latency"`) &&
		strings.Contains(header, `"idle jitter"`) &&
		strings.Contains(header, `"packet loss"`) &&
		strings.Contains(header, `"download"`) &&
		strings.Contains(header, `"upload"`) &&
		strings.Contains(header, `"download bytes"`) &&
		strings.Contains(header, `"upload bytes"`) &&
		strings.Contains(header, `"share url"`) &&
		strings.Contains(header, `"download server count"`) &&
		strings.Contains(header, `"download latency"`) &&
		strings.Contains(header, `"download latency jitter"`) &&
		strings.Contains(header, `"download latency low"`) &&
		strings.Contains(header, `"download latency high"`) &&
		strings.Contains(header, `"upload latency"`) &&
		strings.Contains(header, `"upload latency jitter"`) &&
		strings.Contains(header, `"upload latency low"`) &&
		strings.Contains(header, `"upload latency high"`) &&
		strings.Contains(header, `"idle latency low"`) &&
		strings.Contains(header, `"idle latency high"`))
}

func addTimestampToHeaderLine(header string) string {
	return `"timestamp",` + header
}

func addTimestampToRecordLine(record string) string {

	return `"` + time.Now().String() + `",` + record
}
