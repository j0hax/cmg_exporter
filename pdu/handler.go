package pdu

import (
	"fmt"
	"log"
	"net"

	"github.com/VictoriaMetrics/metrics"
)

func Handler(target string) {
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

	// Gather statistics of both PDUs
	lPower, lEnergy, err := GetStatistics(leftIP.String())
	if err != nil {
		log.Print(err)
		return
	}

	rPower, rEnergy, err := GetStatistics(rightIP.String())
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
	metrics.NewGauge(s, func() float64 {
		return lPower + rPower
	})

	s = fmt.Sprintf(`pdu_left_power{rack="%s"}`, rack)
	metrics.NewGauge(s, func() float64 {
		return lPower
	})

	s = fmt.Sprintf(`pdu_right_power{rack="%s"}`, rack)
	metrics.NewGauge(s, func() float64 {
		return rPower
	})

	s = fmt.Sprintf(`pdu_total_energy{rack="%s"}`, rack)
	c := metrics.NewFloatCounter(s)
	c.Set(lEnergy + rEnergy)

	s = fmt.Sprintf(`pdu_left_energy{rack="%s"}`, rack)
	c = metrics.NewFloatCounter(s)
	c.Set(lEnergy)

	s = fmt.Sprintf(`pdu_right_energy{rack="%s"}`, rack)
	c = metrics.NewFloatCounter(s)
	c.Set(rEnergy)
}
