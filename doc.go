/*
This package is used for go codes to get MegaRaid stat.

Usage

Create a DiskStatus struct by calling diskutil.NewDiskStatus(). You need provide the MegaCli binary path and the count of RAID card in your server.

	ds, err := diskutil.NewDiskStatus(megaPath, adapterCount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DiskStatus New error: %v\n", err)
		return
	}

	err = ds.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DiskStatus Get error: %v\n", err)
		return
	}
	fmt.Println(ds)

After calling ds.Get(), you can visit any stat in the DiskStatus like this:

	for i, ads := range ds.AdapterStats {
		fmt.Printf("adapter #%d \n", i)
		for j, pds := range ads.PhysicalDriveStats {
			pdStatus := pds.FirmwareState
			fmt.Printf("PD%d status: %s\n", j, pdStatus)
		}
		fmt.Printf("\n")
	}

Or print the DiskStatus in json format:

	{
		"adapter_stats": [
			{
				"id": 0,
				"virtual_drive_stats": [
					{
						"virtual_drive": 0,
						"name": "",
						"size": "278.875 GB",
						"state": "Optimal",
						"number_of_drives": 1,
						"encryption_type": "None"
					}
				],
				"physical_drive_stats": [
					{
						"enclosure_device_id": 64,
						"device_id": 8,
						"slot_number": 0,
						"media_error_count": 0,
						"other_error_count": 0,
						"predictive_failure_count": 0,
						"pd_type": "SAS",
						"raw_size": "279.396 GB [0x22ecb25c Sectors]",
						"firmware_state": "Online, Spun Up",
						"brand": "SEAGATE",
						"model": "ST9300605SS",
						"serial_number": "00046XP4MQNJ",
						"drive_emperature": "65C (149.00 F)"
					}
				]
			}
		]
	}

*/

package diskutil
