package diskutil

import "encoding/json"

// VirtualDriveStat is a struct to get the Virtual Drive Stat of a RAID card.
type VirtualDriveStat struct {
	VirtualDrive   int    `json:"virtual_drive"`
	Name           string `json:"name"`
	Size           string `json:"size"`
	State          string `json:"state"`
	NumberOfDrives int    `json:"number_of_drives"`
	Encryptiontype string `json:"encryption_type"`
	RAIDLevel      string `json:"raid_level"`
}

// String is used to get the print string.
func (v *VirtualDriveStat) String() string {
	data, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// ToJson is used to get the json encoded string.
func (v *VirtualDriveStat) ToJson() (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
