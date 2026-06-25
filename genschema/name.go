package main

import (
	"fmt"
	"slices"
	"strings"

	"k8s.io/kube-openapi/pkg/validation/spec"
)

func referenceCallback(name string) spec.Ref {
	group, kind, version := parseNameRef(name)
	return spec.MustCreateRef(fmt.Sprintf(`../%s/%s_%s.json`, group, strings.ToLower(kind), version))
}

// parseNameRef converts an openapi type name into its group, kind, and version components.
//
// Examples:
//   - "k8s.io/kubelet/config/v1/KubeletConfiguration" → group="kubelet.config.k8s.io", version="v1", kind="KubeletConfiguration"
//   - "k8s.io/kubernetes/pkg/apis/core.Sysctl" (internal, no version) → group="api_core", version="", kind="Sysctl"
//
// It finds the first matching APIEntry, replaces the prefix with the entry's Group
// (if non-empty), then splits on "/" to extract the components:
// - The last segment is "version.Kind" (for versioned types) or just "Kind" (for internal types)
// - The remaining segments (reversed and joined by ".") form the group
func parseNameRef(name string) (group, kind, version string) {
	if entry := findAPIEntry(name); entry != nil && entry.Group != "" {
		name = strings.Replace(name, entry.Prefix, entry.Group, 1)
	}
	ss := strings.Split(name, "/")
	vk := ss[len(ss)-1]
	ss = ss[0 : len(ss)-1]
	slices.Reverse(ss)
	group = strings.Join(ss, ".")
	ss = strings.Split(vk, ".")
	if len(ss) == 1 {
		// Internal type without version (e.g. "k8s.io/kubernetes/pkg/apis/core.Sysctl")
		kind = ss[0]
	} else {
		version = ss[0]
		kind = ss[1]
	}
	return
}
