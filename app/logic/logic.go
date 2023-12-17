package logic

import (
	"context"
	"time"

	"chromecast"
	"goodboi"
	"model/device"
	"model/manager"
)

type scanningFunc func(context.Context) ([]device.Device, error)

var (
	supportedDevices = map[string]scanningFunc{
		"goodboi":    goodboi.FindDevices,
		"chromecast": chromecast.FindDevices,
	}

	SupportedExtensions = []string{".mp4", ".srt"}
)

type Logic struct {
	Manager manager.Manager
}

func New(localDiskLibrary string) (*Logic, error) {
	return &Logic{
		Manager: manager.New(),
	}, nil
}

func (l *Logic) FindDevices(timeout time.Duration, types []string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var messages []string

	var typeScanners map[string]scanningFunc
	if len(types) != 0 {
		typeScanners = make(map[string]scanningFunc)
		for _, t := range types {
			f, ok := supportedDevices[t]
			if !ok {
				message := "devices of type " + t + " are not supported"
				messages = append(messages, message)
				continue
			}
			typeScanners[t] = f
		}
		if len(typeScanners) == 0 {
			message := "none of the requested types are supported"
			messages = append(messages, message)
			return messages
		}
	} else {
		typeScanners = supportedDevices
	}

	for t, find := range typeScanners {
		found, err := find(ctx)
		if err != nil {
			message := "failed to find devices of type " + t + ": " + err.Error()
			messages = append(messages, message)
			continue
		}
		for _, device := range found {
			if err := l.Manager.AddDevice(device); err != nil {
				message := "failed to add device " + device.ID() + ": " + err.Error()
				messages = append(messages, message)
				continue
			}
			message := "added device " + device.ID()
			messages = append(messages, message)
		}
	}
	return messages
}
