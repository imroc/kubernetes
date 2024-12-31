/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"k8s.io/kube-openapi/pkg/handler"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

func TestGen(t *testing.T) {
	dummyRef := func(name string) spec.Ref {
		name = strings.ReplaceAll(name, "k8s.io/api/kubelet", "k8s.io/kubelet/config")
		name = strings.ReplaceAll(name, "/", ".")
		ss := strings.Split(name, ".")
		slices.Reverse(ss)
		kind := ss[0]
		version := ss[1]
		group := strings.Join(ss[2:], ".")
		return spec.MustCreateRef(fmt.Sprintf(`../%s/%s_%s.json`, group, strings.ToLower(kind), version))
	}
	for name, value := range GetOpenAPIDefinitions(dummyRef) {
		if name != "k8s.io/api/kubelet/v1.CredentialProvider" {
			continue
		}
		t.Run(name, func(t *testing.T) {
			// TODO(kubernetes/gengo#193): We currently round-trip ints to floats.
			value.Schema = *handler.PruneDefaultsSchema(&value.Schema)
			data, err := json.MarshalIndent(value.Schema, "", "  ")
			if err != nil {
				t.Error(err)
				return
			}
			if err = os.WriteFile("schema.json", data, 0644); err != nil {
				t.Error(err)
			}
		})
	}
}
