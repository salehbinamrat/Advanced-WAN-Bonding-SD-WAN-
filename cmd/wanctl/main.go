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
	switch {
	case len(args) == 1 && args[0] == "interfaces":
		interfaces, err := wan.Inventory(list)
		if err != nil {
			return err
		}
		return json.NewEncoder(out).Encode(interfaces)
	case len(args) == 2 && args[0] == "validate-config":
		f, err := os.Open(args[1])
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := wan.DecodeConfig(f); err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, "configuration valid (not applied)")
		return err
	default:
		return fmt.Errorf("usage: wanctl interfaces | wanctl validate-config FILE")
	}
}

func main() {
	if err := run(os.Args[1:], os.Stdout, net.Interfaces); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
