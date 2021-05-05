package dqlx

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
	Types       []templateDGraphType
	Imports     map[string]bool
}

type templateDGraphType struct {
	Name       string
	Predicates []templateDGraphPredicate
}

type templateDGraphPredicate struct {
	Name     string
	JsonName string
	GoType   string
	IsEdge   bool
}

// GeneratorOption options for the generator
type GeneratorOption struct {
	Path        string
	PackageName string
}

// GenerateTypes given a schema it generates Go structs definitions
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

func getTypeDefinition(schema *SchemaBuilder) ([]templateDGraphType, map[string]bool) {
	types := make([]templateDGraphType, len(schema.Types))
	imports := map[string]bool{}

	for index, dType := range schema.Types {
		templateType := templateDGraphType{
			Name:       dType.name,
			Predicates: nil,
		}

		// Add fields
		templateType.Predicates = append(templateType.Predicates, templateDGraphPredicate{
			Name:     "Uid",
			JsonName: "uid",
			GoType:   "string",
		})

		for _, predicate := range dType.predicates {
			if predicate.ScalarType == ScalarDateTime {
				imports["time"] = true
			}

			predicateType := dqlTypeToGoType(predicate.ScalarType)

			if predicate.List {
				predicateType = fmt.Sprintf("[]%s", predicateType)
			}

			fieldName := predicate.Name

			if strings.Contains(fieldName, ".") {
				parts := strings.Split(fieldName, ".")
				fieldName = parts[len(parts)-1]
			}

			templateType.Predicates = append(templateType.Predicates, templateDGraphPredicate{
				Name:     toCamelCase(fieldName),
				JsonName: predicate.Name,
				GoType:   predicateType,
				IsEdge:   !isKnownScalarType(predicate.ScalarType),
			})
		}

		// Default DType field
		templateType.Predicates = append(templateType.Predicates, templateDGraphPredicate{
			Name:     "DType",
			JsonName: "dgraph.type",
			GoType:   "[]string",
		})

		types[index] = templateType
	}

	return types, imports
}

var link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")

func toCamelCase(str string) string {
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}
