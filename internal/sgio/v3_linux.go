// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package sgio

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/FoxDenHome/goscsi/godefs/sg"
)

type v3 struct{ embeddedDev }

func (v3 v3) Request(cdb, fromdev, todev []byte) error {
	return v3.RequestWithTimeout(cdb, fromdev, todev, time.Second*5)
}

func (v3 v3) RequestWithTimeout(cdb, fromdev, todev []byte, timeout time.Duration) error {
	dir := int32(sg.SG_DXFER_NONE)
	sense := make([]byte, 32)
	if len(todev) > 0 {
		if len(fromdev) > 0 {
			dir = sg.SG_DXFER_TO_FROM_DEV
			copy(fromdev, todev)
			fromdev = fromdev[:len(todev)]
		} else {
			dir = sg.SG_DXFER_TO_DEV
			fromdev = todev
		}
	} else if len(fromdev) > 0 {
		dir = sg.SG_DXFER_FROM_DEV
	}
	var hdr = sg.SgIoHdr{
		InterfaceId:    'S',
		CmdLen:         uint8(len(cdb)),
		MxSbLen:        uint8(len(sense)),
		DxferDirection: dir,
		DxferLen:       uint32(len(fromdev)),
		Dxferp:         &fromdev[0],
		Cmdp:           (*uint8)(&cdb[0]),
		Sbp:            (*uint8)(&sense[0]),
		Timeout:        uint32(timeout.Milliseconds()),
	}
	err := v3.IOCTL(sg.SG_IO, uintptr(unsafe.Pointer(&hdr)))
	if err != nil {
		return err
	}
	if (hdr.Info & sg.SG_INFO_OK_MASK) != sg.SG_INFO_OK {
		if n := int(hdr.SbLenWr); n > 0 {
			err = fmt.Errorf("scsi error: sense: %x",
				sense[:n])
		}
		if hdr.MaskedStatus != 0 {
			err = fmt.Errorf("scsi error: scsi status=%#x",
				hdr.Status)
		}
		if hdr.HostStatus != 0 {
			err = fmt.Errorf("scsi error: host status=%#x",
				hdr.HostStatus)
		}
		if hdr.DriverStatus != 0 {
			err = fmt.Errorf("scsi error: driver status=%#x",
				hdr.DriverStatus)
		}
	}
	return err
}
