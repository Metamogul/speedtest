package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Metamogul/speedtest/speedtester"
)

func main() {
	tester := speedtester.NewSpeedTester("/Users/Jan/output.csv", 1, 0)

	startLookingForTerminationSignal(tester)
	<-performTest(tester)

	tester.Cleanup()
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
