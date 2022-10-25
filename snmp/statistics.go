package snmp

import (
	"time"

	g "github.com/gosnmp/gosnmp"
)

func GetPower(host string) (float64, error) {
	params := &g.GoSNMP{
		Target:        host,
		Port:          161,
		Version:       g.Version3,
		SecurityModel: g.UserSecurityModel,
		MsgFlags:      g.NoAuthNoPriv,
		Timeout:       time.Duration(30) * time.Second,
		SecurityParameters: &g.UsmSecurityParameters{
			UserName: "cmgwatt",
		},
	}

	err := params.Connect()
	if err != nil {
		return 0, err
	}
	defer params.Conn.Close()

	oids := []string{PowerOID}

	elec, err := params.Get(oids)
	if err != nil {
		return 0, err
	}

	return float64(elec.Variables[0].Value.(int)) / 10, nil
}

func GetStatistics(host string) (*Measurement, error) {
	// build our own GoSNMP struct, rather than using g.Default
	params := &g.GoSNMP{
		Target:        host,
		Port:          161,
		Version:       g.Version3,
		SecurityModel: g.UserSecurityModel,
		MsgFlags:      g.NoAuthNoPriv,
		Timeout:       time.Duration(30) * time.Second,
		SecurityParameters: &g.UsmSecurityParameters{
			UserName: "cmgwatt",
		},
	}

	err := params.Connect()
	if err != nil {
		return nil, err
	}
	defer params.Conn.Close()

	oids := []string{CurrentOID, PowerOID, EnergyOID}

	elec, err := params.Get(oids)
	if err != nil {
		return nil, err
	}

	// calculate the average power since the epoch
	elapsed := time.Since(Epoch)
	avg := float64(elec.Variables[2].Value.(int)*100) / elapsed.Hours()

	// Place our electrical values into the data structure
	m := Measurement{
		Current:      float64(elec.Variables[0].Value.(int)) / 100,
		Power:        float64(elec.Variables[1].Value.(int)) / 10,
		Energy:       float64(elec.Variables[2].Value.(int)) / 100,
		AveragePower: avg,
	}

	return &m, nil
}
