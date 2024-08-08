package arp_scan

import (
	"context"
	"encoding/binary"
	"fmt"
	"go-arp-win/color"
	"net"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sys/windows"
)

var target string
var active bool
var colorful bool

var Scanner = &cobra.Command{
	Short: "arp scan tool",
	Use:   "goarp -t ip/slash",
	Run:   run,
}

func init() {
	Scanner.PersistentFlags().StringVarP(&target, "target", "t", "", "ipv4")
	Scanner.PersistentFlags().BoolVarP(&colorful, "color", "", true, "console color")
	Scanner.PersistentFlags().BoolVarP(&active, "active", "", false, "display active device")
}

func run(cmd *cobra.Command, _ []string) {
	if target == "" {
		_ = cmd.Help()
		return
	}

	if _, ipNet, err := net.ParseCIDR(target); err == nil {
		enableColorful()
		console5(color.ColorWhite, "IPv4", "MAC", "TIME", "ERROR")

		Range(ipNet.IP)
		return
	}

	if ip := net.ParseIP(target); ip != nil {
		enableColorful()
		console5(color.ColorWhite, "IPv4", "MAC", "TIME", "ERROR")

		Target(ip)
		return
	}

	fmt.Printf("Invalid target: %s\n", target)
}

func Range(ip net.IP) {
	var limiter = make(chan struct{}, 256)

	var wg sync.WaitGroup
	for _, dstIP := range ips(ip.To4(), ip.DefaultMask()) {
		limiter <- struct{}{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-limiter }()
			var start = time.Now()

			var err error
			var mac net.HardwareAddr
			if mac, err = SendTo(context.Background(), net.IPv4zero.To4(), dstIP.To4()); err != nil {
				if !active {
					console5(color.ColorRed, dstIP.To4().String(), "N/A", time.Since(start).String(), err.Error())
				}
				return
			}
			console5(color.ColorGreen, dstIP.To4().String(), mac.String(), time.Since(start).String(), "N/A")
		}()
	}
	wg.Wait()
}

func Target(ip net.IP) {
	var start = time.Now()

	var err error
	var mac net.HardwareAddr
	if mac, err = SendTo(context.Background(), net.IPv4zero, ip.To4()); err != nil {
		if !active {
			console5(color.ColorRed, ip.To4().String(), "N/A", time.Since(start).String(), err.Error())
		}
		return
	}
	console5(color.ColorGreen, ip.To4().String(), mac.String(), time.Since(start).String(), "N/A")
}

func ips(ip net.IP, mask net.IPMask) (ips []net.IP) {
	var num = binary.BigEndian.Uint32(ip.To4())
	var broadcast = ^binary.BigEndian.Uint32(mask) | num

	for n := num; n <= broadcast; n++ {
		var parseIP [4]byte
		binary.BigEndian.PutUint32(parseIP[:], n)
		ips = append(ips, parseIP[:])
	}
	return
}

func console5(c color.Color, col1, col2, col3, col4 string) {
	if colorful {
		fmt.Printf("%s%-20s %-20s %-20s %-20s%s\n", c, col1, col2, col3, col4, color.Reset)
		return
	}
	fmt.Printf("%-20s %-20s %-20s %-20s\n", col1, col2, col3, col4)
}

func enableColorful() {
	handle := windows.Handle(os.Stdout.Fd())

	var mode uint32
	if err := windows.GetConsoleMode(handle, &mode); err != nil {
		colorful = false
		return
	}

	if err := windows.SetConsoleMode(handle, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err != nil {
		colorful = false
		return
	}
}
