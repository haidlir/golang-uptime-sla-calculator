# Uptime SLA Calculator
[![Build Status](https://travis-ci.org/haidlir/golang-uptime-sla-calculator.svg?branch=master)](https://travis-ci.org/haidlir/golang-uptime-sla-calculator) [![Coverage Status](https://coveralls.io/repos/github/haidlir/golang-uptime-sla-calculator/badge.svg?branch=master)](https://coveralls.io/github/haidlir/golang-uptime-sla-calculator?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/haidlir/golang-uptime-sla-calculator)](https://goreportcard.com/report/github.com/haidlir/golang-uptime-sla-calculator) [![GoDoc](https://godoc.org/github.com/haidlir/golang-uptime-sla-calculator/sla-calculator?status.svg)](https://godoc.org/github.com/haidlir/golang-uptime-sla-calculator/sla-calculator)<br />

## Status
Experimental

## Motivation
To be used by developers in my organization to calculate SLA based on specific formula.
Freely used by others according to the [LICENSE](https://github.com/haidlir/golang-prtg-api-wrapper/blob/master/LICENSE).

## How to Start
```bash
$ go get github.com/haidlir/golang-uptime-sla-calculator
```

## Example
```go
package main

import (
    "log"

    slacalc "github.com/haidlir/golang-uptime-sla-calculator/sla-calculator"
)

type UptimeData struct {
    Timestamp int64
    Value     int
    Exception bool
}

var (
    startTime int64 = 10000
    endTime int64 = 13000
	toleranceDeltaRatio float64 = 0.9
)

var uptimeSeriesData = []UptimeData{
    UptimeData{10100, 0,    false},
    UptimeData{10200, 0,    false},
    UptimeData{10300, 0,    false},
    UptimeData{10400, 0,    false},
    UptimeData{10500, 270,  false},
    UptimeData{10600, 370,  false},
    UptimeData{10700, 470,  false},
    UptimeData{10800, 570,  false},
    UptimeData{10900, 670,  false},
    UptimeData{11000, 40,   false},
    UptimeData{11100, 140,  false},
    UptimeData{11200, 240,  false},
    UptimeData{11300, 340,  false},
    UptimeData{11400, 440,  false},
    UptimeData{11500, 540,  false},
    UptimeData{11600, 0,    false},
    UptimeData{11700, 0,    false},
    UptimeData{11800, 840,  false},
    UptimeData{11900, 940,  false},
    UptimeData{12000, 1040, false},
    UptimeData{12100, 1140, false},
    UptimeData{12200, 0,    false},
    UptimeData{12300, 0,    false},
    UptimeData{12400, 0,    false},
    UptimeData{12500, 10,   false},
    UptimeData{12600, 110,  false},
    UptimeData{12700, 210,  false},
    UptimeData{12800, 0,    true},
    UptimeData{12900, 0,    true},
    UptimeData{13000, 0,    false},
}

func main() {
    // Prepare the data
    uptimeVals := []int{}
    timestamps := []int64{}
    exceptions := []bool{}
    for _, val := range(uptimeSeriesData) {
        uptimeVals = append(uptimeVals, val.Value)
        timestamps = append(timestamps, val.Timestamp)
        exceptions = append(exceptions, val.Exception)
    }
    // Create calculator object
    calc, err := slacalc.NewUptimeSLACalculator(startTime, endTime, timestamps, uptimeVals, toleranceDeltaRatio, exceptions)
    if err != nil {
        log.Fatalf("An Error should not be accoured: %v", err)
    }
    // Print the calculation result
    log.Printf("Connectivity SLA: %.5f", calc.CalculateSNMPAvailability())
    log.Printf("Uptime SLA: %.5f", calc.CalculateUptimeAvailability())
    log.Printf("(Custom) SLA 1: %.5f", calc.CalculateSLA1Availability())
    log.Printf("(Custom) SLA 2: %.5f", calc.CalculateSLA2Availability())
}
```
[More Example...](https://github.com/haidlir/golang-uptime-sla-calculator/tree/master/_example)

## License
It is released under the MIT license. See
[LICENSE](https://github.com/haidlir/golang-uptime-sla-calculator/blob/master/LICENSE).
