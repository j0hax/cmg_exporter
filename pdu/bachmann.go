package pdu

import (
	"github.com/gosnmp/gosnmp"
)

func GetBachmannMetrics(g *gosnmp.GoSNMP) (float64, float64, error) {
	result, err := g.Get([]string{BachmannPower, BachmannEnergy})
	if err != nil {
		return 0, 0, err
	}

	p := gosnmp.ToBigInt(result.Variables[0].Value).Uint64()
	e := gosnmp.ToBigInt(result.Variables[1].Value).Uint64()

	return (float64(p) / 10), (float64(e) / 10), nil
}
