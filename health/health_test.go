package health

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServe_Error(t *testing.T) {
	service := &Server{
		checks: []Check{
			func() error {
				return fmt.Errorf("some err")
			},
		},
	}

	rr := httptest.NewRecorder()
	rq := httptest.NewRequest("GET", "/health", nil)
	service.serve(rr, rq)

	res := rr.Body.String()
	assert.Contains(t, res, "some err")
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestServe_Success(t *testing.T) {
	service := &Server{
		checks: []Check{},
	}

	rr := httptest.NewRecorder()
	rq := httptest.NewRequest("GET", "/health", nil)
	service.serve(rr, rq)

	res := rr.Body.String()
	assert.Empty(t, res)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}
