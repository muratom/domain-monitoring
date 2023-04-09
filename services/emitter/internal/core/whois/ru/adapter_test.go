package ru

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_parseResponse(t *testing.T) {
	file, err := os.Open("./response_sample")
	require.NoError(t, err)
	defer file.Close()

	raw, err := io.ReadAll(file)
	require.NoError(t, err)

	got, err := parseResponse(context.Background(), raw)
	require.NoError(t, err)

	require.Equal(t, "YANDEX.RU", got.domain)
	require.Equal(t, "RU-CENTER-RU", got.registrar)

	require.ElementsMatch(t, []domainState{Registered, Delegated, Verified}, got.state)

	expectedCreated, err := time.Parse(time.RFC3339, "1997-09-23T09:45:07Z")
	require.NoError(t, err)
	require.True(t, got.created.Equal(expectedCreated))

	expectedPaidTill, err := time.Parse(time.RFC3339, "2023-09-30T21:00:00Z")
	require.NoError(t, err)
	require.True(t, got.paidTill.Equal(expectedPaidTill))
}
