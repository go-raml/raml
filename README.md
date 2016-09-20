[![Build Status](https://travis-ci.org/go-raml/raml.svg?branch=v0)](https://travis-ci.org/go-raml/raml)

Looking for an active maintainer!
=================================

September 21st, 2016
--------------

Unfortunately, in the past year I haven't been able to dedicate time to actively support this project, as I've been busy with other private endeavours. Since this project is in use in several projects and has several active forks, I'm looking for someone to take over or assist with maintenance.

The basic features I'd suggest are:

* Support for RAML 1.0 spec
* Validate using jsonschema, full integration
* Implement an integration with https://github.com/astaxie/beego or https://github.com/labstack/echo to auto-generate server code
* Generate go client for API
* Additional features, taking inspiration from https://github.com/go-swagger/go-swagger

If you're interested in taking over this project, the discussion can be found here: https://github.com/go-raml/raml/issues/6

Thanks!

--






--







raml
====

An implementation of a RAML parser for Go. Compliant with RAML 0.8.

Introduction
============

RAML is a YAML-based language that describes RESTful APIs. Together with the
YAML specification, this specification provides all the information necessary
to describe RESTful APIs; to create API client-code and API server-code
generators; and to create API user documentation from RAML API definitions.

The **_raml_** package enables Go programs to parse RAML files and valid RAML API
definitions. It was originally developed within [EverythingMe](https://www.everything.me).

Status
------

The **_raml_** package is currently unstable and does not offer any kind of API
stability guarantees.

Installation
============

The yaml package may be installed by running:

    $ go get gopkg.in/raml.v0

Opening that same URL in a browser will present a nice introductory page
containing links to the documentation, source code, and all versions available
for the given package:

https://gopkg.in/raml.v0

The actual implementation of the package is in GitHub:

https://github.com/go-raml/raml

Contributing to development
---------------------------

Typical installation process for developing purposes:

    $ git clone git@github.com:go-raml/raml.git
    $ cd raml
    $ go build
    $ go install
    $ go test

Usage
=====

Usage is very simple:

	package main
	
	import (
		"fmt"
		raml "gopkg.in/raml.v0"
		"github.com/kr/pretty"
	)
	
	func main() {
	
		fileName := "./samples/congo/api.raml"
	
		if apiDefinition, err := raml.ParseFile(fileName); err != nil {
			fmt.Printf("Failed parsing RAML file %s:\n  %s", fileName, err.Error())
		} else {
			fmt.Printf("Successfully parsed RAML file %s!\n\n", fileName)
			pretty.Printf(apiDefinition)
		}
	}

Getting help
============

* Look up the [RAML 0.8](http://raml.org/spec.html) spec if in doubt
* Contact Alon, the maintainer, directly: diamant.alon@gmail.com

Roadmap
=======

TBD.

Reporting Bugs and Contributing Code
====================================

* Want to report a bug or request a feature? Please open [an issue](https://github.com/go-raml/raml/issues/new).
* Want to contribute to **_raml_**? Fork the project and make a pull request. Cool cool cool.

## License

See [LICENSE](https://github.com/go-raml/raml/blob/v0/LICENSE) file.
