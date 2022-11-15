// Package pdu handles querying of Power Delivery Units
package pdu

import (
	"errors"
	"log"

	"github.com/gosnmp/gosnmp"
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
		oids = append(oids, snmp.RittalWattage, snmp.RittalEnergy, snmp.RittalWattageAlt, snmp.RittalEnergyAlt)
	}

	// get power and energy
	elec, err := s.Get(oids)
	if err != nil {
		return 0, 0, err
	}

	// Hack to filter out nil values from the old Rittal PDUs
	pdu_values := make([]gosnmp.SnmpPDU, 0, len(elec.Variables))
	for _, val := range elec.Variables {
		switch val.Value.(type) {
		case string, int: // add your desired types which will fill newSlice
			pdu_values = append(pdu_values, val)
		}
	}

	power := float64(pdu_values[0].Value.(int)) / 10
	energy := float64(pdu_values[1].Value.(int)) / 10
	return power, energy, nil
}
