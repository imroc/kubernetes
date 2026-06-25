package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/kube-openapi/pkg/common"
)

func writeSchema(outputDir, name string, value common.OpenAPIDefinition) error {
	group, kind, version := parseNameRef(name)
	dir := filepath.Join(outputDir, group)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir %s: %w", dir, err)
	}
	data, err := value.Schema.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal schema for %s: %w", name, err)
	}
	m := make(map[string]any)
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("unmarshal schema for %s: %w", name, err)
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
		return err
	}
	filename := filepath.Join(dir, fmt.Sprintf("%s_%s.json", strings.ToLower(kind), version))
	if group == "meta.apis.pkg.apimachinery.k8s.io" { // not override existed meta api
		if _, err := os.Stat(filename); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			return nil
		}
	}
	log.Println("write", filename)
	return os.WriteFile(filename, data, 0644)
}
