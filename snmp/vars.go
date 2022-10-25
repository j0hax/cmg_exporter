package snmp

import (
	"fmt"
	"time"
)

// InletOID represents all electric readings of the first PDU on the inlet layer
const InletOID = "1.3.6.1.4.1.31770.2.2.8.4.1.5.0.0.255.255.255.255.0"

// The amperes being drawn from the PDU on the Inlet
const CurrentOID = InletOID + ".4"

// The wattage being drawn from the PDU on the Inlet
const PowerOID = InletOID + ".19"

// The sum of watts drawn since the device was installed
const EnergyOID = InletOID + ".36"

// Epoch represents the approximate time servers were started at the Garbsen facility.
//
// In this case, it is April 1st, 2020.
var Epoch = time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC)

// Represents the measurement of a PDU at a single point in time
type Measurement struct {
	Current      float64
	Power        float64
	Energy       float64
	AveragePower float64
}

func (m Measurement) String() string {
	return fmt.Sprintf(
		"Current: %.2f A, Power: %.2f W, Avg: %.2f Wh",
		m.Current,
		m.Power,
		m.AveragePower,
	)
}
