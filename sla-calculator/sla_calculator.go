package slacalculator

import (
	"fmt"
)

const (
	// DEFAULT_FLOAT_VALUE is the default value of each SLA Calculation functions.
	DEFAULT_FLOAT_VALUE = -1.0
	// STATE_DOWN is the string value of down uptime state.
	STATE_DOWN = "down"
	// STATE_UP is the string value of up uptime state.
	STATE_UP = "up"
	// STATE_OPEN is the string value of open uptime state.
	STATE_OPEN = "open"
)

// UptimeSLACalculator calculates Uptime SLA parameters based on specified formulas.
type UptimeSLACalculator struct {
	uptimeValues []int64
	timestamps   []int64
	exceptions   []bool
	startTime    int64
	endTime      int64
}

func checkArguments(startTime, endTime int64, timestamps []int64, uptimeValues []int, exceptions []bool) error {
	// check start and end time
	if startTime < 0 || endTime < 0 {
		return fmt.Errorf("Start or End time is less than 0 (-).")
	}
	if len(timestamps) <= 0 {
		return fmt.Errorf("Timestamp array is empty.")
	}
	if startTime > timestamps[0] {
		return fmt.Errorf("Start time is greater than the first timestamp.")
	}
	if endTime < timestamps[len(timestamps)-1] {
		return fmt.Errorf("End time is less than the last timestamp.")
	}
	// check arrays' length
	if len(timestamps) != len(uptimeValues) {
		return fmt.Errorf("Length of timestamps and uptime value is unmatched.")
	}
	if exceptions != nil {
		if len(timestamps) != len(exceptions) {
			return fmt.Errorf("Length of timestamps and exceptions is unmatched")
		}
	}
	// check timestamp sequence
	for i := range timestamps {
		if i < 1 {
			continue
		}
		if timestamps[i] < timestamps[i-1] {
			return fmt.Errorf("Unordered timestamps is detected.")
		}
	}
	// All is well, and ready to be calculated
	return nil
}

func castSliceIntToInt64(vals []int) []int64 {
	newVals := []int64{}
	for _, val := range vals {
		newVals = append(newVals, int64(val))
	}
	return newVals
}

// NewUptimeSLACalculator returns the uptime calculator object.
func NewUptimeSLACalculator(startTime, endTime int64, timestamps []int64, uptimeValues []int, exceptions []bool) (*UptimeSLACalculator, error) {
	if err := checkArguments(startTime, endTime, timestamps, uptimeValues, exceptions); err != nil {
		return nil, err
	}
	return &UptimeSLACalculator{
		castSliceIntToInt64(uptimeValues),
		timestamps,
		exceptions,
		startTime,
		endTime,
	}, nil
}

// CalculateSNMPAvailability returns the availability value (SLA) based on
// the existence of the data in each timestamp.
func (u *UptimeSLACalculator) CalculateSNMPAvailability() float64 {
	timestamps := u.timestamps
	uptimeValues := u.uptimeValues
	startTime := u.startTime
	endTime := u.endTime
	countedVals := []int64{}
	deltaTimeStamps := []int64{}
	for i := range timestamps {
		if i == 0 {
			deltaTimeStamps = append(deltaTimeStamps, timestamps[i]-startTime)
			if uptimeValues[i] <= 0 {
				countedVals = append(countedVals, 0)
			} else {
				countedVals = append(countedVals, timestamps[i]-startTime)
			}
			continue
		}
		deltaTimeStamps = append(deltaTimeStamps, timestamps[i]-timestamps[i-1])
		if uptimeValues[i] <= 0 {
			countedVals = append(countedVals, 0)
		} else {
			countedVals = append(countedVals, timestamps[i]-timestamps[i-1])
		}
	}
	if delta := endTime - timestamps[len(timestamps)-1]; delta > 0 {
		deltaTimeStamps = append(deltaTimeStamps, delta)
		countedVals = append(countedVals, 0)
	}
	var sumCountedVal, sumDeltaTimestamp int64
	for i := range countedVals {
		sumCountedVal += countedVals[i]
		sumDeltaTimestamp += deltaTimeStamps[i]
	}
	return float64(sumCountedVal) / float64(sumDeltaTimestamp)
}

func transformToSpreadedUptime(startTime, endTime int64, timestamps []int64, uptimeValues []int64) (deltaTimeStamps, countedVals []int64) {
	for i := range timestamps {
		if i == 0 {
			delta := timestamps[i] - startTime
			deltaTimeStamps = append(deltaTimeStamps, delta)
			if uptimeValues[i] <= 0 {
				countedVals = append(countedVals, 0)
			} else {
				if uptimeValues[i] > (delta) {
					countedVals = append(countedVals, delta)
				} else {
					countedVals = append(countedVals, uptimeValues[i])
				}
			}
			continue
		}
		delta := timestamps[i] - timestamps[i-1]
		deltaTimeStamps = append(deltaTimeStamps, delta)
		if uptimeValues[i]-uptimeValues[i-1] <= 0 {
			countedVals = append(countedVals, 0)
		} else {
			countedVals = append(countedVals, uptimeValues[i]-uptimeValues[i-1])
		}
	}

	for i := range deltaTimeStamps {
		if countedVals[i] > deltaTimeStamps[i] {
			for j := i; countedVals[j] > deltaTimeStamps[j]; j-- {
                if j == 0 {
                    if countedVals[j] > deltaTimeStamps[j] {
                        countedVals[j] = deltaTimeStamps[j]
                    }
                    break
                }
				countedVals[j-1] = countedVals[j] - deltaTimeStamps[j]
				countedVals[j] = deltaTimeStamps[j]
			}
		}
	}
	return
}

// CalculateUptimeAvailability returns the availability value (SLA) based on
// the uptime value on each timestamp, it figures the availability of a device is UP
// regardless the connectivity state.
func (u *UptimeSLACalculator) CalculateUptimeAvailability() float64 {
	timestamps := u.timestamps
	uptimeValues := u.uptimeValues
	startTime := u.startTime
	endTime := u.endTime
	deltaTimeStamps, countedVals := transformToSpreadedUptime(startTime, endTime, timestamps, uptimeValues)
	if delta := endTime - timestamps[len(timestamps)-1]; delta > 0 {
		deltaTimeStamps = append(deltaTimeStamps, delta)
		countedVals = append(countedVals, 0)
	}
	var sumCountedVal, sumDeltaTimestamp int64
	for i := range countedVals {
		sumCountedVal += countedVals[i]
		sumDeltaTimestamp += deltaTimeStamps[i]
	}
	return float64(sumCountedVal) / float64(sumDeltaTimestamp)
}

// CalculateSLA1Availability returns the availability value (SLA) based on
// the uptime value on each timestamp, it figures the availability of a connectivity
// while also regarding the device state (it counts as down if device is up but the
// the connectivity is down, other scenarios counted as up)
func (u *UptimeSLACalculator) CalculateSLA1Availability() float64 {
	timestamps := u.timestamps
	uptimeValues := u.uptimeValues
	startTime := u.startTime
	endTime := u.endTime
	deltaTimeStamps, countedVals := transformToSpreadedUptime(startTime, endTime, timestamps, uptimeValues)
	open := true
	for i := range timestamps {
		if open && uptimeValues[i] > 0 {
			open = false
		}
		if open {
			countedVals[i] = 0
			continue
		}
		if (uptimeValues[i] <= 0) && (countedVals[i] > 0) {
			countedVals[i] = 0
			continue
		}
		countedVals[i] = deltaTimeStamps[i]
	}
	// correcting open end
	for i := len(timestamps) - 1; (uptimeValues[i] <= 0) && (i >= 0); i-- {
		countedVals[i] = 0
		if i == 0 {
			break
		}
	}
	if delta := endTime - timestamps[len(timestamps)-1]; delta > 0 {
		deltaTimeStamps = append(deltaTimeStamps, delta)
		countedVals = append(countedVals, 0)
	}
	var sumCountedVal, sumDeltaTimestamp int64
	for i := range countedVals {
		sumCountedVal += countedVals[i]
		sumDeltaTimestamp += deltaTimeStamps[i]
	}
	return float64(sumCountedVal) / float64(sumDeltaTimestamp)
}

// CalculateSLA2Availability returns the availability value (SLA) based on
// the Uptime SLA 1 Availability and also regards the exception (proved by justification
// document in the real life , ie: scheduled maintenance.)
func (u *UptimeSLACalculator) CalculateSLA2Availability() float64 {
	timestamps := u.timestamps
	uptimeValues := u.uptimeValues
	exceptions := u.exceptions
	startTime := u.startTime
	endTime := u.endTime
	deltaTimeStamps, countedVals := transformToSpreadedUptime(startTime, endTime, timestamps, uptimeValues)
	open := true
	for i := range timestamps {
		if exceptions[i] {
			countedVals[i] = deltaTimeStamps[i]
			continue
		}
		if open && uptimeValues[i] > 0 {
			open = false
		}
		if open {
			countedVals[i] = 0
			continue
		}
		if (uptimeValues[i] <= 0) && (countedVals[i] > 0) {
			countedVals[i] = 0
			continue
		}
		countedVals[i] = deltaTimeStamps[i]
	}
	// correcting open end
	for i := len(timestamps) - 1; (uptimeValues[i] <= 0) && (i >= 0); i-- {
		if exceptions[i] {
			continue
		}
		countedVals[i] = 0
		if i == 0 {
			break
		}
	}
	if delta := endTime - timestamps[len(timestamps)-1]; delta > 0 {
		deltaTimeStamps = append(deltaTimeStamps, delta)
		countedVals = append(countedVals, 0)
	}
	var sumCountedVal, sumDeltaTimestamp int64
	for i := range countedVals {
		sumCountedVal += countedVals[i]
		sumDeltaTimestamp += deltaTimeStamps[i]
	}
	return float64(sumCountedVal) / float64(sumDeltaTimestamp)
}

// GetUptimeStateSeriesData returns the state of every uptime series data either up, down, or open.
func (u *UptimeSLACalculator) GetUptimeStateSeriesData() []string {
	timestamps := u.timestamps
	uptimeValues := u.uptimeValues
	startTime := u.startTime
	endTime := u.endTime
	_, countedVals := transformToSpreadedUptime(startTime, endTime, timestamps, uptimeValues)
	// Create States Array
	states := []string{}
	for i := range countedVals {
		if countedVals[i] > 0 {
			states = append(states, STATE_UP)
			continue
		}
		states = append(states, STATE_DOWN)
	}
	for i := len(timestamps) - 1; (uptimeValues[i] <= 0) && (i >= 0); i-- {
		states[i] = STATE_OPEN
		if i == 0 {
			break
		}
	}
	return states
}
