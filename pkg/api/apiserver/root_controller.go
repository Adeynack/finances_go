package apiserver

import (
	"context"
	"database/sql"

	"github.com/adeynack/finances/pkg/api/apimodel"
)

type Service struct {
	DB *sql.DB
}

var _ StrictServerInterface = (*Service)(nil)

// (GET /health)
func (s *Service) GetHealth(context.Context, GetHealthRequestObject) (GetHealthResponseObject, error) {
	return GetHealth200JSONResponse{Status: apimodel.ServerHealthStatusHealthy}, nil
}
