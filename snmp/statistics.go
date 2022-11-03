package snmp

import (
	"errors"
	"strings"

	g "github.com/gosnmp/gosnmp"
)

// Determine the manufacturer of a PDU via SNMP
func GetMfG(params *g.GoSNMP) (Manufacturer, error) {
	oids := []string{ModelOID}

	dat, err := params.Get(oids)
	if err != nil {
		return 0, err
	}

	// Make a string of our data
	model := string(dat.Variables[0].Value.([]byte))

	// Check if the value is known
	if strings.Contains(model, "BlueNet2") {
		return Bachmann, nil
	} else if strings.Contains(model, "Rittal") {
		return Rittal, nil
	}

	return 0, errors.New("could not determine hardware")
}

// GetStatistics requests energy and power via SNMP.
func GetStatistics(params *g.GoSNMP) (float64, float64, error) {
	mfg, err := GetMfG(params)
	if err != nil {
		return 0, 0, err
	}

	oids := []string{}

	if mfg == Bachmann {
		oids = append(oids, BachmannWattage, BachmannEnergy)
	} else if mfg == Rittal {
		oids = append(oids, RittalWattage, RittalEnergy)
	}

	// get power
	elec, err := params.Get(oids)
	if err != nil {
		return 0, 0, err
	}

	power := float64(elec.Variables[0].Value.(int)) / 10
	energy := float64(elec.Variables[1].Value.(int)) / 10
	return power, energy, nil
}
