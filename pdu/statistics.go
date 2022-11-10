// Package pdu handles querying of Power Delivery Units
package pdu

import (
	"errors"
	"log"

	"github.com/j0hax/cmg_exporter/snmp"
)

// GetStatistics requests energy and power via SNMP.
func GetStatistics(target string) (float64, float64, error) {

	// Determine manufacturer of target
	dev, mfg, err := snmp.GetInfo(target)
	if err != nil {
		log.Print(err)
	}

	if dev != snmp.PDU {
		return 0, 0, errors.New("device is not a recognized PDU")
	}

	// connect
	s, err := snmp.Connect(target)
	if err != nil {
		return 0, 0, err
	}
	defer s.Conn.Close()

	// determine appropriate OIDs
	oids := []string{}
	if mfg == snmp.Bachmann {
		oids = append(oids, snmp.BachmannWattage, snmp.BachmannEnergy)
	} else if mfg == snmp.Rittal {
		oids = append(oids, snmp.RittalWattage, snmp.RittalEnergy)
	}

	// get power and energy
	elec, err := s.Get(oids)
	if err != nil {
		return 0, 0, err
	}

	power := float64(elec.Variables[0].Value.(int)) / 10
	energy := float64(elec.Variables[1].Value.(int)) / 10
	return power, energy, nil
}
