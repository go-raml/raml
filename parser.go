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

// This package contains the parser, validator and types that implement the
// RAML specification, as documented here:
// http://raml.org/spec.html
package raml

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	yaml "github.com/advance512/yaml"
	"github.com/kr/pretty"
)

// A RamlError is returned by the ParseFile function when RAML or YAML problems
// are encountered when parsing the RAML document.
// When this error is returned, the value is still parsed partially.
type RamlError struct {
	Errors []string
}

// Parse a RAML file. Returns a raml.APIDefinition value or an error if
// everything is something went wrong.
func ParseFile(filePath string) (*APIDefinition, error) {

	// Get the working directory
	workingDirectory, fileName := filepath.Split(filePath)

	// Read original file contents into a byte array
	mainFileBytes, err := readFileContents(workingDirectory, fileName)

	if err != nil {
		return nil, err
	}

	// Get the contents of the main file
	mainFileBuffer := bytes.NewBuffer(mainFileBytes)

	// Verify the YAML version
	var ramlVersion string
	if firstLine, err := mainFileBuffer.ReadString('\n'); err != nil {
		return nil, fmt.Errorf("Problem reading RAML file (Error: %s)", err.Error())
	} else {

		// We read some data...
		if len(firstLine) >= 10 {
			ramlVersion = firstLine[:10]
		}

		// TODO: Make this smart. We probably won't support multiple RAML
		// versions in the same package - we'll have different branches
		// for different versions. This one is hard-coded to 0.8.
		// Still, would be good to think about this.
		if ramlVersion != "#%RAML 0.8" {
			return nil, errors.New("Input file is not a RAML 0.8 file. Make " +
				"sure the file starts with #%RAML 0.8")
		}
	}

	// Pre-process the original file, following !include directive
	preprocessedContentsBytes, err :=
		preProcess(mainFileBuffer, workingDirectory)

	if err != nil {
		return nil,
			fmt.Errorf("Error preprocessing RAML file (Error: %s)", err.Error())
	}

	pretty.Println(string(preprocessedContentsBytes))

	// Unmarshal into an APIDefinition value
	apiDefinition := new(APIDefinition)
	apiDefinition.RAMLVersion = ramlVersion

	// Go!
	err = yaml.Unmarshal(preprocessedContentsBytes, apiDefinition)

	// Any errors?
	if err != nil {

		return nil, fmt.Errorf("Problems parsing RAML:\n  %s", err.Error())
	}

	// Good.
	return apiDefinition, nil

}

// Reads the contents of a file, returns a bytes buffer
func readFileContents(workingDirectory string, fileName string) ([]byte, error) {

	filePath := filepath.Join(workingDirectory, fileName)

	if fileName == "" {
		return nil, fmt.Errorf("File name cannot be nil: %s", filePath)
	}

	// Read the file
	fileContentsArray, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil,
			fmt.Errorf("Could not read file %s (Error: %s)",
				filePath, err.Error())
	}

	return fileContentsArray, nil
}

// preProcess acts as a preprocessor for a RAML document in YAML format,
// including files referenced via !include. It returns a pre-processed document.
func preProcess(originalContents io.Reader, workingDirectory string) ([]byte, error) {

	// NOTE: Since YAML doesn't support !include directives, and since go-yaml
	// does NOT play nice with !include tags, this has to be done like this.
	// I am considering modifying go-yaml to add custom handlers for specific
	// tags, to add support for !include, but for now - this method is
	// GoodEnough(TM) and since it will only happen once, I am not prematurely
	// optimizing it.

	var preprocessedContents bytes.Buffer

	// Go over each line, looking for !include tags
	scanner := bufio.NewScanner(originalContents)
	var line string

	// Scan the file until we reach EOF or error out
	for scanner.Scan() {
		line = scanner.Text()

		// Did we find an !include directive to handle?
		if idx := strings.Index(line, "!include"); idx != -1 {

			// TODO: Do this better
			includeLength := len("!include ")

			includedFile := line[idx+includeLength:]

			preprocessedContents.Write([]byte(line[:idx]))

			// Get the included file contents
			includedContents, err :=
				readFileContents(workingDirectory, includedFile)

			if err != nil {
				return nil,
					fmt.Errorf("Error including file %s:\n    %s",
						includedFile, err.Error())
			}

			// TODO: Check that you only insert .yaml, .raml, .txt and .md files
			// In case of .raml or .yaml, remove the comments
			// In case of other files, Base64 them first.

			// TODO: Better, step by step checks .. though prolly it'll panic
			// Write text files in the same indentation as the first line
			internalScanner :=
				bufio.NewScanner(bytes.NewBuffer(includedContents))

			// Indent by this much
			firstLine := true
			indentationString := ""

			// Go over each line, write it
			for internalScanner.Scan() {
				internalLine := internalScanner.Text()

				preprocessedContents.WriteString(indentationString)
				if firstLine {
					indentationString = strings.Repeat(" ", idx)
					firstLine = false
				}

				preprocessedContents.WriteString(internalLine)
				preprocessedContents.WriteByte('\n')
			}

		} else {

			// No, just a simple line.. write it
			preprocessedContents.WriteString(line)
			preprocessedContents.WriteByte('\n')
		}
	}

	// Any errors encountered?
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading YAML file: %s", err.Error())
	}

	// Return the preprocessed contents
	return preprocessedContents.Bytes(), nil
}
