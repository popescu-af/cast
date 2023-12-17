package device

import (
	"model/library"
	"time"
)

type Device interface {
	// ID returns the ID of the device.
	ID() string
	// Load loads the given media.
	Load(media, subtitles library.Media) error
	// Play starts playing the current media.
	Play() error
	// Stop stops playing the current media.
	Stop() error
	// Pause pauses the current media.
	Pause() error
	// Seek seeks to the given absolute percentage of the media.
	Seek(absolutePercentage float64) error
	// Fwd seeks forward by the given offset in seconds.
	Fwd(offsetInSeconds float64) error
	// Rev seeks backward by the given offset in seconds.
	Rev(offsetInSeconds float64) error
	// GetPlaybackInfo returns information about the current playback.
	GetPlaybackInfo() (PlaybackInfo, error)
	// Close closes the device.
	Close() error
}

type PlaybackInfo struct {
	// Duration is the duration of the media in seconds.
	Duration float64
	// Position is the current position of the media in seconds.
	Position float64
	// PositionTS is the timestamp of the last position update.
	PositionTS time.Time
	// Playing indicates whether the media is currently playing.
	Playing bool
}
