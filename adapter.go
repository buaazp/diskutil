package diskutil

import (
	"encoding/json"
)

// AdapterStat is a struct to get the Adapter Stat of a RAID card.
// AdapterStat has VirtualDriveStats and PhysicalDriveStats in itself.
type AdapterStat struct {
	// AdapterId          int                 `json:"adapter_id"`
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
