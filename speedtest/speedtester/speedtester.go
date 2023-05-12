package speedtester

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const (
	testIntervalMinutesDefault = 5
	testIntervalMinutesMin     = 1
	testDurationHoursDefault   = 6
)

type SpeedTester struct {
	filePath            string
	file                *os.File
	testIntervalMinutes int
	testDurationHours   int
	ticker              *time.Ticker
	stop                chan bool
}

func NewSpeedTester(filePath string, testIntervalMinutes int, testDurationHours int) *SpeedTester {
	tester := &SpeedTester{
		filePath:            filePath,
		testIntervalMinutes: testIntervalMinutesDefault,
		testDurationHours:   testDurationHoursDefault,
		stop:                make(chan bool, 1),
	}

	if testIntervalMinutes != 0 {
		tester.testIntervalMinutes = testIntervalMinutes
	}

	if tester.testIntervalMinutes < testIntervalMinutesMin {
		tester.testIntervalMinutes = testIntervalMinutesMin
	}

	if testDurationHours != 0 {
		tester.testDurationHours = testDurationHours
	}

	return tester
}

func (t *SpeedTester) Initialize() {
	var err error
	t.file, err = os.OpenFile(t.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Could no open file for writing.")
	}
	log.Printf("Writing all test results to %s\n", t.filePath)
}

func (t *SpeedTester) RunAsync() chan bool {
	var done chan bool = make(chan bool, 1)
	log.Println("Starting test run ...")
	go t.run(done)
	log.Println("Started.")
	return done
}

func (t *SpeedTester) ScheduleStopAsync() {
	go t.scheduleStop()
}

func (t *SpeedTester) run(done chan bool) {
	t.performSingleTest(true)

	t.ticker = time.NewTicker(time.Duration(t.testIntervalMinutes) * time.Minute)
	for {
		select {
		case <-t.ticker.C:
			t.performSingleTest(false)
		case <-t.stop:
			log.Println("Test terminated as requested.")
			done <- true
			return
		}
	}
}

func (t *SpeedTester) performSingleTest(withHeader bool) {
	const path = "/opt/homebrew/bin/speedtest"
	const formatArg = "--format=csv"
	const headerArg = "--output-header"

	log.Println("Performing test ...")

	var result []byte
	var err error
	if withHeader {
		result, err = execCommand(path, formatArg, headerArg).Output()
	} else {
		result, err = execCommand(path, formatArg).Output()
	}

	if err != nil {
		log.Print(err.Error())
		return
	}

	if result != nil {
		t.appendResultToFile(string(result))
	}

	log.Println("Finished performing test, output written to", t.filePath)
}

func execCommand(path string, arg ...string) *exec.Cmd {
	cmd := exec.Command(path, arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	
	return cmd
}

func (t *SpeedTester) appendResultToFile(result string) {
	datawriter := bufio.NewWriter(t.file)

	if _, err := datawriter.WriteString(result); err != nil {
		log.Printf("Error writing test result to %s: %v\n", t.filePath, err)
	}

	datawriter.Flush()
}

func (t *SpeedTester) scheduleStop() {
	// TODO: Correct back to actual time
	time.Sleep(time.Duration(t.testDurationHours) * time.Hour)

	log.Println("Termination requested as scheduled.")
	t.Stop()
}

func (t *SpeedTester) Stop() {
	log.Println("Waiting for test to finish before stopping process ...")

	if t.ticker != nil {
		t.ticker.Stop()
	}
	t.stop <- true
}

func (t *SpeedTester) Cleanup() {
	err := t.file.Close()
	if err != nil {
		fmt.Printf("Error closing output file at %s: %v", t.filePath, t.file)
	} else {
		log.Printf("All test results have ben saved to %s", t.filePath)
	}

}
