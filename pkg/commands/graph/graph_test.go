/// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package graph

import (
	"fmt"
	"testing"

	"sigs.k8s.io/kustomize/pkg/kusttest"
)

func TestFirst(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app")
	th.WriteK("/app/", "")

	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Error creating resource map: %v", err)
	}
}
