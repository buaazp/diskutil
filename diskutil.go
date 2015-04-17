package diskutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	keyExitResult               string = "Exit Code:"
	keyVdVirtualDrive           string = "Virtual Drive"
	keyVdTargetId               string = "Target Id"
	keyVdName                   string = "Name"
	keyVdSize                   string = "Size"
	keyVdState                  string = "State"
	keyVdNumberOfDrives         string = "Number Of Drives"
	keyVdEncryptiontype         string = "Encryption type"
	keyPdEnclosureDeviceId      string = "Enclosure Device ID"
	keyPdSlotNumber             string = "Slot Number"
	keyPdDeviceId               string = "Device Id"
	keyPdMediaErrorCount        string = "Media Error Count"
	keyPdOtherErrorCount        string = "Other Error Count"
	keyPdPredictiveFailureCount string = "Predictive Failure Count"
	keyPdPdtype                 string = "PD type"
	keyPdRawSize                string = "Raw Size"
	keyPdFirmwareState          string = "Firmware state"
	keyPdInquiryData            string = "Inquiry Data"
	keyPdDriveTemperature       string = "Drive Temperature"

	typeString int = iota
	typeInt
	typeUint64
)

// DiskStatus is a struct to get all Adapters' Stat of the server
type DiskStatus struct {
	megacliPath  string
	adapterCount int
	AdapterStats []AdapterStat `json:"adapter_stats"`
}

// String() is used to get the print string.
func (d *DiskStatus) String() string {
	data, err := json.Marshal(d)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// ToJson() is used to get the json encoded string.
func (d *DiskStatus) ToJson() (string, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// NewDiskStatus() use the megaCliPath and apapterCount to build a DiskStatus.
func NewDiskStatus(megaCliPath string, adapterCount int) (*DiskStatus, error) {
	megaCliPath = path.Clean(megaCliPath)
	if !fileExist(megaCliPath) {
		return nil, errors.New("megaCli not exist")
	}
	ds := new(DiskStatus)
	ds.megacliPath = megaCliPath
	ds.adapterCount = adapterCount
	return ds, nil
}

func execCmd(command, args string) (string, error) {
	// fmt.Println("Command: ", command)
	// fmt.Println("Arguments: ", args)
	var argArray []string
	if args != "" {
		argArray = strings.Split(args, " ")
	} else {
		argArray = make([]string, 0)
	}

	cmd := exec.Command(command, argArray...)
	buf, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "The command failed to perform: %s (Command: %s, Arguments: %s)", err, command, args)
		return "", err
	}

	// fmt.Fprintf(os.Stdout, "Result: %s", buf)
	return string(buf), nil
}

// Get() is used to get all the stat of a DiskStatus.
func (d *DiskStatus) Get() error {
	ads := make([]AdapterStat, 0)

	command := d.megacliPath
	for i := 0; i < d.adapterCount; i++ {
		ad := AdapterStat{
			AdapterId: i,
		}
		err := ad.getMegaRaidVdInfo(command)
		if err != nil {
			d.AdapterStats = nil
			return err
		}
		err = ad.getMegaRaidPdInfo(command)
		if err != nil {
			d.AdapterStats = nil
			return err
		}
		ads = append(ads, ad)
	}

	d.AdapterStats = ads
	return nil
}

// GetVirtualDrive() is used to get the VirtualDriveStat of a DiskStatus.
func (d *DiskStatus) GetVirtualDrive() error {
	ads := make([]AdapterStat, 0)

	command := d.megacliPath
	for i := 0; i < d.adapterCount; i++ {
		ad := AdapterStat{
			AdapterId: i,
		}
		err := ad.getMegaRaidVdInfo(command)
		if err != nil {
			d.AdapterStats = nil
			return err
		}
		ads = append(ads, ad)
	}

	d.AdapterStats = ads
	return nil
}

// GetPhysicalDrive() is used to get the PhysicalDriveStat of a DiskStatus.
func (d *DiskStatus) GetPhysicalDrive() error {
	ads := make([]AdapterStat, 0)

	command := d.megacliPath
	for i := 0; i < d.adapterCount; i++ {
		ad := AdapterStat{
			AdapterId: i,
		}
		err := ad.getMegaRaidPdInfo(command)
		if err != nil {
			d.AdapterStats = nil
			return err
		}
		ads = append(ads, ad)
	}

	d.AdapterStats = ads
	return nil
}

// ListBrokenDrive() is used to list the Broken Drives of a DiskStatus.
func (d *DiskStatus) ListBrokenDrive() ([]VirtualDriveStat, []PhysicalDriveStat, error) {
	brokenVds, err := d.ListBrokenVirtualDrive()
	if err != nil {
		return nil, nil, err
	}

	brokenPds, err := d.ListBrokenPhysicalDrive()
	if err != nil {
		return nil, nil, err
	}

	return brokenVds, brokenPds, nil
}

// ListBrokenVirtualDrive() is used to list the Broken Virtual Drives of a DiskStatus.
func (d *DiskStatus) ListBrokenVirtualDrive() ([]VirtualDriveStat, error) {
	err := d.GetVirtualDrive()
	if err != nil {
		return nil, err
	}

	brokenVds := make([]VirtualDriveStat, 0)
	for _, ads := range d.AdapterStats {
		for _, vds := range ads.VirtualDriveStats {
			if vds.State != "Optimal" {
				brokenVds = append(brokenVds, vds)
			}
		}
	}
	return brokenVds, nil
}

// ListBrokenPhysicalDrive() is used to list the Broken Physical Drives of a DiskStatus.
func (d *DiskStatus) ListBrokenPhysicalDrive() ([]PhysicalDriveStat, error) {
	err := d.GetPhysicalDrive()
	if err != nil {
		return nil, err
	}

	brokenPds := make([]PhysicalDriveStat, 0)
	for _, ads := range d.AdapterStats {
		for _, pds := range ads.PhysicalDriveStats {
			if !strings.Contains(pds.FirmwareState, "Online") {
				brokenPds = append(brokenPds, pds)
			}
		}
	}
	return brokenPds, nil
}
