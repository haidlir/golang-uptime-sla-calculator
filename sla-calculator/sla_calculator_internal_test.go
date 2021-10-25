package slacalculator

import (
	"testing"
)

func TestTransformToSpreadedUptime(t *testing.T) {
	startTime := int64(0)
	endTime := int64(20)
	timestamps := []int64{10, 20}
	uptimeVals := []int64{13, 23}
	toleranceDeltaRatio := 0.9
	_, countedVals := transformToSpreadedUptime(startTime, endTime, timestamps, uptimeVals, toleranceDeltaRatio)
	if countedVals[0] != 10 {
		t.Errorf("The first value of counted values is %v, instead of 10", countedVals[0])
	}
}
