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
// Supports two formats:
//  1. Go import path (contains "/"): "k8s.io/kubelet/config/v1.KubeletConfiguration"
//  2. Model name (no "/"): "io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta"
//
// For import path format:
//   - Finds the first matching APIEntry and replaces prefix with Group
//   - Splits on "/" to extract components
//   - Last segment is "version.Kind"
//
// For model name format:
//   - All segments are dot-separated
//   - Last segment is Kind, second-to-last is version
//   - Remaining segments (reversed) form the group
func parseNameRef(name string) (group, kind, version string) {
	if !strings.Contains(name, "/") {
		// Model name format: io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta
		return parseModelName(name)
	}
	// Go import path format
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
		kind = ss[0]
	} else {
		version = ss[0]
		kind = ss[1]
	}
	return
}

// parseModelName handles the model name format (e.g. "io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta").
//
// Model name segments: <domain...>.<path...>.<version>.<Kind>
// - Kind is the last segment
// - Version is the second-to-last segment (e.g. "v1", "v1beta1", or package name like "resource")
// - Remaining segments (reversed) form the group
//
// Examples:
//   - "io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta" → group="meta.apis.pkg.apimachinery.k8s.io", version="v1", kind="ObjectMeta"
//   - "io.k8s.apimachinery.pkg.api.resource.Quantity" → group="api.pkg.apimachinery.k8s.io", version="resource", kind="Quantity"
func parseModelName(name string) (group, kind, version string) {
	parts := strings.Split(name, ".")
	if len(parts) < 3 {
		kind = name
		return
	}
	kind = parts[len(parts)-1]
	version = parts[len(parts)-2]
	groupParts := parts[:len(parts)-2]
	slices.Reverse(groupParts)
	group = strings.Join(groupParts, ".")
	return
}
