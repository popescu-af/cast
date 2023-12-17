package goodboi

import (
	"context"

	"model/device"
)

func FindDevices(context.Context) ([]device.Device, error) {
	return []device.Device{
		&Device{},
	}, nil
}
