package token

import (
	"testing"
	"time"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	"github.com/adiubaidah/syafiiyah-main/pkg/random"
	"github.com/stretchr/testify/require"
)

func TestJwtMaker(t *testing.T) {
	maker, err := NewJWTMaker(random.RandomString(32))
	require.NoError(t, err)

	username := random.RandomString(16)
	roles := []repo.RoleType{repo.RoleTypeAdmin, repo.RoleTypeEmployee}
	role := roles[random.RandomInt(0, int64(len(roles))-1)]
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(&model.User{
		Username: username,
		ID:       int32(random.RandomInt(0, 1000)),
		Role:     role,
	}, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.User.Username)
	require.Equal(t, role, payload.User.Role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}
