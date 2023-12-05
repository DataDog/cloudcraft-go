// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

// Package xtesting provides functions and utilities for testing the Cloudcraft
// SDK.
package xtesting

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/DataDog/cloudcraft-go/internal/xerrors"
)

// ErrGreaterThanZero is returned when a given value is not greater than zero.
const ErrGreaterThanZero xerrors.Error = "value must be greater than zero"

// RandomString returns a random string of length n that is safe for use in a
// URL.
func RandomString(t *testing.T, r io.Reader, n int) (string, error) {
	t.Helper()

	if n <= 0 {
		return "", fmt.Errorf("%w: %d", ErrGreaterThanZero, n)
	}

	b := make([]byte, n) //nolint:makezero // no need for this specific use case

	_, err := io.ReadFull(r, b)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	encoded := base64.URLEncoding.EncodeToString(b)

	return encoded[:n], nil
}

// UniqueName returns an unique name that can be used as a resource name in
// Cloudcraft.
func UniqueName(t *testing.T) string {
	t.Helper()

	suffix, err := RandomString(t, rand.Reader, 16)
	if err != nil {
		t.Fatalf("failed to generate random string for unique name: %v", err)
	}

	return fmt.Sprintf("Go SDK Test (%s)", suffix)
}

// ReadFile reads the named file and returns its contents.
func ReadFile(t *testing.T, name string) []byte {
	t.Helper()

	file, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("failed to read test data file %q: %v", name, err)
	}

	return file
}

// GetEnv returns the value of the environment variable with the given name or
// fails the test if the variable is not set.
func GetEnv(t *testing.T, name string) string {
	t.Helper()

	value, found := os.LookupEnv(name)
	if !found {
		t.Fatalf("environment variable %q is not set; please set it before running the tests", name)
	}

	return value
}

// ParseTime returns the time parsed from the given string or fails the test if
// the string is not a valid time.
func ParseTime(t *testing.T, str string) time.Time {
	t.Helper()

	parsedTime, err := time.Parse(time.RFC3339, str)
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}

	return parsedTime
}
