/*
 * Copyright (c) 2020, 2026, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"
)

const (
	copyright = `///////////////////////////////////////////////////////////////////////////////

    Copyright (c) %s, Oracle and/or its affiliates.
    Licensed under the Universal Permissive License v 1.0 as shown at
    http://oss.oracle.com/licenses/upl.

///////////////////////////////////////////////////////////////////////////////`

	firstParagraph = `///////////////////////////////////////////////////////////////////////////////

NOTE: *** This document must not be manually edited. ***
This document has been generated from the comments in the pkg/api classes.
Any changes should be made by editing the corresponding struct comments.

///////////////////////////////////////////////////////////////////////////////

= Coherence Operator API Docs

A reference guide to the Coherence Operator CRD types.

== Coherence Operator API Docs
This is a reference for the Coherence Operator API types.
These are all the types and fields that are used in the Coherence CRD. 

TIP: This document was generated from comments in the Go structs in the pkg/api/ directory.`

	k8sLink = "https://{k8s-doc-link}/#"
)

var (
	selfLinks = map[string]string{}
)

func main() {
	// hard coded links to types that are not in the k8s docs
	selfLinks["batchv1.CompletionMode"] = "https://pkg.go.dev/k8s.io/api/batch/v1#CompletionMode"
	selfLinks["corev1.DNSPolicy"] = "https://pkg.go.dev/k8s.io/api/core/v1#DNSPolicy"
	selfLinks["corev1.IPFamily"] = "https://pkg.go.dev/k8s.io/api/core/v1#IPFamily"
	selfLinks["corev1.IPFamilyPolicyType"] = "https://pkg.go.dev/k8s.io/api/core/v1#IPFamilyPolicyType"
	selfLinks["corev1.MountPropagationMode"] = "https://pkg.go.dev/k8s.io/api/core/v1#MountPropagationMode"
	selfLinks["corev1.PreemptionPolicy"] = "https://pkg.go.dev/k8s.io/api/core/v1#PreemptionPolicy"
	selfLinks["corev1.PullPolicy"] = "https://pkg.go.dev/k8s.io/api/core/v1#PullPolicy"
	selfLinks["corev1.Protocol"] = "https://pkg.go.dev/k8s.io/api/core/v1#Protocol"
	selfLinks["corev1.RestartPolicy"] = "https://pkg.go.dev/k8s.io/api/core/v1#RestartPolicy"
	selfLinks["corev1.ServiceAffinity"] = "https://pkg.go.dev/k8s.io/api/core/v1#ServiceAffinity"
	selfLinks["corev1.ServiceExternalTrafficPolicyType"] = "https://pkg.go.dev/k8s.io/api/core/v1#ServiceExternalTrafficPolicyType"
	selfLinks["corev1.ServiceType"] = "https://pkg.go.dev/k8s.io/api/core/v1#ServiceType"
	selfLinks["corev1.VolumeSource"] = fmt.Sprintf("%svolume-v1-core", k8sLink)
	selfLinks["intstr.IntOrString"] = "https://pkg.go.dev/k8s.io/apimachinery/pkg/util/intstr#IntOrString"
	selfLinks["corev1.ServiceExternalTrafficPolicy"] = "https://pkg.go.dev/k8s.io/api/core/v1#ServiceExternalTrafficPolicy"
	selfLinks["corev1.ServiceInternalTrafficPolicy"] = "https://pkg.go.dev/k8s.io/api/core/v1#ServiceInternalTrafficPolicy"
	selfLinks["corev1.IPFamilyPolicy"] = "https://pkg.go.dev/k8s.io/api/core/v1#IPFamilyPolicy"

	printAPIDocs(os.Args[1:])
}

func printContents(types []KubeTypes) {
	fmt.Printf("\n=== Table of Contents\n")
	for _, t := range types {
		strukt := t[0]
		if len(t) > 1 {
			fmt.Printf("* <<%s,%s>>\n", strukt.Name, strukt.Name)
		}
	}
}

func printAPIDocs(paths []string) {
	var copyrightYear string
	year := time.Now().Year()
	if year == 2020 {
		copyrightYear = fmt.Sprintf("%d", year)
	} else {
		copyrightYear = fmt.Sprintf("2020, %d", year)
	}

	fmt.Printf(copyright, copyrightYear)
	fmt.Println()
	fmt.Println()
	fmt.Println(firstParagraph)

	types := ParseDocumentationFrom(paths)
	for _, t := range types {
		strukt := t[0]
		selfLinks[strukt.Name] = fmt.Sprintf("<<%s,%s>>", strukt.Name, strukt.Name)
	}

	// we need to parse once more to now add the self links
	types = ParseDocumentationFrom(paths)

	for _, t := range types {
		strukt := t[0]
		if strukt.Name == "Coherence" {
			printType(t, false)
		}
	}

	for _, t := range types {
		strukt := t[0]
		if strukt.Name == "CoherenceJob" {
			printType(t, false)
		}
	}

	sort.Slice(types, func(i, j int) bool {
		first := types[i][0]
		second := types[j][0]
		return first.Name < second.Name
	})

	printContents(types)

	for _, t := range types {
		strukt := t[0]
		if strukt.Name != "Coherence" && strukt.Name != "CoherenceJob" {
			printType(t, true)
		}
	}
}

func printType(t KubeTypes, back bool) {
	strukt := t[0]
	if len(t) > 1 {
		fmt.Printf("\n=== %s\n\n%s\n\n", strukt.Name, strukt.Doc)

		fmt.Println("[cols=\"1,10,1,1\"options=\"header\"]")
		fmt.Println("|===")
		fmt.Println("| Field | Description | Type | Required")
		fields := t[1:]

		for _, f := range fields {
			var d string
			if f.Doc == "" {
				d = "&#160;"
			} else {
				d = strings.ReplaceAll(f.Doc, "\\n", " +\n")
			}
			fmt.Println("m|", f.Name, "|", d, "m|", f.Type, "|", f.Mandatory)
		}
		fmt.Println("|===")
		fmt.Println("")
		if back {
			fmt.Println("<<Table of Contents,Back to TOC>>")
		}
	}
}

// Pair of strings. We keep the name of fields and the doc
type Pair struct {
	Name, Doc, Type string
	Mandatory       bool
}

// KubeTypes is an array to represent all available types in a parsed file. [0] is for the type itself
type KubeTypes []Pair

// ParseDocumentationFrom gets all types' documentation and returns them as an
// array. Each type is again represented as an array (we have to use arrays as we
// need to be sure for the order of the fields). This function returns fields and
// struct definitions that have no documentation as {name, ""}.
func ParseDocumentationFrom(srcs []string) []KubeTypes {
	var docForTypes []KubeTypes

	for _, src := range srcs {
		pkg := astFrom(src)

		for _, kubType := range pkg.Types {
			if structType, ok := kubType.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType); ok {
				var ks KubeTypes
				ks = append(ks, Pair{kubType.Name, fmtRawDoc(kubType.Doc), "", false})
				ks = ProcessFields(structType, ks)
				docForTypes = append(docForTypes, ks)
			}
		}
	}

	return docForTypes
}

func ProcessFields(structType *ast.StructType, ks KubeTypes) KubeTypes {
	for _, field := range structType.Fields.List {
		typeString := fieldType(field.Type)
		fieldMandatory := fieldRequired(field)
		if n := fieldName(field); n != "-" {
			fieldDoc := fmtRawDoc(field.Doc.Text())
			ks = append(ks, Pair{n, fieldDoc, typeString, fieldMandatory})
		} else if strings.Contains(field.Tag.Value, "json:\",inline") {
			if ident, ok := field.Type.(*ast.Ident); ok && ident.Obj != nil {
				if ts, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
					if st, ok := ts.Type.(*ast.StructType); ok {
						ks = ProcessFields(st, ks)
					}
				}
			}
		}
	}
	return ks
}

func astFrom(filePath string) *doc.Package {
	fset := token.NewFileSet()
	m := make(map[string]*ast.File)

	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	m[filePath] = f
	apkg, _ := ast.NewPackage(fset, m, nil, nil)

	return doc.New(apkg, "", 0)
}

func fmtRawDoc(rawDoc string) string {
	var buffer bytes.Buffer
	delPrevChar := func() {
		if buffer.Len() > 0 {
			buffer.Truncate(buffer.Len() - 1) // Delete the last " " or "\n"
		}
	}

	// Ignore all lines after ---
	rawDoc = strings.Split(rawDoc, "---")[0]

	// array to hold and +coh: tags we find in the doc
	var tags []string

	for _, line := range strings.Split(rawDoc, "\n") {
		line = strings.TrimRight(line, " ")
		leading := strings.TrimLeft(line, " ")
		switch {
		case len(line) == 0: // Keep paragraphs
			delPrevChar()
			buffer.WriteString("\n\n")
		case strings.HasPrefix(leading, "+coh:"): // Coherence tag
			tags = append(tags, line)
		case strings.HasPrefix(leading, "TODO"): // Ignore one line TODOs
		case strings.HasPrefix(leading, "+"): // Ignore instructions to go2idl
		default:
			if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
				delPrevChar()
				line = "\n" + line + " +\n" // Replace it with newline. This is useful when we have a line with: "Example:\n\tJSON-someting..."
			} else {
				line += " "
			}
			buffer.WriteString(line)
		}
	}

	postDoc := strings.TrimRight(buffer.String(), "\n")
	// postDoc = strings.Replace(postDoc, "\\\"", "\"", -1) // replace user's \" to "
	// postDoc = strings.Replace(postDoc, "\"", "\\\"", -1) // Escape "
	postDoc = strings.Replace(postDoc, "\n", " +\n", -1)
	postDoc = strings.Replace(postDoc, "\t", "&#160;&#160;&#160;&#160;", -1)
	postDoc = strings.Replace(postDoc, "|", "\\|", -1)

	// process any Coherence tags we found
	postDoc = processCoherenceTags(postDoc, tags)

	return postDoc
}

func processCoherenceTags(doc string, tags []string) string {
	if len(tags) == 0 {
		return doc
	}

	for _, t := range tags {
		if strings.HasPrefix(t, "+coh:doc=") {
			parts := strings.Split(t[9:], ",")
			var link string
			var name string
			if len(parts) == 1 {
				link = parts[0]
				name = parts[0]
			} else {
				link = parts[0]
				name = parts[1]
			}
			doc = fmt.Sprintf("%s +\nsee: <<%s,%s>>", doc, link, name)
		}
	}

	return doc
}

func toLink(typeName string) string {
	selfLink, hasSelfLink := selfLinks[typeName]
	if hasSelfLink {
		return selfLink
	}

	switch {
	case strings.HasPrefix(typeName, "batchv1."):
		return fmt.Sprintf("%s%s-v1-batch[%s]", k8sLink, strings.ToLower(typeName[8:]), escapeTypeName(typeName))
	case strings.HasPrefix(typeName, "corev1."):
		return fmt.Sprintf("%s%s-v1-core[%s]", k8sLink, strings.ToLower(typeName[7:]), escapeTypeName(typeName))
	case strings.HasPrefix(typeName, "metav1."):
		return fmt.Sprintf("%s%s-v1-meta[%s]", k8sLink, strings.ToLower(typeName[7:]), escapeTypeName(typeName))
	}

	return typeName
}

func escapeTypeName(typeName string) string {
	if strings.HasPrefix(typeName, "*") {
		return "&#42;" + typeName[1:]
	}
	return typeName
}

// fieldName returns the name of the field as it should appear in JSON format
// "-" indicates that this field is not part of the JSON representation
func fieldName(field *ast.Field) string {
	jsonTag := ""
	if field.Tag != nil {
		jsonTag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json") // Delete first and last quotation
		if strings.Contains(jsonTag, "inline") {
			return "-"
		}
	}

	jsonTag = strings.Split(jsonTag, ",")[0] // This can return "-"
	if jsonTag == "" {
		if field.Names != nil {
			return field.Names[0].Name
		}
		return field.Type.(*ast.Ident).Name
	}
	return jsonTag
}

// fieldRequired returns whether a field is a required field.
func fieldRequired(field *ast.Field) bool {
	jsonTag := ""
	if field.Tag != nil {
		jsonTag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]).Get("json") // Delete first and last quotation
		return !strings.Contains(jsonTag, "omitempty")
	}

	return false
}

func fieldType(typ ast.Expr) string {
	switch t := typ.(type) {
	case *ast.Ident:
		return toLink(t.Name)
	case *ast.StarExpr:
		return "&#42;" + toLink(fieldType(typ.(*ast.StarExpr).X))
	case *ast.SelectorExpr:
		e := typ.(*ast.SelectorExpr)
		pkg := e.X.(*ast.Ident)
		return toLink(pkg.Name + "." + e.Sel.Name)
	case *ast.ArrayType:
		return "[]" + toLink(fieldType(typ.(*ast.ArrayType).Elt))
	case *ast.MapType:
		mapType := typ.(*ast.MapType)
		return "map[" + toLink(fieldType(mapType.Key)) + "]" + toLink(fieldType(mapType.Value))
	default:
		return ""
	}
}
