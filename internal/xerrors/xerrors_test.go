// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package xerrors_test

import (
	"testing"

	"github.com/DataDog/cloudcraft-go/internal/xerrors"
)

func TestError_ErrorMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		inputError    xerrors.Error
		expectedError string
	}{
		{
			name:          "Empty error",
			inputError:    "",
			expectedError: "",
		},
		{
			name:          "Simple error message",
			inputError:    "something went wrong",
			expectedError: "something went wrong",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualError := tt.inputError.Error()

			if actualError != tt.expectedError {
				t.Fatalf("Expected error: %q, got: %q", tt.expectedError, actualError)
			}
		})
	}
}
