package main

import (
	"time"

	"github.com/gosnmp/gosnmp"
)

// Connect creates a unique connection to a host with standard parameters
func Connect(host string) (*gosnmp.GoSNMP, error) {
	params := &gosnmp.GoSNMP{
		Target:    host,
		Port:      161,
		Community: "public",
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(1) * time.Second,
		Retries:   5,
	}

	err := params.Connect()
	if err != nil {
		return nil, err
	}

	return params, nil
}
