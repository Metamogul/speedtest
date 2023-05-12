package resultfile

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/Clever/csvlint"
)

type ResultFileError string

var HeaderError ResultFileError = "No valid header found in result file."

type MalformedError struct {
	ResultFileError
	ValidationErrors []csvlint.CSVError
}

func (r ResultFileError) Error() string {
	return string(r)
}

type ResultFile struct {
	os.File
	WasEmpty bool
}

func OpenResultFile(filePath string) (*ResultFile, error) {
	fileInfo, err := os.Stat(filePath)

	if err != nil {
		return nil, err
	}

	wasEmpty := fileInfo.Size() == 0

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		return nil, err
	}

	resultFile := ResultFile{File: *file, WasEmpty: wasEmpty}

	if resultFile.WasEmpty {
		return &resultFile, nil
	}

	if _, err := resultFile.isValidCSV(); err != nil {
		return nil, err
	}

	if !resultFile.containsHeaderLine() {
		return nil, HeaderError
	}

	return &resultFile, nil
}

func (c *ResultFile) isValidCSV() (bool, error) {
	c.Seek(0, io.SeekStart)

	csvErrors, isValid, err := csvlint.Validate(&c.File, ',', false)

	if len(csvErrors) > 0 {
		return false, MalformedError{
			ResultFileError:  "No valid result file.",
			ValidationErrors: csvErrors,
		}
	}

	if err != nil {
		return false, err
	}

	c.Seek(0, io.SeekStart)
	return isValid, nil
}

func (c *ResultFile) containsHeaderLine() bool {
	lineScanner := bufio.NewScanner(c)

	for lineScanner.Scan() {
		currentLine := lineScanner.Text()

		if ContainsHeader(currentLine) {
			c.Seek(0, io.SeekStart)
			return true
		}
	}

	return false
}

func ContainsHeader(header string) bool {
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
