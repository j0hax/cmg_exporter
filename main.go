package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/VictoriaMetrics/metrics"
	"github.com/j0hax/cmg_exporter/pdu"
	"github.com/j0hax/cmg_exporter/snmp"
)

// Handler is the basic entrypoint for querying devices.
func Handler(w http.ResponseWriter, req *http.Request) {
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

	// Determine target IP address
	target := m["target"][0]
	log.Printf("%s requests information on %s", req.RemoteAddr, target)

	// Determine manufacturer of target
	dev, _, err := snmp.GetInfo(target)
	if err != nil {
		errmesg := fmt.Sprintf("# ERROR: %s\r\n", err)
		w.Write([]byte(errmesg))
		log.Print(err)
	}

	if dev == snmp.PDU {
		pdu.Handler(target)
	} else if dev == snmp.LCP {

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
