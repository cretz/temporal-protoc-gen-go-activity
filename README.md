# Temporal Activity Protobuf Generator Proof of Concept

This is a protoc plugin for generating easy to use code for calling and implementing activities.

**This is not an official Temporal product**

This generator converts protobuf services to easy to use Temporal activity code. Specifically a struct is generated for
calling activities and an interface is generated for implementing them.

See [the example](example) to see how it can be used.

## Usage

This is similar to the [gRPC quick start](https://grpc.io/docs/languages/go/quickstart/)

This requires [Go](https://golang.org/) installed, and [protoc](https://developers.google.com/protocol-buffers)
installed and on the `PATH`.

If not already installed, the protobuf Go code generator should be installed with:

    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

Similarly, this code generator should also be installed with:

    go install github.com/cretz/temporal-protoc-gen-go-activity/protoc-gen-go-activity@latest

The Go `bin` also needs to be on the `PATH` for those two protoc-gen binaries to be seen by `protoc`. For example, on
Linux:

    export PATH="$PATH:$(go env GOPATH)/bin"

Now `protoc` can be used with both `--go_out` and `--go-activity_out` and Temporal code will be written. See
[the example](example) for a sample of its use.

## Additional Notes

* This is just a proof of concept showing how to use protobuf services. It is not a robust code generator.
* By default, the Temporal Go SDK will use protobuf JSON serialization for activity request/response. While this is much
  more human friendly to external readers, it is less resilient to field name changes than traditional binary protobuf
  serialization which only cares about field numbers. If binary serialization is preferred, the data converter can be
  set to `converter.NewCompositeDataConverter(converter.NewProtoPayloadConverter())` (but this loses fallback converters,
  so instead one may just want to reorder what's in
  [the default](https://github.com/temporalio/sdk-go/blob/master/converter/default_data_converter.go)).
* While the activity struct has synchronous invocations, futures can still be used by still using `ExecuteActivity`
  passing the function reference (e.g. `workflow.ExecuteActivity(ctx, MyService.MyRPC, req)`).