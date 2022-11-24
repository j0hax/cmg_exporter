package vars

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
