// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package cloudcraft

import (
	"net/url"
	"time"

	"github.com/DataDog/cloudcraft-go/internal/xerrors"
	"github.com/DataDog/cloudcraft-go/internal/xos"
)

const (
	// ErrInvalidEndpoint is returned when the endpoint is not a valid URL.
	ErrInvalidEndpoint xerrors.Error = "invalid endpoint"

	// ErrMissingEndpointScheme is returned when a Config is created without a
	// scheme for the endpoint.
	ErrMissingEndpointScheme xerrors.Error = "missing endpoint scheme"

	// ErrMissingEndpointHost is returned when a Config is created without a
	// host for the endpoint.
	ErrMissingEndpointHost xerrors.Error = "missing endpoint host"

	// ErrMissingKey is returned when a Config is created without an API key.
	ErrMissingKey xerrors.Error = "missing API key"

	// ErrInvalidKey is returned when a Config is created with an invalid API
	// key.
	ErrInvalidKey xerrors.Error = "invalid API key; length must be 44"
)

const (
	// DefaultSceme is the default protocol scheme, such as "http" or "https".
	DefaultScheme string = "https"

	// DefaultHost is the default host name or IP address of the Cloudcraft API.
	DefaultHost string = "api.cloudcraft.co"

	// DefaultPort is the default port number of the Cloudcraft API.
	DefaultPort string = "443"

	// DefaultPath is the default path to the Cloudcraft API.
	DefaultPath string = "/"

	// DefaultTimeout is the default timeout for requests made by the Cloudcraft
	// API client.
	DefaultTimeout time.Duration = time.Second * 80
)

// Environment variables used to configure the Config struct.
const (
	EnvScheme  string = "CLOUDCRAFT_PROTOCOL"
	EnvHost    string = "CLOUDCRAFT_HOST"
	EnvPort    string = "CLOUDCRAFT_PORT"
	EnvPath    string = "CLOUDCRAFT_PATH"
	EnvTimeout string = "CLOUDCRAFT_TIMEOUT"
	EnvAPIKey  string = "CLOUDCRAFT_API_KEY" //nolint:gosec // false positive
)

// Config holds the basic configuration for the Cloudcraft API.
type Config struct {
	// endpoint specifies the base URL of the Cloudcraft API for HTTP requests.
	// It is constructed from the Scheme, Host, Port, and Path fields.
	endpoint *url.URL

	// Scheme is the protocol scheme, such as "http" or "https", to use when
	// calling the API.
	//
	// If not set, the value of the CLOUDCRAFT_PROTOCOL environment variable is
	// used. If the environment variable is not set, the default value is
	// "https".
	//
	// This field is optional.
	Scheme string

	// Host is the host name or IP address of the Cloudcraft API.
	//
	// If not set, the value of the CLOUDCRAFT_HOST environment variable is
	// used. If the environment variable is not set, the default value is the
	// public instance of Cloudcraft, "api.cloudcraft.co".
	//
	// This field is optional.
	Host string

	// Port is the port number of the Cloudcraft API.
	//
	// If not set, the value of the CLOUDCRAFT_PORT environment variable is
	// used. If the environment variable is not set, the default value is "443".
	//
	// This field is optional.
	Port string

	// Path is the path to the Cloudcraft API.
	//
	// If not set, the value of the CLOUDCRAFT_PATH environment variable is
	// used. If the environment variable is not set, the default value is "/".
	//
	// This field is optional.
	Path string

	// Key is the API key used to authenticate with the Cloudcraft API.
	//
	// This field is required. [Learn more].
	//
	// [Learn more]: https://developers.cloudcraft.co/#authentication
	Key string

	// Timeout is the time limit for requests made by the Cloudcraft API client
	// to the Cloudcraft API.
	//
	// If not set, the value of the CLOUDCRAFT_TIMEOUT environment variable is
	// used. If the environment variable is not set, the default value is 80
	// seconds.
	//
	// This field is optional.
	Timeout time.Duration
}

// NewConfig returns a new Config with the given API key.
func NewConfig(key string) *Config {
	return &Config{
		Scheme:  DefaultScheme,
		Host:    DefaultHost,
		Port:    DefaultPort,
		Path:    DefaultPath,
		Key:     key,
		Timeout: DefaultTimeout,
	}
}

// NewConfigFromEnv returns a new Config from values set in the environment.
func NewConfigFromEnv() *Config {
	return &Config{
		Scheme:  xos.GetEnv(EnvScheme, DefaultScheme),
		Host:    xos.GetEnv(EnvHost, DefaultHost),
		Port:    xos.GetEnv(EnvPort, DefaultPort),
		Path:    xos.GetEnv(EnvPath, DefaultPath),
		Key:     xos.GetEnv(EnvAPIKey, ""),
		Timeout: xos.GetDurationEnv(EnvTimeout, DefaultTimeout),
	}
}

// Validate checks that the Config is valid.
func (c *Config) Validate() error {
	if c.Scheme == "" {
		return ErrMissingEndpointScheme
	}

	if c.Host == "" {
		return ErrMissingEndpointHost
	}

	if c.Key == "" {
		return ErrMissingKey
	}

	if len(c.Key) != 44 {
		return ErrInvalidKey
	}

	return nil
}
