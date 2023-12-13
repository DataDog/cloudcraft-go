// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package cloudcraft

import "github.com/DataDog/cloudcraft-go/internal/xerrors"

const (
	// ErrNilContext is returned when a nil context is passed to a function.
	ErrNilContext xerrors.Error = "context cannot be nil"

	// ErrNilAccount is returned when a nil account is passed as an argument.
	ErrNilAccount xerrors.Error = "account cannot be nil"

	// ErrAccountsKey is returned when the response from the API to a List call
	// is not a list of AWS or Azure accounts.
	ErrAccountsKey xerrors.Error = "key 'accounts' not found in response"

	// ErrEmptyAccountName is returned when an empty account name is passed as
	// an argument.
	ErrEmptyAccountName xerrors.Error = "account name cannot be empty"

	// ErrMissingAccountID is returned when an empty account ID is passed as an
	// argument.
	ErrEmptyAccountID xerrors.Error = "account ID cannot be empty"

	// ErrEmptyRegion is returned when an empty region is passed as an argument.
	ErrEmptyRegion xerrors.Error = "region cannot be empty"
)
