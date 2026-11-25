# Advanced WAN Bonding & SD-WAN

Foundation prototype for a standalone Linux x86-64 WAN appliance, independent of MikroTik subscriber management. Not production-ready; does not modify networking or forward traffic.

## Architecture and constraints

- Go control plane; Linux kernel data plane. Prefer standard library; justify every external dependency.
- Debian Stable Minimal is the candidate platform; validate against OpenWrt before committing to a bonding engine.
- React/TypeScript for the management UI only when implemented. No empty modules or placeholder services.
- Evaluate OpenMPTCProuter and MQVPN on identical hardware. Peplink, VeloCloud, Fortinet and other documented products are functional references, not code sources. Review upstream licensing before reuse.
- Test first, then implement. Linux Core and Linux Release run in parallel; Release tests precede packaging. Require both checks before merge. Never bypass failures. Keep generated artifacts out of source control.
- Separate management from forwarding; changes require rollback and explicit operator authorization. Preserve MikroTik RouterOS v6/v7 subscriber settings. Do not expose management services on WAN.

## Current slice: read-only interface inventory

`go run ./cmd/wanctl interfaces` prints non-loopback interface inventory as JSON. Administrative status is not carrier or internet health, and no interface is assumed to be a WAN. No elevated privileges needed.

## Local checks

```sh
gofmt -l cmd internal
go vet ./...
go test -race ./...
go build ./cmd/wanctl
```

The module has no external Go dependencies. The Go 1.23 language baseline is compatible with the local test toolchain; CI selects the current stable Go toolchain. Pin a supported exact release and action SHAs before production release.

## Next gates

1. Confirm repository permissions, supported toolchain and license; review baseline CI and enforce branch protection.
2. Design explicit WAN identities and validation; isolate duplicate gateways using Linux routing policy / namespaces in a disposable lab.
3. Add multi-target health monitoring, then weighted per-flow routing and safe failover with rollback.
4. Compare actual tunnel bonding engines with single-flow TCP, UDP, impairment and failure tests. Never represent load balancing as true bonding.
5. Add QoS, Starlink telemetry, secure management, optional MikroTik integration and UI after the networking foundations pass.
6. Validate hardware, recovery, security and long-duration tests before producing installable images. The initial release artifact is only a diagnostic CLI, not a router OS.

## WAN configuration validation

`wanctl validate-config FILE` validates a local JSON document without applying it:

```json
{"management_interface":"eth0","links":[{"id":"wan-1","interface":"eth1","mode":"dhcp","weight":50},{"id":"wan-2","interface":"eth2","mode":"dhcp","weight":50}]}
```

Supports 1–16 explicit WAN identities, unique interfaces, separate management interface, weights 1–100, DHCP and on-link static IPv4 CIDRs. Identical addresses and gateways across different interfaces are permitted deliberately; future isolated routing must be implemented and tested before applying settings. PPPoE, IPv6, duplicate-key rejection, live device validation and network mutation are not implemented in this slice.
