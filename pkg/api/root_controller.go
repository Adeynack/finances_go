package api

import (
	"context"
	"database/sql"
)

type Service struct {
	DB *sql.DB
}

var _ StrictServerInterface = (*Service)(nil)

// (GET /health)
func (s *Service) GetHealth(context.Context, GetHealthRequestObject) (GetHealthResponseObject, error) {
	return GetHealth200JSONResponse{Status: ServerHealthStatusHealthy}, nil
}
