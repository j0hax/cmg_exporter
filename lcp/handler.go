package lcp

import (
	"fmt"

	"github.com/VictoriaMetrics/metrics"
	"github.com/gosnmp/gosnmp"
	"github.com/j0hax/cmg_exporter/vars"
)

// GetMetrics returns metrics for a specific LCP
func GetMetrics(g *gosnmp.GoSNMP) (*LCPInfo, error) {

	// Grab general information
	result, err := g.Get([]string{TempInAvg, TempInSetPoint, TempOutAvg, WaterTempIn, WaterTempOut, WaterFlowRate, WaterValve})
	if err != nil {
		return nil, err
	}

	lcp := &LCPInfo{
		AvgTempIn:      vars.ToFloat(result, 0) / 100,
		TempInSetPoint: vars.ToFloat(result, 1) / 100,
		AvgTempOut:     vars.ToFloat(result, 2) / 100,
		WaterTempIn:    vars.ToFloat(result, 3) / 100,
		WaterTempOut:   vars.ToFloat(result, 4) / 100,
		WaterFlowRate:  vars.ToFloat(result, 5) / 10,
		WaterValve:     gosnmp.ToBigInt(result.Variables[6].Value).Uint64(),
	}

	// Grab the 6 fan speeds
	result, err = g.Get(FanSpeedOIDs)
	if err != nil {
		return nil, err
	}

	for i, r := range result.Variables {
		lcp.Fans[i] = gosnmp.ToBigInt(r.Value).Uint64()
	}

	return lcp, nil
}

// Handler collects data on a PDU and registers power and energy metrics.
func Handler(g *gosnmp.GoSNMP, unit string) {
	lcp, err := GetMetrics(g)
	if err != nil {
		fmt.Print(err)
		return
	}

	s := fmt.Sprintf(`lcp_air_temp_in{unit="%s"}`, unit)
	metrics.NewGauge(s, func() float64 {
		return lcp.AvgTempIn
	})

	s = fmt.Sprintf(`lcp_air_temp_setpoint{unit="%s"}`, unit)
	metrics.NewGauge(s, func() float64 {
		return lcp.TempInSetPoint
	})

	s = fmt.Sprintf(`lcp_air_temp_out{unit="%s"}`, unit)
	metrics.NewGauge(s, func() float64 {
		return lcp.AvgTempOut
	})

	s = fmt.Sprintf(`lcp_water_temp_in{unit="%s"}`, unit)
	metrics.NewGauge(s, func() float64 {
		return lcp.WaterTempIn
	})

	s = fmt.Sprintf(`lcp_water_temp_out{unit="%s"}`, unit)
	metrics.NewGauge(s, func() float64 {
		return lcp.WaterTempOut
	})

	s = fmt.Sprintf(`lcp_water_flow_rate{unit="%s"}`, unit)
	metrics.NewGauge(s, func() float64 {
		return lcp.WaterFlowRate
	})

	s = fmt.Sprintf(`lcp_water_valve{unit="%s"}`, unit)
	metrics.NewGauge(s, func() float64 {
		return float64(lcp.WaterValve)
	})

	s = fmt.Sprintf(`lcp_fan_avg{unit="%s"}`, unit)
	metrics.NewGauge(s, lcp.FanAvg)
}
