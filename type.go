// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package scsi

type Type byte

func (t Type) String() string {
	return map[Type]string{
		0x00: "Direct access block device (e.g., magnetic disk)",
		0x01: "Sequential-access device (e.g., magnetic tape)",
		0x02: "Printer device",
		0x03: "Processor device",
		0x04: "Write-once device (e.g., some optical disks)",
		0x05: "CD/DVD device",
		0x06: "Scanner device (obsolete)",
		0x07: "Optical memory device (e.g., some optical disks)",
		0x08: "Medium changer device (e.g., jukeboxes)",
		0x09: "Communications device (obsolete)",
		0x0A: "Obsolete",
		0x0B: "Obsolete",
		0x0C: "Storage array controller device (e.g., RAID)",
		0x0D: "Enclosure services device",
		0x0E: "Simplified direct-access device (e.g., magnetic disk)",
		0x0F: "Optical card reader/writer device",
		0x10: "Bridge Controller Commands",
		0x11: "Object-based Storage Device",
		0x12: "Automation/Drive Interface",
		0x13: "Reserved",
		0x1D: "Reserved",
		0x1E: "Well known logical unit[b]",
		0x1F: "Unknown or no device type",
	}[t&0x1f]
}
