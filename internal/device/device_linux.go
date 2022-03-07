// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package device

import "syscall"

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

func (dev Dev) IOCTL(cmd, data uintptr) (err error) {
	_, _, errno := syscall.
		RawSyscall(syscall.SYS_IOCTL, uintptr(dev), cmd, data)
	if errno != 0 {
		err = errno
	}
	return
}

func (dev Dev) Read(data []byte) (int, error) {
	return syscall.Read(int(dev), data)
}

func (dev Dev) Write(data []byte) (int, error) {
	return syscall.Write(int(dev), data)
}
