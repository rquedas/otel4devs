package main

import (
	"context"
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmapprovider"
	"go.opentelemetry.io/collector/service"
	"go.opentelemetry.io/collector/service/defaultcomponents"
)

func main(){
	info := component.BuildInfo{
		Command:  "otel-collector-dev",
		Description: "Custom Otel Collector for RQ Dev",
		Version:  "1.0.0",
	}

	factories, err := defaultcomponents.Components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
    }

	configMap := configmapprovider.NewFile("config.yaml")
	
	collectorSettings := service.CollectorSettings{BuildInfo: info, Factories: factories, ConfigMapProvider: configMap}

	collectorApp, err := service.New(collectorSettings)
	
	if err != nil {
		log.Fatal("failed to construct the application: %w", err)
	}

	err = collectorApp.Run(context.TODO())
	if err != nil {
		log.Fatal("application run finished with error: %w", err)
	}
}