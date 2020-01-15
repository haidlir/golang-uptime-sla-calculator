package slacalculator

import (
	"testing"
)

func TestValidateSqfValues(t *testing.T) {
	bakti1UptimeChronology1 := Bakti1UptimeChronology{
		StartTimestamps:     1546300800,
		EndTimestamps:       1546301100,
		UptimeValue:         1600,
		Status:              BaktiRunning,
		LinkFailureDuration: 0,
		RestitutionDuration: 0,
	}
	bakti1UptimeChronology2 := Bakti1UptimeChronology{
		StartTimestamps:     1546301100,
		EndTimestamps:       1546301400,
		UptimeValue:         1900,
		Status:              BaktiRunning,
		LinkFailureDuration: 0,
		RestitutionDuration: 0,
	}
	bakti1UptimeChronology3 := Bakti1UptimeChronology{
		StartTimestamps:     1546301400,
		EndTimestamps:       1546301700,
		UptimeValue:         2200,
		Status:              BaktiRunning,
		LinkFailureDuration: 0,
		RestitutionDuration: 0,
	}
	chronologies := []Bakti1UptimeChronology{
		bakti1UptimeChronology1,
		bakti1UptimeChronology2,
		bakti1UptimeChronology3,
	}
	sqfTimestamps1OK := []int64{1546300800, 1546301100, 1546301400}
	sqfTimestamps1NOK := []int64{1546300800, 1546301110, 1546301400}
	sqfValues1 := []float64{7.1, 7.2, 7.3}
	sqfTimestamps2OK := []int64{1546300800, 1546301100, 1546301400, 1546301700}
	sqfValues2 := []float64{7.1, 7.2, 7.3, 7.4}
	sqfTimestamps3OK := []int64{1546300800, 1546301100}
	sqfValues3 := []float64{7.1, 7.2}

	var err error
	err = validateSqfValues(chronologies, sqfTimestamps1OK, sqfValues1)
	if err != nil {
		t.Errorf("it should be OK: %v", err)
	}
	err = validateSqfValues(chronologies, sqfTimestamps1NOK, sqfValues1)
	if err == nil {
		t.Errorf("it should be NOK")
	}
	err = validateSqfValues(chronologies, sqfTimestamps2OK, sqfValues2)
	if err == nil {
		t.Errorf("it should be NOK")
	}
	err = validateSqfValues(chronologies, sqfTimestamps3OK, sqfValues3)
	if err == nil {
		t.Errorf("it should be NOK")
	}
}
