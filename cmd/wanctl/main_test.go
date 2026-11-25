package main

import (
	"bytes"
	"errors"
	"net"
	"testing"
)

func TestRun(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"interfaces"}, &out, func() ([]net.Interface, error) {
		return []net.Interface{{Index: 2, Name: "eth0", MTU: 1500, Flags: net.FlagUp}}, nil
	})
	if err != nil || out.String() != "[{\"index\":2,\"name\":\"eth0\",\"mtu\":1500,\"administratively_up\":true}]\n" {
		t.Fatalf("got %q, %v", out.String(), err)
	}
}

func TestInvalidCommand(t *testing.T) {
	var out bytes.Buffer
	err := run(nil, &out, func() ([]net.Interface, error) { t.Fatal("unexpected enumeration"); return nil, nil })
	if err == nil || out.Len() != 0 {
		t.Fatal("expected usage error and no output")
	}
}

func TestRunEnumerationError(t *testing.T) {
	var out bytes.Buffer
	sentinel := errors.New("enumeration failed")
	err := run([]string{"interfaces"}, &out, func() ([]net.Interface, error) { return nil, sentinel })
	if !errors.Is(err, sentinel) || out.Len() != 0 {
		t.Fatalf("got %v, %q", err, out.String())
	}
}
