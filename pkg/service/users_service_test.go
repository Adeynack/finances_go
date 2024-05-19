package service

import (
	"testing"

	"github.com/adeynack/finances/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestUsersServicePasswordEncodingAndDecoding(t *testing.T) {
	t.Parallel()
	userService, err := model.NewUsersService()
	require.NoError(t, err)

	t.Run("digest same value will result in different results", func(t *testing.T) {
		first, err := userService.EncodePasswordDigest("testing testing one two")
		require.NoError(t, err)
		second, err := userService.EncodePasswordDigest("testing testing one two")
		require.NoError(t, err)
		require.NotEqual(t, first, second)
	})

	t.Run("Matches digest value", func(t *testing.T) {
		const pwd = "a oh-so-very-secure password"
		digest, err := userService.EncodePasswordDigest(pwd)
		require.NoError(t, err)
		result, err := userService.MatchPassword(digest, pwd)
		require.NoError(t, err)
		require.True(t, result)
	})
}
