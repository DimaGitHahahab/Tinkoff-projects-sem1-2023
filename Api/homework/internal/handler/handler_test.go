package handler

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework/internal/model"
	"homework/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HandlerSuite struct {
	suite.Suite
	service *ServiceMock
	h       Handler
	r       *httptest.ResponseRecorder
}

func (s *HandlerSuite) SetupTest() {
	s.service = NewServiceMock(s.T())
	s.h = NewHandler(s.service)
	s.r = httptest.NewRecorder()
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (s *HandlerSuite) TestHandleCreate() {
	d := model.Device{SerialNum: "12345", Model: "TestModel", IP: "1.1.1.1"}

	payload, _ := json.Marshal(d)
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(payload))

	s.service.CreateDeviceMock.Expect(d).Return(nil)
	s.h.HandleCreate(s.r, req)

	assert.Equal(s.T(), http.StatusCreated, s.r.Code)
}

func (s *HandlerSuite) TestHandleGet() {
	d := model.Device{SerialNum: "12345", Model: "TestModel", IP: "1.1.1.1"}

	req := httptest.NewRequest(http.MethodGet, "/get?num=12345", nil)

	s.service.GetDeviceMock.Expect("12345").Return(d, nil)
	s.h.HandleGet(s.r, req)

	assert.Equal(s.T(), http.StatusOK, s.r.Code)
}

func (s *HandlerSuite) TestHandleDelete() {
	req := httptest.NewRequest(http.MethodDelete, "/delete?num=12345", nil)

	s.service.DeleteDeviceMock.Expect("12345").Return(nil)
	s.h.HandleDelete(s.r, req)

	assert.Equal(s.T(), http.StatusOK, s.r.Code)
}

func (s *HandlerSuite) TestHandleUpdate() {
	d := model.Device{SerialNum: "12345", Model: "TestModel", IP: "1.1.1.1"}

	payload, _ := json.Marshal(d)
	req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(payload))

	s.service.UpdateDeviceMock.Expect(d).Return(nil)
	s.h.HandleUpdate(s.r, req)

	assert.Equal(s.T(), http.StatusOK, s.r.Code)
}

func (s *HandlerSuite) TestHandleCreateInvalidRequest() {
	req := httptest.NewRequest(http.MethodPost, "/create", nil)

	s.h.HandleCreate(s.r, req)

	assert.Equal(s.T(), http.StatusBadRequest, s.r.Code)
}

func (s *HandlerSuite) TestHandleGetInvalidRequest() {
	req := httptest.NewRequest(http.MethodGet, "/get", nil)

	s.service.GetDeviceMock.Expect("").Return(model.Device{}, service.ErrDeviceDoesNotExist)
	s.h.HandleGet(s.r, req)

	assert.Equal(s.T(), http.StatusNotFound, s.r.Code)
}

func (s *HandlerSuite) TestHandleDeleteInvalidRequest() {
	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)

	s.service.DeleteDeviceMock.Expect("").Return(service.ErrDeviceDoesNotExist)
	s.h.HandleDelete(s.r, req)

	assert.Equal(s.T(), http.StatusNotFound, s.r.Code)
}

func (s *HandlerSuite) TestHandleUpdateInvalidRequest() {
	req := httptest.NewRequest(http.MethodPut, "/update", nil)

	s.h.HandleUpdate(s.r, req)

	assert.Equal(s.T(), http.StatusBadRequest, s.r.Code)
}

func (s *HandlerSuite) TestHandleUpdateUnexisting() {
	d := model.Device{SerialNum: "000", Model: "000", IP: "1.1.1.1"}
	payload, _ := json.Marshal(d)
	req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(payload))
	s.service.UpdateDeviceMock.Expect(d).Return(service.ErrDeviceDoesNotExist)
	s.h.HandleUpdate(s.r, req)

	assert.Equal(s.T(), http.StatusNotFound, s.r.Code)
}

func (s *HandlerSuite) TestErrResponse() {
	s.h.ErrResponse(s.r, "test", http.StatusBadRequest)
	assert.Equal(s.T(), http.StatusBadRequest, s.r.Code)
}

func (s *HandlerSuite) TestInvalidIP() {
	d := model.Device{SerialNum: "12345", Model: "TestModel", IP: "1.9999.1"}

	payload, _ := json.Marshal(d)
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(payload))

	s.service.CreateDeviceMock.Expect(d).Return(service.ErrInvalidIPAddress)
	s.h.HandleCreate(s.r, req)

	assert.Equal(s.T(), http.StatusBadRequest, s.r.Code)
}
