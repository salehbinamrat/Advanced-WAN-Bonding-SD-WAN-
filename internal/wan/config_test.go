package wan

import (
	"strings"
	"testing"
)

func TestConfigValidate(t *testing.T) {
	base := func() Config {
		return Config{Management: "eth0", Links: []Link{
			{ID: "wan-1", Interface: "eth1", Mode: "static", Address: "192.168.1.2/24", Gateway: "192.168.1.1", Weight: 50},
			{ID: "wan-2", Interface: "eth2", Mode: "static", Address: "192.168.1.2/24", Gateway: "192.168.1.1", Weight: 50},
		}}
	}
	tests := []struct {
		name   string
		change func(*Config)
		valid  bool
	}{
		{"overlapping addresses intentionally allowed", func(c *Config) {}, true},
		{"dhcp", func(c *Config) { c.Links[0] = Link{ID: "wan-1", Interface: "eth1", Mode: "dhcp", Weight: 1} }, true},
		{"duplicate ID", func(c *Config) { c.Links[1].ID = c.Links[0].ID }, false},
		{"duplicate interface", func(c *Config) { c.Links[1].Interface = c.Links[0].Interface }, false},
		{"management conflict", func(c *Config) { c.Links[0].Interface = c.Management }, false},
		{"loopback", func(c *Config) { c.Links[0].Interface = "lo" }, false},
		{"invalid weight", func(c *Config) { c.Links[0].Weight = 0 }, false},
		{"invalid mode", func(c *Config) { c.Links[0].Mode = "pppoe" }, false},
		{"dhcp with static address", func(c *Config) { c.Links[0].Mode = "dhcp" }, false},
		{"outside gateway", func(c *Config) { c.Links[0].Gateway = "192.168.2.1" }, false},
		{"same gateway", func(c *Config) { c.Links[0].Gateway = "192.168.1.2" }, false},
		{"broadcast gateway", func(c *Config) { c.Links[0].Gateway = "192.168.1.255" }, false},
		{"network address", func(c *Config) { c.Links[0].Address = "192.168.1.0/24" }, false},
		{"point to point", func(c *Config) { c.Links[0].Address = "192.168.1.2/31"; c.Links[0].Gateway = "192.168.1.3" }, true},
		{"IPv6 not supported yet", func(c *Config) { c.Links[0].Address = "2001:db8::2/64" }, false},
		{"empty links", func(c *Config) { c.Links = nil }, false},
		{"invalid management", func(c *Config) { c.Management = "" }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := base()
			tt.change(&c)
			if err := c.Validate(); (err == nil) != tt.valid {
				t.Fatalf("valid=%v; error=%v", tt.valid, err)
			}
		})
	}
}

func TestDecodeConfig(t *testing.T) {
	valid := `{"management_interface":"eth0","links":[{"id":"wan-1","interface":"eth1","mode":"dhcp","weight":1}]}`
	for _, tt := range []struct {
		name  string
		input string
		valid bool
	}{
		{"valid", valid, true},
		{"unknown field", strings.Replace(valid, `"weight"`, `"unknown"`, 1), false},
		{"trailing JSON", valid + " {}", false},
		{"oversize", valid + strings.Repeat(" ", maxConfigBytes), false},
		{"invalid JSON", "{", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeConfig(strings.NewReader(tt.input))
			if (err == nil) != tt.valid {
				t.Fatalf("valid=%v; error=%v", tt.valid, err)
			}
		})
	}
}
