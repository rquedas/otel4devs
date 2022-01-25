

# Introducing... "The Collector"

The OpenTelemetry Collector is the doorway between any MELT datasource and a potential OpenTelemetry consumer.

You can think about "The Collector" (as it's known by close friends and family in the community  ) as a Data bus that can run multiple data pipelines, where each pipeline is responsible to receive MELT data (and additional non-MELT data) from multiple sources and formats, process it, and send it to one or more backends (OTLP compliant or proprietary).

In order to do that, the Collector relies on 3 basic components:

- Receivers: How to get data into the Collector; these can be push or pull based
- Processors: What to do with received data
- Exporters: Where to send received data; these can be push or pull based

Here are how these 3 components work within a typical pipeline configuration:

- ****Step 1****

    Receiver (pull or push) receives data from an external source and basically generates the equivalent OTEL Datasource representation of the data so it can be consumed by the other components within pipeline

- ****Step 2..N****

    Processors can actually be used as pipeline steps that are capable of looking on  what’s coming from the receivers and do something about or with the data.

- ****Step N+1****

    Exporters seat at the end of a pipeline and are focused in exporting the Otel data pipeline to any external format/backend.


Otel Pipelines are datasource driven (Otel currently supports Traces, Logs and Metrics as a datasource), so when you are building a Collector component, it will be focused in providing capabilities to at least 1 of the datasources.

As a developer looking to expand the Collector capabilities, your first step is to understand what datasource(s) are you going to use to build your innovation and contribute to the OpenTelemetry world. Once you picked the datasource(s) you can then decide what components are going to be required to deliver your new Collector’s capability.

Now that you know the basics, here is how the Collector's Tutorials are organized:

## [dev-collector](dev-collector)

 This folder has a tutorial to help you learn how to build your own Collector instance so you can easily build and debug your components

## [builder](builder)

 This folder has a tutorial to help you create a Collector distribution from scratch using the `OpenTelemetry Collector Builder`

## [receiver](receiver)

 This folder has tutorials to help you create receiver components to help ingest different kinds of telemetry into Collector's pipelines so they can be available in the Otel world

