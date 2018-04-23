# here-tracking-go
here-tracking-go is a Go client library for accessing the [HERE Tracking API v2](https://developer.here.com/documentation/tracking/api-reference-swagger.html).

Note: the implementation's interfaces aren't final yet – expect changes

## Credentials
To use the sample application, or any application created with the HERE Tracking Client C Library, you need to:

1. Sign up for a developer account with HERE Tracking.
2. In the vendor role, create some devices with device licenses.
3. In the user role, claim those devices.

The device licenses include device IDs and device secrets.

## Building the Sample Applications
You'll need a valid `$GOPATH` and working Go setup.

### Ingest
Install with `go get`:

```
$ go get -u github.com/choefele/here-tracking-go/cmd/ingest
```

The run the `ingest` providing the device ID and secret:

```
$ ingest <device ID> <device secret>
```

 ## Usage
 ```go
import "github.com/choefele/here-tracking-go/pkg/tracking"
```

Construct a new client, then use the various services on the client to access different parts of the API. For example:

```go
client := tracking.NewClient("", "")
h, err := client.Ingestion.Health(context.Background())

fmt.Printf("Health: %v, error: %v\n", h, err)
```

The services of a client divide the API into logical chunks and correspond to the structure of the [HERE Tracking API v2](https://developer.here.com/documentation/tracking/api-reference-swagger.html.
