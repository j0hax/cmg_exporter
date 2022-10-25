# pdu-exporter

[![Go](https://github.com/j0hax/pdu-exporter/actions/workflows/go.yml/badge.svg)](https://github.com/j0hax/pdu-exporter/actions/workflows/go.yml)
[![Docker](https://github.com/j0hax/pdu-exporter/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/j0hax/pdu-exporter/actions/workflows/docker-publish.yml)

Prometheus Exporter for colocation PDUs

This exporter is _tailor-made_ for tracking power consumption of colocated servers at the Mechanical Engineering Campus of Leibniz University Hannover.

## Method of Operation

**⚠️ This software is experimental:** While its purpose is clear, usage may change substantially until a stable release.

An institute's server rack consists of two PDUs:
- a left one (ending in an even IP)
- a right one (ending in an odd IP)

This exporter
1. takes one of the two possible IPs,
2. calculates the complementing PDUs IP,
3. queries the wattage of both PDUs using SNMPv3
4. provides a prometheus metric

### Example

```console
$ curl 'http://pdu-exporter:1812/metrics?target=10.42.42.42'
pdu_left_power{instance="10.42.42.42"} 307
pdu_right_power{instance="10.42.42.43"} 148
pdu_total_power 455
```

## Installation

Add the following to `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: "pdu_export"
    static_configs:
      - targets:
        - 10.42.42.40 # List of targets to monitor.
        - 10.42.42.42 # Remember to include only every second PDU.
        - 10.42.42.44
        - 10.42.42.46
        - 10.42.42.48
      relabel_configs:
        - source_labels: [__address__]
          target_label: __param_target
        - source_labels: [__param_target]
          target_label: instance
        - target_label: __address__
          replacement: pdu-exporter:1812 # The SNMP exporter's real hostname:port.
```
