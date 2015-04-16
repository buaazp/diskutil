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

	fmt.Println(ds)
}
