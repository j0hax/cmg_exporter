package snmp

type Manufacturer int

// Determine our PDU manufacturers
const (
	Unknown Manufacturer = iota
	Bachmann
	Rittal
)

// ModelOID allows us to diffentiate between PDU manufacturers
const ModelOID = "1.3.6.1.2.1.47.1.2.1.1.2.1"

// BachmannWattage represents the wattage of the first PDU on the inlet layer for a BlueNet2 PDU
const BachmannWattage = "1.3.6.1.4.1.31770.2.2.8.4.1.5.0.0.255.255.255.255.0.19"

// BachmannEnergy represents the energy in kWh measured since the device was activated
const BachmannEnergy = "1.3.6.1.4.1.31770.2.2.8.4.1.5.0.0.255.255.255.255.0.36"

// BachmannWattage represents the "total" wattage for a Rittal PDU
const RittalWattage = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.12"

// RittalEnergy represents the energy in kWh measured since the device was activated
const RittalEnergy = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.20"
