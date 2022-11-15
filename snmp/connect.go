// Package snmp provides common queries and variables for SNMP operations
package snmp

import (
	"time"

	g "github.com/gosnmp/gosnmp"
)

func Connect(host string) (*g.GoSNMP, error) {
	params := &g.GoSNMP{
		Target:    host,
		Port:      161,
		Community: "public",
		Version:   g.Version2c,
		Timeout:   time.Duration(1) * time.Second,
		Retries:   5,
	}

	err := params.Connect()
	if err != nil {
		return nil, err
	}

	return params, nil
}
