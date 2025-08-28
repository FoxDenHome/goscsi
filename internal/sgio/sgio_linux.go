// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package sgio

import (
	"time"
	"unsafe"

	"github.com/FoxDenHome/goscsi/godefs/sg"
	"github.com/FoxDenHome/goscsi/internal/device"
)

type Dev interface {
	Close() error
	Request(cdb, fromdev, todev []byte) error
	RequestWithTimeout(cdb, fromdev, todev []byte, timeout time.Duration) error
}

type embeddedDev interface {
	Close() error
	IOCTL(cmd, data uintptr) error
}

func Open(fn string) (Dev, error) {
	var symver uint32
	symverptr := uintptr(unsafe.Pointer(&symver))
	dev, err := device.Open(fn)
	if err != nil {
		return nil, err
	}
	err = dev.IOCTL(sg.SG_GET_VERSION_NUM, symverptr)
	if err != nil {
		dev.Close()
		return nil, err
	}
	// symver is encoded in X100 segments, e.g.
	//	30536 :: "v3.5.36"
	if symver > 400000 {
		return v4{dev}, nil
	}
	return v3{dev}, nil
}
