/*
 * Copyright (c) 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */
// crd-fixup supplies the existing string zero value for optional list-map keys.
// Kubernetes requires every list-map key to be required or defaulted. Embedded
// upstream API structs can expose optional keys even when the containing CRD list
// is a map. Defaulting preserves the Go zero-value behavior without changing
// server-side apply list semantics.
package main

import (
	"fmt"
	"os"
	"sigs.k8s.io/yaml"
)

func main() {
	for _, name := range os.Args[1:] {
		data, err := os.ReadFile(name)
		if err != nil {
			panic(err)
		}
		var schema map[string]interface{}
		if err := yaml.Unmarshal(data, &schema); err != nil {
			panic(err)
		}
		fix(schema)
		data, err = yaml.Marshal(schema)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(name, data, 0644); err != nil {
			panic(err)
		}
	}
}

func fix(node interface{}) {
	switch value := node.(type) {
	case map[string]interface{}:
		if keys, ok := value["x-kubernetes-list-map-keys"].([]interface{}); ok {
			items, _ := value["items"].(map[string]interface{})
			properties, _ := items["properties"].(map[string]interface{})
			required, _ := items["required"].([]interface{})
			for _, key := range keys {
				property, _ := properties[fmt.Sprint(key)].(map[string]interface{})
				isRequired := false
				for _, name := range required {
					if name == key {
						isRequired = true
					}
				}
				if property != nil && !isRequired && property["default"] == nil && property["type"] == "string" {
					property["default"] = ""
				}
			}
		}
		for _, child := range value {
			fix(child)
		}
	case []interface{}:
		for _, child := range value {
			fix(child)
		}
	}
}
