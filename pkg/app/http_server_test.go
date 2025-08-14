package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adeynack/finances/pkg/api/apiserver"
	"github.com/stretchr/testify/require"
)

func TestHttpServer(t *testing.T) {
	handler := mustCreateHandler()
	// server := httptest.NewServer(handler)
	// t.Cleanup(server.Close)

	t.Run("GET /health", func(t *testing.T) {
		request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/health", nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		require.Equal(t, http.StatusOK, response.Code)
		require.JSONEq(t, `{"status": "healthy"}`, response.Body.String())
	})

	t.Run("GET /books", func(t *testing.T) {
		request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/books", nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		require.Equal(t, http.StatusOK, response.Code, response.Body.String())

		var body apiserver.GetBooks200JSONResponse
		err := json.Unmarshal(response.Body.Bytes(), &body)
		require.NoError(t, err)

		require.Len(t, body.Books, 1)
		require.Equal(t, "Joe's Book", body.Books[0].Name)
		require.Equal(t, "Joe", body.Books[0].OwnerDisplayName)
	})
}
