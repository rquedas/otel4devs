package tailtracer

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/component/componenterror"
)

const (
	typeStr = "tailtracer"
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

	traceRcvr := &tailtracerReceiver{
		logger:       logger,
		nextConsumer: consumer,
		config:       tracemockCfg,
	}
	
	return traceRcvr, nil

}

// NewFactory creates a factory for tracemock receiver.
func NewFactory() component.ReceiverFactory {
	return component.NewReceiverFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesReceiver(createTracesReceiver))
}