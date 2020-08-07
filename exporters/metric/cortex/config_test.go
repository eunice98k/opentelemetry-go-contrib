package cortex_test

import (
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"opentelemetry.io/contrib/exporters/metric/cortex"
)

// This is an example Config struct with mostly default values. The endpoint is not default since
// there is no default endpoint.
var ExampleStandardConfig = cortex.Config{
	Endpoint:      "/api/prom/push",
	Name:          "Standard Config",
	RemoteTimeout: 30 * time.Second,
	PushInterval:  10 * time.Second,
	Client:        http.DefaultClient,
}

var ExampleNoEndpointConfig = cortex.Config{
	Name:          "Standard Config",
	RemoteTimeout: 30 * time.Second,
	PushInterval:  10 * time.Second,
	Client:        http.DefaultClient,
}

// This is an example Config struct with mostly default values. The endpoint is not default since
// there is no default endpoint.
var ExampleTwoBearerTokenConfig = cortex.Config{
	Endpoint:        "/api/prom/push",
	Name:            "Standard Config",
	RemoteTimeout:   30 * time.Second,
	PushInterval:    10 * time.Second,
	BearerToken:     "bearer_token",
	BearerTokenFile: "bearer_token_file",
	Client:          http.DefaultClient,
}

var ExampleTwoPasswordConfig = cortex.Config{
	Endpoint:      "/api/prom/push",
	Name:          "Standard Config",
	RemoteTimeout: 30 * time.Second,
	PushInterval:  10 * time.Second,
	BasicAuth: map[string]string{
		"username":      "user",
		"password":      "password",
		"password_file": "passwordFile",
	},
	Client: http.DefaultClient,
}

// This is an example Config struct with a proxy url. url.Parse returns an error, so the variable it
// is called outside of the struct.
var proxyURL = "/proxy/url"
var parsedProxyURL, err = url.Parse(proxyURL)
var ExampleConfigWithProxy = cortex.Config{
	Endpoint:      "/api/prom/push",
	RemoteTimeout: 30 * time.Second,
	Name:          "Config with proxy",
	ProxyURL:      "/proxy/url",
	PushInterval:  10 * time.Second,
	Client: &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(parsedProxyURL),
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	},
}

// TestValidate checks whether Validate() returns the correct error on different Config structs.
func TestValidate(t *testing.T) {
	tests := []struct {
		testName      string
		config        cortex.Config
		expectedError error
	}{
		{
			"Standard Config",
			ExampleStandardConfig,
			nil,
		},
		{
			"Config with Conflicting Bearer Tokens",
			ExampleTwoBearerTokenConfig,
			cortex.ErrTwoBearerTokens,
		},
		{
			"Config with Conflicting Passwords",
			ExampleTwoPasswordConfig,
			cortex.ErrTwoPasswords,
		},
		{
			"Config with Proxy URL",
			ExampleConfigWithProxy,
			nil,
		},
		{
			"Config with no Endpoint",
			ExampleNoEndpointConfig,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := test.config.Validate()
			require.Equal(t, err, test.expectedError)
		})
	}
}
