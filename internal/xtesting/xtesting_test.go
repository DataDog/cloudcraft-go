// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package xtesting_test

import (
	"crypto/rand"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/DataDog/cloudcraft-go/internal/xtesting"
)

// ErrForced is a mock error used for testing.
var ErrForced = errors.New("forced error")

type errorRandReader struct{}

func (errorRandReader) Read(_ []byte) (_ int, _ error) {
	return 0, ErrForced
}

func TestRandomString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		give     int
		giveRand io.Reader
		want     int
		wantErr  error
	}{
		{
			name:     "PositiveLength",
			give:     10,
			giveRand: rand.Reader,
			want:     10,
		},
		{
			name:     "ZeroLength",
			give:     0,
			giveRand: rand.Reader,
			want:     0,
			wantErr:  xtesting.ErrGreaterThanZero,
		},
		{
			name:     "NegativeLength",
			give:     -5,
			giveRand: rand.Reader,
			want:     0,
			wantErr:  xtesting.ErrGreaterThanZero,
		},
		{
			name:     "RandomReadError",
			give:     10,
			giveRand: errorRandReader{},
			want:     0,
			wantErr:  ErrForced,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := xtesting.RandomString(t, tt.giveRand, tt.give)
			if (err != nil) != (tt.wantErr != nil) {
				t.Fatalf("RandomString() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("RandomString() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(output) != tt.want {
				t.Fatalf("RandomString() length = %d, want %d", len(output), tt.want)
			}
		})
	}
}

func TestUniqueName(t *testing.T) {
	t.Parallel()

	got := xtesting.UniqueName(t)

	if !strings.HasPrefix(got, "Go SDK Test (") || !strings.HasSuffix(got, ")") {
		t.Fatalf("UniqueName() output format is incorrect, got: %s", got)
	}

	// Checking if the length of the unique part (random string) is 16
	// characters. Since the prefix is "Go SDK Test (" and suffix is ")", the
	// unique part starts at 13th character and ends before the last character.
	if len(got) <= 14 || len(got)-14 != 16 {
		t.Fatalf("UniqueName() output does not have the expected length, got: %s", got)
	}
}
