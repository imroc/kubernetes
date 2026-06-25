# Agones support

1.  `// +optional` is added to `Scheduling` field in `GameServerAllocation` Spec.
2.  All functions and test code is been removed.
3.  The agonesv1 dependencies types are copied.
4.  `+k8s:openapi-gen=true` is added to the `doc.go` in both api package.
5.  Code path added to `go.work`.
6.  `go.mod` is added to `agones` folder.
