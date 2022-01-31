# Building a Traces Receiver

 If you are reading this tutorial, you probably already have an idea of what are the concepts behind traces and/or distribute traces, but if you don't you can quickly read through it here: [OpenTelemetry - Data Sources - Traces](https://opentelemetry.io/docs/concepts/data-sources/#traces).

 Here is the definition of those concepts according to OpenTelemetry:

>Traces track the progression of a single request, called a trace, as it is handled by services that make up an application. The request may be initiated by a user or an application. Distributed tracing is a form of tracing that traverses process, network and security boundaries.

Although the definition seems very application centric, you can leverage OpenTelemetry traces as a way to represent a request and quickly understand it's duration and the details about every step involved in completing it.

The trace receiver role is to convert the request telemetry from it's original format into the Otel trace format, so the information can properly processed through the collector's pipelines.

In order to implement a traces receiver you will need the following:

- A `Config` implementation to enable the trace receiver to gather and validate it's configurations within the collector's config.yaml.

- A `ReceiverFactory` implementation so the Collector can properly instantiate the trace receiver component

- A `TracesReceiver` implementation that is responsible to collect the telemetry, convert it to the internal Otel Trace format, and hand the information to the next consumer in the pipeline.

The next sections will guide you through the process of implementing the steps above, so let's get started.



