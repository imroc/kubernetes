package main

import (
	"flag"
	"log"
	"strings"

	"k8s.io/kubernetes/genschema/openapi"
)

// APIEntry defines a type path prefix and its corresponding group replacement.
// Entries are matched in order; the first matching prefix wins.
// Longer prefixes should come before shorter ones to ensure correct matching.
type APIEntry struct {
	Prefix string // type path prefix to match
	Group  string // group name replacement (empty means derive from path)
}

// apis lists all API path prefixes that genschema should generate schemas for.
// Order matters: more specific (longer) prefixes must come before less specific ones.
var apis = []APIEntry{
	// Kubernetes Configuration APIs
	{Prefix: "k8s.io/apimachinery/pkg/api/resource.QuantityValue", Group: "api.pkg.apimachinery.k8s.io/resource.Quantity"},
	{Prefix: "k8s.io/apimachinery/pkg/apis/meta", Group: "meta.apis.pkg.apimachinery.k8s.io"},
	{Prefix: "k8s.io/kubernetes/cmd/kubeadm/app/apis/bootstraptoken", Group: "kubeadm.k8s.io"},
	{Prefix: "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm", Group: "kubeadm.k8s.io"},
	{Prefix: "k8s.io/kubernetes/pkg/apis/core", Group: ""},
	{Prefix: "k8s.io/kubernetes/plugin/pkg/admission/eventratelimit/apis/eventratelimit", Group: "eventratelimit.admission.k8s.io"},
	{Prefix: "k8s.io/api/admission/", Group: "admission.k8s.io/"},
	{Prefix: "k8s.io/api/imagepolicy", Group: "imagepolicy.k8s.io"},
	{Prefix: "k8s.io/api/node_config", Group: "config.node.controllers.cloud-provider.k8s.io"},
	{Prefix: "k8s.io/apiserver/pkg/admission/plugin/webhook/config/apis/webhookadmission", Group: "apiserver.config.k8s.io"},
	{Prefix: "k8s.io/apiserver/pkg/apis/apiserver", Group: "apiserver.config.k8s.io"},
	{Prefix: "k8s.io/apiserver/pkg/apis/audit", Group: "audit.k8s.io"},
	{Prefix: "k8s.io/client-go/pkg/apis/clientauthentication", Group: "client.authentication.k8s.io"},
	{Prefix: "k8s.io/client-go/tools/clientcmd/api", Group: "core.api.k8s.io"},
	{Prefix: "k8s.io/cloud-provider/config", Group: "cloudcontrollermanager.config.k8s.io"},
	{Prefix: "k8s.io/cloud-provider/controllers/service/config", Group: "service.config.controller-manager.k8s.io"},
	{Prefix: "k8s.io/cloud-provider/controllers/node/config", Group: "config.node.controllers.cloud-provider.k8s.io"},
	{Prefix: "k8s.io/component-base/config", Group: "config.component-base.k8s.io"},
	{Prefix: "k8s.io/component-base/logs/api", Group: "api.logs.component-base.k8s.io"},
	{Prefix: "k8s.io/component-base/tracing/api", Group: "api.tracing.component-base.k8s.io"},
	{Prefix: "k8s.io/controller-manager/config", Group: "controllermanager.config.k8s.io"},
	{Prefix: "k8s.io/kube-controller-manager/config", Group: "kubecontrollermanager.config.k8s.io"},
	{Prefix: "k8s.io/kube-proxy/config", Group: "kubeproxy.config.k8s.io"},
	{Prefix: "k8s.io/kube-scheduler/config", Group: "kubescheduler.config.k8s.io"},
	{Prefix: "k8s.io/kubelet/config", Group: "kubelet.config.k8s.io"},
	{Prefix: "k8s.io/kubelet/pkg/apis/credentialprovider", Group: "credentialprovider.kubelet.k8s.io"},

	// External APIs (synced from upstream repos)
	{Prefix: "agones.dev/agones/pkg/apis/allocation", Group: "allocation.agones.dev"},
	{Prefix: "agones.dev/agones/pkg/apis/agones", Group: "agones.apis.pkg.agones.agones.dev"},
	{Prefix: "sigs.k8s.io/gateway-api/apis", Group: "gateway-api.k8s.io"},
}

// findAPIEntry returns the APIEntry whose Prefix matches name, or nil if none match.
// Iterates in slice order so longer/more-specific prefixes take precedence.
func findAPIEntry(name string) *APIEntry {
	for i := range apis {
		if strings.HasPrefix(name, apis[i].Prefix) {
			return &apis[i]
		}
	}
	return nil
}

func main() {
	var outputDir string
	flag.StringVar(&outputDir, "output", "", "output directory for schemas (default: ./schemas, env: OUTPUT_DIR)")
	flag.StringVar(&outputDir, "o", "", "shorthand for --output")
	flag.Parse()

	outputDir = getSchemasDirectory(outputDir)
	for name, value := range openapi.GetOpenAPIDefinitions(referenceCallback) {
		if findAPIEntry(name) == nil {
			continue
		}
		if err := writeSchema(outputDir, name, value); err != nil {
			log.Fatal(err)
		}
	}
}
