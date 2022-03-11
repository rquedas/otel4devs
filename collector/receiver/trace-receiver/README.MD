# Building a Traces Receiver

 If you are reading this tutorial, you probably already have an idea of what are the concepts behind traces and/or distribute traces, but if you don't you can quickly read through it here: [OpenTelemetry - Data Sources - Traces](https://opentelemetry.io/docs/concepts/data-sources/#traces).

 Here is the definition of those concepts according to OpenTelemetry:

>Traces track the progression of a single request, called a trace, as it is handled by services that make up an application. The request may be initiated by a user or an application. Distributed tracing is a form of tracing that traverses process, network and security boundaries.

Although the definition seems very application centric, you can leverage OpenTelemetry traces as a way to represent a request and quickly understand it's duration and the details about every step involved in completing it.

The trace receiver's role is to convert the request telemetry from it's original format into the Otel trace format, so the information can be properly processed through the collector's pipelines.

In order to implement a traces receiver you will need the following:

- A `Config` implementation to enable the trace receiver to gather and validate it's configurations within the collector's config.yaml.

- A `ReceiverFactory` implementation so the Collector can properly instantiate the trace receiver component

- A `TracesReceiver` implementation that is responsible to collect the telemetry, convert it to the internal Otel Trace format, and hand the information to the next consumer in the pipeline.

In this tutorial we will create a sample trace receiver called `tracemock` that simulates a pull operation and generates traces as an outcome of that operation. The next sections will guide you through the process of implementing the steps above in order to create the receiver, so let's get started. 


## Setting up your receiver development and testing environment

First use the tutorial from the [builder](../../builder) folder to create a collector instance named `dev-otelcol`, all you need is to copy the [builder-config.yaml](../../builder/builder-config.yaml) file and make the following changes: 

```yaml
dist:
    module: dev-otelcol # the module name for the new distribution, following Go mod conventions. Optional, but recommended.
    name: dev-otelcol # the binary name. Optional.
    output_path: ./dev-otelcol # the path to write the output (sources and binary). Optional.
```

&nbsp;

As an outcome you should now have a [dev-otelcol](dev-otelcol) folder with your collectors development instance ready to go.

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
```

&nbsp;

Make sure you see the last line, that will confirm that the jaeger exporter has successfully established a connection to your local jaeger instance. Now that we have our environment ready, let's start writing your receiver's code.

First create a folder called `tracemock` under the dev-otelcol so it will host all of our receivers code

```cmd
cd dev-otelcol
mkdir tracemock
```

&nbsp;


## Reading and Validating your Receiver Settings

In order to be instantiated and participate in pipelines the collector needs to identify your receiver and properly load it's settings from within it's configuration file. 

The `tracemock` receiver will have the following settings:

- `interval`: a string representing the time interval (in minutes) between telemetry pull operations 
- `number_of_traces`: the number os mock traces generated for each interval

Here is what the `tracemock` receiver settings will look like:

```yaml
receivers:
  tracemock: #this line represents the ID of your receiver
    interval: 1m
    number_of_traces: 1
```

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
   Interval    string `mapstructure:"interval"`
   NumberOfTraces int `mapstructure:"number_of_traces"`
}

```

>### Reviewing the code
>
>- I imported the `go.opentelemetry.io/collector/config` package, which is where ReceiverSettings is declared.
>- I embedded the `config.ReceiverSettings` as required by the spec.
>- I added the `Interval` and the `NumberOfTraces` fields so I can properly have access to their values from the config.yaml.

&nbsp;

Now that you have access to the settings, you can provide any kind of validation needed for those values by implementing the `Validate` method according to the [validatable](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/config.go#L154) interface.

In this case, the `interval` value will be optional (we will look at generating default values later) but when defined should be at least 1 minute (1m) and the the `number_of_traces` will be a required value. Here is what the config.go looks like after implementing the `Validate` method.

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
   Interval    string `mapstructure:"interval"`
   NumberOfTraces int `mapstructure:"number_of_traces"`
}


// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
    interval, _ := time.ParseDuration(cfg.Interval)
	if (interval.Minutes() < 1){
	   return fmt.Errorf("when defined, the interval has to be set to at least 1 minute (1m)")
	}

	if (cfg.NumberOfTraces < 1){
	   return fmt.Errorf("number_of_traces must be greater or equal to 1")
	}
	return nil
 }
```

>### Reviewing the code
>
>- I imported the `fmt` package, so I can properly format print my error messages.
>- I added the `Validate` method to my Config struct where I am checking if the `interval` setting value is at least 1 minute (1m) and if the `number_of_traces` setting value is greater or equal to 1. If that is not true the Collector will generate an error during it's startup process and display the message accordingly.

&nbsp;


If you want to take a closer look at the structs and interfaces involved in the configuration aspects of a receiver component, take a look at the [config/receiver.go](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/receiver.go) file inside the Collector's GitHub project.

&nbsp;

## Enabling the Collector to instantiate your receiver.

At the beginning of this tutorial, you created your `dev-otelcol` instance, which is bootstrapped with the following components:

* Receivers: Otlp Receiver
* Processors: Batch Processor
* Exporters: Logging and Jaeger Exporters


Go ahead and open the [components.go](dev-otelcol/components.go) file under the [dev-otelcol](dev-otelcol) folder, and let's take a look at the `components()` function. 


```go
func components() (component.Factories, error) {
	var err error
	factories := component.Factories{}

	factories.Extensions, err = component.MakeExtensionFactoryMap(
	)
	if err != nil {
		return component.Factories{}, err
	}

	factories.Receivers, err = component.MakeReceiverFactoryMap(
		otlpreceiver.NewFactory(),
	)
	if err != nil {
		return component.Factories{}, err
	}

	factories.Exporters, err = component.MakeExporterFactoryMap(
		jaegerexporter.NewFactory(),
		loggingexporter.NewFactory(),
	)
	if err != nil {
		return component.Factories{}, err
	}

	factories.Processors, err = component.MakeProcessorFactoryMap(
		batchprocessor.NewFactory(),
	)
	if err != nil {
		return component.Factories{}, err
	}

	return factories, nil
}
```

As you can see, the `components()` function is responsible to provide the Collector the factories for all it's components which is represented by a variable called `factories` of type `component.Factories` (here is the declaration of the [component.Factories](https://github.com/open-telemetry/opentelemetry-collector/blob/main/component/factories.go#L25) struct), which will then be used to instantiate the components that are configured and consumed by the collector's pipelines.

Notice that `factories.Receivers` is the field holding a map to all the receiver factories (instances of `ReceiverFactory`), and it currently has the `otlpreceiver` factory only which is instantiated through the `otlpreceiver.NewFactory()` function call.

The `tracemock` receiver has to provide a `ReceiverFactory` implementation, and although you will find a `ReceiverFactory` interface (you can find it's definition in the [component/receiver.go](https://github.com/open-telemetry/opentelemetry-collector/blob/main/component/receiver.go#L105) file within the Collector's project ), the right way to provide the implementation is by using the functions available within the `go.opentelemetry.io/collector/receiver/receiverhelper` package.

&nbsp;

&nbsp;

### Implementing your ReceiverFactory
&nbsp;

Start by creating a file named factory.go within the [tracemock](dev-otelcol/tracemock) folder

```cmd
cd tracemock
touch factory.go
```

&nbsp;

Now let's follow the convention and add a function named `NewFactory()` that will be responsible to instantiate the `tracemock` factory. Go ahead the add the following code to your `factory.go` file:

```go
package tracemock

import (
	"go.opentelemetry.io/collector/component"
)

// NewFactory creates a factory for tracemock receiver.
func NewFactory() component.ReceiverFactory {

}
```

&nbsp;

In order to instantiate your `tracemock` receiver factory, you will use the  following function from the `receiverhelper` package:

```go
func NewFactory(cfgType config.Type, createDefaultConfig CreateDefaultConfig, options ...FactoryOption) component.ReceiverFactory
```

&nbsp;

The `receiverhelper.NewFactory()` instantiate and returns a `component.ReceiverFactory` and it requires the following parameters:

- `config.Type`: A config.Type instance representing a unique identifier for your receiver across all collector's components.

- `CreateDefaultConfig`: A reference to a function that returns the config.Receiver instance for your receiver.

- `... FactoryOption`: The slice of FactoryOptions that will determine what type of MELT datasource your receiver is capable of processing.

Let's now implement the code to support all the parameters required by `receiverhelper.NewFactory()`

&nbsp;

&nbsp;

### Identifying and Providing default settings for the receiver
&nbsp;

If you take a look at the definition of [config.Type](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/config.go#L151), you will see that it's just a string. So all we need to do is to provide a string constant representing the unique identifier for our receiver.

Previously, we said that the `interval` setting for our `tracemock` receiver would be optional, in that case you will need to provide a default value for it so it can be used as part of the default settings. 

Go ahead and add the following code to your `factory.go` file:

```go
const (
	typeStr = "tracemock"
	defaultInterval = 1 * time.Minute
)
```

&nbsp;

As for default settings, you just need to add a function that returns a config.Receiver holding the default configurations for the `tracemock` receiver.

To accomplish that, go ahead and add the following code to your `factory.go` file:

```go
func createDefaultConfig() config.Receiver {
	return &Config{
		ReceiverSettings:   config.NewReceiverSettings(config.NewComponentID(typeStr)),
		Interval: defaultInterval,
	}
}
```

&nbsp;

After these two changes you will notice a few imports are missing, so here is what your `factory.go` file should look like with the proper imports:

>factory.go
```go
package tracemock

import (
	"time"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

const (
	typeStr = "tracemock"
	defaultInterval = 1 * time.Minute
)

func createDefaultConfig() config.Receiver {
	return &Config{
		ReceiverSettings:   config.NewReceiverSettings(config.NewComponentID(typeStr)),
		Interval: defaultInterval,
	}
}

// NewFactory creates a factory for tracemock receiver.
func NewFactory() component.ReceiverFactory {

}
```

>### Reviewing the code
>
>- Importing the `time` package in order to support the time.Duration type for the defaultInterval
>- Importing the `go.opentelemetry.io/collector/config` package, which is where the Receiver interface and the NewReceiverSettings() and NewComponentID() functions are declared.
>- Added a string constant called `typeStr` to represent the unique identifier (component ID) of the receiver and assigned `tracemock` as it's value. This id is going to be used to fetch the receiver settings from the collector's config.
>- Added a `time.Duration` constant called `defaultInterval` to represent the default value for our receiver's `Interval` setting. We will be setting the default value for 1 minute hence the assignment of `1 * time.Minute` as it's value.
>- Added a function called `createDefaultConfig` which is responsible to return a config.Receiver implementation, which in this case is going to be an instance of our `tracemock.Config` struct.   
>   - The `tracemock.Config.ReceiverSettings` field was initialized using the `config.NewReceiverSettings` function which returns a `config.ReceiverSettings` instance based on a given `config.ComponentID`.  
>   - To provide the proper `config.ComponentID`, we used the function `config.NewComponentID` which returns a `config.ComponentID` for the given `config.Type` which in our case is represented by the variable `typeStr`
>- The `tracemock.Config.Interval` field was initialized with the `defaultInterval` constant.

&nbsp;

If you want to take a closer look at `ReceiverSettings` struct and `NewReceiverSettings` function within the [config/receiver.go](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/receiver.go) file inside the Collector's GitHub project.

If you want to take a closer look at `ComponentID` struct and `NewComponentID` function within the [config/identifiable.go](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/identifiable.go) file inside the Collector's GitHub project.

&nbsp;

&nbsp;

### Enabling the factory to describe the receiver as capable of processing traces
&nbsp;

The same receiver component can process traces, metrics and logs, and the receiver's factory is responsible to describe those capabilities. 

Given that Traces are the subject of the tutorial, that's the only datasource we will enable the `tracemock` receiver to work with. The `receiverhelper` package provides the following function and type to help the factory describe the trace processing capabilities:

```go
func WithTraces(createTracesReceiver CreateTracesReceiver) FactoryOption
type CreateTracesReceiver func(context.Context, component.ReceiverCreateSettings, config.Receiver, consumer.Traces) (component.TracesReceiver, error)
```

&nbsp;

The `receiverhelper.WithTraces()` instantiate and returns a `receiverhelper.FactoryOption` and it requires the following parameters:
- `CreateTracesReceiver`: A reference to a function that matches the `receiverhelper.CreateTracesReceiver` type 

&nbsp;

The `receiverhelper.CreateTracesReceiver` type is a pointer to a function that is responsible to instantiate and return a `component.TraceReceiver` instance and it requires the following parameters:
- `context.Context`: the reference to the collector's `context.Context` so your trace receiver can properly manage it's execution context.
- `component.ReceiverCreateSettings`: the reference to some of the collector's settings under which your receiver is created. 
- `config.Receiver`: the reference for the receiver config settings passed by the collector to the factory so it can properly read it's settings from the collector config.
- `consumer.Traces`: the reference to the `consumer.Traces` responsible to pass the traces generated by your receiver to the next consumer/step (most likely a processor or exporter) configured in the traces pipeline.

&nbsp;

Start by adding the bootstrap code to properly implement the `receiverhelper.CreateTracesReceiver` function pointer. Go ahead and add the following code to your `factory.go` file:

```go
func createTracesReceiver(_ context.Context, params component.ReceiverCreateSettings, baseCfg config.Receiver, consumer consumer.Traces) (component.TracesReceiver, error) {
  return nil,nil
}
```
&nbsp;

You now have all the necessary components to successfully instantiate your receiver factory using the `receiverhelper.NewFactory` function. Go ahead and and update your `NewFactory()` function in your `factory.go` file as follow:

```go
// NewFactory creates a factory for tracemock receiver.
func NewFactory() component.ReceiverFactory {
	return receiverhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		receiverhelper.WithTraces(createTracesReceiver))
}
```
&nbsp;

After these two changes you will notice a few imports are missing, so here is what your `factory.go` file should look like with the proper imports:

> factory.go
```go
package tracemock

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

const (
	typeStr = "tracemock"
	defaultInterval = 1 * time.Minute
)

func createDefaultConfig() config.Receiver {
	return &Config{
		ReceiverSettings:   config.NewReceiverSettings(config.NewComponentID(typeStr)),
		Interval: defaultInterval,
	}
}

func createTracesReceiver(_ context.Context, params component.ReceiverCreateSettings, baseCfg config.Receiver, consumer consumer.Traces) (component.TracesReceiver, error) {
  return nil,nil
}

// NewFactory creates a factory for tracemock receiver.
func NewFactory() component.ReceiverFactory {
	return receiverhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		receiverhelper.WithTraces(createTracesReceiver))
}
```
>### Reviewing the code
>
>- Importing the `context` package in order to support the `context.Context` type referenced in the `createTracesReceiver` function
>- Importing the `go.opentelemetry.io/collector/consumer` package in order to support the `consumer.Traces` type referenced in the `createTracesReceiver` function
>- Importing the `go.opentelemetry.io/collector/receiver/receiverhelper` package which is where the NewFactory() and the WithTrace() functions are declared.
>- Updated the `NewFactory()` function so it returns the `component.ReceiverFactory` generated by the `receiverhelper.NewFactory()` call with the required parameters. The generated receiver factory will be capable of processing traces through the call to `receiverhelper.WithTraces(createTracesReceiver)`

&nbsp;

At this point, you have the `tracemock` factory and config code needed for the Collector to validate the `tracemock` receiver settings if they are defined within the `config.yaml`. You just need to add it to the Collector's initialization process.

&nbsp;

&nbsp;

### Adding the receiver factory to the Collector's initialization
&nbsp;

As explained before, all the Collector components are instantiated by the `components()` function within the `components.go` file.

The `tracemock` receiver factory instance has to be added to the `factories` map so the Collector can load it properly as part of it's initialization process.

Here is what the `components.go` file looks like after making the changes to support that:

>components.go
```go
// Code generated by "go.opentelemetry.io/collector/cmd/builder". DO NOT EDIT.

package main

import (
	"go.opentelemetry.io/collector/component"
	jaegerexporter "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jaegerexporter"
	loggingexporter "go.opentelemetry.io/collector/exporter/loggingexporter"
	batchprocessor "go.opentelemetry.io/collector/processor/batchprocessor"
	otlpreceiver "go.opentelemetry.io/collector/receiver/otlpreceiver"
	"dev-otelcol/tracemock"
)

func components() (component.Factories, error) {
	var err error
	factories := component.Factories{}

	factories.Extensions, err = component.MakeExtensionFactoryMap(
	)
	if err != nil {
		return component.Factories{}, err
	}

	factories.Receivers, err = component.MakeReceiverFactoryMap(
		otlpreceiver.NewFactory(),
		tracemock.NewFactory(),
	)
	if err != nil {
		return component.Factories{}, err
	}

	factories.Exporters, err = component.MakeExporterFactoryMap(
		jaegerexporter.NewFactory(),
		loggingexporter.NewFactory(),
	)
	if err != nil {
		return component.Factories{}, err
	}

	factories.Processors, err = component.MakeProcessorFactoryMap(
		batchprocessor.NewFactory(),
	)
	if err != nil {
		return component.Factories{}, err
	}

	return factories, nil
}
```

>### Reviewing the code
>
>- Importing the `dev-otelcol\tracemock` package which is where the receiver types and function are.
>- Added a call to `tracemock.NewFactory()` as a parameter of the `component.MakeReceiverFactoryMap()` call so your `tracemock` receiver factory is properly added to the `factories` map.

&nbsp;

We added the `tracemock` receiver settings to the `config.yaml` previously, so here is what the beginning of the output for running your collector with `dev-otelcol` command should look like after building it with the current codebase: 

```cmd
dev-otelcol % ./dev-otelcol --config config.yaml
2022-02-24T12:17:41.454-0600    info    service/collector.go:190        Applying configuration...
2022-02-24T12:17:41.454-0600    info    builder/exporters_builder.go:254        Exporter was built.     {"kind": "exporter", "name": "logging"}
2022-02-24T12:17:41.454-0600    info    builder/exporters_builder.go:254        Exporter was built.     {"kind": "exporter", "name": "jaeger"}
2022-02-24T12:17:41.454-0600    info    builder/pipelines_builder.go:222        Pipeline was built.     {"name": "pipeline", "name": "traces"}
2022-02-24T12:17:41.454-0600    info    builder/receivers_builder.go:111        Ignoring receiver as it is not used by any pipeline      {"kind": "receiver", "name": "tracemock"}
2022-02-24T12:17:41.454-0600    info    builder/receivers_builder.go:224        Receiver was built.     {"kind": "receiver", "name": "otlp", "datatype": "traces"}
2022-02-24T12:17:41.454-0600    info    service/service.go:86   Starting extensions...
2022-02-24T12:17:41.454-0600    info    service/service.go:91   Starting exporters...
```

&nbsp;

Look for the log line for "builder/receivers_builder.go:111" (it's the 4th line from the bottom at the snippet showed here), you can see that the Collector found the settings for the `tracemock` receiver, validated them (the current settings are all correct), but ignores the receiver given that it's not used in any pipeline.

Let's check if the `tracemock` factory is validating the receiver settings correctly, the `interval` setting isn't required, so if you remove it from the `config.yaml` and run the command again you should get the same output.

Now, let's test one of the `tracemock` settings validation rules. Remove the `number_of_traces` setting from the `config.yaml`, and here is what the output for running the collector will look like:

```cmd
dev-otelcol % ./dev-otelcol --config config.yaml
Error: invalid configuration: receiver "tracemock" has invalid configuration: number_of_traces must be at least 1
2022/02/24 13:00:20 collector server run finished with error: invalid configuration: receiver "tracemock" has invalid configuration: number_of_traces must be at least 1
```

&nbsp;

The `tracemock` receiver factory and config requirements are done and the Collector is properly loading your component. You can now move to the core of your receiver, the implementation of the component itself.

&nbsp;

&nbsp;

## Implementing the trace receiver component
&nbsp;

In the previous section, I mentioned the fact that a receiver can process any of the OTel supported datasources (traces, metrics or log), and the Collector's API is designed to help you accomplish that.

All the receiver API's responsible to enable the datasources are currently declared in the [component/receiver.go](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/receiver.go) file within the OTel Collector's project in GitHub, open the file and take a minute to browse through all the interfaces declared in it. 

Notice that `component.TracesReceiver` (and it's siblings `MetricsReceiver` and `LogsReceiver`) at this point in time, doesn't describe any specific methods other than the ones it "inherits" from `component.Receiver` which also doesn't describe any specific methods other than the ones it "inherits" from `component.Component`. 

It might feel weird, but remember, the Collector's API was meant to be extensible, and the components and their datasources might evolve in different ways, so the role of those interfaces exist to help support that. 

So, to create a `component.TracesReceiver`, you just need to implement the following methods described by `component.Component` interface:

```go
Start(ctx context.Context, host Host) error
Shutdown(ctx context.Context) error
```
Both methods actually act as event handlers used by the Collector to communicate with it's components as part of their lifecycle.

The `Start()` represents a signal of the Collector telling the component to start it's processing. As part of the event, the Collector will pass the following information:
- `context.Context`: Most of the time, a receiver will be processing a long-running operation, so the recommendation is to ignore this context and actually create a new one from context.Background().
- `Host`:  The host is meant to enable the receiver to communicate with the collector's host once it's up and running.

The `Shutdown()` represents a signal of the Collector telling the component that the service is getting shutdown and as such the component should stop it's processing and make all the necessary cleanup work required:
- `context.Context`: the context passed by the Collector as part of the shutdown operation.

You will start the implementation by creating a new file called `trace-receiver.go` within your project's `tracemock` folder and add the declaration to a type type called `tracemockReceiver` as follow:

```go
type tracemockReceiver struct{

}
```
 Now that you have the `tracemockReceiver` type you can implement the Start() and Shutdown() methods so the receiver type can be compliant with the `component.TraceReceiver` interface.

Here is what the [tracemock/trace-receiver.go](tracemock/trace-receiver.go) file should look like with the methods implementation:

>trace-receiver.go
```go
package tracemock

import (
	"context"
	"go.opentelemetry.io/collector/component"
)

type tracemockReceiver struct {

}

func (tracemockRcvr *tracemockReceiver) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (tracemockRcvr *tracemockReceiver) Shutdown(context.Context) error {
	return nil
}
```
>### Reviewing the code
>
>- Importing the `context` package which is where the `Context` type and functions are declared
>- Importing the `go.opentelemetry.io/collector/component` package which is where the `Host` type is declared
>- Added a bootstrap implementation of the `Start(ctx context.Context, host component.Host)` method to comply with the `component.TraceReceiver` interface. 
>- Added a bootstrap implementation of the `Shutdown(ctx context.Context)` method to comply with the `component.TraceReceiver` interface. 

&nbsp;

The `Start()` method is passing 2 references (`context.Context` and `component.Host`) that your receiver might need to keep so they can be used as part of it's processing operations. 

The `context.Context` reference should be used for creating a new context to support you receiver processing operations, and in that case you will need to decide the best way to handle context cancellation so you can finalize it properly as part of the component's shutdown within the `Shutdown()` method.

The `component.Host` might be useful during the whole lifecycle of the receiver so you might want to keep that reference within your `tracemockReceiver` type.

Here is what the `tracemockReceiver` type declaration will look like after you include the fields for keeping the references suggested above:  

```go
type tracemockReceiver struct {
  host component.Host
  cancel context.CancelFunc
}
```

&nbsp;

Now you need to update the `Start()` methods so the receiver can properly initialize it's own processing context and have the cancellation function kept in the `cancel` field and also initialize it's `host` field value. You will also update the `Stop()` method in order to finalize the context by calling the `cancel` function.

Here is what the `trace-receiver.go` file look like after making the changes above:

>trace-receiver.go
```go
package tracemock

import (
	"context"
	"go.opentelemetry.io/collector/component"
)

type tracemockReceiver struct {
    host component.Host
	cancel context.CancelFunc
}

func (tracemokRcvr *tracemockReceiver) Start(ctx context.Context, host component.Host) error {
    tracemokRcvr.host = host
    ctx = context.Background()
	ctx, tracemokRcvr.cancel = context.WithCancel(ctx)
	return nil
}

func (tracemokRcvr *tracemockReceiver) Shutdown(ctx context.Context) error {
	tracemokRcvr.cancel()
	return nil
}

```
>### Reviewing the code
>- Updated the `Start()` method by adding the initialization to the `host` field with the `component.Host` reference passed by the collector and the `cancel` function field with the cancellation based on a new context created with `context.Background()` (according the Collector's API documentation suggestions).
>- Updated the `Stop()` method by adding a call to the `cancel()` context cancellation function.

&nbsp;

&nbsp;

### Keeping information passed by the receiver's factory
&nbsp;

Now that you have implemented the `component.TraceReceiver` interface methods, your `tracemock` receiver component is ready to be instantiated and returned by it's factory.

Open the [factory.go](dev-otelcol/tracemock/factory.go#L24) file and navigate to the `createTracesReceiver()` function. Notice that the factory will pass references as part of the `createTracesReceiver()` function parameters that your receiver actually requires to work properly like it's configuration settings (`config.Receiver`), the next `Consumer` in the pipeline that will consume the generated traces (`consumer.Traces`) and the Collector's logger so the `tracemock` receiver can add meaningful events to it (`component.ReceiverCreateSettings`).

Given that all this information will be only  made available to the receiver at the moment it's instantiated by the factory, The `tracemockReceiver` type will need fields to keep that information and use it within other stages of it's lifecycle.

Here is what the `trace-receiver.go` file looks like with the updated `tracemockReceiver` type declaration:

>trace-receiver.go
```go
package tracemock

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
	
)

type tracemockReceiver struct {
    host         component.Host
	cancel       context.CancelFunc
	logger       *zap.Logger
	nextConsumer consumer.Traces
	config       *Config
}

func (tracemokRcvr *tracemockReceiver) Start(ctx context.Context, host component.Host) error {
    tracemokRcvr.host = host
    ctx = context.Background()
	ctx, tracemokRcvr.cancel = context.WithCancel(ctx)

	return nil
}

func (tracemokRcvr *tracemockReceiver) Shutdown(ctx context.Context) error {
	tracemokRcvr.cancel()
	return nil
}

```

>### Reviewing the code
>- Importing the `go.opentelemetry.io/collector/consumer` which is where the pipeline's consumer types and interfaces are declared.
>- Importing the `go.uber.org/zap` package, which is what the Collector uses for it's logging capabilities.
>- Added a `zap.Logger` field named `logger` so we can have access to the collector's logger reference from within the receiver.
>- Added a `consumer.Traces` field named `nextConsumer` so we can push the traces generated by the `tracemock` receiver to the next consumer declared in the collector's pipeline.
>- Added a `Config` field named `config` so we can have access to receiver's configuration settings defined within the collector's config.

&nbsp;

The `tracemockReceiver` type is now ready to be instantiated and keep all  meaningful information passed by it's factory.

Open the [factory.go](dev-otelcol/tracemock/factory.go#L24) file and navigate to the `createTracesReceiver()` function. 

The receiver is only instantiated if it's declared as a component within a pipeline and the factory is responsible to make sure the next consumer (either a processor or exporter) in the pipeline is valid otherwise it should generate an error. 

The collector's api provides a package with helper functions for component error handling called `go.opentelemetry.io/collector/component/componenterror`, if you open the [error.go](https://github.com/open-telemetry/opentelemetry-collector/blob/main/component/componenterror/errors.go) you will find the `ErrNilNextConsumer` error, which is what your receiver factory should throw in case the next consumer has an issue and is passed as nil.

The `createTracesReceiver()` function will need a guard clause to make that validation.

You will also need variables to properly initialize the `config` and the `logger` fields of the `tracemockReceiver` instance.

Here is what the `factory.go` file looks like with the updated `createTracesReceiver()` function:

>factory.go
```go
package tracemock

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
	"go.opentelemetry.io/collector/component/componenterror"
)

const (
	typeStr = "tracemock"
	defaultInterval = "1m"
)

func createDefaultConfig() config.Receiver {
	return &Config{
		ReceiverSettings:   config.NewReceiverSettings(config.NewComponentID(typeStr)),
		Interval: defaultInterval,
	}
}

func createTracesReceiver(_ context.Context, params component.ReceiverCreateSettings, baseCfg config.Receiver, consumer consumer.Traces) (component.TracesReceiver, error) {
	if consumer == nil {
		return nil, componenterror.ErrNilNextConsumer
	}
	
	logger := params.Logger
	tracemockCfg := baseCfg.(*Config)

	traceRcvr := &tracemockReceiver{
		logger:       logger,
		nextConsumer: consumer,
		config:       tracemockCfg,
	}
	
	return traceRcvr, nil

}

// NewFactory creates a factory for tracemock receiver.
func NewFactory() component.ReceiverFactory {
	return receiverhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		receiverhelper.WithTraces(createTracesReceiver))
}
```
>### Reviewing the code
>- Importing the `go.opentelemetry.io/collector/component/componenterror` which is where the component errors are declared.
>- Added a variable called `logger` and initialized it with the collector's logger that is available as a field named `Logger` within the `component.ReceiverCreateSettings` reference.
>- Added a variable called `tracemockCfg` and initialized it by casting the `config.Receiver` reference to the `tracemock` receiver `Config`.
>- Added a variable called `traceRcvr` and initialized it with the `tracemockReceiver` instance using the factory information stored within the variables.
>- Updated the return statement to now include the `traceRcvr` instance.

&nbsp;

With the factory fully implemented and instantiating the trace receiver component you are ready to test the receiver as part of a pipeline. Go ahead and add the `tracemock` receiver to your `traces` pipeline in the `config.yaml` as follow:

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp, tracemock]
      processors: []
      exporters: [jaeger, logging]
```

&nbsp;

Here is what the output for running your collector with `dev-otelcol` command should look like after you updated the `traces` pipeline: 

```cmd
dev-otelcol % ./dev-otelcol --config config.yaml
2022-03-03T11:19:50.779-0600    info    service/collector.go:190        Applying configuration...
2022-03-03T11:19:50.780-0600    info    builder/exporters_builder.go:254        Exporter was built.     {"kind": "exporter", "name": "jaeger"}
2022-03-03T11:19:50.780-0600    info    builder/exporters_builder.go:254        Exporter was built.     {"kind": "exporter", "name": "logging"}
2022-03-03T11:19:50.780-0600    info    builder/pipelines_builder.go:222        Pipeline was built.     {"name": "pipeline", "name": "traces"}
2022-03-03T11:19:50.780-0600    info    builder/receivers_builder.go:224        Receiver was built.     {"kind": "receiver", "name": "otlp", "datatype": "traces"}
2022-03-03T11:19:50.780-0600    info    builder/receivers_builder.go:224        Receiver was built.     {"kind": "receiver", "name": "tracemock", "datatype": "traces"}
2022-03-03T11:19:50.780-0600    info    service/service.go:86   Starting extensions...
2022-03-03T11:19:50.780-0600    info    service/service.go:91   Starting exporters...
2022-03-03T11:19:50.780-0600    info    builder/exporters_builder.go:40 Exporter is starting... {"kind": "exporter", "name": "jaeger"}
2022-03-03T11:19:50.781-0600    info    builder/exporters_builder.go:48 Exporter started.       {"kind": "exporter", "name": "jaeger"}
2022-03-03T11:19:50.781-0600    info    jaegerexporter@v0.41.0/exporter.go:186  State of the connection with the Jaeger Collector backend       {"kind": "exporter", "name": "jaeger", "state": "IDLE"}
2022-03-03T11:19:50.781-0600    info    builder/exporters_builder.go:40 Exporter is starting... {"kind": "exporter", "name": "logging"}
2022-03-03T11:19:50.781-0600    info    builder/exporters_builder.go:48 Exporter started.       {"kind": "exporter", "name": "logging"}
2022-03-03T11:19:50.781-0600    info    service/service.go:96   Starting processors...
2022-03-03T11:19:50.781-0600    info    builder/pipelines_builder.go:54 Pipeline is starting... {"name": "pipeline", "name": "traces"}
2022-03-03T11:19:50.781-0600    info    builder/pipelines_builder.go:65 Pipeline is started.    {"name": "pipeline", "name": "traces"}
2022-03-03T11:19:50.781-0600    info    service/service.go:101  Starting receivers...
2022-03-03T11:19:50.781-0600    info    builder/receivers_builder.go:68 Receiver is starting... {"kind": "receiver", "name": "otlp"}
2022-03-03T11:19:50.781-0600    info    otlpreceiver/otlp.go:69 Starting GRPC server on endpoint localhost:55680        {"kind": "receiver", "name": "otlp"}
2022-03-03T11:19:50.783-0600    info    builder/receivers_builder.go:73 Receiver started.       {"kind": "receiver", "name": "otlp"}
2022-03-03T11:19:50.783-0600    info    builder/receivers_builder.go:68 Receiver is starting... {"kind": "receiver", "name": "tracemock"}
2022-03-03T11:19:50.783-0600    info    tracemock/trace-receiver.go:23  I should start processing traces now!   {"kind": "receiver", "name": "tracemock"}
2022-03-03T11:19:50.783-0600    info    builder/receivers_builder.go:73 Receiver started.       {"kind": "receiver", "name": "tracemock"}
2022-03-03T11:19:50.783-0600    info    service/telemetry.go:92 Setting up own telemetry...
2022-03-03T11:19:50.788-0600    info    service/telemetry.go:116        Serving Prometheus metrics      {"address": ":8888", "level": "basic", "service.instance.id": "0ca4907c-6fda-4fe1-b0e9-b73d789354a4", "service.version": "latest"}
2022-03-03T11:19:50.788-0600    info    service/collector.go:239        Starting dev-otelcol... {"Version": "1.0.0", "NumCPU": 12}
2022-03-03T11:19:50.788-0600    info    service/collector.go:135        Everything is ready. Begin running and processing data.
```

&nbsp;

Look for the log line for "builder/receivers_builder.go:68 Receiver is starting... {"kind": "receiver", "name": "tracemock"}", you can see that the Collector found the settings for the `tracemock` receiver within the `traces` pipeline and is now instantiating it and starting it given that next line present the info line we added to the `Start()` method.

Now, go ahead and press `ctrl+c` in your collector's terminal so you want watch the shutdown process happening. Here is what the output should look like:

```cmd
^C2022-03-03T11:20:14.652-0600  info    service/collector.go:166        Received signal from OS {"signal": "interrupt"}
2022-03-03T11:20:14.652-0600    info    service/collector.go:255        Starting shutdown...
2022-03-03T11:20:14.652-0600    info    service/service.go:121  Stopping receivers...
2022-03-03T11:20:14.653-0600    info    tracemock/trace-receiver.go:29  I am done and ready to shutdown!        {"kind": "receiver", "name": "tracemock"}
2022-03-03T11:20:14.653-0600    info    service/service.go:126  Stopping processors...
2022-03-03T11:20:14.653-0600    info    builder/pipelines_builder.go:73 Pipeline is shutting down...    {"name": "pipeline", "name": "traces"}
2022-03-03T11:20:14.653-0600    info    builder/pipelines_builder.go:77 Pipeline is shutdown.   {"name": "pipeline", "name": "traces"}
2022-03-03T11:20:14.653-0600    info    service/service.go:131  Stopping exporters...
2022-03-03T11:20:14.653-0600    info    service/service.go:136  Stopping extensions...
2022-03-03T11:20:14.653-0600    info    service/collector.go:273        Shutdown complete.
```
&nbsp;

As you can see there is an info log line for the `tracemock` receiver which means the component is responding correctly to the `Shutdown()` event.
In the next section you will learn more about the OTel Trace data model so the `tracemock` receiver can finally generate traces!

## The OTel Trace Data Model

You might be familiar with OTel traces by using the SDKs and instrumenting an application so you can see and evaluate your traces within a distributed tracing backend like Jaeger. 

Here is what a trace looks like in Jaeger:

![Jaeger trace](images/Jaeger.jpeg)

Granted, this is a Jaeger trace, but it was generated by a trace pipeline within the collector, therefore you can use it to learn a few things about the OTel trace data model :

- A trace is made of one or multiple spans structured within a hierarchy to represent dependencies.
- The spans can represent operations within a service and/or across services.

Creating a trace within the trace receiver will be slightly different than the way you would do it with the SDKs, so let's start reviewing the high level concepts.

### Trace entities and attributes

In the OTel world, all telemetry is generated by a `Resource`, here is the definition according to the [OTel spec](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/resource/sdk.mds):

>A `Resource` is an immutable representation of the entity producing telemetry as Attributes. For example, a process producing telemetry that is running in a container on Kubernetes has a Pod name, it is in a namespace and possibly is part of a Deployment which also has a name. All three of these attributes can be included in the `Resource`.

Traces are most commonly used to represent a service request (the Services entity described by Jaeger's model), which are normally implemented as processes running in a compute unit, but OTel's API approach to describe a `Resource` through attributes is flexible enough to represent any entity that you may require like ATMs, IoT sensors, the sky is the limit.

So it's safe to say that for a trace to exist, a `Resource` will have to start it. In this tutorial we will simulate a system that has telemetry that demonstrate ATMs making calls to a backend system, therefore we will have to implement code to create the `Resources` representing the ATM and the backend system.

Go ahead and create a file named `model.go` inside the `dev-otelcol/tracemock` folder

```cmd
cd dev-otelcol/tracemock
touch model.go
```

&nbsp;

Now, within the `model.go` file, add the definition for the `ATM` and the `Backendsytem` types as follow:

>model.go
```go
package tracemock

type Atm struct{
    ID           int64
	Version      string
	Name         string
	StateID      string
	SerialNumber string
	ISPNetwork   string
}

type Backendsystem struct{
	Version       string
	ProcessName   string
	OSType        string
    OSVersion     string
	CloudProvider string
	CloudRegion   string	
	ServiceName   string
	Endpoint      string
}
```

These types are meant to represent the entities as they are within the system been observed and they contain information that would be quite meaningful to be added to the traces as part of the `Resource` definition. 

The collector's API provides a package named `pdata` with all the types, interfaces and helper functions required to work with traces (and all the other OTel datasources) within the collector's pipeline components.

Before you can define a `Resource`, you need to create a trace. Traces are propagated within the collector's pipeline by the `pdata.Traces` type, and you can use the helper function `pdata.NewTraces()` to  instantiate it.

Open the `model.go` and add the following function to it:

```go
func generateTraces() pdata.Traces{
	traces := pdata.NewTraces()

	return traces
}
```

The Collector's API provides a type named `ResourceSpans` to represent both the resources and all the operations that they either originated or received in a trace.
