# Building Open Telemetry Receivers

A Receiver is the entry door for any telemetry to be made available to the Otel Collector. It is basically responsible to receive telemetry data on a given format and translate it to the internal Collector's format so it can be processed through a pipeline.

Receivers can either pull telemetry data from a system (normally used when you don't have a way to instrument the system's code) or have data pushed into it (the system can be instrumented and is capable of pushing the telemetry data to an endpoint in the Collector)

Currently, the OpenTelemetry Collector provides apis to create a receiver for each of the supported telemetry datasources. So, if you want to ingest traces into your Otel pipeline, you will build a `TracesReceiver`; for metrics ingestion, you will build a `MetricsReceiver` and last but not least, for logs you will build a `LogsReceiver`.

If you take a look at [receiver.go](https://github.com/open-telemetry/opentelemetry-collector/blob/main/component/receiver.go) you will quickly realize that all the receivers are inherited from the `Component` interface (you can look at [component.go](https://github.com/open-telemetry/opentelemetry-collector/blob/main/component/component.go) for details), which basically enables them to be loaded into a Collector's pipeline.


The Collector's Receiver tutorials are organized as follow:

## [trace-receiver](trace-receiver) 
This folder contains a tutorial that demonstrates how to build a traces receiver to collect telemetry that can meaningfully be translated and analyzed as traces.
