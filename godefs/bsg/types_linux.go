// Copyright © 2016-2022 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

//go:build ignore
// +build ignore

package bsg

/*
#include <asm-generic/param.h>
#include </usr/include/linux/bsg.h>
*/
import "C"

const (
	BSG_PROTOCOL_SCSI = C.BSG_PROTOCOL_SCSI

	BSG_SUB_PROTOCOL_SCSI_CMD       = C.BSG_SUB_PROTOCOL_SCSI_CMD
	BSG_SUB_PROTOCOL_SCSI_TMF       = C.BSG_SUB_PROTOCOL_SCSI_TMF
	BSG_SUB_PROTOCOL_SCSI_TRANSPORT = C.BSG_SUB_PROTOCOL_SCSI_TRANSPORT

	BSG_FLAG_Q_AT_TAIL = C.BSG_FLAG_Q_AT_TAIL
	BSG_FLAG_Q_AT_HEAD = C.BSG_FLAG_Q_AT_HEAD
)

type SgIoV4 C.struct_sg_io_v4
