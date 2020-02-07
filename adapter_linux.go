package diskutil

func (a *AdapterStat) parseMegaRaidVdInfo(info string) error {
	if info == "" {
		return errors.New("mageRaid vd info nil")
	}

	vds := make([]VirtualDriveStat, 0)
	parts := strings.Split(info, keyVdVirtualDriveSpace)

	for _, vdinfo := range parts {
		if strings.Contains(vdinfo, keyVdTargetId) {
			vdinfo = keyVdVirtualDrive + vdinfo
			vd := VirtualDriveStat{}
			lines := strings.Split(vdinfo, "\n")
			for _, line := range lines {
				err := vd.parseLine(line)
				if err != nil {
					continue
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

	err = a.parseMegaRaidVdInfo(parts[0])
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
					continue
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

func (a *AdapterStat) getLog(cmd string) (string, error) {
	args := "-fwtermlog -dsply -a" + strconv.Itoa(a.AdapterId)
	if output, err := execCmd(cmd, args); err != nil {
		return "", err
	} else {
		return output, nil
	}
}
