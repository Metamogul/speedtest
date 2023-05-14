package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Metamogul/speedtest-series/speedtester"
)

const (
	testIntervalMinutesDefault = 5
	testIntervalMinutesMin     = 1
	testDurationHoursDefault   = 6
	usage                      = `Usage of ./speedtest-series:
  -f, --filepath string
        Full path including filename of the result file
  -d, --test-duration-hours int
        Duration after which to terminate the test series. Pass 0 to continue indefenitely (default 6)
  -i, --test-interval-minutes int
        Interval in between single tests, provided in minutes (default 5)`
)

func main() {
	filePath, testIntervalMinutes, testDurationHours, err := parseAndValidateArguments()

	if err != nil {
		log.Printf("Error validating arguments: %v\n", err)
	}

	tester := speedtester.NewSpeedTester(filePath, testIntervalMinutes, testDurationHours)

	startLookingForTerminationSignal(tester)
	<-performTest(tester)

	tester.Cleanup()
}

func parseAndValidateArguments() (filePath string, testIntervalMinutes int, testDurationHours int, err error) {
	flag.StringVar(&filePath, "filepath", "", "Full path including filename of the result file")
	flag.StringVar(&filePath, "f", "", "Full path including filename of the result file")
	flag.IntVar(&testIntervalMinutes, "test-interval-minutes", 5, "Interval in between single tests, provided in minutes")
	flag.IntVar(&testIntervalMinutes, "i", 5, "Interval in between single tests, provided in minutes")
	flag.IntVar(&testDurationHours, "test-duration-hours", 6, "Duration after which to terminate the test series. Pass 0 to continue indefenitely")
	flag.IntVar(&testDurationHours, "d", 6, "Duration after which to terminate the test series. Pass 0 to continue indefenitely")

	flag.Usage = func() {
		fmt.Println(usage)
	}

	flag.Parse()

	if filePath == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return "", 0, 0, err
		}

		filePath = userHome + "/result.csv"
	}

	if testIntervalMinutes == 0 {
		testIntervalMinutes = testIntervalMinutesDefault
	}

	if testIntervalMinutes < testIntervalMinutesMin {
		testIntervalMinutes = testIntervalMinutesMin
	}

	return
}

func performTest(tester *speedtester.SpeedTester) chan bool {
	tester.Initialize()
	tester.ScheduleStopAsync()
	return tester.RunAsync()
}

func startLookingForTerminationSignal(tester *speedtester.SpeedTester) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		s := <-signals
		fmt.Println()
		log.Printf("Termination requested by %s.\n", s)
		tester.Stop()
	}()
}
