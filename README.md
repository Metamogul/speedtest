# speedtest-series

This a little command line utility built on top of [Ookla's speedtest tool](https://www.speedtest.net/de/apps/cli). It executes that tool in regular intervals until a given time has passed and records the output as time series in a single CSV table. This project is still WIP. To modify any options you'll need to change the source code. Take a look at the ToDos below for what's yet to be implemented. 

## Building ##

Clone the repository ( `git clone git@github.com:Metamogul/speedtest.git` ), `cd` to the speedtest-series subfolder and run `go build` or `go run main.go` to build and run the tool right away.

### Dependencies ###

For it's functionality the app depends on `speedtest` CLI by Ookla. This dependency is not included. Please follow the [offical directions](https://www.speedtest.net/de/apps/cli) on how to install.

## Usage ##

Usage of speedtest-series:

	-f, --filepath string
		Full path including filename of the result file
	-d, --test-duration-hours int
		Duration after which to terminate the test series. Pass 0 to continue indefenitely (default 6)
	-i, --test-interval-minutes int
		Interval in between single tests, provided in minutes (default 5)
 	
	--help
		Display this help

## ToDos ##

- Option for additinal log file instead of outputting logs only to stdout
- Tests
- Create `brew` package for MacOS
- Include binary(s) in release
- Create distributable format for other OSes
- Option to use stdout for result output instead of outputting logs

## Contributing ##

Feel free to contribute by creating a fork and issuing a pull request. When issuing a pull request, it would be nice if you could relate it to an open ticket so there's documentation later on.

## Reporting a bug ##

To report a bug, please [create an issue ticket](https://github.com/Metamogul/speedtest-series/issues) for it. In the ticket please provide a description of the state of the app, the action you've been performing, the expected outcome and the actual outcome. Also include the go version you've been using as reported by `go version` as well as any other information that seems relevant to you.

## License ##

This project is distributed in it's entirety under the permissive [Apache 2.0 license](https://github.com/Metamogul/speedtest-series/blob/main/LICENSE) as included.
