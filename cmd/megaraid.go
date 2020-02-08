package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	flag.Parse()

	ds, err := GetDiskStatus()
	if err != nil {
		log.Fatalf("DiskStatus error: %v\n", err)
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

	r.GET("/broken/:type/", func(c *gin.Context) {
		if err := ds.Get(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {
			if volumes, drives, err := ds.ListBrokenDrive(); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			} else {
				switch strings.ToLower(c.Param("type")) {
				case "volumes":
					c.JSON(http.StatusOK, volumes)
				case "drives":
					c.JSON(http.StatusOK, drives)
				default:
					c.JSON(http.StatusNotFound, nil)
				}
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

	// for prometheus exporter
	prometheus.MustRegister(NewMegaCollector(ds))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	_ = r.Run(listenAddress)
}
