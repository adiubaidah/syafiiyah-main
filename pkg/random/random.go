package random

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt(min, max int64) int64 {
	return min + seededRand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[seededRand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomEmail() string {
	return RandomString(8) + "@" + RandomString(5) + ".com"
}

func RandomBool() bool {
	return seededRand.Intn(2) == 1
}

func RandomTimeStamp() time.Time {
	min := time.Date(2024, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2025, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := RandomInt(0, delta)
	return time.Unix(min+sec, 0)
}

func RandomTimeOnly() time.Time {
	hour := rand.Intn(24)
	minute := rand.Intn(60)
	second := rand.Intn(60)

	return time.Date(0, 1, 1, hour, minute, second, 0, time.UTC)
}
