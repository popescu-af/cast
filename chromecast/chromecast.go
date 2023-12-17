package chromecast

import (
	"context"
	"net"
	"strings"

	"github.com/grandcat/zeroconf"

	"logger"
	"model/device"
)

func FindDevices(ctx context.Context) ([]device.Device, error) {
	resolver, err := zeroconf.NewResolver()
	if err != nil {
		return nil, err
	}
	entriesChannel := make(chan *zeroconf.ServiceEntry, 5)
	err = resolver.Browse(ctx, "_googlecast._tcp", "local", entriesChannel)
	if err != nil {
		return nil, err
	}
	var entries []*zeroconf.ServiceEntry
	for {
		select {
		case <-ctx.Done():
			goto CreateDevices
		case e := <-entriesChannel:
			entries = append(entries, e)
		}
	}
CreateDevices:
	devices := make([]device.Device, 0, len(entries))
	for _, entry := range entries {
		if !strings.Contains(entry.Service, "_googlecast.") {
			logger.Log.Printf("skipping '%s' as it does not contain '_googlecast.'", entry.Service)
			continue
		}
		var ip net.IP
		if len(entry.AddrIPv6) > 0 {
			ip = entry.AddrIPv6[0]
		} else if len(entry.AddrIPv4) > 0 {
			ip = entry.AddrIPv4[0]
		}
		device, err := NewDevice(entry.Instance, ip, entry.Port, entry.Text)
		if err != nil {
			logger.Log.Printf("failed to instantiate '%s' with error %v", entry.Service, err)
			continue
		}
		devices = append(devices, device)
	}
	return devices, nil
}
