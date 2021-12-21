package main

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service/defaultcomponents"
	"go.opentelemetry.io/collector/config/configmapprovider"
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
	
}