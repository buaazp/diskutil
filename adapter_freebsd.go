package diskutil

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	matchPdInfo = `(?m)(\d+)\s+\(\s*(.+)\)\s+([\w\s]+)\s+<([\w\-\s]+)\s+(\w+)\s+serial=(\w+)>\s+(\w+)\s+(\w+)$`
	matchVdInfo = `(?m)(\w+) \(\s+(\w+)\)\s+([\w\-]+)\s+(\w+) ([\w ]+)(Disabled|Enabled)\s*(\<(\w+)\>)?`
)

func (a *AdapterStat) getMegaRaidVdInfo(command string) error {
	args := "show volumes"
	if output, err := execCmd(command, args); err != nil {
		return err
	} else {
		if err := a.parseMegaRaidVdInfo(output); err != nil {
			return err
		}
	}

	return nil
}

func (a *AdapterStat) getMegaRaidPdInfo(command string) error {
	args := "show drives"
	if output, err := execCmd(command, args); err != nil {
		return err
	} else {
		if err := a.parseMegaRaidPdInfo(output); err != nil {
			return err
		}
	}

	return nil
}

func (a *AdapterStat) getLog(cmd string) (string, error) {
	if output, err := execCmd(cmd, "show events -c debug"); err != nil {
		return "", err
	} else {
		return output, nil
	}
}

func (a *AdapterStat) parseMegaRaidPdInfo(info string) error {
	matcher := regexp.MustCompile(matchPdInfo)
	pds := make([]PhysicalDriveStat, 0)

	for _, item := range strings.Split(info, "\n") {
		match := matcher.FindStringSubmatch(item)
		if match != nil && len(matcher.SubexpNames()) > 0 {
			deviceId, _ := strconv.Atoi(match[1])
			pds = append(pds, PhysicalDriveStat{
				DeviceId:         deviceId,
				SlotNumber:       deviceId,
				RawSize:          strings.ToLower(match[2]),
				FirmwareState:    strings.TrimSpace(match[3]),
				Brand:            match[4],
				Model:            match[5],
				SerialNumber:     match[6],
				Pdtype:           match[7],
				DriveTemperature: "N/A",
			})
		}
	}

	a.PhysicalDriveStats = pds
	return nil
}

func (a *AdapterStat) parseMegaRaidVdInfo(info string) error {
	matcher := regexp.MustCompile(matchVdInfo)
	vds := make([]VirtualDriveStat, 0)

	for _, item := range strings.Split(info, "\n") {
		match := matcher.FindStringSubmatch(item)
		if match != nil && len(matcher.SubexpNames()) > 0 {
			name := strings.TrimSpace(match[8])
			if len(name) <= 0 {
				name = match[1]
			}

			vds = append(vds, VirtualDriveStat{
				Name:           name,
				Size:           strings.ToLower(match[2]),
				State:          strings.ToLower(strings.TrimSpace(match[5])),
				NumberOfDrives: 0,
				Encryptiontype: "N/A",
				RAIDLevel:      match[3],
			})
		}
	}

	a.VirtualDriveStats = vds
	return nil
}
