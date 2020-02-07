package diskutil

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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

func execCmd(command, args string) (string, error) {
	var argArray []string
	if args != "" {
		argArray = strings.Split(args, " ")
	} else {
		argArray = make([]string, 0)
	}

	cmd := exec.Command(command, argArray...)
	buf, err := cmd.Output()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "The command failed to perform: %s (Command: %s, Arguments: %s)", err, command, args)
		return "", err
	}

	return string(buf), nil
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
			if strings.ToLower(vds.State) != "optimal" {
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
			if !strings.Contains(strings.ToLower(pds.FirmwareState), "online") {
				brokenPds = append(brokenPds, pds)
			}
		}
	}
	return brokenPds, nil
}
