package session

import (
	"github.com/google/uuid"

	"model/device"
	"model/library"
)

type Session interface {
	// ID returns the ID of the session.
	ID() string
	// Play resumes playing the current media.
	Play(fromBeginning bool) error
	// Pause pauses the current media.
	Pause() error
	// Stop stops the current media.
	Stop() error
	// Seek seeks to the given absolute percentage of the media.
	Seek(absolutePercentage float64) error
	// Fwd seeks forward by the given offset in seconds.
	Fwd(offsetInSeconds float64) error
	// Rev seeks backward by the given offset in seconds.
	Rev(offsetInSeconds float64) error
	// GetPlaybackInfo returns info about the current playback.
	GetPlaybackInfo() (device.PlaybackInfo, error)
}

func New(device device.Device, media, subtitles library.Media) (Session, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	s := &session{
		id:        id.String(),
		device:    device,
		media:     media,
		subtitles: subtitles,
	}
	if err := s.device.Load(s.media, s.subtitles); err != nil {
		return nil, err
	}
	return s, nil
}

type session struct {
	id        string
	device    device.Device
	media     library.Media
	subtitles library.Media
}

func (s *session) ID() string {
	return s.id
}

func (s *session) Play(fromBeginning bool) error {
	if fromBeginning {
		if err := s.device.Seek(0); err != nil {
			return err
		}
	}
	return s.device.Play()
}

func (s *session) Pause() error {
	return s.device.Pause()
}

func (s *session) Stop() error {
	return s.device.Stop()
}

func (s *session) Seek(absolutePercentage float64) error {
	return s.device.Seek(absolutePercentage)
}

func (s *session) Fwd(offsetInSeconds float64) error {
	return s.device.Fwd(offsetInSeconds)
}

func (s *session) Rev(offsetInSeconds float64) error {
	return s.device.Rev(offsetInSeconds)
}

func (s *session) GetPlaybackInfo() (device.PlaybackInfo, error) {
	return s.device.GetPlaybackInfo()
}
