package repository

import "database/sql"

func strPtr(source sql.NullString) *string {
	if source.Valid {
		return &source.String
	}

	return nil
}

func int64Ptr(source sql.NullInt64) *int64 {
	if source.Valid {
		return &source.Int64
	}

	return nil
}
