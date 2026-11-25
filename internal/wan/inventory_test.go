package wan

import (
	"errors"
	"net"
	"reflect"
	"testing"
)

func TestInventory(t *testing.T) {
	input := []net.Interface{
		{Index: 4, Name: "eth1", MTU: 1500, Flags: net.FlagUp},
		{Index: 1, Name: "lo", MTU: 65536, Flags: net.FlagLoopback | net.FlagUp},
		{Index: 2, Name: "eth0", MTU: 1500},
	}
	got, err := Inventory(func() ([]net.Interface, error) { return input, nil })
	want := []Interface{{Index: 2, Name: "eth0", MTU: 1500}, {Index: 4, Name: "eth1", MTU: 1500, Up: true}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, %v; want %+v", got, err, want)
	}
	if input[0].Name != "eth1" {
		t.Fatal("input mutated")
	}
}

func TestInventoryError(t *testing.T) {
	sentinel := errors.New("unavailable")
	got, err := Inventory(func() ([]net.Interface, error) { return nil, sentinel })
	if got != nil || !errors.Is(err, sentinel) {
		t.Fatalf("got %+v, %v", got, err)
	}
}

func TestInventoryEmpty(t *testing.T) {
	got, err := Inventory(func() ([]net.Interface, error) { return nil, nil })
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("got %+v, %v", got, err)
	}
}
