package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/buaazp/diskutil"
)

var (
	mfiutilPath   string
	listenAddress string
)

func init() {
	flag.StringVar(&mfiutilPath, "mfiutil-path", "/usr/sbin/mfiutil", "the mfiutil binary path")
	flag.StringVar(&listenAddress, "listen-address", ":9101", "server listen address")
}

func GetDiskStatus() (*diskutil.DiskStatus, error) {
	if _, err := os.Stat(mfiutilPath); err != nil || os.IsExist(err) {
		return nil, fmt.Errorf("%s is not exists", mfiutilPath)
	}

	return &diskutil.DiskStatus{
		MfiutilPath: mfiutilPath,
	}, nil
}
