package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/buaazp/diskutil"
	"github.com/gin-gonic/gin"
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

func main() {
	flag.Parse()

	ds, err := diskutil.NewDiskStatus(megaPath, adapterCount)
	if err != nil {
		log.Fatalf("DiskStatus New error: %v\n", err)
	}

	r := gin.Default()
	if true {
		gin.SetMode(gin.DebugMode)
	}

	r.GET("/physical-drive-stats", func(c *gin.Context) {
		if err := ds.Get(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {
			for _, ads := range ds.AdapterStats {
				c.JSON(http.StatusOK, ads.PhysicalDriveStats)
			}
		}
	})

	r.GET("/virtual-drive-stats", func(c *gin.Context) {
		if err := ds.Get(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {
			for _, ads := range ds.AdapterStats {
				c.JSON(http.StatusOK, ads.VirtualDriveStats)
			}
		}
	})

	r.GET("/log", func(c *gin.Context) {
		if logs, err := ds.Log(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {
			c.Header("Content-Type", "text/plain")
			for _, tmp := range logs {
				c.String(http.StatusOK, tmp)
			}
		}
	})

	_ = r.Run(listenAddress)
}
