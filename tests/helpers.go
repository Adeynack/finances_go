package tests

import (
	"net/http"
	"sync"
	"testing"

	"github.com/adeynack/finances/pkg/app"
	"github.com/adeynack/finances/pkg/repository"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

var internalTestDbConnection = sync.OnceValue(func() *bun.DB {
	return app.MustConnectDatabase()
})

// GetTestDB starts a transaction on the internal test dabase connection
// that will automatically be rollbacked after the test is performed.
func GetTestDB(t testing.TB) repository.DB {
	tx, err := internalTestDbConnection().BeginTx(t.Context(), nil)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = tx.Rollback()
	})

	return tx
}

func CreateTestAPIHandler(t testing.TB) http.Handler {
	return app.MustCreateHandler(GetTestDB(t))
}
