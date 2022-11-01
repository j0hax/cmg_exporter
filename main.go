package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	g "github.com/gosnmp/gosnmp"

	"github.com/j0hax/pdu-exporter/snmp"

	"github.com/VictoriaMetrics/metrics"
)

func connectSNMP(host string) (*g.GoSNMP, error) {
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

func handler(w http.ResponseWriter, req *http.Request) {
	// parse URL parameters
	u, err := url.Parse(req.RequestURI)
	if err != nil {
		log.Print(err)
		return
	}

	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Print(err)
		return
	}

	target := m["target"][0]

	log.Printf("%s requests information on %s", req.RemoteAddr, target)

	leftIP := net.ParseIP(target).To4()
	if leftIP == nil {
		log.Print("Target is not a valid IP adress")
		return
	}

	// Ensure the IP is really the one of a left rack (last byte should be even)
	if leftIP[3]%2 == 1 {
		leftIP[3]--
	}

	// IP of right PDU is simply the last field increased by 1
	rightIP := make(net.IP, len(leftIP))
	copy(rightIP, leftIP)
	rightIP[3]++

	// Create the struct describing our SNMP state
	leftSNMP, err := connectSNMP(leftIP.String())
	if err != nil {
		log.Print(err)
	}

	defer leftSNMP.Conn.Close()

	// Gather statistics of both PDUs
	left, err := snmp.GetPower(leftSNMP)
	if err != nil {
		log.Print(err)
		return
	}

	// Create the struct describing our SNMP state
	rightSNMP, err := connectSNMP(leftIP.String())
	if err != nil {
		log.Print(err)
	}

	defer rightSNMP.Conn.Close()

	right, err := snmp.GetPower(rightSNMP)
	if err != nil {
		log.Print(err)
		return
	}

	// Get the name of the rack
	name, err := net.LookupAddr(rightIP.String())
	if err != nil {
		log.Print(err)
		return
	}

	rack := name[0][0:3]

	s := fmt.Sprintf(`pdu_total_power{rack="%s"}`, rack)
	metrics.GetOrCreateGauge(s, func() float64 {
		return left + right
	})

	s = fmt.Sprintf(`pdu_left_power{rack="%s"}`, rack)
	metrics.GetOrCreateGauge(s, func() float64 {
		return left
	})

	s = fmt.Sprintf(`pdu_right_power{rack="%s"}`, rack)
	metrics.GetOrCreateGauge(s, func() float64 {
		return right
	})

	metrics.WritePrometheus(w, true)
	metrics.UnregisterAllMetrics()
}

func main() {
	// Expose the registered metrics at `/metrics` path.
	http.HandleFunc("/metrics", handler)
	http.ListenAndServe(":1812", nil)
	log.Print("PDU Exporter running.")
}
