package api

import (
	"context"
)

type Implementation struct{}

var _ StrictServerInterface = (*Implementation)(nil)

// (GET /health)
func (s *Implementation) GetHealth(context.Context, GetHealthRequestObject) (GetHealthResponseObject, error) {
	return GetHealth200Response{}, nil
}
