// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package sgio

import (
	"syscall"
	"unsafe"

	"github.com/platinasystems/scsi/internal/godefs/sg"
)

type Dev int

func Open(fn string) (Dev, error) {
	const (
		mode = syscall.O_RDONLY | syscall.O_NONBLOCK | syscall.O_CLOEXEC
		perm = 0
	)
	fd, err := syscall.Open(fn, mode, perm)
	return Dev(fd), err
}

func (dev Dev) Close() error {
	return syscall.Close(int(dev))
}
func (dev Dev) IOCTL(cmd, v uintptr) (err error) {
	fd := uintptr(dev)
	_, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, fd, cmd, v)
	if errno != 0 {
		err = errno
	}
	return
}

// Version returns symver ecoded with X100 segments, e.g. 30536 :: "v3.5.36"
func (dev Dev) Version() (ver int32) {
	err := dev.IOCTL(sg.SG_GET_VERSION_NUM, uintptr(unsafe.Pointer(&ver)))
	if err != nil {
		ver = -1
	}
	return
}
