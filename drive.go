package diskutil

import "encoding/json"

// PhysicalDriveStat is a struct to get the Physical Drive Stat of a RAID card.
type PhysicalDriveStat struct {
	EnclosureDeviceId      int    `json:"enclosure_device_id"`
	DeviceId               int    `json:"device_id"`
	SlotNumber             int    `json:"slot_number"`
	MediaErrorCount        int    `json:"media_error_count"`
	OtherErrorCount        int    `json:"other_error_count"`
	PredictiveFailureCount int    `json:"predictive_failure_count"`
	Pdtype                 string `json:"pd_type"`
	RawSize                string `json:"raw_size"`
	FirmwareState          string `json:"firmware_state"`
	Brand                  string `json:"brand"`
	Model                  string `json:"model"`
	SerialNumber           string `json:"serial_number"`
	DriveTemperature       string `json:"drive_emperature"`
}

// String is used to get the print string.
func (p *PhysicalDriveStat) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// ToJson is used to get the json encoded string.
func (p *PhysicalDriveStat) ToJson() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
