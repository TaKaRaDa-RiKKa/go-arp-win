package arp_scan_test

import (
	"context"
	"encoding/binary"
	"go-arp-win/arp_scan"
	"net"
	"sync"
	"testing"
)

func TestSendARPV4(t *testing.T) {
	var wg sync.WaitGroup
	for _, ip := range ips(net.IPv4Mask(192, 168, 1, 0)) {
		wg.Add(1)
		go func(ip net.IP) {
			defer wg.Done()

			var err error
			var hardwareAddr net.HardwareAddr
			if hardwareAddr, err = arp_scan.SendTo(context.Background(), net.IPv4zero, ip); err != nil {
				return
			}

			names, _ := net.LookupAddr(ip.To4().String())
			t.Logf("IP: %s \t MAC: %s \t HostName: %s", ip.To4().String(), hardwareAddr.String(), names)
		}(ip)
	}
	wg.Wait()
}

func ips(ipMask net.IPMask) (ips []net.IP) {
	for i := 1; i < 255; i++ {
		var num = binary.BigEndian.Uint32(ipMask) + uint32(i)

		var buf = [4]byte{}
		binary.BigEndian.PutUint32(buf[0:], num)
		ips = append(ips, buf[:])
	}
	return
}
