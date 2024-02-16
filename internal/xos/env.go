// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package xos

import (
	"os"
	"strconv"
	"time"
)

// GetEnv returns the string value of the environment variable named by the key.
//
// If the variable is present in the environment the value (which may be empty)
// or if the variable is unset, a fallback value is returned.
func GetEnv(key, fallback string) string {
	value, found := os.LookupEnv(key)
	if !found || value == "" {
		return fallback
	}

	return value
}

// GetIntEnv returns the integer value of the environment variable named by the
// key.
//
// If the variable is present in the environment the value (which may be empty)
// or if the variable is unset, a fallback value is returned.
func GetIntEnv(key string, fallback int) int {
	value, found := os.LookupEnv(key)
	if !found || value == "" {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return intValue
}

// GetDurationEnv returns the time.Duration value of the environment variable
// named by the key.
//
// If the variable is present in the environment the value (which may be empty)
// or if the variable is unset, a fallback value is returned.
func GetDurationEnv(key string, fallback time.Duration) time.Duration {
	value, found := os.LookupEnv(key)
	if !found || value == "" {
		return fallback
	}

	durationValue, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return durationValue
}
