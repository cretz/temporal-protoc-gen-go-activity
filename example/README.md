# Greeting Example

This example is the equivalent of the
[samples-go greeting example](https://github.com/temporalio/samples-go/tree/main/greetings) but using a
[proto file](greetingspb/greetings.proto).

## Generation

Assuming all dependencies of `go`, `protoc`, `protoc-gen-go`, and `protoc-gen-go-activity` are on the `PATH`, run:

    protoc --go_out=paths=source_relative:. --go-activity_out=paths=source_relative:. greetingspb/greetings.proto

## Running

    go run .