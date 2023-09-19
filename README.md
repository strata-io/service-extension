# Service Extension Development Library
Integrating identity systems requires extreme configurability. Service Extensions 
give administrators the ability to customize the behavior of the Maverics Orchestrator
to suit the particular needs of their integration. This library aids in the 
development of Service Extensions.

## Local Setup
To develop Service Extensions on your local machine, follow the instructions below. 
The steps listed below are only to aid in local development and are not required to run
the Maverics platform. 

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

## Documentation
For library documentation, please visit the [Godoc site](https://pkg.go.dev/github.com/strata-io/service-extension).

For Maverics specific documentation, please visit the [product doc site](https://docs.strata.io/).
