package tracemock

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type tracemockReceiver struct {
    host component.Host
	cancel context.CancelFunc
	logger       *zap.Logger
	nextConsumer consumer.Traces
	config       *Config
}

func (tracemokRcvr *tracemockReceiver) Start(ctx context.Context, host component.Host) error {
    tracemokRcvr.host = host
    ctx = context.Background()
	ctx, tracemokRcvr.cancel = context.WithCancel(ctx)
	tracemokRcvr.logger.Info("I should start processing traces now!")
	return nil
}

func (tracemokRcvr *tracemockReceiver) Shutdown(ctx context.Context) error {
	tracemokRcvr.cancel()
	tracemokRcvr.logger.Info("I am done and ready to shutdown!")
	return nil
}


