package goodboi

import (
	"time"

	"model/device"
	"model/library"
)

var deviceID = "goodboi-device"

type Device struct {
}

func (d *Device) ID() string {
	return deviceID
}

func (d *Device) Load(media, subtitles library.Media) error {
	println("Loaded media", media.ID(), "successfully")
	return nil
}

func (d *Device) Play() error {
	println("Playing media successfully")
	return nil
}

func (d *Device) Stop() error {
	println("Stopped media successfully")
	return nil
}

func (d *Device) Pause() error {
	println("Paused media successfully")
	return nil
}

func (d *Device) Seek(absolutePercentage float64) error {
	println("Seeked media successfully at", absolutePercentage, "absolute percentage")
	return nil
}

func (d *Device) Fwd(offsetInSeconds float64) error {
	println("Forwarded media successfully by", offsetInSeconds, "seconds")
	return nil
}

func (d *Device) Rev(offsetInSeconds float64) error {
	println("Reversed media successfully by", offsetInSeconds, "seconds")
	return nil
}

func (d *Device) GetPlaybackInfo() (device.PlaybackInfo, error) {
	println("Getting media status successfully")
	return device.PlaybackInfo{
		Duration:   0,
		Position:   0,
		PositionTS: time.Now(),
		Playing:    false,
	}, nil
}

func (d *Device) Close() error {
	println("Closed device successfully")
	return nil
}
