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