package tests

import (
	"net/http"
	"sync"
	"testing"

	"github.com/adeynack/finances/pkg/app"
	"github.com/adeynack/finances/pkg/repository"
	"github.com/stretchr/testify/require"
)

var internalTestDbConnection = sync.OnceValue(func() repository.DB {
	return repository.NewDB(app.MustConnectDatabase())
})

// GetTestDB starts a transaction on the internal test dabase connection
// that will automatically be rollbacked after the test is performed.
func GetTestDB(t testing.TB) repository.DB {
	db, txc, err := internalTestDbConnection().BeginTx(t.Context())
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = txc.Rollback()
	})

	return db
}

func CreateTestAPIHandler(t testing.TB) http.Handler {
	return app.MustCreateHandler(GetTestDB(t))
}
