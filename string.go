// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package goscsi

import "strings"

type String []byte

func (s String) String() string {
	b := []byte(s)
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			b = b[:i]
			break
		}
	}
	return strings.TrimSpace(string(b))
}
