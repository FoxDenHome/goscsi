// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	scsi "github.com/FoxDenHome/goscsi"
)

func main() {
	usage := fmt.Sprint("usage: ", filepath.Base(os.Args[0]), " DEVICE")
	if len(os.Args) != 2 {
		fatal(usage)
	}
	arg1 := os.Args[1]
	if arg1 == "-h" || strings.TrimLeft(arg1, "-") == "help" {
		fmt.Println(usage)
		return
	}
	dev, err := scsi.Open(arg1)
	if err != nil {
		fatal(err)
	}
	defer dev.Close()

	pages, err := dev.VPDs()
	if err != nil {
		fatal(err)
	}

	for _, pc := range pages {
		switch pc {
		case scsi.Page00:
			if std, err := dev.StandardInquiry(); err != nil {
				warn(err)
			} else {
				fmt.Println("Vendor:", std.VendorID())
				fmt.Println("Type:", std.Type())
				fmt.Println("Model:", std.ProductID())
				fmt.Println("Rev:", std.ProductRev())
			}
		case scsi.Page80:
			if inq, err := dev.SerialNumber(); err != nil {
				warn(err)
			} else {
				fmt.Println("S/N:", inq)
			}
		case scsi.Page83:
			if dscs, err := dev.IDs(); err != nil {
				warn(err)
			} else if len(dscs) > 0 {
				fmt.Println("ID:")
				for _, dsc := range dscs {
					l := strings.Fields(dsc.String())
					fmt.Printf("  %d: %s\n", dsc.Type(),
						strings.Join(l, " "))
				}
			}
		}
	}
}

func fatal(args ...interface{}) {
	warn(args...)
	os.Exit(1)
}

func warn(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}
