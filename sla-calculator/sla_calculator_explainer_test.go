package slacalculator_test

import (
	"math"
	"testing"

	slacalc "github.com/haidlir/golang-uptime-sla-calculator/sla-calculator"
)

var (
	startTimeBakti int64 = 10000
	endTimeBakti   int64 = 13000
)

var uptimeBaktiSeriesData = []UptimeData{
	{10000, 0, false}, // open
	{10100, 0, false}, // open
	{10200, 0, false}, // power failure
	{10300, 0, false}, // link failure
	{10400, 0, false}, // link failure
	{10500, 270, false},
	{10600, 370, false},
	{10700, 470, false},
	{10800, 570, false},
	{10900, 670, false},
	{11000, 40, false},
	{11100, 140, false},
	{11200, 240, false},
	{11300, 340, false},
	{11400, 440, false},
	{11500, 540, false},
	{11600, 0, false}, // link failure
	{11700, 0, false}, // link failure
	{11800, 840, false},
	{11900, 940, false},
	{12000, 1040, false},
	{12100, 1140, false},
	{12200, 0, false}, // power failure
	{12300, 0, false}, // power failure
	{12400, 0, false}, // power failure
	{12500, 10, false},
	{12600, 110, false},
	{12700, 210, false},
	{12800, 0, true}, // open
	{12900, 0, true}, // open
	{13000, 0, true}, // no next
}

func TestExplainBakti1Uptime(t *testing.T) {
	// form each array
	{
		endTime := endTimeBakti
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range uptimeBaktiSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		calc, err := slacalc.NewUptimeSLACalculator(startTimeBakti, endTime, timestamps, uptimeVals, exceptions)
		if err != nil {
			t.Fatalf("An Error should not be accoured: %v", err)
		}
		uptimeChronologies := calc.ExplainBakti1Uptime()
		if uptimeChronologies == nil {
			t.Fatalf("It should be not nil")
		}
		stateLen := map[int]int{}
		for _, uptimeChronology := range uptimeChronologies {
			state := uptimeChronology.Status
			if _, ok := stateLen[state]; ok {
				stateLen[state]++
			} else {
				stateLen[state] = 1
			}
		}
		if _, ok := stateLen[slacalc.BaktiRunning]; !ok {
			t.Fatalf("No running state in the states array")
		}
		if stateLen[slacalc.BaktiRunning] != 18 {
			t.Errorf("The amount of running state is %v instead of 18", stateLen[slacalc.BaktiRunning])
		}
		if _, ok := stateLen[slacalc.BaktiLinkFailure]; !ok {
			t.Fatalf("No link failure state in the states array")
		}
		if stateLen[slacalc.BaktiLinkFailure] != 4 {
			t.Errorf("The amount of link failure state is %v instead of 4", stateLen[slacalc.BaktiLinkFailure])
		}
		if _, ok := stateLen[slacalc.BaktiPowerFailure]; !ok {
			t.Fatalf("No power failure state in the states array")
		}
		if stateLen[slacalc.BaktiPowerFailure] != 4 {
			t.Errorf("The amount of power failure state is %v instead of 4", stateLen[slacalc.BaktiPowerFailure])
		}
		if _, ok := stateLen[slacalc.BaktiOpen]; !ok {
			t.Fatalf("No open state in the states array")
		}
		if stateLen[slacalc.BaktiOpen] != 4 {
			t.Errorf("The amount of open state is %v instead of 4", stateLen[slacalc.BaktiPowerFailure])
		}
	}
}

func TestCalcBakti1Uptime(t *testing.T) {
	// uptimeBaktiSeriesData
	{
		endTime := endTimeBakti
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range uptimeBaktiSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		calc, err := slacalc.NewUptimeSLACalculator(startTimeBakti, endTime, timestamps, uptimeVals, exceptions)
		if err != nil {
			t.Fatalf("An Error should not be accoured: %v", err)
		}
		bakti1Availability := calc.CalcBakti1Uptime()
		if bakti1Availability == nil {
			t.Fatalf("It should be not nil")
		}
		if math.Abs(bakti1Availability.Availability-.733333) > ACCURACY {
			t.Errorf("Availability for uptime data is %.2f instead of 0.733333", bakti1Availability.Availability)
		}
	}
	// AllUp
	{
		endTime := endTimeBakti
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range allUpUptimeSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		calc, err := slacalc.NewUptimeSLACalculator(timestamps[0], endTime, timestamps, uptimeVals, exceptions)
		if err != nil {
			t.Fatalf("An Error should not be accoured: %v", err)
		}
		bakti1Availability := calc.CalcBakti1Uptime()
		if bakti1Availability == nil {
			t.Fatalf("It should be not nil")
		}
		if math.Abs(bakti1Availability.Availability-1.) > ACCURACY {
			t.Errorf("Availability for all up data is %.2f instead of 1.0", bakti1Availability.Availability)
		}
	}
	// AllDown
	{
		endTime := endTimeBakti
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range allDownUptimeSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		calc, err := slacalc.NewUptimeSLACalculator(startTimeBakti, endTime, timestamps, uptimeVals, exceptions)
		if err != nil {
			t.Fatalf("An Error should not be accoured: %v", err)
		}
		bakti1Availability := calc.CalcBakti1Uptime()
		if bakti1Availability == nil {
			t.Fatalf("It should be not nil")
		}
		if math.Abs(bakti1Availability.Availability-0.0) > ACCURACY {
			t.Errorf("Availability for all down data is %.2f instead of 0.0", bakti1Availability.Availability)
		}
	}
}

func TestCalcBakti1UptimeTrimmed(t *testing.T) {
	// uptimeBaktiSeriesData
	{
		// Trim Time
		var startTrimmedTime int64 = 10250
		var endTrimmedTime int64 = 12650
		endTime := endTimeBakti
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range uptimeBaktiSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		calc, err := slacalc.NewUptimeSLACalculator(startTimeBakti, endTime, timestamps, uptimeVals, exceptions)
		if err != nil {
			t.Fatalf("An Error should not be accoured: %v", err)
		}
		bakti1Availability := calc.CalcBakti1UptimeTrimmed(startTrimmedTime, endTrimmedTime, false)
		if bakti1Availability == nil {
			t.Fatalf("It should be not nil")
		}
		if math.Abs(bakti1Availability.Availability-1.00) > ACCURACY {
			t.Errorf("Availability for uptime data is %.2f instead of 1.00", bakti1Availability.Availability)
		}
	}
	// AllUp
	{
		// Trim Time
		var startTrimmedTime int64 = 10550
		var endTrimmedTime int64 = 11450
		endTime := endTimeBakti
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range uptimeBaktiSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		calc, err := slacalc.NewUptimeSLACalculator(startTimeBakti, endTime, timestamps, uptimeVals, exceptions)
		if err != nil {
			t.Fatalf("An Error should not be accoured: %v", err)
		}
		bakti1Availability := calc.CalcBakti1UptimeTrimmed(startTrimmedTime, endTrimmedTime, false)
		if bakti1Availability == nil {
			t.Fatalf("It should be not nil")
		}
		if math.Abs(bakti1Availability.Availability-1.0) > ACCURACY {
			t.Errorf("Availability for uptime data is %.2f instead of 1.0", bakti1Availability.Availability)
		}
	}
	// Final no Open
	{
		// Trim Time
		var startTrimmedTime int64 = 10000
		var endTrimmedTime int64 = 13000
		endTime := endTimeBakti
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range uptimeBaktiSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		calc, err := slacalc.NewUptimeSLACalculator(startTimeBakti, endTime, timestamps, uptimeVals, exceptions)
		if err != nil {
			t.Fatalf("An Error should not be accoured: %v", err)
		}
		bakti1Availability := calc.CalcBakti1UptimeTrimmed(startTrimmedTime, endTrimmedTime, true)
		if bakti1Availability == nil {
			t.Fatalf("It should be not nil")
		}
		if bakti1Availability.OpenDuration != 0 {
			t.Fatalf("It should be zero")
		}
	}
	// With Open
	{
		// Trim Time
		var startTrimmedTime int64 = 10000
		var endTrimmedTime int64 = 13000
		endTime := endTimeBakti
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range uptimeBaktiSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		calc, err := slacalc.NewUptimeSLACalculator(startTimeBakti, endTime, timestamps, uptimeVals, exceptions)
		if err != nil {
			t.Fatalf("An Error should not be accoured: %v", err)
		}
		bakti1Availability := calc.CalcBakti1UptimeTrimmed(startTrimmedTime, endTrimmedTime, false)
		if bakti1Availability == nil {
			t.Fatalf("It should be not nil")
		}
		if bakti1Availability.OpenDuration == 0 {
			t.Fatalf("It should be not zero")
		}
	}
}
