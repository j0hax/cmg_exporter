package vars

import "github.com/gosnmp/gosnmp"

type Device int

// Enumerates the types of devices.
//
// PDUController and PDUMan are PDUs manufactured by Rittal.
// CMCIII-PU LCP is a liquid cooling package by Rittal.
// BlueNet2 is a PDU manufactured by Bachmann.
const (
	Unknown Device = iota
	Pdu
	Lcp
)

// TypeOID allows for a diffentiation between manufacturers and device types
const TypeOID = ".1.3.6.1.2.1.47.1.2.1.1.2.1"

// Small helper function to convert an SNMP Result into a usable float
func ToFloat(result *gosnmp.SnmpPacket, index int) float64 {
	pdu := result.Variables[index]
	i := gosnmp.ToBigInt(pdu.Value).Uint64()
	return float64(i)
}
