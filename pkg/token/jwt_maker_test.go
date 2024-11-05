package token

import (
	"testing"
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/random"
	"github.com/stretchr/testify/require"
)

func TestJwtMaker(t *testing.T) {
	maker, err := NewJWTMaker(random.RandomString(32))
	require.NoError(t, err)

	username := random.RandomString(16)
	role := random.RandomString(4)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}
