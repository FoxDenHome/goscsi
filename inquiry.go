// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package goscsi

import (
	"fmt"

	"github.com/FoxDenHome/goscsi/godefs/scsi"
	"github.com/FoxDenHome/goscsi/internal/vpd"
)

const (
	Page00  = 0x00
	Page80  = 0x80
	Page83  = 0x83
	DataLen = 0xfc
)

func (dev Dev) Inquiry(vpd vpd.Request, pagecode byte, data []byte) error {
	cdb := []byte{
		scsi.INQUIRY,
		byte(vpd),
		byte(pagecode),
		0x00,
		byte(len(data)),
		0x00,
	}
	return dev.Request(cdb, data, nil)
}

/*
VPDs() returns the `SUPPORTED PAGE LIST` from the page 0x00 EVPD inquiry
response.

<https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf>

	Table 483 Supported Vital Product Data pages
	┏━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┓
	┃  Bit┃   7    ┃   6    ┃   5    ┃   4    ┃   3    ┃   2    ┃   1    ┃   0    ┃
	┃Byte ┃        ┃        ┃        ┃        ┃        ┃        ┃        ┃        ┃
	┣━━━━━╇━━━━━━━━┻━━━━━━━━┻━━━━━━━━╇━━━━━━━━┻━━━━━━━━┻━━━━━━━━┻━━━━━━━━┻━━━━━━━━┫
	┃ 0   │ PERIPHERAL QUALIFIER     │ PERIPHERAL DEVICE TYPE                     ┃
	┠─────┼──────────────────────────┴────────────────────────────────────────────┨
	┃ 1   │ PAGE CODE (00h)                                                       ┃
	┠─────┼───────────────────────────────────────────────────────────────────────┨
	┃ 2   │ Reserved                                                              ┃
	┠─────┼───────────────────────────────────────────────────────────────────────┨
	┃ 3   │ PAGE LENGTH (n-3)                                                     ┃
	┠─────┼───────────────────────────────────────────────────────────────────────┨
	┃ 4   │                                                                       ┃
	┠─────┤                                                                       ┃
	┃ ... │ SUPPORTED PAGE LIST                                                   ┃
	┠─────┤                                                                       ┃
	┃ n   │                                                                       ┃
	┗━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
*/
func (dev Dev) VPDs() (l []byte, err error) {
	data := make([]byte, DataLen)
	if err = dev.Inquiry(vpd.Enable, Page00, data); err != nil {
		err = fmt.Errorf("page00evpd: %w", err)
	} else if data[1] != Page00 {
		err = fmt.Errorf("page00evpd: code mismatch: %#x", data[1])
	} else if n := int(data[3]) + 1; n > len(data) {
		err = fmt.Errorf("page00evpd: excessive length: %d", n)
	} else {
		l = data[4:n]
	}
	return
}

func (dev Dev) StandardInquiry() (std StdInqData, err error) {
	data := make([]byte, DataLen)
	if err = dev.Inquiry(vpd.Disable, Page00, data); err != nil {
		err = fmt.Errorf("page0: %w", err)
	} else if data[1] != Page00 {
		err = fmt.Errorf("page0: code mismatch: %#x", data[1])
	} else {
		std = StdInqData(data)
	}
	return
}

/*
SerialNumber() returns the `Product Serial Number` field from the page 0x80
EVPD inquiry response.

<https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf>

		Table 484 Unit Serial Number page (80h)
		┏━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┓
		┃  Bit┃   7    ┃   6    ┃   5    ┃   4    ┃   3    ┃   2    ┃   1    ┃   0    ┃
		┃Byte ┃        ┃        ┃        ┃        ┃        ┃        ┃        ┃        ┃
		┣━━━━━╇━━━━━━━━┻━━━━━━━━┻━━━━━━━━╇━━━━━━━━┻━━━━━━━━┻━━━━━━━━┻━━━━━━━━┻━━━━━━━━┫
		┃ 0   │ PERIPHERAL QUALIFIER     │ PERIPHERAL DEVICE TYPE                     ┃
		┠─────┼──────────────────────────┴────────────────────────────────────────────┨
		┃ 1   │ PAGE CODE (80h)                                                       ┃
		┠─────┼───────────────────────────────────────────────────────────────────────┨
		┃ 2   │ Reserved                                                              ┃
		┠─────┼───────────────────────────────────────────────────────────────────────┨
		┃ 3   │ PAGE LENGTH                                                           ┃
		┠─────┼───────────────────────────────────────────────────────────────────────┨
	        ┃ 4   │                                                                       ┃
	        ┠─────┤                                                                       ┃
	        ┃ ... │ Product Serial Number                                                 ┃
	        ┠─────┤                                                                       ┃
		┃ n   │                                                                       ┃
		┗━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
*/
func (dev Dev) SerialNumber() (sn String, err error) {
	data := make([]byte, DataLen)
	if err = dev.Inquiry(vpd.Enable, Page80, data); err != nil {
		err = fmt.Errorf("page80evpd: %w", err)
	} else if data[1] != Page80 {
		err = fmt.Errorf("page80evpd: code mismatch: %#x", data[1])
	} else if n := int(data[3]); n+4 > len(data) {
		err = fmt.Errorf("page80evpd: excessive length: %d", n)
	} else {
		sn = String(data[4 : n+4])
	}
	return
}

/*
IDs() returns the identification descriptors from the page 0x83 EVPD inquiry
response.

<https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf>

		Table 459 Device Identification VPD page
		┏━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┓
		┃  Bit┃   7    ┃   6    ┃   5    ┃   4    ┃   3    ┃   2    ┃   1    ┃   0    ┃
		┃Byte ┃        ┃        ┃        ┃        ┃        ┃        ┃        ┃        ┃
		┣━━━━━╇━━━━━━━━┻━━━━━━━━┻━━━━━━━━╇━━━━━━━━┻━━━━━━━━┻━━━━━━━━┻━━━━━━━━┻━━━━━━━━┫
		┃ 0   │ PERIPHERAL QUALIFIER     │ PERIPHERAL DEVICE TYPE                     ┃
		┠─────┼──────────────────────────┴────────────────────────────────────────────┨
		┃ 1   │ PAGE CODE (83h)                                                       ┃
		┠─────┼───────────────────────────────────────────────────────────────────────┨
	        ┃ 2   │ (MSB)                                                                 ┃
	        ┠─────┤          PAGE LENGTH (n-3)                                            ┃
		┃ 3   │                                                                 (LSB) ┃
		┠─────┴───────────────────────────────────────────────────────────────────────┨
		┃ Identification descriptor list ...                                          ┃
		┠─────┬───────────────────────────────────────────────────────────────────────┨
		┃ 4   │ First Identification Descriptor (see table 460)                       ┃
		┠─────┼───────────────────────────────────────────────────────────────────────┨
		┃ n   │ Last Identification Descriptor                                        ┃
		┗━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
*/
func (dev Dev) IDs() (dscs []IdDsc, err error) {
	data := make([]byte, DataLen)
	if err = dev.Inquiry(vpd.Enable, Page83, data); err != nil {
		err = fmt.Errorf("page83evpd: %w", err)
	} else if data[1] != Page83 {
		err = fmt.Errorf("page83evpd: code mismatch: %#x", data[1])
	} else if n := 4 + (int(data[2])<<8 | int(data[3])); n > len(data) {
		err = fmt.Errorf("page83evpd: excessive length: %d", n)
	} else {
		for i, j := 4, 4; i < n; i = j {
			j += 4 + int(data[i+3])
			if j > n {
				j = n
			}
			dscs = append(dscs, IdDsc(data[i:j]))
		}
	}
	return
}

/*
StdIngData provides methods that parse the respective fields of the page 0x00
inquiry response.

<https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf>

	Table 59 Standard INQUIRY data format
	┏━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┓
	┃  Bit┃   7    ┃   6    ┃   5    ┃   4    ┃   3    ┃   2    ┃   1    ┃   0    ┃
	┃Byte ┃        ┃        ┃        ┃        ┃        ┃        ┃        ┃        ┃
	┣━━━━━╇━━━━━━━━┻━━━━━━━━┻━━━━━━━━╇━━━━━━━━┻━━━━━━━━┻━━━━━━━━┻━━━━━━━━┻━━━━━━━━┫
	┃ 0   │ PERIPHERAL QUALIFIER     │ PERIPHERAL DEVICE TYPE                     ┃
	┠─────┼────────┬─────────────────┴────────────────────────────────────────────┨
	┃ 1   │  RMB   │                  Reserved                                    ┃
	┠─────┼────────┴──────────────────────────────────────────────────────────────┨
	┃ 2   │                              Version                                  ┃
	┠─────┼────────┬────────┬─────────────────┬───────────────────────────────────┫
	┃ 3   │Obsolete│Obsolete│     Reserved    │         Response Data Format      ┃
	┠─────┼────────┴────────┴─────────────────┴───────────────────────────────────┫
	┃ 4   │                           Additional Length (n-4)                     ┃
	┠─────┼───────────────────────────────────────────────────────────────────────┫
	┃ 5   │                           Reserved                                    ┃
	┠─────┼───────────────────────────────────────────────────────────────────────┫
	┃ 6   │                           Reserved                                    ┃
	┠─────┼────────┬────────┬────────┬────────┬────────┬────────┬────────┬────────┫
	┃ 7   │Obsolete│Obsolete│Obsolete│Obsolete│Obsolete│Obsolete│ CMDQUE │   VS   ┃
	┠─────┼────────┴────────┴────────┴────────┴────────┴────────┴────────┴────────┫
	┃ 8   │ (MSB)                                                                 ┃
	┠─────┤                           VENDOR IDENTIFICATION                       ┃
	┃ 15  │                                                                 (LSB) ┃
	┠─────┼───────────────────────────────────────────────────────────────────────┫
	┃ 16  │ (MSB)                                                                 ┃
	┠─────┤                           PRODUCT IDENTIFICATION                      ┃
	┃ 31  │                                                                 (LSB) ┃
	┠─────┼───────────────────────────────────────────────────────────────────────┫
	┃ 32  │ (MSB)                                                                 ┃
	┠─────┤                           PRODUCT REVISION LEVEL                      ┃
	┃ 35  │                                                                 (LSB) ┃
	┠─────┼───────────────────────────────────────────────────────────────────────┫
	┃ 36  │                                                                       ┃
	┠─────┤                           SERIAL NUMBER                               ┃
	┃ 43  │                                                                       ┃
	┠─────┼───────────────────────────────────────────────────────────────────────┫
	┃ 56  │                                                                       ┃
	┠─────┤                           Reserved                                    ┃
	┃ 95  │                                                                       ┃
	┠─────┴───────────────────────────────────────────────────────────────────────┫
	┃                             Vendor-Specific Parameters                      ┃
	┠─────┬───────────────────────────────────────────────────────────────────────┫
	┃ 96  │                                                                       ┃
	┠─────┤                           Vendor Specific                             ┃
	┃ n   │                                                                       ┃
	┗━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
*/
type StdInqData []byte

func (data StdInqData) Type() Type {
	return Type(data[0] & 0x1f)
}

func (data StdInqData) VendorID() String {
	return String(data[8:16])
}

func (data StdInqData) ProductID() String {
	return String(data[16:32])
}

func (data StdInqData) ProductRev() String {
	return String(data[32:36])
}

func (data StdInqData) SerialNumber() String {
	return String(data[36:44])
}

/*
<https://www.seagate.com/files/staticfiles/support/docs/manual/Interface%20manuals/100293068j.pdf>

		Table 460 Identification Descriptor
		┏━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┳━━━━━━━━┓
		┃  Bit┃   7    ┃   6    ┃   5    ┃   4    ┃   3    ┃   2    ┃   1    ┃   0    ┃
		┃Byte ┃        ┃        ┃        ┃        ┃        ┃        ┃        ┃        ┃
		┣━━━━━╇━━━━━━━━┻━━━━━━━━┻━━━━━━━━┻━━━━━━━━╇━━━━━━━━┻━━━━━━━━┻━━━━━━━━┻━━━━━━━━┫
		┃ 0   │ PROTOCOL IDENTIFER                │ CODE SET                          ┃
		┠─────┼────────┬────────┬─────────────────┼───────────────────────────────────┨
		┃ 1   │ PIV    │Reserved│ ASSOCIATION     │ IDENTIFIER TYPE                   ┃
		┠─────┼────────┴────────┴─────────────────┴───────────────────────────────────┨
		┃ 2   │ Reserved                                                              ┃
		┠─────┼───────────────────────────────────────────────────────────────────────┨
		┃ 3   │ IDENTIFIER LENGTH (n-3)                                               ┃
		┠─────┼───────────────────────────────────────────────────────────────────────┨
	        ┃ 4   │                                                                       ┃
	        ┠─────┤                                                                       ┃
	        ┃ ... │ IDENTIFIER                                                            ┃
	        ┠─────┤                                                                       ┃
		┃ n   │                                                                       ┃
		┗━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
*/
type IdDsc []byte

func (dsc IdDsc) Ascii() bool {
	return 2 == (0xf & dsc[0])
}

func (dsc IdDsc) Protocol() uint8 {
	return 0xf & (uint8(dsc[0]) >> 4)
}

func (dsc IdDsc) Type() uint8 {
	return 0xf & uint8(dsc[1])
}

func (dsc IdDsc) Association() uint8 {
	return 0x3 & (uint8(dsc[1]) >> 4)
}

func (dsc IdDsc) PIV() bool {
	return 1 == (dsc[1] & (1 << 7))
}

func (dsc IdDsc) Len() int {
	n := int(dsc[3])
	if n > len(dsc) {
		n = len(dsc)
	}
	return n
}

func (dsc IdDsc) ID() []byte {
	return []byte(dsc[4:])
}

func (dsc IdDsc) String() string {
	if dsc.Ascii() {
		return string(dsc.ID())
	}
	return fmt.Sprintf("%#x", dsc.ID())
}
