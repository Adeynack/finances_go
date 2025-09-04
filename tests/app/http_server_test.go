package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/adeynack/finances/pkg/api/apimodel"
	"github.com/adeynack/finances/pkg/api/apiserver"
	"github.com/adeynack/finances/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.llib.dev/testcase"
)

func TestHttpServer(t *testing.T) {
	s := testcase.NewSpec(t)

	handler := testcase.Let(s, func(t *testcase.T) http.Handler {
		return tests.CreateTestAPIHandler(t.TB)
	})

	method := testcase.Var[string]{ID: "method"}
	path := testcase.Var[string]{ID: "path"}
	requestBody := testcase.LetValue[io.Reader](s, nil)

	response := testcase.Let(s, func(t *testcase.T) httptest.ResponseRecorder {
		request := httptest.NewRequestWithContext(t.Context(), method.Get(t), path.Get(t), requestBody.Get(t))
		innerResponse := httptest.NewRecorder()
		handler.Get(t).ServeHTTP(innerResponse, request)

		return *innerResponse
	})

	s.Describe("/health", func(s *testcase.Spec) {
		path.LetValue(s, "/health")

		s.When("GET", func(s *testcase.Spec) {
			method.LetValue(s, http.MethodGet)

			s.Then("responds with 200 OK and body", func(t *testcase.T) {
				r := response.Get(t)
				require.Equal(t, http.StatusOK, r.Code)
				require.JSONEq(t, `{"status": "healthy"}`, r.Body.String())
			})
		})
	})

	s.Describe("/books", func(s *testcase.Spec) {
		path.LetValue(s, "/books")

		s.When("POST", func(s *testcase.Spec) {
			method.LetValue(s, http.MethodPost)

			s.When("the owner ID does not exist", func(s *testcase.Spec) {
				requestBody.Let(s, func(t *testcase.T) io.Reader {
					return strings.NewReader(`{
							"book": {
								"name": "My all new shiny book",
								"owner_id": "095005a0-aa18-43f6-b578-66800aa6aae8",
								"default_currency_iso_code": "CAD"
							}
						}
					`)
				})

				s.Then("it fails with a 422 Unprocessable Entity", func(t *testcase.T) {
					r := response.Get(t)
					require.Equal(t, http.StatusUnprocessableEntity, r.Code)
					var jsonError apimodel.Error
					require.NoError(t, json.Unmarshal(r.Body.Bytes(), &jsonError))
					require.Equal(t, "validation error: owner does not exist", jsonError.Title)
				})
			})

			s.When("the request contains a valid book to create", func(s *testcase.Spec) {
				requestBody.Let(s, func(t *testcase.T) io.Reader {
					return strings.NewReader(`{
							"book": {
								"name": "My all new shiny book",
								"owner_id": "569bcfdd-4056-42cd-af9c-285fa5ce92c8",
								"default_currency_iso_code": "CAD"
							}
						}`)
				})

				s.Then("it creates a book and responds 201 Created with the new book as body", func(t *testcase.T) {
					startTime := time.Now()
					r := response.Get(t)

					require.Equal(t, http.StatusCreated, r.Code, r.Body)

					var body apiserver.CreateBook201JSONResponse
					require.NoError(t, json.Unmarshal(r.Body.Bytes(), &body), r.Body)
					require.NotZero(t, body.Book.Id, "expecting an ID to be set")
					require.GreaterOrEqual(t, body.Book.CreatedAt, startTime, "expecting CreatedAt to be set to book creation time")
					require.GreaterOrEqual(t, body.Book.UpdatedAt, startTime, "expecting CreatedAt to be set to book creation time")
					expectedBody := apiserver.CreateBook201JSONResponse{
						Book: apimodel.Book{
							Id:                     body.Book.Id,
							CreatedAt:              body.Book.CreatedAt,
							UpdatedAt:              body.Book.UpdatedAt,
							Name:                   "My all new shiny book",
							OwnerId:                uuid.MustParse("569bcfdd-4056-42cd-af9c-285fa5ce92c8"),
							OwnerDisplayName:       "Joe",
							DefaultCurrencyIsoCode: "CAD",
						},
					}
					require.Equal(t, expectedBody, body, r.Body)
				})
			})
		})

		s.When("GET", func(s *testcase.Spec) {
			method.LetValue(s, http.MethodGet)

			s.Then("it responds 200 OK with a list of books", func(t *testcase.T) {
				r := response.Get(t)

				require.Equal(t, http.StatusOK, r.Code, r.Body)
				const expectedBody = `{
							"books": [
								{
									"created_at": "2025-08-24T00:20:00Z",
									"default_currency_iso_code": "USD",
									"id": "5da6e20f-eecd-456b-a8dd-ae1a63d0268e",
									"name": "Foo, the Book",
									"owner_display_name": "Joe",
									"owner_id": "569bcfdd-4056-42cd-af9c-285fa5ce92c8",
									"updated_at": "2025-08-24T00:20:00Z"
								},
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
				require.JSONEq(t, expectedBody, r.Body.String(), r.Body)
			})
		})
	})

	s.Describe("/books/:id", func(s *testcase.Spec) {
		bookId := testcase.Let[string](s, nil)
		path.Let(s, func(t *testcase.T) string { return fmt.Sprintf("/books/%s", bookId.Get(t)) })

		s.When("GET", func(s *testcase.Spec) {
			method.LetValue(s, http.MethodGet)

			s.When("the book ID exists", func(s *testcase.Spec) {
				bookId.LetValue(s, "8d8666c0-016f-49fb-8f59-4150a822ffb2")

				s.Then("it responds with 200 OK with the book", func(t *testcase.T) {
					r := response.Get(t)
					require.Equal(t, http.StatusOK, r.Code, r.Body)
					const expectedBody = `{
						"book": {
							"created_at": "2025-08-14T21:47:58.211393Z",
							"default_currency_iso_code": "EUR",
							"id": "8d8666c0-016f-49fb-8f59-4150a822ffb2",
							"name": "Joe's Book",
							"owner_display_name": "Joe",
							"owner_id": "569bcfdd-4056-42cd-af9c-285fa5ce92c8",
							"updated_at": "2025-08-14T21:47:58.211393Z"
						}
					}`
					require.JSONEq(t, expectedBody, r.Body.String(), r.Body)
				})
			})
		})
	})

	s.Describe("/exchanges", func(s *testcase.Spec) {
		path.LetValue(s, "/exchanges")

		s.When("GET", func(s *testcase.Spec) {
			method.LetValue(s, http.MethodGet)

			s.Then("responds 200 OK with a list of exchanges", func(t *testcase.T) {
				r := response.Get(t)
				require.Equal(t, http.StatusOK, r.Code, r.Body)
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
												"counterpart_amount": -10000,
												"created_at": "2025-08-14T21:54:36.523784Z",
												"destination_register_id": "f01ac27a-89f4-4aab-b534-0cf77cee659c",
												"exchange_id": "4c703f3b-7505-4785-9c63-14b38b9e1129",
												"id": "61dad97b-f7a0-45eb-b325-e6fa5f601937",
												"status": "uncleared",
												"updated_at": "2025-08-14T21:54:36.523784Z"
											},
											{
												"amount": 1000,
												"counterpart_amount": -1000,
												"created_at": "2025-08-14T21:54:36.523784Z",
												"destination_register_id": "af903dfd-3e65-41d2-83f8-db1422b5159c",
												"exchange_id": "4c703f3b-7505-4785-9c63-14b38b9e1129",
												"id": "118dd813-e2c0-4eed-8dc4-773b87de9b93",
												"status": "uncleared",
												"updated_at": "2025-08-14T21:54:36.523784Z"
											},
											{
												"amount": 5012,
												"counterpart_amount": -5012,
												"created_at": "2025-08-14T21:54:36.523784Z",
												"destination_register_id": "df622be3-dff4-4541-b2c7-40947bc68d5f",
												"exchange_id": "4c703f3b-7505-4785-9c63-14b38b9e1129",
												"id": "88c2dfe2-5589-4695-bd6f-2ccd24bda109",
												"status": "uncleared",
												"updated_at": "2025-08-14T21:54:36.523784Z"
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
												"counterpart_amount": -12356,
												"created_at": "2025-08-14T21:51:40.606261Z",
												"destination_register_id": "f7baf4ac-52f8-494c-87b3-5178b416411f",
												"exchange_id": "e6469f9d-7388-4ea0-b6f3-e1b085ae1f7f",
												"id": "2fce88f8-2977-4cb6-89ae-af21ff108e35",
												"status": "uncleared",
												"updated_at": "2025-08-14T21:51:40.606261Z"
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
												"counterpart_amount": -3000,
												"created_at": "2025-08-14T21:54:36.523784Z",
												"destination_register_id": "af903dfd-3e65-41d2-83f8-db1422b5159c",
												"exchange_id": "862207f1-1d28-40b5-903f-d3dbb312e4b9",
												"id": "e4261caa-4924-491d-9b3c-2be680f83859",
												"status": "uncleared",
												"updated_at": "2025-08-14T21:54:36.523784Z"
											},
											{
												"amount": 1000,
												"counterpart_amount": -1000,
												"created_at": "2025-08-14T21:54:36.523784Z",
												"destination_register_id": "f01ac27a-89f4-4aab-b534-0cf77cee659c",
												"exchange_id": "862207f1-1d28-40b5-903f-d3dbb312e4b9",
												"id": "34f1be6c-9595-4aa4-a508-0ef1228a7536",
												"status": "uncleared",
												"updated_at": "2025-08-14T21:54:36.523784Z"
											},
											{
												"amount": 2000,
												"counterpart_amount": -2000,
												"created_at": "2025-08-14T21:54:36.523784Z",
												"destination_register_id": "df622be3-dff4-4541-b2c7-40947bc68d5f",
												"exchange_id": "862207f1-1d28-40b5-903f-d3dbb312e4b9",
												"id": "9b861a0a-3637-471e-b24c-d121f5a273ca",
												"status": "uncleared",
												"updated_at": "2025-08-14T21:54:36.523784Z"
											}
										],
										"status": "uncleared",
										"updated_at": "0001-01-01T00:00:00Z"
									}
								]
							}`
				require.JSONEq(t, expectedBody, r.Body.String(), r.Body)
			})
		})
	})
}
