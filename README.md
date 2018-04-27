# here-tracking-go
here-tracking-go is a Go client library for accessing the [HERE Tracking API v2](https://developer.here.com/documentation/tracking/api-reference-swagger.html).

Note: the implementation's interfaces aren't final yet – expect changes

This libary provides a client for devices to send new location data and a client for programmatically accessing the admin API. The services of each client divide the API into logical chunks and correspond to the structure of the [HERE Tracking API v2](https://developer.here.com/documentation/tracking/api-reference-swagger.html).

## Device Client

### Credentials
To use the device client:

1. Sign up for a [developer account with HERE Tracking](https://app.tracking.here.com).
2. In the vendor role, create some devices with device licenses.
3. In the user role, claim those devices.

The device license includes device ID and device secret.

### Building the Sample Application
You'll need a valid `$GOPATH` and working Go setup. Then install with `go get`:

```
$ go get -u github.com/choefele/here-tracking-go/cmd/ingest
```

Run `ingest` providing the device ID and secret:

```
$ ingest <device ID> <device secret>
```

This will send test location data to the HERE Tracking Service, which you can monitor in the [admin console](https://app.tracking.here.com/).

### Usage
 ```go
import "github.com/choefele/here-tracking-go/pkg/tracking"
```

Construct a new client, then use the various services on the client to access different parts of the API. For example:

```go
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ingest device_id device_secret")
		os.Exit(-1)
	}

	client := tracking.NewDeviceClient(os.Args[1], os.Args[2])
	dr := &tracking.DataRequest{
		Timestamp: tracking.Time{Time: time.Now()},
		Position: &tracking.Position{
			Lat:      52,
			Lng:      13,
			Accuracy: 100,
		},
	}
	err := client.Ingestion.Send(context.Background(), []*tracking.DataRequest{dr})
	fmt.Printf("Send: done, error: %v\n", err)
}
```

## Admin Client

### Building the Sample Application
You'll need a valid `$GOPATH` and working Go setup. Then install with `go get`:

```
$ go get -u github.com/choefele/here-tracking-go/cmd/admin
```

Run `admin` providing your email and password:

```
$ admin <email> <password>
```

This will list your devices.

### Usage
 ```go
import "github.com/choefele/here-tracking-go/pkg/tracking"
```

Construct a new client, then use the various services on the client to access different parts of the API. For example:

```go
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: admin email password")
		os.Exit(-1)
	}

	client := tracking.NewAdminClient(os.Args[1], os.Args[2])
	err := client.User.ListDevices(context.Background())
	fmt.Printf("ListDevices: done, error: %v\n", err)
}
```
