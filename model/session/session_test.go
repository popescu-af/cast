package session_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"model/device"
	"model/mocks"
	"model/session"
)

func TestSession(t *testing.T) {
	// Create a new device.
	deviceID := uuid.Must(uuid.NewRandom())
	d := &mocks.Device{}
	d.On("ID").Return(deviceID)

	// Create a new media.
	mediaID := uuid.Must(uuid.NewRandom())
	media := &mocks.Media{}
	media.On("ID").Return(mediaID)
	subtitleID := uuid.Must(uuid.NewRandom())
	subtitle := &mocks.Media{}
	subtitle.On("ID").Return(subtitleID)

	var errDevice = errors.New("device error")

	// Create a new session with load error fails.
	d.On("Load", media, subtitle).Return(errDevice)
	_, err := session.New(d, media, subtitle)
	require.Equal(t, errDevice, err)
	d.ExpectedCalls = []*mock.Call{}

	// Create a new session.
	d.On("Load", media, subtitle).Return(nil)
	s, err := session.New(d, media, subtitle)
	require.NoError(t, err)

	// Test session methods.
	d.On("Play").Return(errDevice)
	require.Equal(t, errDevice, s.Play(false))
	d.ExpectedCalls = []*mock.Call{}

	d.On("Play").Return(nil)
	require.NoError(t, s.Play(false))
	d.ExpectedCalls = []*mock.Call{}

	d.On("Seek", 0.0).Return(errDevice)
	require.Equal(t, errDevice, s.Play(true))
	d.ExpectedCalls = []*mock.Call{}

	d.On("Seek", 0.0).Return(nil)
	d.On("Play").Return(errDevice)
	require.Equal(t, errDevice, s.Play(true))
	d.ExpectedCalls = []*mock.Call{}

	d.On("Seek", 0.0).Return(nil)
	d.On("Play").Return(nil)
	require.NoError(t, s.Play(true))
	d.ExpectedCalls = []*mock.Call{}

	d.On("Pause").Return(errDevice)
	require.Equal(t, errDevice, s.Pause())
	d.ExpectedCalls = []*mock.Call{}

	d.On("Pause").Return(nil)
	require.NoError(t, s.Pause())
	d.ExpectedCalls = []*mock.Call{}

	d.On("Stop").Return(errDevice)
	require.Equal(t, errDevice, s.Stop())
	d.ExpectedCalls = []*mock.Call{}

	d.On("Stop").Return(nil)
	require.NoError(t, s.Stop())
	d.ExpectedCalls = []*mock.Call{}

	d.On("Seek", 50.1).Return(errDevice)
	require.Equal(t, errDevice, s.Seek(50.1))
	d.ExpectedCalls = []*mock.Call{}

	d.On("Seek", 50.1).Return(nil)
	require.NoError(t, s.Seek(50.1))
	d.ExpectedCalls = []*mock.Call{}

	d.On("Fwd", 10.4).Return(errDevice)
	require.Equal(t, errDevice, s.Fwd(10.4))
	d.ExpectedCalls = []*mock.Call{}

	d.On("Fwd", 10.4).Return(nil)
	require.NoError(t, s.Fwd(10.4))
	d.ExpectedCalls = []*mock.Call{}

	d.On("Rev", 10.4).Return(errDevice)
	require.Equal(t, errDevice, s.Rev(10.4))
	d.ExpectedCalls = []*mock.Call{}

	d.On("Rev", 10.4).Return(nil)
	require.NoError(t, s.Rev(10.4))
	d.ExpectedCalls = []*mock.Call{}

	d.On("GetPlaybackInfo").Return(device.PlaybackInfo{}, errDevice)
	_, err = s.GetPlaybackInfo()
	require.Equal(t, errDevice, err)
	d.ExpectedCalls = []*mock.Call{}

	d.On("GetPlaybackInfo").Return(device.PlaybackInfo{}, nil)
	_, err = s.GetPlaybackInfo()
	require.NoError(t, err)
	d.ExpectedCalls = []*mock.Call{}
}
