package deku

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"html/template"
	"io/fs"
	"io/ioutil"
	"regexp"
	"strings"
)

//go:embed templates/types.tmpl
var typesTmpl string

type templateTypesVariables struct {
	PackageName string
	Types       []TemplateType
	Imports     map[string]bool
}

type TemplateType struct {
	Name   string
	Fields []TemplateField
}

type TemplateField struct {
	Name     string
	JsonName string
	GoType   string
}

type GeneratorOption struct {
	Path        string
	PackageName string
}

func GenerateTypes(schema *SchemaBuilder, options GeneratorOption) error {
	tmpl, err := template.New("types").Parse(typesTmpl)

	if err != nil {
		return err
	}

	out := bytes.Buffer{}

	typeDefinitions, imports := getTypeDefinition(schema)

	err = tmpl.Execute(&out, templateTypesVariables{
		PackageName: options.PackageName,
		Types:       typeDefinitions,
		Imports:     imports,
	})

	if err != nil {
		return err
	}

	formattedCode, err := format.Source(out.Bytes())

	if err != nil {
		return err
	}

	return ioutil.WriteFile(options.Path, formattedCode, fs.ModePerm)
}

func getTypeDefinition(schema *SchemaBuilder) ([]TemplateType, map[string]bool) {
	types := make([]TemplateType, len(schema.Types))
	imports := map[string]bool{}

	for index, dType := range schema.Types {
		templateType := TemplateType{
			Name:   dType.name,
			Fields: nil,
		}

		// Add fields
		for _, predicate := range dType.predicates {
			if predicate.ScalarType == ScalarDateTime {
				imports["time"] = true
			}

			predicateType := dgraphScalarToGoType(predicate.ScalarType)

			if predicate.List {
				predicateType = fmt.Sprintf("[]%s", predicateType)
			}

			fieldName := predicate.Name

			if strings.Contains(fieldName, ".") {
				parts := strings.Split(fieldName, ".")
				fieldName = parts[len(parts)-1]
			}

			templateType.Fields = append(templateType.Fields, TemplateField{
				Name:     toCamelCase(fieldName),
				JsonName: predicate.Name,
				GoType:   predicateType,
			})
		}

		// Default DType field
		templateType.Fields = append(templateType.Fields, TemplateField{
			Name:     "DType",
			JsonName: "dgraph.type",
			GoType:   "[]string",
		})

		types[index] = templateType
	}

	return types, imports
}

func dgraphScalarToGoType(value DGraphScalar) string {
	switch value {
	case ScalarPassword, ScalarString, ScalarUID:
		return "string"
	case ScalarInt:
		return "int64"
	case ScalarDateTime:
		return "time.Time"
	case ScalarBool:
		return "bool"
	case ScalarFloat:
		return "float64"
	}

	return string(value)
}

var link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")

func toCamelCase(str string) string {
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}
