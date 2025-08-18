package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
		const expectedBody = `{
			"books": [
				{
					"created_at": "2025-08-14T21:47:58.211393Z",
					"default_currency_iso_code": "EUR",
					"id": "8d8666c0-016f-49fb-8f59-4150a822ffb2",
					"name": "Joe's Book",
					"owner_display_name": "Joe",
					"owner_id": "569bcfdd-4056-42cd-af9c-285fa5ce92c8",
					"updated_at": "2025-08-14T21:47:58.211393Z"
				}
			]
		}`
		require.JSONEq(t, expectedBody, response.Body.String(), response.Body.String())
	})

	t.Run("GET /exchanges", func(t *testing.T) {
		request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/exchanges", nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		require.Equal(t, http.StatusOK, response.Code, response.Body.String())
		const expectedBody = `{
				"exchanges": [
					{
						"created_at": "2025-08-14T21:53:51.988009Z",
						"date": "2025-08-14",
						"description": "Grocery",
						"id": "4c703f3b-7505-4785-9c63-14b38b9e1129",
						"register_id": "7625b655-732d-49e0-a86b-43994ea89359",
						"splits": [
							{
								"amount": 10000,
								"created_at": "0001-01-01T00:00:00Z",
								"destination_register_id": "",
								"exchange_id": "",
								"id": "",
								"memo": "",
								"status": "",
								"updated_at": "0001-01-01T00:00:00Z"
							},
							{
								"amount": 1000,
								"created_at": "0001-01-01T00:00:00Z",
								"destination_register_id": "",
								"exchange_id": "",
								"id": "",
								"memo": "",
								"status": "",
								"updated_at": "0001-01-01T00:00:00Z"
							},
							{
								"amount": 5012,
								"created_at": "0001-01-01T00:00:00Z",
								"destination_register_id": "",
								"exchange_id": "",
								"id": "",
								"memo": "",
								"status": "",
								"updated_at": "0001-01-01T00:00:00Z"
							}
						],
						"status": "uncleared",
						"updated_at": "0001-01-01T00:00:00Z"
					},
					{
						"created_at": "2025-08-14T21:50:45.43635Z",
						"date": "2025-08-14",
						"description": "Transfer to credit card",
						"id": "e6469f9d-7388-4ea0-b6f3-e1b085ae1f7f",
						"register_id": "7625b655-732d-49e0-a86b-43994ea89359",
						"splits": [
							{
								"amount": 12356,
								"created_at": "0001-01-01T00:00:00Z",
								"destination_register_id": "",
								"exchange_id": "",
								"id": "",
								"memo": "",
								"status": "",
								"updated_at": "0001-01-01T00:00:00Z"
							}
						],
						"status": "uncleared",
						"updated_at": "0001-01-01T00:00:00Z"
					},
					{
						"created_at": "2025-08-14T21:53:51.988009Z",
						"date": "2025-08-21",
						"description": "Grocery",
						"id": "862207f1-1d28-40b5-903f-d3dbb312e4b9",
						"register_id": "7625b655-732d-49e0-a86b-43994ea89359",
						"splits": [
							{
								"amount": 3000,
								"created_at": "0001-01-01T00:00:00Z",
								"destination_register_id": "",
								"exchange_id": "",
								"id": "",
								"memo": "",
								"status": "",
								"updated_at": "0001-01-01T00:00:00Z"
							},
							{
								"amount": 1000,
								"created_at": "0001-01-01T00:00:00Z",
								"destination_register_id": "",
								"exchange_id": "",
								"id": "",
								"memo": "",
								"status": "",
								"updated_at": "0001-01-01T00:00:00Z"
							},
							{
								"amount": 2000,
								"created_at": "0001-01-01T00:00:00Z",
								"destination_register_id": "",
								"exchange_id": "",
								"id": "",
								"memo": "",
								"status": "",
								"updated_at": "0001-01-01T00:00:00Z"
							}
						],
						"status": "uncleared",
						"updated_at": "0001-01-01T00:00:00Z"
					}
				]
			}`
		require.JSONEq(t, expectedBody, response.Body.String(), response.Body.String())
	})
}
