# Building a Traces Receiver

 If you are reading this tutorial, you probably already have an idea of what are the concepts behind traces and/or distribute traces, but if you don't you can quickly read through it here: [OpenTelemetry - Data Sources - Traces](https://opentelemetry.io/docs/concepts/data-sources/#traces).

 Here is the definition of those concepts according to OpenTelemetry:

>Traces track the progression of a single request, called a trace, as it is handled by services that make up an application. The request may be initiated by a user or an application. Distributed tracing is a form of tracing that traverses process, network and security boundaries.

Although the definition seems very application centric, you can leverage OpenTelemetry traces as a way to represent a request and quickly understand it's duration and the details about every step involved in completing it.

The trace receiver role is to convert the request telemetry from it's original format into the Otel trace format, so the information can be properly processed through the collector's pipelines.

In order to implement a traces receiver you will need the following:

- A `Config` implementation to enable the trace receiver to gather and validate it's configurations within the collector's config.yaml.

- A `ReceiverFactory` implementation so the Collector can properly instantiate the trace receiver component

- A `TracesReceiver` implementation that is responsible to collect the telemetry, convert it to the internal Otel Trace format, and hand the information to the next consumer in the pipeline.

In this tutorial we will create a sample trace receiver called `tracemock` that simulates a pull operation and generates traces as an outcome of that operation. The next sections will guide you through the process of implementing the steps above in order to create the receiver, so let's get started. 


### Setting up your receiver development and testing environment

First use the tutorial from the [builder](../../builder) folder to create a collector instance named `dev-otelcol`, all you need is to copy the [builder-config.yaml](../../builder/builder-config.yaml) file and make the following changes: 

```yaml
dist:
    name: dev-otelcol # the binary name. Optional.
    output_path: ./dev-otelcol # the path to write the output (sources and binary). Optional.
```

&nbsp;

As an outcome you should now have a `dev-otelcol` folder with your collectors development instance ready to go.

In order to properly test your trace receiver, you will need a distributed tracing backend so the collector can send the telemetry to it. We will be using [JAEGER](https://www.jaegertracing.io/docs/1.30/getting-started), if you don't have a `JAEGER` instance running, you can easily stand one using docker with the following command:

```cmd
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.29
```

&nbsp;

Now, create a config.yaml file so you can setup your collector's components. 

```cmd
cd dev-otelcol
touch config.yaml
```

&nbsp;

For now, you just need a basic traces pipeline with the `otlp` receiver, the `jaeger` and `logging` exporters, here is what your `config.yaml` file should look like:

>config.yaml
```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:55680

processors:

exporters:
  logging:
    logLevel: debug
  jaeger:
    endpoint: localhost:14250
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [jaeger, logging]
```

&nbsp;


In order to verify that your initial pipeline is properly setup, you should have the following output after running your `dev-otelcol` command: 

```cmd
dev-otelcol % ./dev-otelcol --config config.yaml
2022-02-01T09:29:36.205-0600    info    service/collector.go:190        Applying configuration...
2022-02-01T09:29:36.206-0600    info    builder/exporters_builder.go:254        Exporter was built.  {"kind": "exporter", "name": "logging"}
2022-02-01T09:29:36.207-0600    info    builder/exporters_builder.go:254        Exporter was built.  {"kind": "exporter", "name": "jaeger"}
2022-02-01T09:29:36.207-0600    info    builder/pipelines_builder.go:222        Pipeline was built.  {"name": "pipeline", "name": "traces"}
2022-02-01T09:29:36.207-0600    info    builder/receivers_builder.go:224        Receiver was built.  {"kind": "receiver", "name": "otlp", "datatype": "traces"}
2022-02-01T09:29:36.207-0600    info    service/service.go:86   Starting extensions...
2022-02-01T09:29:36.207-0600    info    service/service.go:91   Starting exporters...
2022-02-01T09:29:36.207-0600    info    builder/exporters_builder.go:40 Exporter is starting... {"kind": "exporter", "name": "logging"}
2022-02-01T09:29:36.207-0600    info    builder/exporters_builder.go:48 Exporter started.       {"kind": "exporter", "name": "logging"}
2022-02-01T09:29:36.207-0600    info    builder/exporters_builder.go:40 Exporter is starting... {"kind": "exporter", "name": "jaeger"}
2022-02-01T09:29:36.207-0600    info    builder/exporters_builder.go:48 Exporter started.       {"kind": "exporter", "name": "jaeger"}
2022-02-01T09:29:36.207-0600    info    service/service.go:96   Starting processors...
2022-02-01T09:29:36.207-0600    info    builder/pipelines_builder.go:54 Pipeline is starting... {"name": "pipeline", "name": "traces"}
2022-02-01T09:29:36.207-0600    info    jaegerexporter@v0.41.0/exporter.go:186  State of the connection with the Jaeger Collector backend     {"kind": "exporter", "name": "jaeger", "state": "IDLE"}
2022-02-01T09:29:36.207-0600    info    builder/pipelines_builder.go:65 Pipeline is started.    {"name": "pipeline", "name": "traces"}
2022-02-01T09:29:36.208-0600    info    service/service.go:101  Starting receivers...
2022-02-01T09:29:36.208-0600    info    builder/receivers_builder.go:68 Receiver is starting... {"kind": "receiver", "name": "otlp"}
2022-02-01T09:29:36.208-0600    info    otlpreceiver/otlp.go:69 Starting GRPC server on endpoint localhost:55680      {"kind": "receiver", "name": "otlp"}
2022-02-01T09:29:36.212-0600    info    builder/receivers_builder.go:73 Receiver started.       {"kind": "receiver", "name": "otlp"}
2022-02-01T09:29:36.212-0600    info    service/telemetry.go:92 Setting up own telemetry...
2022-02-01T09:29:36.219-0600    info    service/telemetry.go:116        Serving Prometheus metrics   {"address": ":8888", "level": "basic", "service.instance.id": "b9424849-d205-45ed-aa51-768c606f9b12", "service.version": "latest"}
2022-02-01T09:29:36.219-0600    info    service/collector.go:239        Starting dev-otelcol... {"Version": "1.0.0", "NumCPU": 12}
2022-02-01T09:29:36.219-0600    info    service/collector.go:135        Everything is ready. Begin running and processing data.
2022-02-01T09:29:37.208-0600    info    jaegerexporter@v0.41.0/exporter.go:186  State of the connection with the Jaeger Collector backend     {"kind": "exporter", "name": "jaeger", "state": "READY"}
````

&nbsp;

Make sure you see the last line, that will confirm that the jaeger exporter has successfully established a connection to your local jaeger instance. Now that we have our environment ready, let's start writing your receiver's code.

First create a folder called `tracermock` under the dev-otelcol so it will host all of our receivers code

```cmd
cd dev-otelcol
mkdir tracermock
````

&nbsp;

### Reading and Validating your Receiver Settings

In order to be instantiated and participate in pipelines the collector needs to identify your receiver and properly load it's settings from within it's configuration file. 

The `tracemock` receiver will have the following settings:

- `interval`: the time interval (in minutes) between telemetry pull operations 
- `numberOfTraces`: the number os mock traces generated for each interval

Here is what the `tracemock` receiver settings will look like:

```yaml
receivers:
  tracemock: #this line represents the ID of your receiver
    interval: 1
    numberOfTraces: 1
````

&nbsp;

Under the `tracemock` folder, create a file named `config.go` where you will write all the code to support your receiver settings.

&nbsp;

```cmd
cd tracemock
touch config.go
```

&nbsp;

To implement the configuration aspects of a receiver you need create a `Config` struct, so go ahead the add the following code to your `config.go` file:


```go
package tracemock

type Config struct{

}
```

&nbsp;

In order to be able to give your receiver access to it's settings the `Config` struct must:

- embed the [config.ReceiverSettings](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/receiver.go#L32) struct or a struct that extends it.

- Add a field for each of the receiver's settings.

Here is what your config.go file should look like after you implemented the requirements above


>config.go
```go
package tracemock

import (
  "go.opentelemetry.io/collector/config"
)

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
   config.ReceiverSettings `mapstructure:",squash"`
   Interval    int `mapstructure:"interval"`
   NumberOfTraces int `mapstructure:"numberOfTraces"`
}

```

>### Reviewing the code
>
>- I imported the `go.opentelemetry.io/collector/config` package, which is where ReceiverSettings is declared.
>- I embedded the `config.ReceiverSettings` as required by the spec.
>- I added the `Interval` and the `NumberOfTraces` fields so I can properly have access to their values from the config.yaml.

&nbsp;

Now that you have access to the settings, you can provide any kind of validation needed for those values by implementing the `Validate` method according to the [validatable](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/config.go#L154) interface.

In this case, the `interval` value will be optional (we will look at generating default values later) and the the `numberOfTraces` will be a required value. Here is what the config.go looks like after implementing the `Validate` method.

>config.go
```go
package tracemock

import (
  "fmt"
  "go.opentelemetry.io/collector/config"
)

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
   config.ReceiverSettings `mapstructure:",squash"`
   Interval    int `mapstructure:"interval"`
   NumberOfTraces int `mapstructure:"numberOfTraces"`
}


// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	if (cfg.NumberOfTraces >= 1){
	   return fmt.Errorf("numberOfTraces must be greater or equal to 1")
	}
	return nil
 }
```

>### Reviewing the code
>
>- I imported the `fmt` package, so I can properly format print my error messages.
>- I added the `Validate` method to my Config struct where I am checking if the `numberOfTraces` setting value is greater or equal to 1. If that is not true the Collector will generate an error during it's startup process and display the message accordingly.

&nbsp;


If you want to take a closer look at the structs and interfaces involved in the configuration aspects of a receiver component, take a look at the [config/receiver.go](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/receiver.go) file inside the Collector's GitHub project.














