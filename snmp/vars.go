package snmp

// Device allows us to keep track of devices.
//
// It influences how statistics are queried
type Device int

// Determine our types of devices
//
// PDUs are Power Distribution Units,
// LCPs are Liquiod Cooling Packages.
const (
	PDU Device = iota
	LCP
)

// Manufacturer allows us to keep track of manufacturers.
//
// It influences which specific OIDs are used to query common parameters
type Manufacturer int

// Determine our equipment manufacturers
const (
	Bachmann Manufacturer = iota
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

// RittalWattageAlt is an alternative OID to RittalWattage for older devices
//
// Unfortunately, it only differs by one decimal
const RittalWattageAlt = "1.3.6.1.4.1.2606.7.4.2.2.1.11.1.12"

// RittalEnergy represents the energy in kWh measured since the device was activated
const RittalEnergy = "1.3.6.1.4.1.2606.7.4.2.2.1.11.2.20"

// RittalEnergyAlt is an alternative OID to RittalEnergy for older devices.
//
// Unfortunately, it only differs by one decimal
const RittalEnergyAlt = "1.3.6.1.4.1.2606.7.4.2.2.1.11.1.20"
