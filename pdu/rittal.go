package pdu

import (
	"errors"

	"github.com/gosnmp/gosnmp"
)

// GetRittalType determines the type of PDU we are dealing with:
//
// PDU-Controller, used by most institutes, or PDU-MAN, an older variant.
func GetRittalType(g *gosnmp.GoSNMP) (RittalType, error) {
	result, err := g.Get([]string{RittalTypeOID})
	if err != nil {
		return 0, err
	}

	v := string(result.Variables[0].Value.([]byte))

	switch v {
	case "PDU-Controller":
		return Controller, nil
	case "PDU-MAN":
		return Man, nil
	default:
		return 0, errors.New("could not determine if device is PDU-Controller or PDU-MAN")
	}
}

// GetControllerMetrics returns the power and energy of a Rittal PDU-Controller
func GetControllerMetrics(g *gosnmp.GoSNMP) (float64, float64, error) {
	result, err := g.Get([]string{RittalPower, RittalEnergy})
	if err != nil {
		return 0, 0, err
	}

	p := gosnmp.ToBigInt(result.Variables[0].Value).Uint64()
	e := gosnmp.ToBigInt(result.Variables[1].Value).Uint64()

	return float64(p), (float64(e) / 10), nil
}

// GetManMetrics returns the power and energy of a Rittal PDU-MAN
func GetManMetrics(g *gosnmp.GoSNMP) (float64, float64, error) {
	result, err := g.Get([]string{RittalAltPower, RittalAltEnergy})
	if err != nil {
		return 0, 0, err
	}

	p := gosnmp.ToBigInt(result.Variables[0].Value).Uint64()
	e := gosnmp.ToBigInt(result.Variables[1].Value).Uint64()

	return float64(p), (float64(e) / 10), nil
}

func GetRittalMetrics(g *gosnmp.GoSNMP) (float64, float64, error) {
	// Determine if the PDU is from Rittal or Bachmann
	m, err := GetRittalType(g)
	if err != nil {
		return 0, 0, err
	}

	switch m {
	case Controller:
		return GetControllerMetrics(g)
	case Man:
		return GetManMetrics(g)
	default:
		return 0, 0, errors.New("could not differentiate between Rittal PDU type")
	}
}
