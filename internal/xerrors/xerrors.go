// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

// Package xerrors provides helper functions and types for error handling.
package xerrors

// Error is an [imuutable error] type.
//
// [imuutable error]: https://dave.cheney.net/2016/04/07/constant-errors
type Error string

// Error implements the error interface for Error.
func (e Error) Error() string {
	return string(e)
}
