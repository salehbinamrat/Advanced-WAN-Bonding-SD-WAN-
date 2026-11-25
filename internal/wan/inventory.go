// Package wan provides read-only network interface discovery.
package wan

import (
	"fmt"
	"net"
	"sort"
)

// Interface describes an interface, not a verified WAN uplink.
type Interface struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
	MTU   int    `json:"mtu"`
	Up    bool   `json:"administratively_up"`
}

// Inventory returns non-loopback interfaces in stable index order.
// Administrative up status does not imply carrier or Internet connectivity.
func Inventory(list func() ([]net.Interface, error)) ([]Interface, error) {
	interfaces, err := list()
	if err != nil {
		return nil, fmt.Errorf("enumerate interfaces: %w", err)
	}
	result := make([]Interface, 0, len(interfaces))
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		result = append(result, Interface{Index: iface.Index, Name: iface.Name, MTU: iface.MTU, Up: iface.Flags&net.FlagUp != 0})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Index < result[j].Index })
	return result, nil
}
