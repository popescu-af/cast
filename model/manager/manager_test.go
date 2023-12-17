package manager_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"model/library"
	"model/manager"
	"model/mocks"
)

func TestDeviceManagement(t *testing.T) {
	// Create a new device.
	deviceID := uuid.Must(uuid.NewRandom()).String()
	d := &mocks.Device{}
	d.On("ID").Return(deviceID)

	// Create a new manager and add the device.
	m := manager.New()
	require.NoError(t, m.AddDevice(d))
	require.Equal(t, manager.ErrDeviceAlreadyExists, m.AddDevice(d))

	// Verify that the device was added.
	devices := m.ListDevices()
	require.Len(t, devices, 1)

	// Verify that the device has the correct ID.
	_, ok := devices[deviceID]
	require.True(t, ok)
}

func TestLibraryAndMediaManagement(t *testing.T) {
	// Create some media.
	mediaIDs := []string{
		uuid.Must(uuid.NewRandom()).String(),
		uuid.Must(uuid.NewRandom()).String(),
		uuid.Must(uuid.NewRandom()).String(),
	}
	media := make([]library.Media, 3)
	for i := range media {
		m := &mocks.Media{}
		m.On("ID").Return(mediaIDs[i])
		media[i] = m
	}

	// Create some libraries.
	libID1 := uuid.Must(uuid.NewRandom()).String()
	lib1 := &mocks.Library{}
	lib1.On("ID").Return(libID1)
	lib1.On("List").Return(media[:2], nil)

	libID2 := uuid.Must(uuid.NewRandom()).String()
	lib2 := &mocks.Library{}
	lib2.On("ID").Return(libID2)
	lib2.On("List").Return(media[2:], nil)

	// Create a new manager and add the libraries.
	m := manager.New()
	m.AddLibrary(lib1)
	m.AddLibrary(lib2)

	// Verify that the libraries were added.
	libraries := m.ListLibraries()
	require.Len(t, libraries, 2)
	_, ok := libraries[libID1]
	require.True(t, ok)
	_, ok = libraries[libID2]
	require.True(t, ok)

	// Verify that listing media for an inexistent library fails.
	_, err := m.ListMedia(uuid.Must(uuid.NewRandom()).String())
	require.Equal(t, manager.ErrLibraryNotFound, err)

	// Verify that the correct media are returned for each library.
	media1, err := m.ListMedia(libID1)
	require.NoError(t, err)
	require.Len(t, media1, 2)
	require.Equal(t, mediaIDs[0], media1[0].ID())
	require.Equal(t, mediaIDs[1], media1[1].ID())

	media2, err := m.ListMedia(libID2)
	require.NoError(t, err)
	require.Len(t, media2, 1)
	require.Equal(t, mediaIDs[2], media2[0].ID())

	// Verify that getting inexistent media fails.
	_, err = m.GetMedia(uuid.Must(uuid.NewRandom()).String())
	require.Equal(t, manager.ErrMediaNotFound, err)

	// Verify that the correct media are returned.
	for _, id := range mediaIDs {
		m, err := m.GetMedia(id)
		require.NoError(t, err)
		require.Equal(t, id, m.ID())
	}
}

func TestSessionManagement(t *testing.T) {
	// Create a new device.
	deviceID := uuid.Must(uuid.NewRandom()).String()
	d := &mocks.Device{}
	d.On("ID").Return(deviceID)

	// Create a new media.
	mediaID := uuid.Must(uuid.NewRandom()).String()
	media := &mocks.Media{}
	media.On("ID").Return(mediaID)

	// Create a new subtitle.
	subtitleID := uuid.Must(uuid.NewRandom()).String()
	subtitle := &mocks.Media{}
	subtitle.On("ID").Return(subtitleID)

	// Create a new library.
	libID := uuid.Must(uuid.NewRandom()).String()
	l := &mocks.Library{}
	l.On("ID").Return(libID)

	// Create a new manager and add the device and library.
	m := manager.New()
	m.AddDevice(d)
	m.AddLibrary(l)

	var errLibrary = errors.New("library error")
	var errDevice = errors.New("device error")
	sessionInputAndExpectedOutput := []struct {
		deviceID   string
		mediaID    string
		subtitleID string
		libraryErr error
		deviceErr  error
		err        error
	}{
		{uuid.Must(uuid.NewRandom()).String(), mediaID, subtitleID, nil, nil, manager.ErrDeviceNotFound},
		{uuid.Must(uuid.NewRandom()).String(), mediaID, "", nil, nil, manager.ErrDeviceNotFound},
		{deviceID, uuid.Must(uuid.NewRandom()).String(), subtitleID, nil, nil, manager.ErrMediaNotFound},
		{deviceID, uuid.Must(uuid.NewRandom()).String(), "", nil, nil, manager.ErrMediaNotFound},
		{deviceID, mediaID, uuid.Must(uuid.NewRandom()).String(), nil, nil, manager.ErrMediaNotFound},
		{deviceID, mediaID, subtitleID, errLibrary, nil, errLibrary},
		{deviceID, mediaID, "", errLibrary, nil, errLibrary},
		{deviceID, mediaID, subtitleID, nil, errDevice, errDevice},
		{deviceID, mediaID, "", nil, errDevice, errDevice},
		{deviceID, mediaID, subtitleID, nil, nil, nil},
		{deviceID, mediaID, "", nil, nil, nil},
	}

	// Create a new session.
	for _, input := range sessionInputAndExpectedOutput {
		d.On("ID").Return(deviceID)
		if input.deviceID == deviceID {
			l.On("List").Return([]library.Media{media, subtitle}, input.libraryErr)
			media.On("ID").Return(mediaID)
			if input.subtitleID != "" {
				subtitle.On("ID").Return(subtitleID)
			}
		}
		if input.subtitleID == "" {
			d.On("Load", media, nil).Return(input.deviceErr)
		} else {
			d.On("Load", media, subtitle).Return(input.deviceErr)
		}
		_, err := m.CreateSession(input.deviceID, input.mediaID, input.subtitleID)
		require.Equal(t, input.err, err)

		// Reset mock expectations.
		media.ExpectedCalls = []*mock.Call{}
		l.ExpectedCalls = []*mock.Call{}
		d.ExpectedCalls = []*mock.Call{}
	}
}
