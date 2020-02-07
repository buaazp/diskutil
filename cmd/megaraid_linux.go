package main

import (
	"flag"

	"github.com/buaazp/diskutil"
)

var (
	megaPath      string
	adapterCount  int
	listenAddress string
)

func init() {
	flag.StringVar(&megaPath, "mega-path", "/opt/MegaRAID/MegaCli/MegaCli64", "megaCli binary path")
	flag.IntVar(&adapterCount, "adapter", 0, "adapter count in your server")
	flag.StringVar(&listenAddress, "listen-address", "0.0.0.0:9101", "server listen address")
}

// GetDiskStatus for handling disk status object
func GetDiskStatus() (*diskutil.DiskStatus, error) {
	return diskutil.NewDiskStatus(diskutil.DiskStatusConfig{
		MegaCliPath:  megaPath,
		AdapterCount: adapterCount,
	})
}
