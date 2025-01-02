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
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"k8s.io/kube-openapi/pkg/validation/spec"
)

const OUTPUT_DIR = "~/dotfiles/kubeschemas"

func parseNameRef(name string) (group, kind, version string) {
	for prefix, group := range APIs {
		if group == "" {
			continue
		}
		if strings.HasPrefix(name, prefix) {
			name = strings.Replace(name, prefix, group, 1)
		}
	}
	ss := strings.Split(name, "/")
	vk := ss[len(ss)-1]
	ss = ss[0 : len(ss)-1]
	slices.Reverse(ss)
	group = strings.Join(ss, ".")
	ss = strings.Split(vk, ".")
	version = ss[0]
	kind = ss[1]
	return
}

var APIs = map[string]string{
	"k8s.io/kubelet/config":                                 "kubelet.config.k8s.io",
	"k8s.io/client-go/tools/clientcmd/api":                  "core.api.k8s.io",
	"k8s.io/kube-controller-manager/config":                 "kubecontrollermanager.config.k8s.io",
	"k8s.io/component-base/tracing/api":                     "",
	"k8s.io/component-base/logs/api":                        "",
	"k8s.io/component-base/config":                          "",
	"k8s.io/apimachinery/pkg/api/resource.QuantityValue":    "api.pkg.apimachinery.k8s.io/resource.Quantity",
	"k8s.io/cloud-provider/controllers/service/config":      "",
	"k8s.io/cloud-provider/config":                          "cloudcontrollermanager.config.k8s.io",
	"k8s.io/controller-manager/config":                      "controllermanager.config.k8s.io",
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm":        "kubeadm.k8s.io",
	"k8s.io/kubernetes/cmd/kubeadm/app/apis/bootstraptoken": "",
	"k8s.io/api/admission/":                                 "admission.k8s.io/",
	"k8s.io/apimachinery/pkg/apis/meta":                     "",
	"k8s.io/apiserver/pkg/apis/audit":                       "audit.k8s.io",
	"k8s.io/apiserver/pkg/apis/apiserver":                   "apiserver.config.k8s.io",
	// "k8s.io/apiserver/pkg/apis/apiserver": "apiserver.k8s.io",
	"k8s.io/api/node_config":                     "config.node.controllers.cloud-provider.k8s.io",
	"k8s.io/kube-proxy/config":                   "kubeproxy.config.k8s.io",
	"k8s.io/kube-scheduler/config":               "kubescheduler.config.k8s.io",
	"k8s.io/kubelet/pkg/apis/credentialprovider": "credentialprovider.kubelet.k8s.io",
	"k8s.io/apiserver/pkg/admission/plugin/webhook/config/apis/webhookadmission": "apiserver.config.k8s.io",
	"k8s.io/client-go/pkg/apis/clientauthentication":                             "client.authentication.k8s.io",
	"k8s.io/kubernetes/plugin/pkg/admission/eventratelimit/apis/eventratelimit":  "eventratelimit.admission.k8s.io",
	"k8s.io/api/imagepolicy": "imagepolicy.k8s.io",
}

func TestWriteSchema(t *testing.T) {
	outputDir := OUTPUT_DIR
	if strings.Contains(outputDir, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Error(err)
			return
		}
		outputDir = strings.Replace(outputDir, "~", homeDir, 1)
	}
	outputDir, err := filepath.EvalSymlinks(outputDir)
	if err != nil {
		t.Error(err)
		return
	}

	ref := func(name string) spec.Ref {
		group, kind, version := parseNameRef(name)
		return spec.MustCreateRef(fmt.Sprintf(`../%s/%s_%s.json`, group, strings.ToLower(kind), version))
	}

	for name, value := range GetOpenAPIDefinitions(ref) {
		fmt.Println(name)
		allowed := false
		for prefix := range APIs {
			if strings.HasPrefix(name, prefix) {
				allowed = true
				break
			}
		}
		if !allowed {
			continue
		}
		group, kind, version := parseNameRef(name)
		dir := filepath.Join(outputDir, group)
		err := os.MkdirAll(dir, 0644)
		if err != nil {
			t.Error(err)
			return
		}
		data, err := value.Schema.MarshalJSON()
		if err != nil {
			t.Error(err)
			return
		}
		m := make(map[string]any)
		err = json.Unmarshal(data, &m)
		if err != nil {
			t.Error(err)
			return
		}
		props, ok := m["properties"].(map[string]any)
		if ok {
			propApiVersion, ok1 := props["apiVersion"].(map[string]any)
			propKind, ok2 := props["kind"].(map[string]any)
			if ok1 && ok2 {
				apiVersion := fmt.Sprintf("%s/%s", group, version)
				apiVersion = strings.Replace(apiVersion, "core.api.k8s.io/", "", 1)
				propApiVersion["enum"] = []string{apiVersion}
				propKind["enum"] = []string{kind}
				m["required"] = []string{"apiVersion", "kind"}
			}
		}
		data, err = json.MarshalIndent(&m, "", "  ")
		if err != nil {
			t.Error(err)
			return
		}
		filename := filepath.Join(dir, fmt.Sprintf("%s_%s.json", strings.ToLower(kind), version))
		if group == "meta.apis.pkg.apimachinery.k8s.io" { // not override existed meta api
			if _, err := os.Stat(filename); err != nil {
				if !os.IsNotExist(err) {
					t.Error(err)
					return
				}
			} else {
				continue
			}
		}
		fmt.Println("write", filename)
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			t.Error(err)
		}
	}
}
