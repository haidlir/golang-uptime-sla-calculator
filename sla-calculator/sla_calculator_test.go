package slacalculator_test

import (
	// "fmt"
	// "log"
	"math"
	"testing"

	slacalc "github.com/haidlir/golang-uptime-sla-calculator/sla-calculator"
)

const (
	ACCURACY = 0.001
)

type UptimeData struct {
	Timestamp int64
	Value     int
	Exception bool
}

var (
	startTime           int64   = 10000
	endTime             int64   = 13000
	toleranceDeltaRatio float64 = 0.9
)

var uptimeSeriesData = []UptimeData{
	{10100, 0, false},
	{10200, 0, false},
	{10300, 0, false},
	{10400, 0, false},
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
	{11600, 0, false},
	{11700, 0, false},
	{11800, 840, false},
	{11900, 940, false},
	{12000, 1040, false},
	{12100, 1140, false},
	{12200, 0, false},
	{12300, 0, false},
	{12400, 0, false},
	{12500, 10, false},
	{12600, 110, false},
	{12700, 210, false},
	{12800, 0, true},
	{12900, 0, true},
	{13000, 0, false},
}

var allDownUptimeSeriesData = []UptimeData{
	{10100, 0, false},
	{10200, 0, false},
	{10300, 0, false},
	{10400, 0, false},
	{10500, 0, false},
	{10600, 0, false},
	{10700, 0, false},
	{10800, 0, false},
	{10900, 0, false},
	{11000, 0, false},
	{11100, 0, false},
	{11200, 0, false},
	{11300, 0, false},
	{11400, 0, false},
	{11500, 0, false},
	{11600, 0, false},
	{11700, 0, false},
	{11800, 0, false},
	{11900, 0, false},
	{12000, 0, false},
	{12100, 0, false},
	{12200, 0, false},
	{12300, 0, false},
	{12400, 0, false},
	{12500, 0, false},
	{12600, 0, false},
	{12700, 0, false},
	{12800, 0, false},
	{12900, 0, false},
	{13000, 0, false},
}

var allUpUptimeSeriesData = []UptimeData{
	{10100, 100, false},
	{10200, 200, false},
	{10300, 300, false},
	{10400, 400, false},
	{10500, 500, false},
	{10600, 600, false},
	{10700, 700, false},
	{10800, 800, false},
	{10900, 900, false},
	{11000, 1000, false},
	{11100, 1100, false},
	{11200, 1200, false},
	{11300, 1300, false},
	{11400, 1400, false},
	{11500, 1500, false},
	{11600, 1600, false},
	{11700, 1700, false},
	{11800, 1800, false},
	{11900, 1900, false},
	{12000, 2000, false},
	{12100, 2100, false},
	{12200, 2200, false},
	{12300, 2300, false},
	{12400, 2400, false},
	{12500, 2500, false},
	{12600, 2600, false},
	{12700, 2700, false},
	{12800, 2800, false},
	{12900, 2900, false},
	{13000, 3000, false},
}

var (
	expetedSNMPAvailability              float64 = 0.58
	expetedUptimeAvailability            float64 = 0.6516
	expetedSLA1Availability              float64 = 0.6774
	expetedSLA2Availability              float64 = 0.7419
	expetedSLA2AvailabilityWithException float64 = 0.9

	expectedAllDown float64 = 0.0
	expectedAllUp   float64 = 1.0
)

func TestUptimeSLACalculator(t *testing.T) {
	// form each array
	{
		endTime := endTime + 100
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range uptimeSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		calc, err := slacalc.NewUptimeSLACalculator(startTime, endTime, timestamps, uptimeVals, toleranceDeltaRatio, exceptions)
		if err != nil {
			t.Fatalf("An Error should not be accoured: %v", err)
		}
		t.Run("CalculateSNMPAvailability", func(t *testing.T) {
			snmpAvai := calc.CalculateSNMPAvailability()
			if math.Abs(snmpAvai-expetedSNMPAvailability) >= ACCURACY {
				t.Errorf("The calculated SNMP Availability value is %v, instead of %v", snmpAvai, expetedSNMPAvailability)
			}
		})
		t.Run("CalculateUptimeAvailability", func(t *testing.T) {
			uptimeAvai := calc.CalculateUptimeAvailability()
			if math.Abs(uptimeAvai-expetedUptimeAvailability) >= ACCURACY {
				t.Errorf("The calculated Uptime Availability value is %v, instead of %v", uptimeAvai, expetedUptimeAvailability)
			}
		})
		t.Run("CalculateSLA1Availability", func(t *testing.T) {
			sla1Avai := calc.CalculateSLA1Availability()
			if math.Abs(sla1Avai-expetedSLA1Availability) >= ACCURACY {
				t.Errorf("The calculated SLA 1 Availability value is %v, instead of %v", sla1Avai, expetedSLA1Availability)
			}
		})
		t.Run("CalculateSLA2Availability", func(t *testing.T) {
			sla2Avai := calc.CalculateSLA2Availability()
			if math.Abs(sla2Avai-expetedSLA2Availability) >= ACCURACY {
				t.Errorf("The calculated SLA 2 Availability value is %v, instead of %v", sla2Avai, expetedSLA2Availability)
			}
		})
		t.Run("GetUptimeStateSeriesData", func(t *testing.T) {
			states := calc.GetUptimeStateSeriesData()
			stateLen := map[string]int{}
			for _, state := range states {
				if _, ok := stateLen[state]; ok {
					stateLen[state]++
				} else {
					stateLen[state] = 1
				}
			}
			if _, ok := stateLen["up"]; !ok {
				t.Fatalf("No UP state in the states array")
			}
			if stateLen["up"] != 22 {
				t.Errorf("The amount of UP state is %v instead of 22", stateLen["up"])
			}
			if _, ok := stateLen["down"]; !ok {
				t.Fatalf("No DOWN state in the states array")
			}
			if stateLen["down"] != 5 {
				t.Errorf("The amount of DOWN state is %v instead of 5", stateLen["down"])
			}
			if _, ok := stateLen["open"]; !ok {
				t.Fatalf("No OPEN state in the states array")
			}
			if stateLen["open"] != 3 {
				t.Errorf("The amount of OPEN state is %v instead of 3", stateLen["open"])
			}
		})
		t.Run("GetTotalUptimeAndDowntime", func(t *testing.T) {
			uptime, downtime := calc.GetTotalUptimeAndDowntime()
			if math.Abs((float64(uptime)/float64(uptime+downtime))-expetedUptimeAvailability) >= ACCURACY {
				t.Errorf("The calculated Uptime Availability value is %v, instead of %v", float64(uptime)/float64(uptime+downtime), expetedUptimeAvailability)
			}
		})
	}

	{
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range allDownUptimeSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		calc, err := slacalc.NewUptimeSLACalculator(startTime, endTime, timestamps, uptimeVals, toleranceDeltaRatio, exceptions)
		if err != nil {
			t.Fatalf("An Error should not be accoured: %v", err)
		}
		t.Run("CalculateSNMPAvailability: All Down", func(t *testing.T) {
			snmpAvai := calc.CalculateSNMPAvailability()
			if math.Abs(snmpAvai-expectedAllDown) >= ACCURACY {
				t.Errorf("The calculated SNMP Availability value is %v, instead of %v", snmpAvai, expectedAllDown)
			}
		})
		t.Run("CalculateUptimeAvailability: All Down", func(t *testing.T) {
			uptimeAvai := calc.CalculateUptimeAvailability()
			if math.Abs(uptimeAvai-expectedAllDown) >= ACCURACY {
				t.Errorf("The calculated Uptime Availability value is %v, instead of %v", uptimeAvai, expectedAllDown)
			}
		})
		t.Run("CalculateSLA1Availability: All Down", func(t *testing.T) {
			sla1Avai := calc.CalculateSLA1Availability()
			if math.Abs(sla1Avai-expectedAllDown) >= ACCURACY {
				t.Errorf("The calculated SLA 1 Availability value is %v, instead of %v", sla1Avai, expectedAllDown)
			}
		})
		t.Run("CalculateSLA2Availability: All Down", func(t *testing.T) {
			sla2Avai := calc.CalculateSLA2Availability()
			if math.Abs(sla2Avai-expectedAllDown) >= ACCURACY {
				t.Errorf("The calculated SLA 2 Availability value is %v, instead of %v", sla2Avai, expectedAllDown)
			}
		})
		t.Run("CalculateSLA2Availability: All Down with Exception", func(t *testing.T) {
			exceptions[len(exceptions)-1] = true
			calc, err := slacalc.NewUptimeSLACalculator(startTime, endTime, timestamps, uptimeVals, toleranceDeltaRatio, exceptions)
			if err != nil {
				t.Fatalf("An Error should not be accoured: %v", err)
			}
			expected := 0.0333
			sla2Avai := calc.CalculateSLA2Availability()
			if math.Abs(sla2Avai-expected) >= ACCURACY {
				t.Errorf("The calculated SLA 2 Availability value is %v, instead of %v", sla2Avai, expected)
			}
		})
		t.Run("CalculateSLA2Availability: All Down with Exception 2", func(t *testing.T) {
			exceptions[len(exceptions)-1] = false
			exceptions[0] = true
			calc, err := slacalc.NewUptimeSLACalculator(startTime, endTime, timestamps, uptimeVals, toleranceDeltaRatio, exceptions)
			if err != nil {
				t.Fatalf("An Error should not be accoured: %v", err)
			}
			expected := 0.0333
			sla2Avai := calc.CalculateSLA2Availability()
			if math.Abs(sla2Avai-expected) >= ACCURACY {
				t.Errorf("The calculated SLA 2 Availability value is %v, instead of %v", sla2Avai, expected)
			}
		})
		t.Run("GetUptimeStateSeriesData: All Down", func(t *testing.T) {
			states := calc.GetUptimeStateSeriesData()
			stateLen := map[string]int{}
			for _, state := range states {
				if _, ok := stateLen[state]; ok {
					stateLen[state]++
				} else {
					stateLen[state] = 1
				}
			}
			if _, ok := stateLen["up"]; ok {
				t.Fatalf("UP state exists in the states array")
			}
			if stateLen["up"] != 0 {
				t.Errorf("The amount of UP state is %v instead of 0", stateLen["up"])
			}
			if _, ok := stateLen["down"]; ok {
				t.Fatalf("DOWN state exists in the states array")
			}
			if stateLen["down"] != 0 {
				t.Errorf("The amount of DOWN state is %v instead of 0", stateLen["down"])
			}
			if _, ok := stateLen["open"]; !ok {
				t.Fatalf("No OPEN state in the states array")
			}
			if stateLen["open"] != 30 {
				t.Errorf("The amount of OPEN state is %v instead of 30", stateLen["open"])
			}
		})
		t.Run("GetTotalUptimeAndDowntime", func(t *testing.T) {
			uptime, downtime := calc.GetTotalUptimeAndDowntime()
			if math.Abs((float64(uptime)/float64(uptime+downtime))-expectedAllDown) >= ACCURACY {
				t.Errorf("The calculated Uptime Availability value is %v, instead of %v", float64(uptime)/float64(uptime+downtime), expectedAllDown)
			}
		})
	}

	{
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range allUpUptimeSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		calc, err := slacalc.NewUptimeSLACalculator(startTime, endTime, timestamps, uptimeVals, toleranceDeltaRatio, exceptions)
		if err != nil {
			t.Fatalf("An Error should not be accoured: %v", err)
		}
		t.Run("CalculateSNMPAvailability: All Up", func(t *testing.T) {
			snmpAvai := calc.CalculateSNMPAvailability()
			if math.Abs(snmpAvai-expectedAllUp) >= ACCURACY {
				t.Errorf("The calculated SNMP Availability value is %v, instead of %v", snmpAvai, expectedAllUp)
			}
		})
		t.Run("CalculateUptimeAvailability: All Up", func(t *testing.T) {
			uptimeAvai := calc.CalculateUptimeAvailability()
			if math.Abs(uptimeAvai-expectedAllUp) >= ACCURACY {
				t.Errorf("The calculated Uptime Availability value is %v, instead of %v", uptimeAvai, expectedAllUp)
			}
		})
		t.Run("CalculateSLA1Availability: All Up", func(t *testing.T) {
			sla1Avai := calc.CalculateSLA1Availability()
			if math.Abs(sla1Avai-expectedAllUp) >= ACCURACY {
				t.Errorf("The calculated SLA 1 Availability value is %v, instead of %v", sla1Avai, expectedAllUp)
			}
		})
		t.Run("CalculateSLA2Availability: All Up", func(t *testing.T) {
			sla2Avai := calc.CalculateSLA2Availability()
			if math.Abs(sla2Avai-expectedAllUp) >= ACCURACY {
				t.Errorf("The calculated SLA 2 Availability value is %v, instead of %v", sla2Avai, expectedAllUp)
			}
		})
		t.Run("GetUptimeStateSeriesData", func(t *testing.T) {
			states := calc.GetUptimeStateSeriesData()
			stateLen := map[string]int{}
			for _, state := range states {
				if _, ok := stateLen[state]; ok {
					stateLen[state]++
				} else {
					stateLen[state] = 1
				}
			}
			if _, ok := stateLen["up"]; !ok {
				t.Fatalf("No UP state in the states array")
			}
			if stateLen["up"] != 30 {
				t.Errorf("The amount of UP state is %v instead of 30", stateLen["up"])
			}
			if _, ok := stateLen["down"]; ok {
				t.Fatalf("DOWN state exists in the states array")
			}
			if stateLen["down"] != 0 {
				t.Errorf("The amount of DOWN state is %v instead of 0", stateLen["down"])
			}
			if _, ok := stateLen["open"]; ok {
				t.Fatalf("OPEN state exists in the states array")
			}
			if stateLen["open"] != 0 {
				t.Errorf("The amount of OPEN state is %v instead of 0", stateLen["open"])
			}
		})
		t.Run("GetTotalUptimeAndDowntime", func(t *testing.T) {
			uptime, downtime := calc.GetTotalUptimeAndDowntime()
			if math.Abs((float64(uptime)/float64(uptime+downtime))-expectedAllUp) >= ACCURACY {
				t.Errorf("The calculated Uptime Availability value is %v, instead of %v", float64(uptime)/float64(uptime+downtime), expectedAllUp)
			}
		})
	}

	// Error Case
	{
		uptimeVals := []int{}
		timestamps := []int64{}
		exceptions := []bool{}
		for _, val := range allUpUptimeSeriesData {
			uptimeVals = append(uptimeVals, val.Value)
			timestamps = append(timestamps, val.Timestamp)
			exceptions = append(exceptions, val.Exception)
		}
		t.Run("Start Time is (-)", func(t *testing.T) {
			_, err := slacalc.NewUptimeSLACalculator(-1, endTime, timestamps, uptimeVals, toleranceDeltaRatio, exceptions)
			if err == nil {
				t.Fatalf("Error should be occured.")
			}
		})
		t.Run("Timestamp length unmatches to uptime length", func(t *testing.T) {
			_, err := slacalc.NewUptimeSLACalculator(startTime, endTime, timestamps[:len(timestamps)-2], uptimeVals, toleranceDeltaRatio, exceptions)
			if err == nil {
				t.Fatalf("Error should be occured.")
			}
		})
		t.Run("Uptime length unmatches to timestamp  length", func(t *testing.T) {
			_, err := slacalc.NewUptimeSLACalculator(startTime, endTime, timestamps, uptimeVals[:len(timestamps)-2], toleranceDeltaRatio, exceptions)
			if err == nil {
				t.Fatalf("Error should be occured.")
			}
		})
		t.Run("Tolerance Delta Ratio is not between 0 to 1", func(t *testing.T) {
			_, err := slacalc.NewUptimeSLACalculator(startTime, endTime, timestamps, uptimeVals, -1, exceptions)
			if err == nil {
				t.Fatalf("Error should be occured.")
			}
		})
		t.Run("SLA 2 Nil Exception", func(t *testing.T) {
			_, err := slacalc.NewUptimeSLACalculator(startTime, endTime, timestamps, uptimeVals, toleranceDeltaRatio, exceptions[:len(exceptions)-2])
			if err == nil {
				t.Fatalf("Error should be occured.")
			}
		})
		t.Run("Older Start Time", func(t *testing.T) {
			_, err := slacalc.NewUptimeSLACalculator(timestamps[0]+1, endTime, timestamps, uptimeVals, toleranceDeltaRatio, exceptions[:len(exceptions)-2])
			if err == nil {
				t.Fatalf("Error should be occured.")
			}
		})
		t.Run("Earlier End Time", func(t *testing.T) {
			_, err := slacalc.NewUptimeSLACalculator(startTime, timestamps[len(timestamps)-1]-1, timestamps, uptimeVals, toleranceDeltaRatio, exceptions[:len(exceptions)-2])
			if err == nil {
				t.Fatalf("Error should be occured.")
			}
		})
		t.Run("Nil Timestamp", func(t *testing.T) {
			_, err := slacalc.NewUptimeSLACalculator(startTime, endTime, nil, uptimeVals, toleranceDeltaRatio, exceptions[:len(exceptions)-2])
			if err == nil {
				t.Fatalf("Error should be occured.")
			}
		})
		t.Run("Unordered Timestamp", func(t *testing.T) {
			timestamps[0], timestamps[len(timestamps)-1] = timestamps[len(timestamps)-1], timestamps[0]
			_, err := slacalc.NewUptimeSLACalculator(startTime, endTime, timestamps, uptimeVals, toleranceDeltaRatio, exceptions)
			if err == nil {
				t.Fatalf("Error should be occured.")
			}
		})
	}
}
