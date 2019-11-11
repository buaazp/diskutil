package main

import (
	"github.com/buaazp/diskutil"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type MegaCollector struct {
	status       *diskutil.DiskStatus
	virtualDesc  *prometheus.Desc
	physicalDesc *prometheus.Desc
}

func (collector *MegaCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.virtualDesc
	ch <- collector.physicalDesc
}

func (collector *MegaCollector) Collect(ch chan<- prometheus.Metric) {
	if err := collector.status.Get(); err != nil {
		return
	}

	for _, ads := range collector.status.AdapterStats {
		for _, vd := range ads.VirtualDriveStats {
			ch <- prometheus.MustNewConstMetric(collector.virtualDesc, prometheus.GaugeValue,
				float64(vd.VirtualDrive), strconv.Itoa(vd.NumberOfDrives), vd.State, vd.Size, vd.RAIDLevel)
		}

		for _, pv := range ads.PhysicalDriveStats {
			ch <- prometheus.MustNewConstMetric(collector.physicalDesc, prometheus.GaugeValue,
				float64(pv.DeviceId), pv.SerialNumber, strconv.Itoa(pv.SlotNumber), pv.Brand, pv.Model, pv.RawSize, pv.FirmwareState)
		}
	}
}

func NewMegaCollector(ds *diskutil.DiskStatus) *MegaCollector {
	return &MegaCollector{
		status: ds,
		virtualDesc: prometheus.NewDesc("mega_virtual_drives", "Virtual drive status",
			[]string{"number_of_drives", "state", "size", "raid_level"}, nil),
		physicalDesc: prometheus.NewDesc("mega_physical_drives", "Physical drive status",
			[]string{"serial", "slot", "brand", "model", "size", "state"}, nil),
	}
}
