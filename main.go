package main

import (
	"github.com/spf13/cobra"
	"go-arp-win/arp_scan"
)

func main() {
	root := &cobra.Command{}

	root.AddCommand(arp_scan.Scanner)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
