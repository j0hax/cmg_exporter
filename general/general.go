// Package general represents generic information that can be requested via SNMP
package general

import (
	"fmt"
	"log"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/beevik/ntp"
	"github.com/gosnmp/gosnmp"
	"github.com/j0hax/cmg_exporter/vars"
)

// RFC3418 defines SNMPv2-MIB::sysUpTime
const SystemUptimeOID = ".1.3.6.1.2.1.1.3.0"

// RFC2790 defines HOST-RESOURCES-MIB::hrSystemDate.0
const SystemDateOID = ".1.3.6.1.2.1.25.1.2.0"

// RFC2790 defines HOST-RESOURCES-MIB::hrSystemProcesses.0
const SystemProcessesOID = ".1.3.6.1.2.1.25.1.6.0"

// The preferred NTP Server for calculating the differences in clock times
var NTPServer = "time1.uni-hannover.de"

func getUptime(g *gosnmp.GoSNMP) (float64, error) {
	result, err := g.Get([]string{SystemUptimeOID})
	if err != nil {
		return 0, err
	}

	// Divide by 100, as TimeTicks are represented as centiseconds
	return vars.ToFloat(result, 0) / 100, nil
}

func getNumProcesses(g *gosnmp.GoSNMP) (uint64, error) {
	result, err := g.Get([]string{SystemProcessesOID})
	if err != nil {
		return 0, err
	}

	// Divide by 100, as TimeTicks are represented as centiseconds
	return gosnmp.ToBigInt(result.Variables[0].Value).Uint64(), nil
}

// SnmpTimeString parses SNMP Data into a RFC3339 compatible timestamp
func SnmpTimeString(c []byte) string {
	year := (int(c[0]) << 8) | int(c[1])
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d%c%02d:%02d", year, c[2], c[3], c[4], c[5], c[6], c[8], c[9], c[10])
}

// getTimeDeviation calculates the difference between g's onboard clock and the time provided by an NTP server
func getTimeDeviation(g *gosnmp.GoSNMP) (float64, error) {
	refTime, err := ntp.Time(NTPServer)
	if err != nil {
		return 0, err
	}

	result, err := g.Get([]string{SystemDateOID})
	if err != nil {
		return 0, err
	}

	dateStr := SnmpTimeString(result.Variables[0].Value.([]byte))

	sysTime, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return 0, err
	}

	diff := sysTime.Sub(refTime)

	return diff.Seconds(), nil
}

func Handler(g *gosnmp.GoSNMP, unit string) {
	s := fmt.Sprintf(`device_uptime_seconds{unit="%s"}`, unit)
	uptime, err := getUptime(g)
	if err != nil {
		log.Panic(err)
	}

	cnt := metrics.NewFloatCounter(s)
	cnt.Set(uptime)

	d, err := getTimeDeviation(g)
	if err != nil {
		log.Panic(err)
	}

	s = fmt.Sprintf(`device_clock_drift_seconds{unit="%s"}`, unit)
	metrics.NewGauge(s, func() float64 {
		return d
	})

	p, err := getNumProcesses(g)
	if err != nil {
		log.Panic(err)
	}

	s = fmt.Sprintf(`device_num_processes{unit="%s"}`, unit)
	metrics.NewGauge(s, func() float64 {
		return float64(p)
	})
}
