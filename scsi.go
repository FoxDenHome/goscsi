// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

// Package scsi provides a device command interface.
package scsi

type Dev struct{ EmbeddedDev }

type EmbeddedDev interface {
	Close() error
	Request(cdb, fromdev []byte, todev ...byte) error
}

func Open(fn string) (Dev, error) {
	dev, err := open(fn)
	return Dev{dev}, err
}
