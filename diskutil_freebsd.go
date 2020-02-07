package diskutil

type DiskStatus struct {
	MfiutilPath  string
	AdapterStats []AdapterStat `json:"adapter_stat"`
}

func (d *DiskStatus) Log() ([]string, error) {
	ad := AdapterStat{}
	if log, err := ad.getLog(d.MfiutilPath); err != nil {
		return []string{}, err
	} else {
		return []string{log}, nil
	}
}

// GetPhysicalDrive is used to get the PhysicalDriveStat of a DiskStatus.
func (d *DiskStatus) GetPhysicalDrive() error {

	return nil
}

// GetVirtualDrive is used to get the VirtualDriveStat of a DiskStatus.
func (d *DiskStatus) GetVirtualDrive() error {

	return nil
}

// Get() is used to get all the stat of a DiskStatus.
func (d *DiskStatus) Get() error {
	ad := AdapterStat{}

	if err := ad.getMegaRaidPdInfo(d.MfiutilPath); err != nil {
		d.AdapterStats = nil
		return err
	}

	if err := ad.getMegaRaidVdInfo(d.MfiutilPath); err != nil {
		d.AdapterStats = nil
		return err
	}

	d.AdapterStats = []AdapterStat{ad}
	return nil
}
