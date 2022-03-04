// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package scsi

import "github.com/platinasystems/scsi/internal/sgio"

func Open(name string) (Dev, error) {
	sg, err := sgio.Open(name)
	if err != nil {
		return Dev{}, err
	}
	if sg.Version() > 400000 {
		return Dev{sgio.V4{sg}}, nil
	}
	return Dev{sgio.V3{sg}}, nil
}
