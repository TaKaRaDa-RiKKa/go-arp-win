package arp_scan_test

import (
	"go-arp-win/arp_scan"
	"net"
	"testing"
)

func TestRange(t *testing.T) {
	_, ipnet, err := net.ParseCIDR("192.168.1.1/24")
	if err != nil {
		t.Fatal(err)
	}
	arp_scan.Range(ipnet.IP.To4())
}

func TestIP(t *testing.T) {
	arp_scan.Target(net.ParseIP("192.168.1.1").To4())
}
