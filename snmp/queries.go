package snmp

import (
	"errors"
	"strings"
)

// Determine the manufacturer of a PDU via SNMP
func GetInfo(target string) (Device, Manufacturer, error) {
	s, err := Connect(target)
	if err != nil {
		return 0, 0, err
	}
	defer s.Conn.Close()

	oids := []string{ModelOID}

	dat, err := s.Get(oids)
	if err != nil {
		return 0, 0, err
	}

	// Make a string of our data
	model := string(dat.Variables[0].Value.([]byte))

	// Check if the value is known
	if strings.Contains(model, "BlueNet2") {
		return PDU, Bachmann, nil
	} else if strings.Contains(model, "Rittal") {
		if strings.Contains(model, "PDU") {
			return PDU, Rittal, nil
		} else if strings.Contains(model, "LCP") {
			return LCP, Rittal, nil
		}
	}

	return 0, 0, errors.New("could not determine hardware")
}
