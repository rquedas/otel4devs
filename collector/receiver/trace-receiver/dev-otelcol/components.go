// Code generated by "go.opentelemetry.io/collector/cmd/builder". DO NOT EDIT.

package main

import (
	"go.opentelemetry.io/collector/component"
	jaegerexporter "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/jaegerexporter"
	loggingexporter "go.opentelemetry.io/collector/exporter/loggingexporter"
	batchprocessor "go.opentelemetry.io/collector/processor/batchprocessor"
	otlpreceiver "go.opentelemetry.io/collector/receiver/otlpreceiver"
	tailtracer "github.com/rquedas/otel4devs/collector/receiver/trace-receiver/tailtracer"
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
		tailtracer.NewFactory(),
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
