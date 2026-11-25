package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/salehbinamrat/Advanced-WAN-Bonding-SD-WAN-/internal/wan"
)

func run(args []string, out io.Writer, list func() ([]net.Interface, error)) error {
	if len(args) != 1 || args[0] != "interfaces" {
		return fmt.Errorf("usage: wanctl interfaces")
	}
	interfaces, err := wan.Inventory(list)
	if err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(interfaces)
}

func main() {
	if err := run(os.Args[1:], os.Stdout, net.Interfaces); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
