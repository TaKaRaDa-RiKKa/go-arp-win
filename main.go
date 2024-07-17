package main

import (
	"go-arp-win/arp_scan"
)

func main() {
	if err := arp_scan.Scanner.Execute(); err != nil {
		panic(err)
	}
}
