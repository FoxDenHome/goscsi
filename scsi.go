// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

// package goscsi provides a device command interface.
package goscsi

import "time"

type Dev struct{ EmbeddedDev }

type EmbeddedDev interface {
	Close() error
	Request(cdb, fromdev, todev []byte) error
	RequestWithTimeout(cdb, fromdev, todev []byte, timeout time.Duration) error
}

func Open(fn string) (Dev, error) {
	dev, err := open(fn)
	return Dev{dev}, err
}
