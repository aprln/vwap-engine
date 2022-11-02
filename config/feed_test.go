package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFeed(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    Feed
	}{
		{
			name: "no env vars",
			want: Feed{
				Name:            FeedNameCoinbase,
				WSConnectionURL: deftFeedWSConnectionURL,
			},
		},
		{
			name: "with env vars",
			envVars: map[string]string{
				"FEED_NAME":              "banana",
				"FEED_WS_CONNECTION_URL": "monkey",
			},
			want: Feed{
				Name:            "banana",
				WSConnectionURL: "monkey",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				for k, v := range tt.envVars {
					t.Setenv(k, v)
				}
				got := NewFeed()
				assert.Equal(t, tt.want, got)
			},
		)
	}
}
