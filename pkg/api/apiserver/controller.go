package apiserver

import (
	"context"

	"github.com/adeynack/finances/pkg/api/apimodel"
	"github.com/adeynack/finances/pkg/repository"
)

type Service struct {
	Repo repository.R
}

var _ StrictServerInterface = (*Service)(nil)

// (GET /health)
func (s *Service) GetHealth(context.Context, GetHealthRequestObject) (GetHealthResponseObject, error) {
	return GetHealth200JSONResponse{Status: apimodel.ServerHealthStatusHealthy}, nil
}

func (s *Service) GetBooks(ctx context.Context, request GetBooksRequestObject) (GetBooksResponseObject, error) {
	booksForAPIResponse, err := s.Repo.GetBooks(ctx)
	if err != nil {
		return nil, err
	}

	return GetBooks200JSONResponse{Books: booksForAPIResponse}, nil
}

// GetExchanges implements StrictServerInterface.
func (s *Service) GetExchanges(ctx context.Context, request GetExchangesRequestObject) (GetExchangesResponseObject, error) {
	exchangesForAPIResponse, err := s.Repo.GetExchangesWithSplits(ctx)
	if err != nil {
		return nil, err
	}

	return GetExchanges200JSONResponse{Exchanges: exchangesForAPIResponse}, nil
}
