// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package xhttp

import (
	"crypto/tls"
	"net/http"
	"time"
)

const (
	// DefaultMaxIddleConns is the default maximum number of idle connections in
	// the pool.
	DefaultMaxIddleConns int = 100

	// DefaultMaxIddleConnsPerHost is the default maximum number of idle connections in
	// the pool per host.
	DefaultMaxIddleConnsPerHost int = 10

	// DefaultLRUClientSessionCacheCapacity is the default capacity of the LRU client session cache.
	DefaultLRUClientSessionCacheCapacity int = 64
)

// NewClient creates a new HTTP client with sane defaults given the provided
// timeout.
func NewClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				MinVersion:             tls.VersionTLS13,
				SessionTicketsDisabled: false,
				ClientSessionCache:     tls.NewLRUClientSessionCache(DefaultLRUClientSessionCacheCapacity),
			},
			MaxIdleConns:        DefaultMaxIddleConns,
			MaxIdleConnsPerHost: DefaultMaxIddleConnsPerHost,
			DisableCompression:  true,
			ForceAttemptHTTP2:   true,
		},
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
