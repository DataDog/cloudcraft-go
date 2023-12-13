// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package xhttp

import (
	"fmt"
	"io"
	"net/http"

	"github.com/DataDog/cloudcraft-go/internal/xerrors"
)

const (
	// ErrCannotDrainResponse is returned when a response body cannot be drained.
	ErrCannotDrainResponse xerrors.Error = "cannot drain response body"

	// ErrCannotCloseResponse is returned when a response body cannot be closed.
	ErrCannotCloseResponse xerrors.Error = "cannot close response body"
)

// DrainResponseBody reads and discards the remaining content of the response
// body until EOF, then closes it. If an error occurs while draining or closing
// the response body, an error is returned.
func DrainResponseBody(resp *http.Response) error {
	_, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotDrainResponse, err)
	}

	if err = resp.Body.Close(); err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCloseResponse, err)
	}

	return nil
}
