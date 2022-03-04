// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

// Vital Product Data request bit
package vpd

type Request byte

const (
	Disable Request = 0
	Enable  Request = 1
)
