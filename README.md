# cmg_exporter

[![Go](https://github.com/j0hax/pdu-exporter/actions/workflows/go.yml/badge.svg)](https://github.com/j0hax/pdu-exporter/actions/workflows/go.yml)
[![Docker](https://github.com/j0hax/pdu-exporter/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/j0hax/pdu-exporter/actions/workflows/docker-publish.yml)

Prometheus Exporter for colocation PDUs

This exporter is _tailor-made_ for tracking power consumption of colocated servers at the Mechanical Engineering Campus of Leibniz University Hannover.

## Method of Operation

**⚠️ This software is experimental:** While its purpose is clear, usage may change substantially until a stable release.

At the moment this service provides power and energy statistics for Rittal PDU-Controller, PDU-Man and Bachmann BlueNet2 PDUs. Rittal LCP support coming soon.

### Statistics reported

To avoid confusion keep in mind:

> ...a kilowatt is a unit of power but a kilowatt-hour (1 kilowatt times 1 hour) is a unit of energy. 

More information can be found at the Website for [Energy Education](https://energyeducation.ca/encyclopedia/Energy_vs_power).

#### Energy

The _total_ wattage drawn through the PDU in its service time, in kWh.

#### Power

The current wattage being drawn through the PDU.

### Example

```console
$ curl 'http://pdu-exporter:1812/metrics?target=10.42.42.42'
pdu_total_energy{rack="s12"} 16554.4
pdu_total_power{rack="s12"} 707
```

## Installation

Add the following to `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: "pdu_export"
    static_configs:
      - targets:
        - 10.42.42.40 # List of targets to monitor.
        - 10.42.42.41
        - 10.42.42.42
        - 10.42.42.43
        - 10.42.42.44
        - 10.42.42.45
        # etc.
      relabel_configs:
        - source_labels: [__address__]
          target_label: __param_target
        - source_labels: [__param_target]
          target_label: instance
        - target_label: __address__
          replacement: pdu-exporter:1812 # The SNMP exporter's real hostname:port.
      metric_relabel_configs:
      - source_labels: [rack]
        regex: s01 # The PDU Exporter gathers Rack number from the hostname.
        target_label: institute
        replacement: IMES # This can be used to assign entities to a rack via RegEx
      - source_labels: [rack]
        regex: s(02|12|21)
        target_label: institute
        replacement: IDS/IKM

```
