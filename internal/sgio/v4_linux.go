// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package sgio

import (
	"time"
	"unsafe"

	"github.com/FoxDenHome/goscsi/godefs/bsg"
	"github.com/FoxDenHome/goscsi/godefs/sg"
)

type v4 struct{ embeddedDev }

func (v4 v4) Request(cdb, fromdev, todev []byte) error {
	return v4.RequestWithTimeout(cdb, fromdev, todev, time.Second*5)
}

func (v4 v4) RequestWithTimeout(cdb, fromdev, todev []byte, timeout time.Duration) error {
	sense := make([]byte, 32)
	hdr := bsg.SgIoV4{
		Guard:          'Q',
		Protocol:       bsg.BSG_PROTOCOL_SCSI,
		Subprotocol:    bsg.BSG_SUB_PROTOCOL_SCSI_CMD,
		RequestLen:     uint32(len(cdb)),
		Request:        uint64(uintptr(unsafe.Pointer(&cdb[0]))),
		MaxResponseLen: uint32(len(sense)),
		Response:       uint64(uintptr(unsafe.Pointer(&sense[0]))),
		Timeout:        uint32(timeout.Milliseconds()),
	}
	if len(fromdev) > 0 {
		hdr.DinXferLen = uint32(len(fromdev))
		hdr.DinXferp = uint64(uintptr(unsafe.Pointer(&fromdev[0])))
	}
	if len(todev) > 0 {
		hdr.DoutXferLen = uint32(len(todev))
		hdr.DoutXferp = uint64(uintptr(unsafe.Pointer(&todev[0])))
	}
	return v4.IOCTL(sg.SG_IO, uintptr(unsafe.Pointer(&hdr)))
}
