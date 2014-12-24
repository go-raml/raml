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

// This file contains all of the RAML types.

// TODO: We don't support !include of non-text files. RAML supports including
//       of many file types.

// "Any" type, for our convenience
type Any interface{}

// For extra clarity
type HTTPCode int      // e.g. 200
type HTTPHeader string // e.g. Content-Length

// The RAML Specification uses collections of named parameters for the
// following properties: URI parameters, query string parameters, form
// parameters, request bodies (depending on the media type), and request
// and response headers.
//
// Some fields are pointers to distinguish Zero values and no values
type NamedParameter struct {

	// NOTE: We currently do not support Named Parameters With Multiple Types.
	// TODO: Add support for Named Parameters With Multiple Types. Should be
	// done sort of like the DefinitionChoice type.

	// The name of the Parameter, as defined by the type containing it.
	Name string
	// TODO: Fill this during the post-processing phase

	// A friendly name used only for display or documentation purposes.
	// If displayName is not specified, it defaults to the property's key
	DisplayName string `yaml:"displayName"` // TODO: Auto-fill this

	// The intended use or meaning of the parameter
	Description string

	// The primitive type of the parameter's resolved value. Can be:
	//
	// Type	Description
	// string	- Value MUST be a string.
	// number	- Value MUST be a number. Indicate floating point numbers as defined by YAML.
	// integer	- Value MUST be an integer. Floating point numbers are not allowed. The integer type is a subset of the number type.
	// date		- Value MUST be a string representation of a date as defined in RFC2616 Section 3.3 [RFC2616]. See Date Representations.
	// boolean	- Value MUST be either the string "true" or "false" (without the quotes).
	// file		- (Applicable only to Form properties) Value is a file. Client generators SHOULD use this type to handle file uploads correctly.
	Type string
	// TODO: Verify the enum options

	// If the enum attribute is defined, API clients and servers MUST verify
	// that a parameter's value matches a value in the enum array
	Enum []Any `yaml:",flow"`

	// The pattern attribute is a regular expression that a parameter of type
	// string MUST match. Regular expressions MUST follow the regular
	// expression specification from ECMA 262/Perl 5. (string only)
	Pattern *string

	// The minLength attribute specifies the parameter value's minimum number
	// of characters (string only)
	MinLength *int `yaml:"minLength"`
	// TODO: go-yaml doesn't raise an error when the minLength isn't an integer!
	// find out why and fix it.

	// The maxLength attribute specifies the parameter value's maximum number
	// of characters (string only)
	MaxLength *int `yaml:"maxLength"`

	// The minimum attribute specifies the parameter's minimum value. (numbers
	// only)
	Minimum *float64

	// The maximum attribute specifies the parameter's maximum value. (numbers
	// only)
	Maximum *float64

	// An example value for the property. This can be used, e.g., by
	// documentation generators to generate sample values for the property.
	Example string

	// The repeat attribute specifies that the parameter can be repeated,
	// i.e. the parameter can be used multiple times
	Repeat *bool // TODO: What does this mean?

	// Whether the parameter and its value MUST be present when a call is made.
	// In general, parameters are optional unless the required attribute is
	// included and its value set to 'true'.
	// For a URI parameter, its default value is 'true'.
	Required bool

	// The default value to use for the property if the property is omitted or
	// its value is not specified
	Default Any

	format Any `ramlFormat:"Named parameters must be mappings. Example: userId: {displayName: 'User ID', description: 'Used to identify the user.', type: 'integer', minimum: 1, example: 5}"`
}

// Headers used in Methods and other types
type Header NamedParameter

// All documentation of the API is of this format.
type Documentation struct {
	Title   string `yaml:"title"`
	Content string `yaml:"content"`
}

// Some method verbs expect the resource to be sent as a request body.
// For example, to create a resource, the request must include the details of
// the resource to create.
// Resources CAN have alternate representations. For example, an API might
// support both JSON and XML representations.
type Body struct {
	mediaType string `yaml:"mediaType"`
	// TODO: Fill this during the post-processing phase

	// The structure of a request or response body MAY be further specified
	// by the schema property under the appropriate media type.
	// The schema key CANNOT be specified if a body's media type is
	// application/x-www-form-urlencoded or multipart/form-data.
	// All parsers of RAML MUST be able to interpret JSON Schema [JSON_SCHEMA]
	// and XML Schema [XML_SCHEMA].
	// Alternatively, the value of the schema field MAY be the name of a schema
	// specified in the root-level schemas property
	Schema string `yaml:"schema"`

	// Brief description
	Description string `yaml:"description"`

	// Example attribute to generate example invocations
	Example string `yaml:"example"`

	// Web forms REQUIRE special encoding and custom declaration.
	// If the API's media type is either application/x-www-form-urlencoded or
	// multipart/form-data, the formParameters property MUST specify the
	// name-value pairs that the API is expecting.
	// The formParameters property is a map in which the key is the name of
	// the web form parameter, and the value is itself a map the specifies
	// the web form parameter's attributes
	FormParameters map[string]NamedParameter `yaml:"formParameters"`
	// TODO: This doesn't make sense in response bodies.. separate types for
	// request and response body?

	Headers map[HTTPHeader]Header `yaml:"headers"`
}

// Container of Body types, necessary because of technical reasons.
type Bodies struct {

	// Instead of using a simple map[HTTPHeader]Body for the body
	// property of the Response and Method, we use the Bodies struct. Why?
	// Because some RAML APIs don't use the MIMEType part, instead relying
	// on the mediaType property in the APIDefinition.
	// So, you might see:
	//
	// responses:
	//   200:
	//     body:
	//       example: "some_example" : "123"
	//
	// and also:
	//
	// responses:
	//   200:
	//     body:
	//       application/json:
	//         example: |
	//           {
	//             "some_example" : "123"
	//           }

	// As in the Body type.
	DefaultSchema string `yaml:"schema"`

	// As in the Body type.
	DefaultDescription string `yaml:"description"`

	// As in the Body type.
	DefaultExample string `yaml:"example"`

	// As in the Body type.
	DefaultFormParameters map[string]NamedParameter `yaml:"formParameters"`

	// TODO: Is this ever used? I think I put it here by mistake.
	//Headers               map[HTTPHeader]Header     `yaml:"headers"`

	// Resources CAN have alternate representations. For example, an API
	// might support both JSON and XML representations. This is the map
	// between MIME-type and the body definition related to it.
	ForMIMEType map[string]Body `yaml:",regexp:.*"`

	// TODO: For APIs without a priori knowledge of the response types for
	// their responses, "*/*" MAY be used to indicate that responses that do
	// not matching other defined data types MUST be accepted. Processing
	// applications MUST match the most descriptive media type first if
	// "*/*" is used.
}

// Resource methods MAY have one or more responses.
type Response struct {

	// HTTP status code of the response
	HTTPCode HTTPCode
	// TODO: Fill this during the post-processing phase

	// Clarifies why the response was emitted. Response descriptions are
	// particularly useful for describing error conditions.
	Description string

	// An API's methods may support custom header values in responses
	Headers map[HTTPHeader]Header `yaml:"headers"`

	// TODO: API's may include the the placeholder token {?} in a header name
	// to indicate that any number of headers that conform to the specified
	// format can be sent in responses. This is particularly useful for
	// APIs that allow HTTP headers that conform to some naming convention
	// to send arbitrary, custom data.

	// Each response MAY contain a body property. Responses that can return
	// more than one response code MAY therefore have multiple bodies defined.
	Bodies Bodies `yaml:"body"`
}

// A ResourceType/Trait/SecurityScheme choice contains the name of a
// ResourceType/Trait/SecurityScheme as well as the parameters used to create
// an instance of it.
// Parameters MUST be of type string.
type DefinitionParameters map[string]string
type DefinitionChoice struct {
	Name string

	// The definitions of resource types and traits MAY contain parameters,
	// whose values MUST be specified when applying the resource type or trait,
	// UNLESS the parameter corresponds to a reserved parameter name, in which
	// case its value is provided by the processing application.
	// Same goes for security schemes.
	Parameters DefinitionParameters
}

// Unmarshal a node which MIGHT be a simple string or a
// map[string]DefinitionParameters
func (dc *DefinitionChoice) UnmarshalYAML(unmarshaler func(interface{}) error) error {

	simpleDefinition := new(string)
	parameterizedDefinition := make(map[string]DefinitionParameters)

	var err error

	// Unmarshal into a string
	if err = unmarshaler(simpleDefinition); err == nil {
		dc.Name = *simpleDefinition
		dc.Parameters = nil
	} else if err = unmarshaler(parameterizedDefinition); err == nil {
		// Didn't work? Now unmarshal into a map
		for choice, params := range parameterizedDefinition {
			dc.Name = choice
			dc.Parameters = params
		}
	}

	// Still didn't work? Panic

	return err
}

// A trait is a partial method definition that, like a method, can provide
// method-level properties such as description, headers, query string
// parameters, and responses. Methods that use one or more traits inherit
// those traits' properties.
type Trait struct {

	// TODO: Parameters MUST be indicated in resource type and trait definitions
	// by double angle brackets (double chevrons) enclosing the parameter name;
	// for example, "<<tokenName>>".

	// TODO: Auto-fill the methodName parameter

	// TODO: In trait definitions, there is one reserved parameter name,
	// methodName, in addition to the resourcePath and resourcePathName. The
	// processing application MUST set the value of the methodName parameter
	// to the inheriting method's name. The processing application MUST set
	// the values of the resourcePath and resourcePathName parameters the same
	// as in resource type definitions.

	// TODO: Parameter values MAY further be transformed by applying one of
	// the following functions:
	// * The !singularize function MUST act on the value of the parameter
	// by a locale-specific singularization of its original value. The only
	// locale supported by this version of RAML is United States English.
	// * The !pluralize function MUST act on the value of the parameter by a
	// locale-specific pluralization of its original value. The only locale
	// supported by this version of RAML is United States English.

	Name string
	// TODO: Fill this during the post-processing phase

	// The usage property of a resource type or trait is used to describe how
	// the resource type or trait should be used
	Usage string

	// Briefly describes what the method does to the resource
	Description string

	// As in Method.
	Bodies Bodies `yaml:"body"`

	// As in Method.
	Headers map[HTTPHeader]Header `yaml:"headers"`

	// As in Method.
	Responses map[HTTPCode]Response `yaml:"responses"`

	// As in Method.
	QueryParameters map[string]NamedParameter `yaml:"queryParameters"`

	// As in Method.
	Protocols []string `yaml:"protocols"`

	// When defining resource types and traits, it can be useful to capture
	// patterns that manifest several levels below the inheriting resource or
	// method, without requiring the creation of the intermediate levels.
	// For example, a resource type definition may describe a body parameter
	// that will be used if the API defines a post method for that resource,
	// but the processing application should not create the post method itself.
	//
	// This optional structure key indicates that the value of the property
	// should be applied if the property name itself (without the question
	// mark) is already defined (whether explicitly or implicitly) at the
	// corresponding level in that resource or method.
	OptionalBodies          Bodies                    `yaml:"body?"`
	OptionalHeaders         map[HTTPHeader]Header     `yaml:"headers?"`
	OptionalResponses       map[HTTPCode]Response     `yaml:"responses?"`
	OptionalQueryParameters map[string]NamedParameter `yaml:"queryParameters?"`
}

// Method that is part of a ResourceType. DIfferentiated from Traits since it
// doesn't contain Usage, optional fields etc.
type ResourceTypeMethod struct {
	Name string
	// TODO: Fill this during the post-processing phase

	// Briefly describes what the method does to the resource
	Description string

	// As in Method.
	Bodies Bodies `yaml:"body"`
	// TODO: Check - how does the mediaType play play here? What it do?

	// As in Method.
	Headers map[HTTPHeader]Header `yaml:"headers"`

	// As in Method.
	Responses map[HTTPCode]Response `yaml:"responses"`

	// As in Method.
	QueryParameters map[string]NamedParameter `yaml:"queryParameters"`

	// As in Method.
	Protocols []string `yaml:"protocols"`
}

// Resource and method declarations are frequently repetitive. For example, if
// an API requires OAuth authentication, the API definition must include the
// access_token query string parameter (which is defined by the queryParameters
// property) in all the API's resource method declarations.
//
// Moreover, there are many advantages to reusing patterns across multiple
// resources and methods. For example, after defining a collection-type
// resource's characteristics, that definition can be applied to multiple
// resources. This use of patterns encouraging consistency and reduces
// complexity for both servers and clients.
//
// A resource type is a partial resource definition that, like a resource, can
// specify a description and methods and their properties. Resources that use
// a resource type inherit its properties, such as its methods.
type ResourceType struct {

	// TODO: Auto-fill the resourcePath and resourcePathName parameters
	// Remove mediaTypeExtension.

	// TODO: Parameters MUST be indicated in resource type and trait definitions
	// by double angle brackets (double chevrons) enclosing the parameter name;
	// for example, "<<tokenName>>".

	// TODO: In resource type definitions, there are two reserved parameter
	// names: resourcePath and resourcePathName. The processing application
	// MUST set the values of these reserved parameters to the inheriting
	// resource's path (for example, "/users") and the part of the path
	// following the rightmost "/" (for example, "users"), respectively.
	// Processing applications MUST also omit the value of any
	// mediaTypeExtension found in the resource's URI when setting
	// resourcePath and resourcePathName.

	// TODO: Parameter values MAY further be transformed by applying one of
	// the following functions:
	// * The !singularize function MUST act on the value of the parameter
	// by a locale-specific singularization of its original value. The only
	// locale supported by this version of RAML is United States English.
	// * The !pluralize function MUST act on the value of the parameter by a
	// locale-specific pluralization of its original value. The only locale
	// supported by this version of RAML is United States English.

	// Name of the resource type
	Name string
	// TODO: Fill this during the post-processing phase

	// The usage property of a resource type or trait is used to describe how
	// the resource type or trait should be used
	Usage string

	// Briefly describes what the resource type
	Description string

	// As in Resource.
	UriParameters map[string]NamedParameter `yaml:"uriParameters"`

	// As in Resource.
	BaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters"`

	// In a RESTful API, methods are operations that are performed on a
	// resource. A method MUST be one of the HTTP methods defined in the
	// HTTP version 1.1 specification [RFC2616] and its extension,
	// RFC5789 [RFC5789].
	Get    *ResourceTypeMethod `yaml:"get"`
	Head   *ResourceTypeMethod `yaml:"head"`
	Post   *ResourceTypeMethod `yaml:"post"`
	Put    *ResourceTypeMethod `yaml:"put"`
	Delete *ResourceTypeMethod `yaml:"delete"`
	Patch  *ResourceTypeMethod `yaml:"patch"`

	// When defining resource types and traits, it can be useful to capture
	// patterns that manifest several levels below the inheriting resource or
	// method, without requiring the creation of the intermediate levels.
	// For example, a resource type definition may describe a body parameter
	// that will be used if the API defines a post method for that resource,
	// but the processing application should not create the post method itself.
	//
	// This optional structure key indicates that the value of the property
	// should be applied if the property name itself (without the question
	// mark) is already defined (whether explicitly or implicitly) at the
	// corresponding level in that resource or method.
	OptionalUriParameters     map[string]NamedParameter `yaml:"uriParameters?"`
	OptionalBaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters?"`
	OptionalGet               *ResourceTypeMethod       `yaml:"get?"`
	OptionalHead              *ResourceTypeMethod       `yaml:"head?"`
	OptionalPost              *ResourceTypeMethod       `yaml:"post?"`
	OptionalPut               *ResourceTypeMethod       `yaml:"put?"`
	OptionalDelete            *ResourceTypeMethod       `yaml:"delete?"`
	OptionalPatch             *ResourceTypeMethod       `yaml:"patch?"`
}

// A trait-like structure to a security scheme mechanism so as to extend
// the mechanism, such as specifying response codes, HTTP headers or custom
// documentation.
type SecuritySchemeMethod struct {
	Bodies          Bodies                    `yaml:"body"`
	Headers         map[HTTPHeader]Header     `yaml:"headers"`
	Responses       map[HTTPCode]Response     `yaml:"responses"`
	QueryParameters map[string]NamedParameter `yaml:"queryParameters"`
}

// Most REST APIs have one or more mechanisms to secure data access, identify
// requests, and determine access level and data visibility.
type SecurityScheme struct {
	Name string
	// TODO: Fill this during the post-processing phase

	// Briefly describes the security scheme
	Description string

	// The type attribute MAY be used to convey information about
	// authentication flows and mechanisms to processing applications
	// such as Documentation Generators and Client generators.
	Type string
	// TODO: Verify that it is of the values accepted: "OAuth 1.0",
	// "OAuth 2.0", "Basic Authentication", "Digest Authentication",
	// "x-{other}"

	// The describedBy attribute MAY be used to apply a trait-like structure
	// to a security scheme mechanism so as to extend the mechanism, such as
	// specifying response codes, HTTP headers or custom documentation.
	// This extension allows API designers to describe security schemes.
	// As a best practice, even for standard security schemes, API designers
	// SHOULD describe the security schemes' required artifacts, such as
	// headers, URI parameters, and so on.
	// Including the security schemes' description completes an API's documentation.
	DescribedBy SecuritySchemeMethod

	// The settings attribute MAY be used to provide security schema-specific
	// information. Depending on the value of the type parameter, its attributes
	// can vary.
	Settings map[string]Any
	// TODO: Verify OAuth 1.0, 2.0 settings
	// TODO: Add to documentaiotn

	// If the scheme's type is x-other, API designers can use the properties
	// in this mapping to provide extra information to clients that understand
	// the x-other type.
	Other map[string]string
}

// Methods are operations that are performed on a resource
type Method struct {
	Name string
	// TODO: Fill this during the post-processing phase

	// Briefly describes what the method does to the resource
	Description string

	// Applying a securityScheme definition to a method overrides whichever
	// securityScheme has been defined at the root level. To indicate that
	// the method is protected using a specific security scheme, the method
	// MUST be defined by using the securedBy attribute
	// Custom parameters can be provided to the security scheme.
	SecuredBy []DefinitionChoice `yaml:"securedBy"`
	// TODO: To indicate that the method may be called without applying any
	// securityScheme, the method may be annotated with the null securityScheme.

	// The method's non-standard HTTP headers. The headers property is a map
	// in which the key is the name of the header, and the value is itself a
	// map specifying the header attributes.
	Headers map[HTTPHeader]Header `yaml:"headers"`
	// TODO: Examples for headers are REQUIRED.
	// TODO: If the header name contains the placeholder token {*}, processing
	// applications MUST allow requests to send any number of headers that
	// conform to the format specified, with {*} replaced by 0 or more valid
	// header characters, and offer a way for implementations to add an
	// arbitrary number of such headers. This is particularly useful for APIs
	// that allow HTTP headers that conform to custom naming conventions to
	// send arbitrary, custom data.

	// A RESTful API method can be reached HTTP, HTTPS, or both.
	// A method can override an API's protocols value for that single method
	// by setting a different value for the fields.
	Protocols []string `yaml:"protocols"`

	// The queryParameters property is a map in which the key is the query
	// parameter's name, and the value is itself a map specifying the query
	//  parameter's attributes
	QueryParameters map[string]NamedParameter `yaml:"queryParameters"`

	// Some method verbs expect the resource to be sent as a request body.
	// A method's body is defined in the body property as a hashmap, in which
	// the key MUST be a valid media type.
	Bodies Bodies `yaml:"body"`
	// TODO: Check - how does the mediaType play play here? What it do?

	// Resource methods MAY have one or more responses. Responses MAY be
	// described using the description property, and MAY include example
	// attributes or schema properties.
	// Responses MUST be a map of one or more HTTP status codes, where each
	// status code itself is a map that describes that status code.
	Responses map[HTTPCode]Response `yaml:"responses"`

	// Methods may specify one or more traits from which they inherit using the
	// is property
	Is []DefinitionChoice `yaml:"is"`
	// TODO: Add support for inline traits?
}

// A resource is the conceptual mapping to an entity or set of entities.
type Resource struct {

	// Resources are identified by their relative URI, which MUST begin with
	// a slash (/).
	URI string
	// TODO: Fill this during the post-processing phase

	// A resource defined as a child property of another resource is called a
	// nested resource, and its property's key is its URI relative to its
	// parent resource's URI. If this is not nil, then this resource is a
	// child resource.
	Parent *Resource
	// TODO: Fill this during the post-processing phase

	// A friendly name to the resource
	DisplayName string

	// Briefly describes the resource
	Description string

	// A securityScheme may also be applied to a resource by using the
	// securedBy key, which is equivalent to applying the securityScheme to
	// all methods of this Resource.
	// Custom parameters can be provided to the security scheme.
	SecuredBy []DefinitionChoice `yaml:"securedBy"`
	// TODO: To indicate that the method may be called without applying any
	// securityScheme, the method may be annotated with the null securityScheme.

	// A resource or a method can override a base URI template's values.
	// This is useful to restrict or change the default or parameter selection
	// in the base URI. The baseUriParameters property MAY be used to override
	// any or all parameters defined at the root level baseUriParameters
	// property, as well as base URI parameters not specified at the root level.
	// In a resource structure of resources and nested resources with their
	// methods, the most specific baseUriParameter fully overrides any
	// baseUriParameter definition made before
	BaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters"`

	// Template URIs containing URI parameters can be used to define a
	// resource's relative URI when it contains variable elements.
	// The values matched by URI parameters cannot contain slash (/) characters
	UriParameters map[string]NamedParameter `yaml:"uriParameters"`

	// TODO: If a URI parameter in a resource's relative URI is not explicitly
	// described in a uriParameters property for that resource, it MUST still
	// be treated as a URI parameter with defaults as specified in the Named
	// Parameters section of this specification. Its type is "string", it is
	// required, and its displayName is its name (i.e. without the surrounding
	// curly brackets [{] and [}]). In the example below, the top-level
	// resource has two URI parameters, "folderId" and "fileId

	// TOOD: A special uriParameter, mediaTypeExtension, is a reserved
	// parameter. It may be specified explicitly in a uriParameters property
	// or not specified explicitly, but its meaning is reserved: it is used
	// by a client to specify that the body of the request or response be of
	// the associated media type. By convention, a value of .json is
	// equivalent to an Accept header of application/json and .xml is
	// equivalent to an Accept header of text/xml.

	// Resources may specify the resource type from which they inherit using
	// the type property. The resource type may be defined inline as the value
	// of the type property (directly or via an !include), or the value of
	// the type property may be the name of a resource type defined within
	// the root-level resourceTypes property.
	// NOTE: inline not currently supported.
	Type *DefinitionChoice `yaml:"type"`
	// TODO: Add support for inline ResourceTypes

	// A resource may use the is property to apply the list of traits to all
	// its methods.
	Is []DefinitionChoice `yaml:"is"`
	// TODO: Add support for inline traits?

	// In a RESTful API, methods are operations that are performed on a
	// resource. A method MUST be one of the HTTP methods defined in the
	// HTTP version 1.1 specification [RFC2616] and its extension,
	// RFC5789 [RFC5789].
	Get    *Method `yaml:"get"`
	Head   *Method `yaml:"head"`
	Post   *Method `yaml:"post"`
	Put    *Method `yaml:"put"`
	Delete *Method `yaml:"delete"`
	Patch  *Method `yaml:"patch"`

	// A resource defined as a child property of another resource is called a
	// nested resource, and its property's key is its URI relative to its
	// parent resource's URI.
	Nested map[string]*Resource `yaml:",regexp:/.*"`
}

// TODO: Resource.GetBaseURIParameter --> includeds APIDefinition BURIParams..
// TODO: Resource.GetAbsoluteURI

// The API Definition describes the basic information of an API, such as its
// title and base URI, and describes how to define common schema references.
type APIDefinition struct {

	// RAML 0.8
	RAMLVersion string `yaml:"raml_version"`

	// The title property is a short plain text description of the RESTful API.
	// The title property's value SHOULD be suitable for use as a title for the
	// contained user documentation
	Title string `yaml:"title"`

	// If RAML API definition is targeted to a specific API version, it should
	// be noted here
	Version string `yaml:"version"`

	// A RESTful API's resources are defined relative to the API's base URI.
	// If the baseUri value is a Level 1 Template URI, the following reserved
	// base URI parameters are available for replacement:
	//
	// version - The content of the version field.
	BaseUri string
	// TODO: If a URI template variable in the base URI is not explicitly
	// described in a baseUriParameters property, and is not specified in a
	// resource-level baseUriParameters property, it MUST still be treated as
	// a base URI parameter with defaults as specified in the Named Parameters
	//  section of this specification. Its type is "string", it is required,
	// and its displayName is its name (i.e. without the surrounding curly
	// brackets [{] and [}]).

	// A resource or a method can override a base URI template's values.
	// This is useful to restrict or change the default or parameter selection
	// in the base URI. The baseUriParameters property MAY be used to override
	// any or all parameters defined at the root level baseUriParameters
	// property, as well as base URI parameters not specified at the root level.
	// In a resource structure of resources and nested resources with their
	// methods, the most specific baseUriParameter fully overrides any
	// baseUriParameter definition made before
	BaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters"`
	// TODO: Generate these also from the baseUri

	// Level 1 URI custom parameters, which are useful in a variety of scenario.
	// URI parameters can be further defined by using the uriParameters
	// property. The use of uriParameters is OPTIONAL. The uriParameters
	// property MUST be a map in which each key MUST be the name of the URI
	// parameter as defined in the baseUri property. The uriParameters CANNOT
	// contain a key named version because it is a reserved URI parameter name.
	UriParameters map[string]NamedParameter `yaml:"uriParameters"`

	// A RESTful API can be reached HTTP, HTTPS, or both
	Protocols []string `yaml:"protocols"`

	// The media types returned by API responses, and expected from API
	// requests that accept a body, MAY be defaulted by specifying the
	// mediaType property.
	// The property's value MAY be a single string with a valid media type:
	//
	// One of the following YAML media types:
	// * text/yaml
	// * text/x-yaml
	// * application/yaml
	// * application/x-yaml*
	//
	// Any type from the list of IANA MIME Media Types,
	// http://www.iana.org/assignments/media-types
	// A custom type that conforms to the regular expression:
	// * "application\/[A-Za-z.-0-1]*+?(json|xml)"
	MediaType string `yaml:"mediaType"`

	// The schemas property specifies collections of schemas that could be
	// used anywhere in the API definition.
	// The value of the schemas property is an array of maps; in each map,
	// the keys are the schema name, and the values are schema definitions:
	// []map[SchemaName]SchemaString
	Schemas []map[string]string
	// TODO: Flatten the arrays of maps here.

	// The securitySchemes property MUST be used to specify an API's security
	// mechanisms, including the required settings and the authentication
	// methods that the API supports.
	// []map[SchemeName]SecurityScheme
	SecuritySchemes []map[string]SecurityScheme `yaml:"securitySchemes"`
	// TODO: Flatten the arrays of maps here.

	// To apply a securityScheme definition to every method in an API, the
	// API MAY be defined using the securedBy attribute. This specifies that
	// all methods in the API are protected using that security scheme.
	// Custom parameters can be provided to the security scheme.
	SecuredBy []DefinitionChoice `yaml:"securedBy"`

	// The API definition can include a variety of documents that serve as a
	// user guides and reference documentation for the API. Such documents can
	// clarify how the API works or provide business context.
	// All the sections are in the order in which the documentation is declared.
	Documentation []Documentation `yaml:"documentation"`

	// To apply a trait definition to a method, so that the method inherits the
	// trait's characteristics, the method MUST be defined by using the is
	// attribute. The value of the is attribute MUST be an array of any number
	// of elements, each of which MUST be a) one or more trait keys (names)
	// included in the traits declaration, or b) one or more trait definition
	// maps.
	// []map[TraitName]Trait
	Traits []map[string]Trait `yaml:"traits"`
	// TODO: Flatten the arrays of maps here.

	// The resourceTypes and traits properties are declared at the API
	// definition's root level with the resourceTypes and traits property keys,
	// respectively. The value of each of these properties is an array of maps;
	// in each map, the keys are resourceType or trait names, and the values
	// are resourceType or trait definitions, respectively.
	// []map[ResourceTypeName]ResourceType
	ResourceTypes []map[string]ResourceType `yaml:"resourceTypes"`
	// TODO: Flatten the arrays of maps here.

	// Resources are identified by their relative URI, which MUST begin with a
	// slash (/). A resource defined as a root-level property is called a
	// top-level resource. Its property's key is the resource's URI relative
	// to the baseUri. A resource defined as a child property of another
	// resource is called a nested resource, and its property's key is its
	// URI relative to its parent resource's URI.
	Resources map[string]Resource `yaml:",regexp:/.*"`
}

// This function receives a path, splits it and traverses the resource
// tree to find the appropriate resource
func (r *APIDefinition) GetResource(path string) *Resource {
	return nil
}
