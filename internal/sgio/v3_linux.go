// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package sgio

import (
	"fmt"
	"unsafe"

	"github.com/platinasystems/scsi/internal/godefs/sg"
)

const (
	Millisecond = 1
	Second      = 1000 * Millisecond
	Timeout     = 5 * Second
)

type V3 struct{ Dev }

func (v3 V3) Request(cdb, data []byte, todev ...byte) error {
	dir := int32(sg.SG_DXFER_NONE)
	sense := make([]byte, 32)
	if len(todev) > 0 {
		if len(data) > 0 {
			dir = sg.SG_DXFER_TO_FROM_DEV
			copy(data, todev)
			data = data[:len(todev)]
		} else {
			dir = sg.SG_DXFER_TO_DEV
			data = todev
		}
	} else if len(data) > 0 {
		dir = sg.SG_DXFER_FROM_DEV
	}
	var hdr = sg.SgIoHdr{
		InterfaceId:    'S',
		CmdLen:         uint8(len(cdb)),
		MxSbLen:        uint8(len(sense)),
		DxferDirection: dir,
		DxferLen:       uint32(len(data)),
		Dxferp:         &data[0],
		Cmdp:           (*uint8)(&cdb[0]),
		Sbp:            (*uint8)(&sense[0]),
		Timeout:        Timeout,
	}
	err := v3.IOCTL(sg.SG_IO, uintptr(unsafe.Pointer(&hdr)))
	if err != nil {
		return err
	}
	if (hdr.Info & sg.SG_INFO_OK_MASK) != sg.SG_INFO_OK {
		err = fmt.Errorf("INQUIRY !OK:")
		if n := int(hdr.SbLenWr); n > 0 {
			err = fmt.Errorf("%v\n\tsense: %x",
				err, sense[:n])
		}
		if hdr.MaskedStatus != 0 {
			err = fmt.Errorf("%v\n\tscsi status=%#x",
				err, hdr.Status)
		}
		if hdr.HostStatus != 0 {
			err = fmt.Errorf("%v\n\thost status=%#x",
				err, hdr.HostStatus)
		}
		if hdr.DriverStatus != 0 {
			err = fmt.Errorf("%v\n\tdriver status=%#x",
				err, hdr.DriverStatus)
		}
	}
	return err
}
