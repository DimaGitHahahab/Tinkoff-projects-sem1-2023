package service

import (
	"github.com/stretchr/testify/assert"
	"homework/internal/model"
	"net"
	"strconv"
	"testing"
)

func TestSafeMapAdd(t *testing.T) {
	m := NewStorage()

	d := model.Device{SerialNum: "1", Model: "model1", IP: "1.1.1.1"}

	ok := m.Add(d)
	assert.True(t, ok)

	ok = m.Add(d)
	assert.False(t, ok)

	d = model.Device{SerialNum: "2", Model: "model1", IP: "1.1.1.1"}

	ok = m.Add(d)
	assert.True(t, ok)
}

func TestSafeMapGet(t *testing.T) {
	m := NewStorage()

	d := model.Device{SerialNum: "1", Model: "model1", IP: "1.1.1.1"}

	m.Add(d)

	gotDevice, ok := m.Get(d.SerialNum)
	assert.True(t, ok)
	assert.Equal(t, d, gotDevice)

	gotDevice, ok = m.Get("2")
	assert.False(t, ok)
	assert.Equal(t, model.Device{}, gotDevice)
}

func TestSafeMapDel(t *testing.T) {
	m := NewStorage()

	d := model.Device{SerialNum: "1", Model: "model1", IP: "1.1.1.1"}

	m.Add(d)

	ok := m.Del(d.SerialNum)
	assert.True(t, ok)

	ok = m.Del(d.SerialNum)
	assert.False(t, ok)
}

func FuzzVerifyDeviceData(f *testing.F) {
	f.Fuzz(func(t *testing.T, serialNum, model, ip string) {
		d := model.Device{SerialNum: serialNum, Model: model, IP: ip}
		res := verifyDeviceData(d)
		if res != nil {
			if d.Model == "" {
				assert.Equal(t, ErrInvalidModel, res)
			} else if d.SerialNum == "" {
				assert.Equal(t, ErrInvalidSerialNumber, res)
			} else if net.ParseIP(d.IP) == nil {
				assert.Equal(t, ErrInvalidIPAddress, res)
			}
		}
	})
}

func TestCreateDevice(t *testing.T) {
	storage := NewStorageMock(t)
	s := NewService(storage)

	wantDevice := model.Device{SerialNum: "123", Model: "model1", IP: "1.1.1.1"}

	storage.AddMock.Expect(wantDevice).Return(true)

	err := s.CreateDevice(wantDevice)
	assert.Nil(t, err)

	storage.GetMock.Expect(wantDevice.SerialNum).Return(wantDevice, true)

	gotDevice, err := s.GetDevice(wantDevice.SerialNum)
	assert.Nil(t, err)

	assert.Equal(t, wantDevice, gotDevice)
}

func TestCreateMultipleDevices(t *testing.T) {
	storage := NewStorageMock(t)
	s := NewService(storage)

	devices := []model.Device{
		{
			SerialNum: "1", Model: "model1", IP: "1.1.1.1",
		},
		{
			SerialNum: "2", Model: "model1", IP: "1.1.1.1",
		},
		{
			SerialNum: "3", Model: "model2", IP: "1.1.1.1",
		},
		{
			SerialNum: "3", Model: "model3", IP: "123.123.123.123",
		},
	}

	for _, d := range devices {
		storage.AddMock.Expect(d).Return(true)
		err := s.CreateDevice(d)
		assert.Nil(t, err)
	}

	for _, d := range devices {
		storage.GetMock.Expect(d.SerialNum).Return(d, true)
		gotDevice, err := s.GetDevice(d.SerialNum)
		assert.Nil(t, err)
		assert.Equal(t, d, gotDevice)
	}

}

func TestCreateInvalidSerialNum(t *testing.T) {
	storage := NewStorageMock(t)
	s := NewService(storage)

	invalidDevice := model.Device{SerialNum: "", Model: "model1", IP: "1.1.1.1"}

	storage.AddMock.Expect(invalidDevice).Return(true)

	err := s.CreateDevice(invalidDevice)
	assert.NotNil(t, err)
}
func TestCreateInvalidIP(t *testing.T) {
	storage := NewStorageMock(t)
	s := NewService(storage)

	invalidDevice := model.Device{SerialNum: "1", Model: "model1", IP: "1.99999.1.1"}

	storage.AddMock.Expect(invalidDevice).Return(true)

	err := s.CreateDevice(invalidDevice)
	assert.NotNil(t, err)
}

func TestCreateDuplicate(t *testing.T) {
	storage := NewStorageMock(t)
	s := NewService(storage)

	d := model.Device{SerialNum: "123", Model: "model1", IP: "1.1.1.1"}

	storage.AddMock.Expect(d).Return(true)
	err := s.CreateDevice(d)
	assert.Nil(t, err)

	storage.AddMock.Expect(d).Return(false)
	err = s.CreateDevice(d)
	assert.NotNil(t, err)
}

func TestGetDeviceUnexisting(t *testing.T) {
	storage := NewStorageMock(t)
	s := NewService(storage)

	serialNum := "000"

	storage.GetMock.Expect(serialNum).Return(model.Device{}, false)

	d, err := s.GetDevice(serialNum)
	assert.Equal(t, model.Device{}, d)
	assert.NotNil(t, err)

}

func TestDeleteDevice(t *testing.T) {
	storage := NewStorageMock(t)
	s := NewService(storage)

	d := model.Device{SerialNum: "1", Model: "model1", IP: "1.1.1.1"}

	storage.AddMock.Expect(d).Return(true)

	_ = s.CreateDevice(d)

	storage.DelMock.Expect(d.SerialNum).Return(true)

	err := s.DeleteDevice(d.SerialNum)
	assert.Nil(t, err)
}

func TestDeleteDeviceUnexisting(t *testing.T) {
	storage := NewStorageMock(t)
	s := NewService(storage)

	serialNum := "000"

	storage.DelMock.Expect(serialNum).Return(false)

	err := s.DeleteDevice(serialNum)
	assert.NotNil(t, err)
}

func TestUpdateDevice(t *testing.T) {
	storage := NewStorageMock(t)
	s := NewService(storage)

	d := model.Device{SerialNum: "1", Model: "model1", IP: "1.1.1.1"}

	storage.AddMock.Expect(d).Return(true)

	_ = s.CreateDevice(d)

	updDevice := model.Device{SerialNum: "1", Model: "model2 pro max", IP: "1.1.1.1"}

	storage.GetMock.Expect(updDevice.SerialNum).Return(d, true)
	storage.AddMock.Expect(updDevice).Return(true)

	err := s.UpdateDevice(updDevice)
	assert.Nil(t, err)
}

func TestUpdateDeviceUnexsting(t *testing.T) {
	storage := NewStorageMock(t)
	s := NewService(storage)

	updDevice := model.Device{SerialNum: "1", Model: "model1", IP: "1.1.1.1"}

	storage.GetMock.Expect(updDevice.SerialNum).Return(model.Device{}, false)

	err := s.UpdateDevice(updDevice)
	assert.NotNil(t, err)

}

// BenchmarkServiceUpdateDevice benchmarks the UpdateDevice method of Service.
func BenchmarkServiceUpdateDevice(b *testing.B) {
	storage := NewStorage()
	s := NewService(storage)

	d := model.Device{SerialNum: "123", Model: "model1", IP: "1.1.1.1"}

	_ = s.CreateDevice(d)

	for i := 2; i < b.N; i++ {
		newModel := d.Model[:len(d.Model)-1] + strconv.Itoa(i)
		d.Model = newModel

		err := s.UpdateDevice(d)
		assert.Nil(b, err)
	}
}
