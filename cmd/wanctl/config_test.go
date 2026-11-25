package main

import (
	"bytes"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateConfigCommand(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
		valid bool
	}{
		{"valid", `{"management_interface":"eth0","links":[{"id":"wan-1","interface":"eth1","mode":"dhcp","weight":1}]}`, true},
		{"invalid", `{"management_interface":"eth0","links":[]}`, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			if err := os.WriteFile(path, []byte(tt.input), 0600); err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			err := run([]string{"validate-config", path}, &out, func() ([]net.Interface, error) {
				t.Fatal("config validation must not inspect live interfaces")
				return nil, nil
			})
			if (err == nil) != tt.valid {
				t.Fatalf("valid=%v; error=%v", tt.valid, err)
			}
			if tt.valid && out.String() != "configuration valid (not applied)\n" {
				t.Fatalf("unexpected output %q", out.String())
			}
			if !tt.valid && out.Len() != 0 {
				t.Fatalf("unexpected success output %q", out.String())
			}
		})
	}
}
