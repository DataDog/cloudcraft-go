// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

// Package endpoint provides a function to parse fragments of an URL into an
// *url.URL.
package endpoint

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/DataDog/cloudcraft-go/internal/xerrors"
)

const (
	_Slash           string = "/"
	_ColonSlashSlash string = "://"
	_Colon           string = ":"
)

const (
	ErrMissingFragment xerrors.Error = "missing scheme or host"
	ErrInvalidScheme   xerrors.Error = "invalid URL scheme"
)

// Parse parses fragments of an URL into an *url.URL.
func Parse(scheme, host, port, path string) (*url.URL, error) {
	if scheme == "" || host == "" {
		return nil, ErrMissingFragment
	}

	if scheme != "https" && scheme != "http" {
		return nil, fmt.Errorf("%w", ErrInvalidScheme)
	}

	if path == "" {
		path = _Slash
	}

	var builder strings.Builder

	builder.Grow(len(scheme) + len(host) + len(port) + len(path) + 4)

	builder.WriteString(scheme)
	builder.WriteString(_ColonSlashSlash)
	builder.WriteString(host)

	if port != "" {
		builder.WriteString(_Colon)
		builder.WriteString(port)
	}

	builder.WriteString(path)

	uri, err := url.Parse(builder.String())
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return uri, nil
}
