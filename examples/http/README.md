# HTTP clients and requests making in service extension

This example demonstrates how to make HTTP requests in service extension
using a shared HTTP client.

## Setup
Please reference the [maverics.yaml](maverics.yaml) configuration file and update 
the file path to the `serveSE` service extension appropriately. 
It should point to the [http.go](http.go) file.

## Testing
1. Start the Orchestrator with the `maverics.yaml` configuration file.
2. You should see the following log message in the Orchestrator logs:
```
   ID: 1, Phone: 1-770-736-8031 x56442
   ID: 2, ...
```