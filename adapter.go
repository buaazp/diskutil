package diskutil

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// AdapterStat is a struct to get the Adapter Stat of a RAID card.
// AdapterStat has VirtualDriveStats and PhysicalDriveStats in itself.
type AdapterStat struct {
	AdapterId          int                 `json:"adapter_id"`
	VirtualDriveStats  []VirtualDriveStat  `json:"virtual_drive_stats"`
	PhysicalDriveStats []PhysicalDriveStat `json:"physical_drive_stats"`
}

// String() is used to get the print string.
func (a *AdapterStat) String() string {
	data, err := json.Marshal(a)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// ToJson() is used to get the json encoded string.
func (a *AdapterStat) ToJson() (string, error) {
	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (a *AdapterStat) parseMegaRaidVdInfo(info string) error {
	if info == "" {
		return errors.New("mageRaid vd info nil")
	}

	vds := make([]VirtualDriveStat, 0)

	parts := strings.Split(info, keyVdVirtualDrive)
	for _, vdinfo := range parts {
		if strings.Contains(vdinfo, keyVdTargetId) {
			vdinfo = keyVdVirtualDrive + vdinfo
			vd := VirtualDriveStat{}
			lines := strings.Split(vdinfo, "\n")
			for _, line := range lines {
				err := vd.parseLine(line)
				if err != nil {
					return err
				}
			}
			vds = append(vds, vd)
		}
	}

	a.VirtualDriveStats = vds
	return nil
}

func (a *AdapterStat) getMegaRaidVdInfo(command string) error {
	args := "-ldinfo -lall -a" + strconv.Itoa(a.AdapterId)

	output, err := execCmd(command, args)
	if err != nil {
		return err
	}
	parts := strings.SplitN(output, keyExitResult, 2)
	if len(parts) != 2 {
		return errors.New("megaCli output illegal")
	}
	result := strings.TrimSpace(parts[1])
	if result != "0x00" {
		return errors.New("megaCli return error: " + result)
	}

	err = a.parseMegaRaidVdInfo(output)
	if err != nil {
		return err
	}
	return nil
}

func (a *AdapterStat) parseMegaRaidPdInfo(info string) error {
	if info == "" {
		return errors.New("mageRaid pd info nil")
	}

	pds := make([]PhysicalDriveStat, 0)

	parts := strings.Split(info, keyPdEnclosureDeviceId)
	for _, pdinfo := range parts {
		if strings.Contains(pdinfo, keyPdSlotNumber) {
			pdinfo = keyPdEnclosureDeviceId + pdinfo
			pd := PhysicalDriveStat{}
			lines := strings.Split(pdinfo, "\n")
			for _, line := range lines {
				err := pd.parseLine(line)
				if err != nil {
					return err
				}
			}
			pds = append(pds, pd)
		}
	}

	a.PhysicalDriveStats = pds
	return nil
}

func (a *AdapterStat) getMegaRaidPdInfo(command string) error {
	args := "-pdlist -a" + strconv.Itoa(a.AdapterId)

	output, err := execCmd(command, args)
	if err != nil {
		return err
	}
	parts := strings.SplitN(output, keyExitResult, 2)
	if len(parts) != 2 {
		return errors.New("megaCli output illegal")
	}
	result := strings.TrimSpace(parts[1])
	if result != "0x00" {
		return errors.New("megaCli return error: " + result)
	}

	err = a.parseMegaRaidPdInfo(output)
	if err != nil {
		return err
	}
	return nil
}
