package server

import (
	"fmt"
	"github.com/PatrickKvartsaniy/print-me-at/errors"
	"github.com/PatrickKvartsaniy/print-me-at/models"
	"github.com/PatrickKvartsaniy/print-me-at/server/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	msg = "hello"
	ts  = "2020-09-02T17:57:15%2B03:00"
)

type TestSuite struct {
	suite.Suite
	ctrl     *gomock.Controller
	handler  http.Handler
	repoMock *mock.MockRepository
	server   *Server
	recorder *httptest.ResponseRecorder
}

func (s *TestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.server = &Server{repo: mock.NewMockRepository(s.ctrl)}
	s.repoMock = mock.NewMockRepository(s.ctrl)
	s.handler = http.HandlerFunc(s.server.messageScheduler)
	s.recorder = httptest.NewRecorder()
}

func (s *TestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s TestSuite) Test_ServerMessageScheduler(t *testing.T) {
	tcs := []struct {
		url          string
		expectedCode int
	}{
		{
			url:          "/printMeAt",
			expectedCode: http.StatusBadRequest,
		},
		{
			url:          fmt.Sprintf("/printMeAt?msg=%s&ts=%s", msg, ts),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tcs {
		s.repoMock.EXPECT().AddNewTask(gomock.Any()).Return(gomock.Eq(nil))
		req, err := http.NewRequest(http.MethodGet, tc.url, nil)
		assert.NoError(t, err)

		s.handler.ServeHTTP(s.recorder, req)
		assert.Equal(t, tc.expectedCode, s.recorder.Code)
	}
}

func Test_parseQueryParams(t *testing.T) {
	tcs := []struct {
		url string
		out models.Message
		err error
	}{
		{
			url: "/printMeAt",
			out: models.Message{},
			err: fmt.Errorf("%w: message is missing", errors.InvalidParameters),
		},
		{
			url: fmt.Sprintf("/printMeAt?msg=%s&ts=%s", msg, ts),
			out: models.Message{
				Value:     msg,
				Timestamp: 1599058635,
			},
			err: nil,
		},
	}
	for _, tc := range tcs {
		req, err := http.NewRequest(http.MethodGet, tc.url, nil)
		assert.NoError(t, err)
		res, err := parseQueryParams(req)
		if tc.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, tc.err.Error())
		}
		assert.Equal(t, tc.out.Timestamp, res.Timestamp)
		assert.Equal(t, tc.out.Value, res.Value)
	}
}
