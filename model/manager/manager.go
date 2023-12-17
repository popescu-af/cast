package manager

import (
	"errors"
	"sync"

	"logger"
	"model/device"
	"model/library"
	"model/session"
)

var (
	ErrDeviceAlreadyExists = errors.New("device already exists")
	ErrDeviceNotFound      = errors.New("device not found")
	ErrLibraryNotFound     = errors.New("library not found")
	ErrMediaNotFound       = errors.New("media not found")
	ErrSessionNotFound     = errors.New("session not found")
)

type Manager interface {
	// AddDevice adds a new device.
	AddDevice(device device.Device) error
	// ListDevices returns the list of devices.
	ListDevices() map[string]device.Device
	// AddLibrary adds a new library.
	AddLibrary(library library.Library)
	// ListLibraries returns the list of libraries.
	ListLibraries() map[string]library.Library
	// ListMedia returns the list of media in the library.
	ListMedia(libraryID string) ([]library.Media, error)
	// GetMedia returns the media with the given ID.
	GetMedia(mediaID string) (library.Media, error)
	// CreateSession creates a new session.
	CreateSession(deviceID, mediaID, subtitleID string) (session.Session, error)
}

func New() Manager {
	return &manager{
		devices:   make(map[string]device.Device),
		libraries: make(map[string]library.Library),
	}
}

type manager struct {
	mtx       sync.RWMutex
	devices   map[string]device.Device
	libraries map[string]library.Library
}

func (m *manager) AddDevice(device device.Device) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if oldDevice, ok := m.devices[device.ID()]; ok {
		if err := oldDevice.Close(); err != nil {
			logger.Log.Printf("failed to close device '%s': %s", oldDevice.ID(), err)
		}
	}
	m.devices[device.ID()] = device
	return nil
}

func (m *manager) ListDevices() map[string]device.Device {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return copyMap(m.devices)
}

func (m *manager) AddLibrary(library library.Library) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.libraries[library.ID()] = library
}

func (m *manager) ListLibraries() map[string]library.Library {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return copyMap(m.libraries)
}

func (m *manager) GetMedia(mediaID string) (library.Media, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	for _, library := range m.libraries {
		mediaList, err := library.List()
		if err != nil {
			return nil, err
		}
		for _, m := range mediaList {
			if m.ID() == mediaID {
				return m, nil
			}
		}
	}
	return nil, ErrMediaNotFound
}

func (m *manager) ListMedia(libraryID string) ([]library.Media, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	for _, library := range m.libraries {
		if library.ID() == libraryID {
			return library.List()
		}
	}
	return nil, ErrLibraryNotFound
}

func (m *manager) CreateSession(deviceID, mediaID, subtitleID string) (session.Session, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	var device device.Device
	for _, d := range m.devices {
		if d.ID() == deviceID {
			device = d
			break
		}
	}
	if device == nil {
		return nil, ErrDeviceNotFound
	}

	var media, subtitle library.Media
	for _, library := range m.libraries {
		mediaList, err := library.List()
		if err != nil {
			return nil, err
		}
		for _, m := range mediaList {
			if m.ID() == mediaID {
				media = m
			}
			if subtitleID != "" && m.ID() == subtitleID {
				subtitle = m
			}
			if media != nil && (subtitleID == "" || subtitle != nil) {
				goto MediaFound
			}
		}
	}
MediaFound:
	if media == nil {
		return nil, ErrMediaNotFound
	}
	if subtitleID != "" && subtitle == nil {
		return nil, ErrMediaNotFound
	}
	session, err := session.New(device, media, subtitle)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func copyMap[T any](m map[string]T) map[string]T {
	result := make(map[string]T, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}
