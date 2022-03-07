// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package sgio

import (
	"unsafe"

	"github.com/platinasystems/scsi/internal/godefs/bsg"
	"github.com/platinasystems/scsi/internal/godefs/sg"
)

type v4 struct{ embeddedDev }

func (v4 v4) Request(cdb, data []byte, todev ...byte) error {
	sense := make([]byte, 32)
	hdr := bsg.SgIoV4{
		Guard:          'Q',
		Protocol:       bsg.BSG_PROTOCOL_SCSI,
		Subprotocol:    bsg.BSG_SUB_PROTOCOL_SCSI_CMD,
		RequestLen:     uint32(len(cdb)),
		Request:        uint64(uintptr(unsafe.Pointer(&cdb[0]))),
		MaxResponseLen: uint32(len(sense)),
		Response:       uint64(uintptr(unsafe.Pointer(&sense[0]))),
	}
	if len(data) > 0 {
		hdr.DinXferLen = uint32(len(data))
		hdr.DinXferp = uint64(uintptr(unsafe.Pointer(&data[0])))
	}
	if len(todev) > 0 {
		hdr.DoutXferLen = uint32(len(todev))
		hdr.DoutXferp = uint64(uintptr(unsafe.Pointer(&todev[0])))
	}
	return v4.IOCTL(sg.SG_IO, uintptr(unsafe.Pointer(&hdr)))
}
