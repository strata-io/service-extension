# Service Extension Development Library
Integrating identity systems requires extreme configurability. Service Extensions 
give administrators the ability to customize the behavior of the Maverics Orchestrator
to suit the particular needs of their integration. This library aids in the 
development of Service Extensions.

## Local Setup
To develop Service Extensions on your local machine, follow the instructions below. 
The steps listed below are only to aid in local development and are not required to run
the Maverics platform.

### Project Initialization
Download the latest version of Go. [Instructions](https://go.dev/doc/install) can be 
found on the Go website.

After Go is downloaded, create a project directory. For organizational purpose, it is
recommended to create this project alongside other Maverics configuration files. For
example, in an `/etc/maverics/extensions` directory for Linux users.

Next run `go mod init example.com/extensions` to initialize a Go module which will be
used for tracking dependencies. You can replace `example.com` with the name of your 
company. Your directory structure look similar to the below.
```
etc
└── maverics
    ├── extensions
    │         ├── go.mod
    │         └── go.sum
    └── maverics.yaml
```

After the `go.mod` file has been successfully created, next run 
`go get github.com/strata-io/service-extension` to add this library as a dependency
to your Go project. Adding this library as a dependency to your Go project will enable
the library to be imported in your Service Extensions.

You are now able to import and use the Service Extension library! The code snippet 
below demonstrates how the library can be imported.

`/etc/maverics/extensions/auth.go`
```go
package main

import (
	"net/http"

	"github.com/strata-io/service-extension/orchestrator"
)

func IsAuthenticated(api orchestrator.Orchestrator, rw http.ResponseWriter, req *http.Request) bool {
	return false
}
```

### Testing, Vetting and Formatting
Service extensions can be programmatically tested. To run unit tests, run the 
standard `go test ./...` command in the root of your project directory. To learn more
about how to test service extensions, see the Go testing [docs](https://pkg.go.dev/cmd/go#hdr-Test_packages)
or run `go help test`.

Vetting service extensions can help catch common programming errors. To vet your 
extensions, run the `go vet ./...` command in the root of your project directory. To 
learn more about vetting service extensions, see the Go vet 
[docs](https://pkg.go.dev/cmd/vet) or run `go help vet`.

Formatting your code can help make your code more readable and idiomatic. To format 
service extensions, run the `go fmt ./...` command in the root of your project 
directory. To learn more about formatting service extensions, see the Go formatting
[docs](https://pkg.go.dev/cmd/gofmt) or run `go help fmt`.

Most integrated development environments (IDEs) have built-in support for testing,
vetting, and formatting code. See the below list for how to leverage these features 
in popular IDEs:
- [VSCode](https://learn.microsoft.com/en-us/azure/developer/go/configure-visual-studio-code)
- [GoLand](https://www.jetbrains.com/help/go/quick-start-guide-goland.html)

## Documentation
For library documentation, please visit the [Godoc site](https://pkg.go.dev/github.com/strata-io/service-extension).

For Maverics specific documentation, please visit the [product doc site](https://docs.strata.io/).
