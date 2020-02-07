package diskutil

const (
	keyExitResult               string = "Exit Code:"
	keyVdVirtualDrive           string = "Virtual Drive"
	keyVdVirtualDriveSpace      string = "\n\n\n"
	keyVdVirtualDriveType       string = "Virtual Drive Type"
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
	keyRAIDLevel                string = "RAID Level"

	typeString int = iota
	typeInt
	typeUint64
)

type DiskStatusConfig struct {
	MegaCliPath  string
	AdapterCount int
}

// DiskStatus is a struct to get all Adapters' Stat of the server
type DiskStatus struct {
	megacliPath  string
	adapterCount int
	AdapterStats []AdapterStat `json:"adapter_stats"`
}

// NewDiskStatus use the megaCliPath and apapterCount to build a DiskStatus.
func NewDiskStatus(config DiskStatusConfig) (*DiskStatus, error) {
	megaCliPath := path.Clean(config.MegaCliPath)
	if !fileExist(megaCliPath) {
		return nil, errors.New("megaCli not exist")
	}
	ds := new(DiskStatus)
	ds.megacliPath = megaCliPath
	ds.adapterCount = config.AdapterCount
	return ds, nil
}

func (d *DiskStatus) Log() ([]string, error) {
	logs := make([]string, 0)

	for i := 0; i <= d.adapterCount; i++ {
		ad := AdapterStat{
			AdapterId: i,
		}
		if log, err := ad.getLog(d.megacliPath); err == nil {
			logs = append(logs, log)
		}
	}

	return logs, nil
}

// GetPhysicalDrive is used to get the PhysicalDriveStat of a DiskStatus.
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

// GetVirtualDrive is used to get the VirtualDriveStat of a DiskStatus.
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

// Get() is used to get all the stat of a DiskStatus.
func (d *DiskStatus) Get() error {
	ads := make([]AdapterStat, 0)

	command := d.megacliPath
	for i := 0; i <= d.adapterCount; i++ {
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
