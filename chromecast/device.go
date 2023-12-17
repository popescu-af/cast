package chromecast

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"model/device"
	"model/library"
)

type Device struct {
	id         string
	properties map[string]string
	transport  *transport
	state      *state
	mtx        sync.Mutex
}

var (
	ErrFailedToConnect = fmt.Errorf("failed to connect to device")
)

func NewDevice(id string, ip net.IP, port int, kvSlice []string) (*Device, error) {
	properties := make(map[string]string, len(kvSlice))
	for _, v := range kvSlice {
		s := strings.SplitN(v, "=", 2)
		if len(s) == 2 {
			properties[s[0]] = s[1]
		}
	}
	ctx := context.Background()
	t, err := newTransport(ctx, ip, port)
	if err != nil {
		return nil, err
	}
	s, err := newState(ctx, t)
	if err != nil {
		return nil, err
	}
	t.startReceiving(s.handleReceivedMessage)
	return &Device{
		id:         id,
		properties: properties,
		transport:  t,
		state:      s,
	}, nil
}

func (d *Device) Close() error {
	_, err := d.synchronizeCall(
		5*time.Second,
		func() ([]interface{}, error) {
			if err := d.state.close(); err != nil {
				return nil, err
			}
			return nil, d.transport.close()
		},
		"request_stop",
	)
	return err
}

func (d *Device) ID() string {
	return d.id
}

func (d *Device) Load(media, subtitles library.Media) error {
	d.Stop()
	return d.synchronizeSimpleCommand(10*time.Second, "request_cast", media, subtitles)
}

func (d *Device) Play() error {
	return d.synchronizeSimpleCommand(5*time.Second, "request_control", "PLAY")
}

func (d *Device) Stop() error {
	return d.synchronizeSimpleCommand(5*time.Second, "request_stop")
}

func (d *Device) Pause() error {
	return d.synchronizeSimpleCommand(5*time.Second, "request_control", "PAUSE")
}

func (d *Device) Seek(absolutePercentage float64) error {
	return d.synchronizeSimpleCommand(10*time.Second, "request_control", "SEEK", seekTypeAbsolutePercentage, absolutePercentage)
}

func (d *Device) Fwd(offsetInSeconds float64) error {
	return d.synchronizeSimpleCommand(10*time.Second, "request_control", "SEEK", seekTypeRelative, offsetInSeconds)
}

func (d *Device) Rev(offsetInSeconds float64) error {
	return d.synchronizeSimpleCommand(10*time.Second, "request_control", "SEEK", seekTypeRelative, offsetInSeconds)
}

func (d *Device) GetPlaybackInfo() (device.PlaybackInfo, error) {
	retvals, err := d.synchronizeCall(
		5*time.Second,
		func() ([]interface{}, error) {
			return []interface{}{
				d.state.playbackInfo,
			}, nil
		},
		"request_control",
		"GET_STATUS",
	)
	if err != nil {
		return device.PlaybackInfo{}, err
	}
	return retvals[0].(device.PlaybackInfo), nil
}

func (d *Device) synchronizeSimpleCommand(
	timeout time.Duration,
	event string,
	eventArgs ...interface{},
) error {
	_, err := d.synchronizeCall(
		timeout,
		func() ([]interface{}, error) {
			return nil, nil
		},
		event,
		eventArgs...,
	)
	return err
}

func (d *Device) synchronizeCall(
	timeout time.Duration,
	afterFunc func() ([]interface{}, error),
	event string,
	eventArgs ...interface{},
) (
	[]interface{},
	error,
) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	d.state.commandSync.setWaiting(true)
	if err := d.state.emitEvent(event, eventArgs...); err != nil {
		d.state.commandSync.setWaiting(false)
		return nil, err
	}
	if err := d.state.commandSync.wait(timeout); err != nil {
		return nil, err
	}
	return afterFunc()
}
