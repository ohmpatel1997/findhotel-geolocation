package repository_test

import (
	"errors"
	"testing"

	"github.com/ohmpatel1997/findhotel-geolocation/integration/repository"
	"github.com/stretchr/testify/require"
)

func TestPGConfigString(t *testing.T) {
	assert := require.New(t)

	testCases := []struct {
		pgc            repository.PGConfig
		expectedString string
		expectedError  error
	}{
		{repository.PGConfig{}, "", errors.New("Host must be set")},
		{repository.PGConfig{Host: "foo"}, "", errors.New("Port must be set")},
		{repository.PGConfig{Host: "foo", Port: 9}, "", errors.New("User must be set")},
		{repository.PGConfig{Host: "foo", Port: 9, User: "geddy"}, "", errors.New("Password must be set")},
		{repository.PGConfig{Host: "foo", Port: 9, User: "geddy", Password: "gg"}, "", errors.New("DBName must be set")},
		{repository.PGConfig{Host: "foo", Port: 9, User: "geddy", Password: "gg", DBName: "lee"}, "", errors.New("SSLMode must be set")},
		{repository.PGConfig{Host: "foo", Port: 9, User: "geddy", Password: "gg", DBName: "lee", SSLMode: "peart"}, "host=foo port=9 user=geddy password=gg dbname=lee sslmode=peart", nil},
	}

	for _, tc := range testCases {
		s, err := tc.pgc.String()

		assert.Equal(tc.expectedString, s)
		assert.Equal(tc.expectedError, err)
	}
}
