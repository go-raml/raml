// Copyright 2014 DoAT. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation and/or
//    other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED “AS IS” WITHOUT ANY WARRANTIES WHATSOEVER.
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
// THE IMPLIED WARRANTIES OF NON INFRINGEMENT, MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE ARE HEREBY DISCLAIMED. IN NO EVENT SHALL DoAT OR CONTRIBUTORS
// BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// // THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
// EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// The views and conclusions contained in the software and documentation are those of
// the authors and should not be interpreted as representing official policies,
// either expressed or implied, of DoAT.

package raml

// This file contains all code related to YAML and RAML errors.

import (
	"fmt"
	"strings"

	yaml "github.com/advance512/yaml"
)

// A RamlError is returned by the ParseFile function when RAML or YAML problems
// are encountered when parsing the RAML document.
type RamlError struct {
	Errors []string
}

func (e *RamlError) Error() string {
	return fmt.Sprintf("Error parsing RAML:\n  %s\n",
		strings.Join(e.Errors, "\n  "))
}

// Populate the RAML error value with converted YAML error strings (with
// additional context)
func populateRAMLError(ramlError *RamlError,
	yamlErrors *yaml.TypeError) {

	// Go over the errors
	for _, currErr := range yamlErrors.Errors {

		// Create the RAML errors
		ramlError.Errors =
			append(ramlError.Errors, convertYAMLError(currErr))
	}
}

// Convert a YAML error string into RAML error string, with more context
func convertYAMLError(yamlError string) string {

	if strings.Contains(yamlError, "cannot unmarshal") {

		yamlErrorParts := strings.Split(yamlError, " ")

		if len(yamlErrorParts) >= 7 {

			fmt.Println(yamlError)

			var ok bool
			var source string
			var target string
			var targetName string
			line := yamlErrorParts[1]
			line = line[:len(line)-1]

			// TODO: support more complex types:
			// map[string]raml.NamedParameter -->
			// detect map, format to:
			//   "mapping of %s to %s", ramlTypeNames["string"], ramlTypeNames["raml.NamedParameter"]
			// if "string" is not found, use the key, i.e. "string" in this case.
			// so the output would be:
			//   mapping of string to named parameter

			// TODO: instead of having string in the key of some mappings,
			// perhaps use a type alias:
			//   type Name string
			//   map[Name]NamedParameter
			// would output:
			//   mapping of name string to named parameter

			if source, ok = yamlTypeToName[yamlErrorParts[4]]; !ok {
				source = yamlErrorParts[4]
			}
			fmt.Println("source: ", source)

			if source == "string" {
				source = fmt.Sprintf("string (got %s)", yamlErrorParts[5])
				target = yamlErrorParts[7]
			} else {
				target = yamlErrorParts[6]

			}
			if targetName, ok = ramlTypeNames[target]; !ok {
				targetName = target
			}

			target, _ = ramlTypes[target]

			return fmt.Sprintf("line %s: %s cannot be of "+
				"type %s, must be %s", line, targetName, source, target)

		}
	}

	// Otherwise
	return fmt.Sprintf("YAML error, %s", yamlError)
}

var yamlTypeToName map[string]string = map[string]string{
	"!!seq":       "sequence",
	"!!map":       "mapping",
	"!!int":       "integer",
	"!!str":       "string",
	"!!null":      "null",
	"!!bool":      "boolean",
	"!!float":     "float",
	"!!timestamp": "timestamp",
	"!!binary":    "binary",
	"!!merge":     "merge",
}

var ramlTypeNames map[string]string = map[string]string{
	"string": "string value",
	"int":    "numeric value",
	"raml.NamedParameter":       "named parameter",
	"raml.HTTPCode":             "HTTP code",
	"raml.HTTPHeader":           "HTTP header",
	"raml.Header":               "header",
	"raml.Documentation":        "documentation",
	"raml.Body":                 "body",
	"raml.Response":             "response",
	"raml.DefinitionParameters": "definition parameters",
	"raml.DefinitionChoice":     "definition choice",
	"raml.Trait":                "trait",
	"raml.ResourceTypeMethod":   "resource type method",
	"raml.ResourceType":         "resource type",
	"raml.SecuritySchemeMethod": "security scheme method",
	"raml.SecurityScheme":       "security scheme",
	"raml.Method":               "method",
	"raml.Resource":             "resource",
	"raml.APIDefinition":        "API definition",
}

var ramlTypes map[string]string = map[string]string{
	"string": "string",
	"int":    "integer",
	"raml.NamedParameter":       "mapping",
	"raml.HTTPCode":             "integer",
	"raml.HTTPHeader":           "string",
	"raml.Header":               "mapping",
	"raml.Documentation":        "mapping",
	"raml.Body":                 "mapping",
	"raml.Response":             "mapping",
	"raml.DefinitionParameters": "mapping",
	"raml.DefinitionChoice":     "string or mapping",
	"raml.Trait":                "mapping",
	"raml.ResourceTypeMethod":   "mapping",
	"raml.ResourceType":         "mapping",
	"raml.SecuritySchemeMethod": "mapping",
	"raml.SecurityScheme":       "mapping",
	"raml.Method":               "mapping",
	"raml.Resource":             "mapping",
	"raml.APIDefinition":        "mapping",
}
