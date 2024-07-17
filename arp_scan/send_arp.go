package arp_scan

import (
	"context"
	"encoding/binary"
	"net"
	"syscall"
	"unsafe"
)

var (
	// https://learn.microsoft.com/zh-cn/windows/win32/api/iphlpapi/
	modiphlpapi = syscall.NewLazyDLL("iphlpapi.dll")

	// https://learn.microsoft.com/zh-cn/windows/win32/api/iphlpapi/nf-iphlpapi-sendarp
	sendARP = modiphlpapi.NewProc("SendARP")
)

func SendTo(ctx context.Context, srcIP, dstIP net.IP) (_ net.HardwareAddr, err error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var mac [6]byte
		var macLen = uint32(len(mac))

		if code, _, _ := sendARP.Call(
			uintptr(binary.LittleEndian.Uint32(dstIP.To4()[:])),
			uintptr(binary.LittleEndian.Uint32(srcIP.To4()[:])),
			uintptr(unsafe.Pointer(&mac[0])),
			uintptr(unsafe.Pointer(&macLen)),
		); code != 0 {
			return nil, syscall.Errno(code)
		}
		return mac[:], nil
	}
}
