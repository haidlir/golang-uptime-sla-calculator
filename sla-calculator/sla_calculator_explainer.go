package slacalculator

const (
	// BaktiRunning is the status definition of running services
	BaktiRunning = 1
	// BaktiLinkFailure is the status definition of down services caused by link failure
	BaktiLinkFailure = 2
	// BaktiPowerFailure is the status definition of down services caused by power failure
	BaktiPowerFailure = 3
	// BaktiOpen is the status definition of down services caused by no unkown reason (open)
	BaktiOpen = 4
)

// Bakti1UptimeChronology explains the chronolgy of each interval uptime data
type Bakti1UptimeChronology struct {
	StartTimestamps int64
	EndTimestamps   int64
	UptimeValue     int64
	Status          int
}

func transformSomething(timestamps []int64, uptimeValues []int64) (countedVals []int64) {
	// Copy uptimeValues to countedVals
	countedVals = append(countedVals, uptimeValues...)
	maxIteration := len(uptimeValues)
	var i, j int = 0, 0
	for (i < maxIteration) && (j < maxIteration) {
		if uptimeValues[i] == 0 {
			i++
			continue
		}
		if j != i {
			for k := i; countedVals[k] > 0 && k > j; k-- {
				countedVals[k-1] = countedVals[k] - (timestamps[k] - timestamps[k-1])
			}
			// normalisasi index i dan j
			j = i
		}
		i++
		j++
	}
	return countedVals
}

// ExplainBakti1Uptime explains the status of service refers to uptime data from BAKTI series.
func (u *UptimeSLACalculator) ExplainBakti1Uptime() []Bakti1UptimeChronology {
	timestamps := []int64{}
	uptimeValues := []int64{}
	if u.startTime != u.timestamps[0] {
		timestamps = append(timestamps, u.startTime)
		uptimeValues = append(uptimeValues, int64(0))
	}
	timestamps = append(timestamps, u.timestamps...)
	uptimeValues = append(uptimeValues, u.uptimeValues...)
	if u.endTime != u.timestamps[len(u.timestamps)-1] {
		timestamps = append(timestamps, u.endTime)
		uptimeValues = append(uptimeValues, int64(0))
	}
	countedVals := transformSomething(timestamps, uptimeValues)
	chronologies := []Bakti1UptimeChronology{}
	for i := 0; i < len(countedVals); i++ {
		if i == len(countedVals)-1 {
			break
		}
		var status int
		if uptimeValues[i] > 0 {
			status = BaktiRunning
		} else if countedVals[i] > 0 && uptimeValues[i] == 0 {
			status = BaktiLinkFailure
		} else if countedVals[i] <= 0 {
			status = BaktiPowerFailure
		}
		chronology := Bakti1UptimeChronology{
			StartTimestamps: timestamps[i],
			EndTimestamps:   timestamps[i+1],
			UptimeValue:     uptimeValues[i],
			Status:          status,
		}
		chronologies = append(chronologies, chronology)
	}
	// Open Status Correction from start
	for i := 0; i < len(chronologies) && countedVals[i] == 0; i++ {
		chronologies[i].Status = BaktiOpen
	}
	// Open Status Correction from end
	for i := len(chronologies) - 1; i > 0 && countedVals[i] == 0; i-- {
		chronologies[i].Status = BaktiOpen
	}
	return chronologies
}

// Bakti1Availability explains the chronolgy of each interval uptime data
type Bakti1Availability struct {
	Availability        float64
	LinkFailureDuration int64
	OpenDuration        int64
	Chronologies        []Bakti1UptimeChronology
}

// CalcBakti1Uptime returns the SLA and explains the status of service refers to uptime data from BAKTI series.
func (u *UptimeSLACalculator) CalcBakti1Uptime() *Bakti1Availability {
	chronologies := u.ExplainBakti1Uptime()
	// Calc Availability
	periodDuration := u.endTime - u.startTime
	var linkFailureDuration int64
	var openDuration int64
	for _, chronology := range chronologies {
		if chronology.Status == BaktiLinkFailure {
			linkFailureDuration += chronology.EndTimestamps - chronology.StartTimestamps
		} else if chronology.Status == BaktiOpen {
			openDuration += chronology.EndTimestamps - chronology.StartTimestamps
		}
	}
	availability := 1.0 - (float64(linkFailureDuration+openDuration) / float64(periodDuration))
	baktiAvailability := Bakti1Availability{
		Availability:        availability,
		LinkFailureDuration: linkFailureDuration,
		OpenDuration:        openDuration,
		Chronologies:        chronologies,
	}
	return &baktiAvailability
}

// CalcBakti1UptimeTrimmed returns the SLA and explains the status of service refers to uptime data from BAKTI series.
func (u *UptimeSLACalculator) CalcBakti1UptimeTrimmed(sTrimDate, eTrimDate int64, isFinalCalc bool) *Bakti1Availability {
	chronologies := u.ExplainBakti1Uptime()
	trimmedChronology := trimChronology(chronologies, sTrimDate, eTrimDate, isFinalCalc)
	// Calc Availability
	periodDuration := eTrimDate - sTrimDate
	var linkFailureDuration int64
	var openDuration int64
	for _, chronology := range trimmedChronology {
		if chronology.Status == BaktiLinkFailure {
			linkFailureDuration += chronology.EndTimestamps - chronology.StartTimestamps
		} else if chronology.Status == BaktiOpen {
			openDuration += chronology.EndTimestamps - chronology.StartTimestamps
		}
	}
	availability := 1.0 - (float64(linkFailureDuration+openDuration) / float64(periodDuration))
	baktiAvailability := Bakti1Availability{
		Availability:        availability,
		LinkFailureDuration: linkFailureDuration,
		OpenDuration:        openDuration,
		Chronologies:        trimmedChronology,
	}
	return &baktiAvailability
	// return nil
}

func trimChronology(chronologies []Bakti1UptimeChronology, sTrimDate, eTrimDate int64, isFinalCalc bool) []Bakti1UptimeChronology {
	iStart := 0
	iEnd := len(chronologies)
	// Check frontside
	for i := 0; i < len(chronologies); i++ {
		if chronologies[i].StartTimestamps >= sTrimDate {
			iStart = i
			break
		}
	}
	// Check backside
	for i := len(chronologies) - 1; i >= 0; i-- {
		if chronologies[i].EndTimestamps <= eTrimDate {
			iEnd = i
			break
		}
	}
	chronologies[iStart].StartTimestamps = sTrimDate
	chronologies[iEnd].EndTimestamps = eTrimDate
	newChronologies := chronologies[iStart : iEnd+1]
	// Normalize Open to PowerFailure if isFinalCalc true
	if !isFinalCalc {
		return newChronologies
	}
	for i := 0; i < len(newChronologies); i++ {
		if newChronologies[i].Status != BaktiOpen {
			break
		} else {
			newChronologies[i].Status = BaktiPowerFailure
		}
	}
	for i := len(newChronologies) - 1; i >= 0; i-- {
		if newChronologies[i].Status != BaktiOpen {
			break
		} else {
			newChronologies[i].Status = BaktiPowerFailure
		}
	}
	return newChronologies
}
