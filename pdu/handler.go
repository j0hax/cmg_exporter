package pdu

import (
	"errors"
	"fmt"

	"github.com/VictoriaMetrics/metrics"
	"github.com/gosnmp/gosnmp"
	"github.com/j0hax/cmg_exporter/vars"
)

// GetManufacturer determines the manufactuerer of a PDU
func GetManufacturer(g *gosnmp.GoSNMP) (Manufacturer, error) {
	result, err := g.Get([]string{vars.TypeOID})
	if err != nil {
		return 0, err
	}

	v := string(result.Variables[0].Value.([]byte))

	switch v {
	case "Rittal PDU":
		return Rittal, nil
	case "BlueNet2":
		return Bachmann, nil
	default:
		return 0, errors.New("could not determine device type")
	}
}

// GetMetrics returns the power and energy using manufacturer/model specific methods
func GetMetrics(g *gosnmp.GoSNMP) (float64, float64, error) {
	// Determine if the PDU is from Rittal or Bachmann
	m, err := GetManufacturer(g)
	if err != nil {
		return 0, 0, err
	}

	switch m {
	case Bachmann:
		return GetBachmannMetrics(g)
	case Rittal:
		return GetRittalMetrics(g)
	default:
		return 0, 0, errors.New("could not determine PDU manufacturer")
	}
}

// Handler collects data on a PDU and registers power and energy metrics.
func Handler(g *gosnmp.GoSNMP, rack string) {
	p, e, err := GetMetrics(g)
	if err != nil {
		fmt.Print(err)
		return
	}

	s := fmt.Sprintf(`pdu_total_power{rack="%s"}`, rack)
	metrics.NewGauge(s, func() float64 {
		return p
	})

	s = fmt.Sprintf(`pdu_total_energy{rack="%s"}`, rack)
	c := metrics.NewFloatCounter(s)
	c.Set(e)
}
