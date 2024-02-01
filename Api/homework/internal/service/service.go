package service

import (
	"errors"
	"homework/internal/model"
	"net"
	"sync"
)

var (
	ErrDeviceAlreadyExists = errors.New("device already exists")
	ErrDeviceDoesNotExist  = errors.New("device doesn't exist")
	ErrInvalidModel        = errors.New("invalid model")
	ErrInvalidSerialNumber = errors.New("invalid serial number")
	ErrInvalidIPAddress    = errors.New("invalid IP address")
)

type Service interface {
	GetDevice(string) (model.Device, error)
	CreateDevice(model.Device) error
	DeleteDevice(string) error
	UpdateDevice(model.Device) error
}

func NewService(s Storage) Service {
	return &storageService{devices: s}
}

type storageService struct {
	devices Storage
}

func (s *storageService) GetDevice(num string) (model.Device, error) {
	d, ok := s.devices.Get(num)
	if !ok {
		return model.Device{}, ErrDeviceDoesNotExist
	}
	return d, nil
}

func (s *storageService) CreateDevice(d model.Device) error {
	ok := s.devices.Add(d)
	if !ok {
		return ErrDeviceAlreadyExists
	}

	if err := verifyDeviceData(d); err != nil {
		return err
	}

	return nil
}
func verifyDeviceData(d model.Device) error {
	if d.Model == "" {
		return ErrInvalidModel
	}
	if d.SerialNum == "" {
		return ErrInvalidSerialNumber
	}
	if net.ParseIP(d.IP) == nil {
		return ErrInvalidIPAddress
	}
	return nil
}

func (s *storageService) DeleteDevice(num string) error {
	ok := s.devices.Del(num)
	if !ok {
		return ErrDeviceDoesNotExist
	}
	return nil
}

func (s *storageService) UpdateDevice(updDev model.Device) error {
	_, ok := s.devices.Get(updDev.SerialNum)
	if !ok {
		return ErrDeviceDoesNotExist
	}
	if err := verifyDeviceData(updDev); err != nil {
		return err
	}
	s.devices.Add(updDev)
	return nil
}

type Storage interface {
	Add(d model.Device) bool
	Get(num string) (model.Device, bool)
	Del(num string) bool
}

func NewStorage() Storage {
	return &SafeMap{devices: make(map[string]model.Device), mu: sync.RWMutex{}}
}

type SafeMap struct {
	devices map[string]model.Device
	mu      sync.RWMutex
}

func (m *SafeMap) Add(d model.Device) bool {
	m.mu.RLock()
	_, ok := m.devices[d.SerialNum]
	m.mu.RUnlock()
	m.mu.Lock()
	m.devices[d.SerialNum] = d
	m.mu.Unlock()
	return !ok
}

func (m *SafeMap) Get(num string) (model.Device, bool) {
	m.mu.RLock()
	d, ok := m.devices[num]
	m.mu.RUnlock()
	if !ok {
		return model.Device{}, false
	}
	return d, true
}

func (m *SafeMap) Del(num string) bool {
	m.mu.RLock()
	_, ok := m.devices[num]
	m.mu.RUnlock()
	if !ok {
		return false
	}
	m.mu.Lock()
	delete(m.devices, num)
	m.mu.Unlock()
	return true
}
