package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"

	"github.com/VictoriaMetrics/metrics"
	"github.com/gosnmp/gosnmp"
	"github.com/j0hax/cmg_exporter/lcp"
	"github.com/j0hax/cmg_exporter/pdu"
	"github.com/j0hax/cmg_exporter/vars"
)

// GetName returns the device name (prefix and number) via reverse DNS lookup.
// This is used to determine the rack or unit name.
//
// For example, "s02-pdu-links.serverraum.mgmt.mb.uni-hannover.de" results in "s02", and
// lcp1.serverraum.mgmt.mb.uni-hannover.de" results in "lcp1"
func GetName(target string) (string, error) {
	name, err := net.LookupAddr(target)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("[a-zA-Z]+[0-9]+")
	return re.FindString(name[0]), nil
}

// GetType determines if the device is a PDU or LCP
func GetType(g *gosnmp.GoSNMP) (vars.Device, error) {
	result, err := g.Get([]string{vars.TypeOID})
	if err != nil {
		return vars.Unknown, err
	}

	v := string(result.Variables[0].Value.([]byte))

	switch v {
	case "Rittal PDU":
		return vars.Pdu, nil
	case "BlueNet2":
		return vars.Pdu, nil
	case "Rittal LCP":
		return vars.Lcp, nil
	default:
		return vars.Unknown, errors.New("could not determine device type")
	}
}

// Handler is the basic entrypoint for querying devices.
func Handler(w http.ResponseWriter, req *http.Request) {
	// parse URL parameters
	u, err := url.Parse(req.RequestURI)
	if err != nil {
		log.Print(err)
		return
	}

	// Determine target IP address
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Print(err)
		return
	}

	target := m["target"][0]
	log.Printf("%s requests information on %s", req.RemoteAddr, target)

	g, err := Connect(target)
	if err != nil {
		log.Print(err)
		return
	}

	t, err := GetType(g)
	if err != nil {
		log.Print(err)
		return
	}

	name, err := GetName(target)
	if err != nil {
		log.Print(err)
		return
	}

	switch t {
	case vars.Pdu:
		pdu.Handler(g, name)
	case vars.Lcp:
		lcp.Handler(g, name)
	}

	metrics.WritePrometheus(w, false)
	metrics.UnregisterAllMetrics()
}
func main() {
	// Expose the registered metrics at `/metrics` path.
	http.HandleFunc("/metrics", Handler)
	http.ListenAndServe(":1812", nil)
	log.Print("PDU Exporter running.")
}
