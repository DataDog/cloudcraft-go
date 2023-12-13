// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package cloudcraft

import (
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestImageExportParams_query(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give ImageExportParams
		want url.Values
	}{
		{
			name: "Empty parameters",
			want: url.Values{},
		},
		{
			name: "All parameters set",
			give: ImageExportParams{
				PaperSize:   "A4",
				Grid:        true,
				Transparent: true,
				Landscape:   true,
				Scale:       1.5,
				Width:       1024,
				Height:      768,
			},
			want: url.Values{
				"paperSize":   []string{"A4"},
				"grid":        []string{"true"},
				"transparent": []string{"true"},
				"landscape":   []string{"true"},
				"scale":       []string{strconv.FormatFloat(1.5, 'f', -1, 32)},
				"width":       []string{"1024"},
				"height":      []string{"768"},
			},
		},
		{
			name: "Only paperSize and transparent",
			give: ImageExportParams{
				PaperSize:   "A3",
				Transparent: true,
			},
			want: url.Values{
				"paperSize":   []string{"A3"},
				"transparent": []string{"true"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.give.query(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ImageExportParams.query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBudgetExportParams_query(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give BudgetExportParams
		want url.Values
	}{
		{
			name: "Empty parameters",
			want: url.Values{},
		},
		{
			name: "All parameters set",
			give: BudgetExportParams{
				Currency: "USD",
				Period:   "monthly",
				Rate:     "standard",
			},
			want: url.Values{
				"currency": []string{"USD"},
				"period":   []string{"monthly"},
				"rate":     []string{"standard"},
			},
		},
		{
			name: "Only currency and period",
			give: BudgetExportParams{
				Currency: "EUR",
				Period:   "yearly",
			},
			want: url.Values{
				"currency": []string{"EUR"},
				"period":   []string{"yearly"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.give.query(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("BudgetExportParams.query() = %v, want %v", got, tt.want)
			}
		})
	}
}
