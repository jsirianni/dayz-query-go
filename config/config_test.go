package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name     string
		env      map[string]string
		expected *Config
		err      error
	}{
		{
			"valid configuration",
			map[string]string{
				EnvServerList: "10.99.1.10:5000,10.99.1.11:5001",
			},
			&Config{
				ServerList: []ServerEndpoint{
					{
						host: "10.99.1.10",
						port: "5000",
					},
					{
						host: "10.99.1.11",
						port: "5001",
					},
				},
			},
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables for this iteration
			// Cleanup with defer.
			for k, v := range tc.env {
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("failed to set environment variable %q: %v", k, err)
				}
				defer os.Unsetenv(k)
			}

			c, err := New()
			if tc.err != nil {
				require.Error(t, err, "an error was expected")
				require.ErrorAs(t, err, &tc.err)
				return
			}
			require.NoError(t, err, "an error was not expected")

			require.Equal(t, tc.expected, c, "config mismatch")
		})
	}
}
