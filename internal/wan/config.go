package wan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/netip"
	"regexp"
)

const maxConfigBytes = 64 * 1024

var (
	idPattern    = regexp.MustCompile(`^[a-z][a-z0-9-]{0,31}$`)
	ifacePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,14}$`)
)

// Config describes desired settings only; it never applies network changes.
type Config struct {
	Management string `json:"management_interface"`
	Links      []Link `json:"links"`
}

type Link struct {
	ID        string `json:"id"`
	Interface string `json:"interface"`
	Mode      string `json:"mode"`
	Address   string `json:"address,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Weight    int    `json:"weight"`
}

func DecodeConfig(r io.Reader) (Config, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxConfigBytes+1))
	if err != nil {
		return Config{}, fmt.Errorf("read configuration: %w", err)
	}
	if len(data) > maxConfigBytes {
		return Config{}, fmt.Errorf("configuration exceeds %d bytes", maxConfigBytes)
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	var c Config
	if err := d.Decode(&c); err != nil {
		return Config{}, fmt.Errorf("decode configuration: %w", err)
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return Config{}, fmt.Errorf("expected exactly one JSON document")
	}
	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}

func (c Config) Validate() error {
	if !validInterface(c.Management) {
		return fmt.Errorf("invalid management interface")
	}
	if len(c.Links) < 1 || len(c.Links) > 16 {
		return fmt.Errorf("expected 1 to 16 WAN links")
	}
	ids := make(map[string]bool)
	interfaces := map[string]bool{c.Management: true}
	for i, link := range c.Links {
		if !idPattern.MatchString(link.ID) || ids[link.ID] {
			return fmt.Errorf("link %d: invalid or duplicate ID", i+1)
		}
		if !validInterface(link.Interface) || interfaces[link.Interface] {
			return fmt.Errorf("link %s: invalid, duplicate or management interface", link.ID)
		}
		if link.Weight < 1 || link.Weight > 100 {
			return fmt.Errorf("link %s: weight must be 1 to 100", link.ID)
		}
		if err := link.validateAddress(); err != nil {
			return fmt.Errorf("link %s: %w", link.ID, err)
		}
		ids[link.ID] = true
		interfaces[link.Interface] = true
	}
	return nil
}

func validInterface(name string) bool {
	return name != "lo" && ifacePattern.MatchString(name)
}

func (l Link) validateAddress() error {
	switch l.Mode {
	case "dhcp":
		if l.Address != "" || l.Gateway != "" {
			return fmt.Errorf("DHCP must not specify static addressing")
		}
	case "static":
		prefix, err := netip.ParsePrefix(l.Address)
		if err != nil || !prefix.Addr().Is4() || prefix.Bits() == 0 || prefix.Bits() == 32 || !prefix.Addr().IsGlobalUnicast() {
			return fmt.Errorf("expected usable IPv4 interface CIDR with prefix length 1 to 31")
		}
		gateway, err := netip.ParseAddr(l.Gateway)
		if err != nil || !gateway.Is4() || !gateway.IsGlobalUnicast() || !prefix.Contains(gateway) || gateway == prefix.Addr() {
			return fmt.Errorf("expected distinct IPv4 gateway within interface subnet")
		}
		if prefix.Bits() < 31 {
			network := prefix.Masked().Addr().As4()
			broadcast := network
			for bit := prefix.Bits(); bit < 32; bit++ {
				broadcast[bit/8] |= 1 << (7 - bit%8)
			}
			for _, addr := range []netip.Addr{prefix.Addr(), gateway} {
				if addr == netip.AddrFrom4(network) || addr == netip.AddrFrom4(broadcast) {
					return fmt.Errorf("network and broadcast addresses are not valid endpoints")
				}
			}
		}
	default:
		return fmt.Errorf("mode must be dhcp or static")
	}
	return nil
}
