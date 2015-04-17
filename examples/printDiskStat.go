package main

import (
	"flag"
	"fmt"
	"github.com/buaazp/diskutil"
	"os"
)

var (
	megaPath     string
	adapterCount int
)

func init() {
	flag.StringVar(&megaPath, "mega-path", "/opt/MegaRAID/MegaCli/MegaCli64", "megaCli binary path")
	flag.IntVar(&adapterCount, "adapter-count", 1, "adapter count in your server")
}

func main() {
	flag.Parse()
	ds, err := diskutil.NewDiskStatus(megaPath, adapterCount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DiskStatus New error: %v\n", err)
		return
	}

	err = ds.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DiskStatus Get error: %v\n", err)
		return
	}

	for i, ads := range ds.AdapterStats {
		fmt.Printf("adapter #%d \n", i)
		for j, pds := range ads.PhysicalDriveStats {
			pdStatus := pds.FirmwareState
			fmt.Printf("PD%d status: %s\n", j, pdStatus)
		}
		fmt.Printf("\n")
	}

	jsonStatus, err := ds.ToJson()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DiskStatus ToJson error: %v\n", err)
		return
	}
	fmt.Println(jsonStatus)
}
