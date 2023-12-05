// Unless explicitly stated otherwise all files in this repository are licensed under the Apache-2.0 License.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-Present Datadog, Inc.

package integration_test

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Call flag.Parse explicitly to prevent testing.Short() from panicking.
	flag.Parse()

	if testing.Short() {
		os.Exit(0)
	}

	os.Exit(m.Run())
}
