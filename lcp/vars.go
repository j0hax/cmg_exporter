package lcp

type LCPInfo struct {
	AvgTempIn     float64
	AvgTempOut    float64
	InBot         float64
	Fans          [6]uint64
	WaterTempIn   float64
	WaterTempOut  float64
	WaterFlowRate float64
}

// FanAvg returns the average percentage of all running fans
func (l *LCPInfo) FanAvg() float64 {
	total := float64(len(l.Fans))

	var sum uint64
	for _, v := range l.Fans {
		sum += v
	}

	return float64(sum) / total
}

// ActiveFans returns the number of fans that are running
func (l *LCPInfo) ActiveFans() int {
	var sum int

	for _, v := range l.Fans {
		if v > 0 {
			sum++
		}
	}

	return sum
}

// NOTE: all integer temperatures must be divided by 100

/* Device 2 / Variable 17: 'Air Temp.Server In.x.Value'
const TempInTop = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.17"
const TempInMid = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.26"
const TempInBot = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.35"
*/

// Device 2 / Variable 45: 'Air Temp.Server In.Average.Value'
const TempInAvg = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.45"

/* Device 2 / Variable 54: 'Air Temp.Server Out.x.Value'
const TempOutTop = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.54"
const TempOutMid = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.63"
const TempOutBot = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.72"
*/

// Device 2 / Variable 81: 'Air Temp.Server Out.Average.Value'
const TempOutAvg = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.81"

// Fan Value OIDs in percent
var FanSpeedOIDs = []string{
	"1.3.6.1.4.1.2606.7.4.2.2.1.11.2.90",
	"1.3.6.1.4.1.2606.7.4.2.2.1.11.2.94",
	"1.3.6.1.4.1.2606.7.4.2.2.1.11.2.98",
	"1.3.6.1.4.1.2606.7.4.2.2.1.11.2.102",
	"1.3.6.1.4.1.2606.7.4.2.2.1.11.2.106",
	"1.3.6.1.4.1.2606.7.4.2.2.1.11.2.110",
}

// Water Values
const WaterTempIn = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.135"
const WaterTempOut = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.144"
const WaterFlowRate = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.153"
