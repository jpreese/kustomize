// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package graph

import (
	"fmt"
	"testing"

	"sigs.k8s.io/kustomize/pkg/kusttest"
)

func TestFirst(t *testing.T) {

	const kustomizationFile = `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - single.yaml
`

	const singleFile = `
apiVersion: apps/v1
metadata:
  name: dply1
kind: Deployment
`

	th := kusttest_test.NewKustTestHarness(t, "/app")

	th.WriteK("/app/", kustomizationFile)
	th.WriteF("/app/single.yaml", singleFile)

	th.WriteK("/app/another/", "")
	//th.WriteF("/app/another/moo.yaml", "")

	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Error creating resource map: %v", err)
	}

	for key, value := range m {
		fmt.Printf("key: %v\n", key)
		fmt.Printf("value: %v\n", value)
		fmt.Printf("refby: %v\n\n", value.GetRefBy())
	}
}
