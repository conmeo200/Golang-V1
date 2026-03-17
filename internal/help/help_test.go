package help

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMethod(t *testing.T) {
	// Handler giả lập để kiểm tra xem nó có được gọi không
	handlerWasCalled := false
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerWasCalled = true
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Allowed Method", func(t *testing.T) {
		handlerWasCalled = false
		middleware := Method(http.MethodPost, mockHandler)
		
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()
		
		middleware.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, handlerWasCalled)
	})

	t.Run("Not Allowed Method", func(t *testing.T) {
		handlerWasCalled = false
		middleware := Method(http.MethodPost, mockHandler)
		
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		
		middleware.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.False(t, handlerWasCalled)
		assert.Contains(t, w.Body.String(), "method not allowed")
	})
}
