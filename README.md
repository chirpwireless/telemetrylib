# TelemetryLib

`telemetrylib` is a Go library designed to streamline the integration of OpenTelemetry into your Go applications. It
provides a unified approach to setting up observability, including metrics, logs, and traces, helping you maintain and
troubleshoot systems effectively.

## Getting started

### Installation

To install telemetrylib, use the following go get command:

```shell
go get gitlab.com/chirpwireless/backend/telemetrylib
```

### Basic Setup

To integrate `telemetrylib` into your application, place the following code at the beginning of your `main` function:

```go
otelShutdown, err := telemetrylib.SetupOTelSDK(ctx)
if err != nil {
    panic("Error initializing Otel SDK " + err.Error())
}
defer func () {
    err = errors.Join(err, otelShutdown(ctx))
}()
```

**Explanation**

- `SetupOTelSDK`: This function initializes the OpenTelemetry SDK, setting up metrics, logs, and trace collectors.
- `otelShutdown`: Ensures that all telemetry data is flushed and resources are cleaned up when your application shuts
  down.

### Gin Integration

If you're using the Gin framework, you can integrate OpenTelemetry by adding the following middleware:

- Install the library:

```shell
go get github.com/Cyprinus12138/otelgin
```

- Add the middleware after creating a new instance of Gin:

```go
root.Use(otelgin.Middleware("<service_name>"))
```

**Explanation**

- `otelgin.Middleware`: This middleware enables traces and metrics for HTTP requests handled by Gin, allowing you to
  monitor incoming and outgoing requests.

### HTTP Client Integration

To instrument external HTTP service calls, follow these steps:

- Install the library:

```shell
go get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
```

- Set up the HTTP client after initializing the OpenTelemetry SDK:

```go
http.DefaultClient = &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
```

**Explanation**

- `otelhttp.NewTransport`: This wraps the default HTTP transport, enabling trace and metrics propagation for HTTP client
  requests.

### gRPC Client Integration

For gRPC client instrumentation, use the following steps:

- Install the library:

```shell
go get go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
```

- Include the interceptor in your gRPC client:

```go
opts = append(opts, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
```

**Explanation**

- `otelgrpc.NewClientHandler`: This handler allows you to trace gRPC client requests, providing insights into outgoing
  gRPC traffic.

### gRPC Server Integration

To instrument a gRPC server, follow these steps:

- Install the library:

```shell
go get go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
```

- Initialize your gRPC server with the appropriate handler:

```go
grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
```

**Explanation**

- `otelgrpc.NewServerHandler`: This handler enables tracing of incoming gRPC requests, allowing you to monitor
  server-side gRPC interactions.

### SQL Database Integration

For SQL database instrumentation, use the following steps:

Install the library:

```shell
go get github.com/XSAM/otelsql
```

Initialize your database instance:

```go
driverName, err := otelsql.Register("postgres", otelsql.WithSpanOptions(otelsql.SpanOptions{
        OmitConnResetSession: true,
    }),
    otelsql.WithSQLCommenter(true),
)
if err != nil {
    return err
}

db, err := sql.Open(driverName, settings.Instance().PostgresUrl)
if err != nil {
    return err
}

dbinst := sqlx.NewDb(db, "postgres")
```

**Explanation**

- `otelsql.Register`: Registers a new SQL driver with tracing capabilities, allowing you to capture SQL query execution
  details.

## Limitations

Currently, Go lacks robust auto-instrumentation capabilities compared to languages like Java or Node.js. The manual
steps outlined above are necessary to achieve comprehensive instrumentation. However, the Go auto-instrumentation
feature is in alpha, and future updates may simplify this process once it becomes stable.

For more information, refer to
the [OpenTelemetry Go Instrumentation GitHub repository](https://github.com/open-telemetry/opentelemetry-go-instrumentation).

