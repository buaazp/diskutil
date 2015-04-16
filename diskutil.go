package diskutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const (
	KeyVdVirtualDrive           string = "Virtual Drive"
	KeyVdTargetId               string = "Target Id"
	KeyVdName                   string = "Name"
	KeyVdSize                   string = "Size"
	KeyVdState                  string = "State"
	KeyVdNumberOfDrives         string = "Number Of Drives"
	KeyVdEncryptionType         string = "Encryption Type"
	KeyPdEnclosureDeviceID      string = "Enclosure Device ID"
	KeyPdSlotNumber             string = "Slot Number"
	KeyPdDeviceId               string = "Device Id"
	KeyPdMediaErrorCount        string = "Media Error Count"
	KeyPdOtherErrorCount        string = "Other Error Count"
	KeyPdPredictiveFailureCount string = "Predictive Failure Count"
	KeyPdPdType                 string = "PD Type"
	KeyPdRawSize                string = "Raw Size"
	KeyPdFirmwareState          string = "Firmware state"
	KeyPdInquiryData            string = "Inquiry Data"
	KeyPdDriveTemperature       string = "Drive Temperature"

	TypeString int = iota
	TypeInt
	TypeUint64
)

type VirtualDriveStat struct {
	VirtualDrive   int    `json:"virtual_drive"`
	Name           string `json:"name"`
	Size           string `json:"size"`
	State          string `json:"state"`
	NumberOfDrives int    `json:"number_of_drives"`
	EncryptionType string `json:"encryption_type"`
}

type PhysicalDriveStat struct {
	EnclosureDeviceID      int    `json:"enclosure_device_id"`
	DeviceId               int    `json:"device_id"`
	SlotNumber             int    `json:"slot_number"`
	MediaErrorCount        int    `json:"media_error_count"`
	OtherErrorCount        int    `json:"other_error_count"`
	PredictiveFailureCount int    `json:"predictive_failure_count"`
	PdType                 string `json:"pd_type"`
	RawSize                string `json:"raw_size"`
	FirmwareState          string `json:"firmware_state"`
	Brand                  string `json:"brand"`
	Model                  string `json:"model"`
	SerialNumber           string `json:"serial_number"`
	DriveTemperature       string `json:"drive_emperature"`
}

type AdapterStat struct {
	ID                 int                 `json:"id"`
	VirtualDriveStats  []VirtualDriveStat  `json:"virtual_drive_stats"`
	PhysicalDriveStats []PhysicalDriveStat `json:"physical_drive_stats"`
}

type DiskStatus struct {
	megacliPath  string
	AdapterCount int           `json:"adapter_count"`
	AdapterStats []AdapterStat `json:"adapter_stats"`
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func NewDiskStatus(mp string, ac int) (*DiskStatus, error) {
	mp = path.Clean(mp)
	if !Exist(mp) {
		return nil, errors.New("megaCli not exist")
	}
	ds := new(DiskStatus)
	ds.megacliPath = mp
	ds.AdapterCount = ac
	return ds, nil
}

func (d *DiskStatus) String() string {
	data, err := json.Marshal(d)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (d *DiskStatus) ToJson() (string, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(data), nil
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

func parseFiled(line, filed string, targetType int) (interface{}, error) {
	fileds := strings.SplitN(line, ":", 2)
	if len(fileds) != 2 {
		return nil, errors.New("format illegal: " + line)
	}

	data := strings.TrimSpace(fileds[1])
	if targetType == TypeString {
		return data, nil
	} else if targetType == TypeInt {
		value, err := strconv.ParseInt(data, 10, 0)
		if err != nil {
			return nil, err
		}
		return int(value), nil
	} else if targetType == TypeUint64 {
		value, err := strconv.ParseUint(data, 10, 0)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
	return nil, errors.New("type not supported")
}

func (v *VirtualDriveStat) parseLine(line string) error {
	if strings.HasPrefix(line, KeyVdVirtualDrive) {
		parts := strings.SplitN(line, "(", 2)
		if len(parts) != 2 {
			return errors.New("format illegal: " + line)
		}
		virtualDrive, err := parseFiled(parts[0], KeyVdVirtualDrive, TypeInt)
		if err != nil {
			return err
		}
		v.VirtualDrive = virtualDrive.(int)
	} else if strings.HasPrefix(line, KeyVdName) {
		name, err := parseFiled(line, KeyVdName, TypeString)
		if err != nil {
			return err
		}
		v.Name = name.(string)
	} else if strings.HasPrefix(line, KeyVdSize) {
		size, err := parseFiled(line, KeyVdSize, TypeString)
		if err != nil {
			return err
		}
		v.Size = size.(string)
	} else if strings.HasPrefix(line, KeyVdState) {
		state, err := parseFiled(line, KeyVdState, TypeString)
		if err != nil {
			return err
		}
		v.State = state.(string)
	} else if strings.HasPrefix(line, KeyVdNumberOfDrives) {
		numberOfDrives, err := parseFiled(line, KeyVdNumberOfDrives, TypeInt)
		if err != nil {
			return err
		}
		v.NumberOfDrives = numberOfDrives.(int)
	} else if strings.HasPrefix(line, KeyVdEncryptionType) {
		encryptionType, err := parseFiled(line, KeyVdEncryptionType, TypeString)
		if err != nil {
			return err
		}
		v.EncryptionType = encryptionType.(string)
	}
	return nil
}

func (a *AdapterStat) parseMegaRaidVdInfo(info string) error {
	if info == "" {
		return errors.New("mageRaid vd info nil")
	}

	vds := make([]VirtualDriveStat, 0)

	parts := strings.Split(info, KeyVdVirtualDrive)
	for _, vdinfo := range parts {
		if strings.Contains(vdinfo, KeyVdTargetId) {
			vdinfo = KeyVdVirtualDrive + vdinfo
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
	args := "-ldinfo -lall -a" + strconv.Itoa(a.ID)

	output, err := execCmd(command, args)
	if err != nil {
		return err
	}
	parts := strings.SplitN(output, "Exit Code:", 2)
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

func (p *PhysicalDriveStat) parseLine(line string) error {
	if strings.HasPrefix(line, KeyPdEnclosureDeviceID) {
		enclosureDeviceID, err := parseFiled(line, KeyPdEnclosureDeviceID, TypeInt)
		if err != nil {
			return err
		}
		p.EnclosureDeviceID = enclosureDeviceID.(int)
	} else if strings.HasPrefix(line, KeyPdDeviceId) {
		deviceId, err := parseFiled(line, KeyPdDeviceId, TypeInt)
		if err != nil {
			return err
		}
		p.DeviceId = deviceId.(int)
	} else if strings.HasPrefix(line, KeyPdSlotNumber) {
		slotNumber, err := parseFiled(line, KeyPdSlotNumber, TypeInt)
		if err != nil {
			return err
		}
		p.SlotNumber = slotNumber.(int)
	} else if strings.HasPrefix(line, KeyPdMediaErrorCount) {
		mediaErrorCount, err := parseFiled(line, KeyPdMediaErrorCount, TypeInt)
		if err != nil {
			return err
		}
		p.MediaErrorCount = mediaErrorCount.(int)
	} else if strings.HasPrefix(line, KeyPdOtherErrorCount) {
		otherErrorCount, err := parseFiled(line, KeyPdOtherErrorCount, TypeInt)
		if err != nil {
			return err
		}
		p.OtherErrorCount = otherErrorCount.(int)
	} else if strings.HasPrefix(line, KeyPdPredictiveFailureCount) {
		predictiveFailureCount, err := parseFiled(line, KeyPdPredictiveFailureCount, TypeInt)
		if err != nil {
			return err
		}
		p.PredictiveFailureCount = predictiveFailureCount.(int)
	} else if strings.HasPrefix(line, KeyPdPdType) {
		pdType, err := parseFiled(line, KeyPdPdType, TypeString)
		if err != nil {
			return err
		}
		p.PdType = pdType.(string)
	} else if strings.HasPrefix(line, KeyPdRawSize) {
		rawSize, err := parseFiled(line, KeyPdRawSize, TypeString)
		if err != nil {
			return err
		}
		p.RawSize = rawSize.(string)
	} else if strings.HasPrefix(line, KeyPdFirmwareState) {
		firmwareState, err := parseFiled(line, KeyPdFirmwareState, TypeString)
		if err != nil {
			return err
		}
		p.FirmwareState = firmwareState.(string)
	} else if strings.HasPrefix(line, KeyPdInquiryData) {
		inquiryData, err := parseFiled(line, KeyPdInquiryData, TypeString)
		if err != nil {
			return err
		}
		inquiryStr := inquiryData.(string)
		parts := strings.Fields(inquiryStr)
		if len(parts) == 3 {
			p.SerialNumber = parts[2]
			p.Model = parts[1]
			p.Brand = parts[0]
		} else if len(parts) == 2 {
			p.SerialNumber = parts[1]
			p.Model = parts[0]
		} else if len(parts) == 1 {
			p.SerialNumber = parts[0]
		}
	} else if strings.HasPrefix(line, KeyPdDriveTemperature) {
		driveTemperature, err := parseFiled(line, KeyPdDriveTemperature, TypeString)
		if err != nil {
			return err
		}
		p.DriveTemperature = driveTemperature.(string)
	}
	return nil
}

func (a *AdapterStat) parseMegaRaidPdInfo(info string) error {
	if info == "" {
		return errors.New("mageRaid pd info nil")
	}

	pds := make([]PhysicalDriveStat, 0)

	parts := strings.Split(info, KeyPdEnclosureDeviceID)
	for _, pdinfo := range parts {
		if strings.Contains(pdinfo, KeyPdSlotNumber) {
			pdinfo = KeyPdEnclosureDeviceID + pdinfo
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
	args := "-pdlist -a" + strconv.Itoa(a.ID)

	output, err := execCmd(command, args)
	if err != nil {
		return err
	}
	parts := strings.SplitN(output, "Exit Code:", 2)
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

func (d *DiskStatus) Get() error {
	// command := "/opt/MegaRAID/MegaCli/MegaCli64"
	command := d.megacliPath
	ads := make([]AdapterStat, 0)

	for i := 0; i < d.AdapterCount; i++ {
		ad := AdapterStat{
			ID: i,
		}
		err := ad.getMegaRaidVdInfo(command)
		if err != nil {
			return err
		}
		err = ad.getMegaRaidPdInfo(command)
		if err != nil {
			return err
		}
		ads = append(ads, ad)
	}

	d.AdapterStats = ads
	return nil
}
