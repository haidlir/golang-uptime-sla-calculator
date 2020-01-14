package slacalculator

import "fmt"

const (
	// BaktiRunning is the status definition of running services
	BaktiRunning = 1
	// BaktiLinkFailure is the status definition of down services caused by link failure
	BaktiLinkFailure = 2
	// BaktiPowerFailure is the status definition of down services caused by power failure
	BaktiPowerFailure = 3
	// BaktiOpen is the status definition of down services caused by no unkown reason (open)
	BaktiOpen = 4
	// SqfGte71 is the status of SQF SLA greater than or equal to 7.1
	SqfGte71 = 101
	// Sqflt3 is the status of SQF SLA less than 3
	Sqflt3 = 102
	// Sqfbt713Quota between 7.1 and 3 in quota
	Sqfbt713Quota = 103
	// Sqfbt713NonQuota between 7.1 and 3 non quota
	Sqfbt713NonQuota = 104
)

// Bakti1UptimeChronology explains the chronolgy of each interval uptime data
type Bakti1UptimeChronology struct {
	StartTimestamps     int64
	EndTimestamps       int64
	UptimeValue         int64
	Status              int
	LinkFailureDuration int64
	RestitutionDuration int64
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
			StartTimestamps:     timestamps[i],
			EndTimestamps:       timestamps[i+1],
			UptimeValue:         uptimeValues[i],
			Status:              status,
			RestitutionDuration: 0,
		}
		if status == BaktiLinkFailure {
			chronology.LinkFailureDuration = chronology.EndTimestamps - chronology.StartTimestamps
		} else {
			chronology.LinkFailureDuration = 0
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
	RestitutionDuration int64
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

func calcRestitutionPerPeriod(chronologies []Bakti1UptimeChronology) []Bakti1UptimeChronology {
	// Tolerate 5 Mins every link failure series
	var iStart, iEnd int
	isLinkFailure := false
	var LinkFailureTolerance int64 = 300
	for i := 0; i < len(chronologies); i++ {
		if chronologies[i].Status == BaktiLinkFailure {
			if !isLinkFailure {
				isLinkFailure = true
				iStart = i
			}
			iEnd = i
		} else {
			if isLinkFailure {
				isLinkFailure = false
				tolerance := LinkFailureTolerance
				for j := iEnd; tolerance > 0 && j >= iStart; j-- {
					restitution := chronologies[j].LinkFailureDuration - tolerance
					if (restitution) <= 0 {
						chronologies[j].RestitutionDuration = 0
						tolerance -= LinkFailureTolerance
					} else {
						chronologies[j].RestitutionDuration = restitution
					}
				}
			}
		}
	}
	return chronologies
}

// CalcBakti1UptimeTrimmed returns the SLA and explains the status of service refers to uptime data from BAKTI series.
func (u *UptimeSLACalculator) CalcBakti1UptimeTrimmed(sTrimDate, eTrimDate int64, isFinalCalc bool) *Bakti1Availability {
	chronologies := u.ExplainBakti1Uptime()
	chronologies = calcRestitutionPerPeriod(chronologies)
	trimmedChronology := trimChronology(chronologies, sTrimDate, eTrimDate, isFinalCalc)
	// Calc Availability
	periodDuration := eTrimDate - sTrimDate
	var linkFailureDuration int64
	var restitutionDuration int64
	var openDuration int64
	for _, chronology := range trimmedChronology {
		if chronology.Status == BaktiLinkFailure {
			restitutionDuration += chronology.RestitutionDuration
			linkFailureDuration += chronology.LinkFailureDuration
		} else if chronology.Status == BaktiOpen {
			openDuration += chronology.EndTimestamps - chronology.StartTimestamps
		}
	}
	availability := 1.0 - (float64(restitutionDuration+openDuration) / float64(periodDuration))
	baktiAvailability := Bakti1Availability{
		Availability:        availability,
		LinkFailureDuration: linkFailureDuration,
		RestitutionDuration: restitutionDuration,
		OpenDuration:        openDuration,
		Chronologies:        trimmedChronology,
	}
	return &baktiAvailability
}

func trimChronology(chronologies []Bakti1UptimeChronology, sTrimDate, eTrimDate int64, isFinalCalc bool) []Bakti1UptimeChronology {
	iStart := 0
	iEnd := len(chronologies) - 1
	// Check frontside
	for i := 0; i < len(chronologies); i++ {
		if chronologies[i].StartTimestamps >= sTrimDate && chronologies[i].EndTimestamps > sTrimDate {
			iStart = i
			break
		}
	}
	// Check backside
	for i := len(chronologies) - 1; i >= 0; i-- {
		if chronologies[i].EndTimestamps <= eTrimDate && chronologies[i].StartTimestamps < eTrimDate {
			// Make sure not overlap
			if i < iStart {
				break
			}
			iEnd = i
			break
		}
	}
	if iStart > iEnd {
		iStart, iEnd = iEnd, iStart
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

// BaktiSqfAvailability explains the chronolgy of each interval uptime data considering SQF data
type BaktiSqfAvailability struct {
	Availability        float64
	LinkFailureDuration int64
	RestitutionDuration int64
	OpenDuration        int64
	Chronologies        []BaktiSqfChronology
	RainQuotaUsed       int64
}

// BaktiSqfChronology explains the chronolgy of each interval uptime data
type BaktiSqfChronology struct {
	Bakti1UptimeChronology
	SqfStatus int
	SqfValue  float64
	RainQuota int64
}

// CalcBaktiSqf returns the SLA and explains the status of service refers to uptime data from BAKTI series.
func CalcBaktiSqf(bakti1Chronologies []Bakti1UptimeChronology, sqfTimestamps []int64, sqfValues []float64, rainQuota int64) (*BaktiSqfAvailability, int64, error) {
	// SQF Value Validation
	if len(bakti1Chronologies) <= 0 || len(sqfTimestamps) <= 0 || len(sqfValues) <= 0 {
		return nil, -1, fmt.Errorf("one of inputted array is empty: %v, %v, %v", len(bakti1Chronologies), len(sqfTimestamps), len(sqfValues))
	}
	err := validateSqfValues(bakti1Chronologies, sqfTimestamps, sqfValues)
	if err != nil {
		return nil, -1, fmt.Errorf("failed on sqf values validation: %v", err)
	}
	chronologies := []BaktiSqfChronology{}
	for i, chronology := range bakti1Chronologies {
		newChronology := BaktiSqfChronology{
			SqfValue:  sqfValues[i],
			RainQuota: -1,
		}
		newChronology.StartTimestamps = chronology.StartTimestamps
		newChronology.EndTimestamps = chronology.EndTimestamps
		newChronology.UptimeValue = chronology.UptimeValue
		newChronology.Status = chronology.Status
		newChronology.LinkFailureDuration = chronology.LinkFailureDuration
		newChronology.RestitutionDuration = chronology.RestitutionDuration
		chronologies = append(chronologies, newChronology)
	}
	// Backup Rain Quota
	oldRainQuota := rainQuota
	for i := 0; i < len(chronologies); i++ {
		if chronologies[i].Status != BaktiLinkFailure {
			chronologies[i].RainQuota = rainQuota
		} else if chronologies[i].SqfValue >= 7.1 {
			chronologies[i].SqfStatus = SqfGte71
			chronologies[i].RestitutionDuration = 0
		} else if chronologies[i].SqfValue < 3 {
			chronologies[i].SqfStatus = Sqflt3
		} else if rainQuota > 0 {
			newRainQuota := rainQuota - chronologies[i].RestitutionDuration
			if newRainQuota >= 0 {
				rainQuota = newRainQuota
				chronologies[i].RestitutionDuration = 0
				chronologies[i].RainQuota = newRainQuota
				chronologies[i].SqfStatus = Sqfbt713Quota
			} else {
				chronologies[i].RestitutionDuration -= rainQuota
				chronologies[i].RainQuota = 0
				chronologies[i].SqfStatus = Sqfbt713Quota
				rainQuota = 0
			}
		} else {
			chronologies[i].SqfStatus = Sqfbt713NonQuota
			chronologies[i].RainQuota = 0
		}
	}
	// Calc Availability
	periodDuration := chronologies[len(chronologies)-1].EndTimestamps - chronologies[0].StartTimestamps
	var linkFailureDuration int64
	var restitutionDuration int64
	var openDuration int64
	for _, chronology := range chronologies {
		if chronology.Status == BaktiLinkFailure {
			linkFailureDuration += chronology.LinkFailureDuration
			// sqfQuota status is included to calculate while quota unable to cover all the restituion
			if chronology.SqfStatus == Sqfbt713NonQuota || chronology.SqfStatus == Sqflt3 || chronology.SqfStatus == Sqfbt713Quota {
				restitutionDuration += chronology.RestitutionDuration
			}
		} else if chronology.Status == BaktiOpen {
			openDuration += chronology.EndTimestamps - chronology.StartTimestamps
		}
	}
	availability := 1.0 - (float64(restitutionDuration+openDuration) / float64(periodDuration))
	baktiSqfAvailability := BaktiSqfAvailability{
		Availability:        availability,
		LinkFailureDuration: linkFailureDuration,
		RestitutionDuration: restitutionDuration,
		OpenDuration:        openDuration,
		Chronologies:        chronologies,
	}
	baktiSqfAvailability.RainQuotaUsed = oldRainQuota - rainQuota
	return &baktiSqfAvailability, rainQuota, nil
}

func validateSqfValues(bakti1UptimeChronologies []Bakti1UptimeChronology, sqfTimestamps []int64, sqfValues []float64) error {
	if len(sqfTimestamps) != len(sqfValues) {
		return fmt.Errorf("len of sqf timestamps and sqf values not same")
	}
	if len(bakti1UptimeChronologies) != len(sqfValues) {
		return fmt.Errorf("len of chronologies and sqf values not same")
	}
	for i := 0; i < len(bakti1UptimeChronologies); i++ {
		if bakti1UptimeChronologies[i].StartTimestamps != sqfTimestamps[i] {
			return fmt.Errorf("timestamps index %v of chronologies and sqf no same: %v and %v", i,
				bakti1UptimeChronologies[i].StartTimestamps, sqfTimestamps[i])
		}
	}
	return nil
}
