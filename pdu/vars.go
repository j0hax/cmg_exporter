package pdu

type Manufacturer int

// Enumerates the types of PDU manufactuerers
const (
	Bachmann Manufacturer = iota
	Rittal
)

type RittalType int

// Enumerates the types of PDU manufactuerers
const (
	Man RittalType = iota
	Controller
)

// RittalTypeOID allows for a differentiation between the two types of Rital PDUs
//
// The two options are "PDU-Controller" and "PDU-MAN"
const RittalTypeOID = ".1.3.6.1.4.1.2606.7.4.1.2.1.2.1"

// BachmannPower represents the wattage of the first PDU on the inlet layer for a BlueNet2 PDU
const BachmannPower = ".1.3.6.1.4.1.31770.2.2.8.4.1.5.0.0.255.255.255.255.0.19"

// BachmannEnergy represents the energy in kWh measured since the device was activated
const BachmannEnergy = ".1.3.6.1.4.1.31770.2.2.8.4.1.5.0.0.255.255.255.255.0.36"

// RittalControllerPower represents the "total" wattage for a Rittal PDU-Controller
const RittalPower = ".1.3.6.1.2.1.99.1.1.1.4.101003"

// RittalManPower represents the "total" wattage for a Rittal PDU-MAN
const RittalAltPower = ".1.3.6.1.2.1.99.1.1.1.4.1003"

// RittalControllerEnergy represents the energy in kWh measured since the device was activated
const RittalEnergy = ".1.3.6.1.2.1.99.1.1.1.4.101004"

// RittalManEnergy is an alternative OID to RittalEnergy for older devices.
const RittalAltEnergy = ".1.3.6.1.2.1.99.1.1.1.4.1004"
