// Package general represents generic information that can be requested via SNMP
package general

import (
	"fmt"
	"log"

	"github.com/VictoriaMetrics/metrics"
	"github.com/gosnmp/gosnmp"
	"github.com/j0hax/cmg_exporter/vars"
)

// RFC3418 defines SNMPv2-MIB::sysUpTime
const UptimeOID = ".1.3.6.1.2.1.1.3.0"

func getUptime(g *gosnmp.GoSNMP) (float64, error) {
	result, err := g.Get([]string{UptimeOID})
	if err != nil {
		return 0, err
	}

	// Divide by 100, as TimeTicks are represented as centiseconds
	return vars.ToFloat(result, 0) / 100, nil
}

func Handler(g *gosnmp.GoSNMP, unit string) {
	s := fmt.Sprintf(`device_uptime_seconds{unit="%s"}`, unit)

	uptime, err := getUptime(g)
	if err != nil {
		log.Panic(err)
	}

	cnt := metrics.NewFloatCounter(s)

	cnt.Set(uptime)
}
