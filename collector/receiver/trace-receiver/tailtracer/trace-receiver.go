package tailtracer

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type tailtracerReceiver struct {
    host component.Host
	cancel context.CancelFunc
	logger       *zap.Logger
	nextConsumer consumer.Traces
	config       *Config
}

func (tracemokRcvr *tailtracerReceiver) Start(ctx context.Context, host component.Host) error {
    tracemokRcvr.host = host
    ctx = context.Background()
	ctx, tracemokRcvr.cancel = context.WithCancel(ctx)
 
	interval, _ := time.ParseDuration(tracemokRcvr.config.Interval)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				tracemokRcvr.logger.Info("I should start processing traces now!")
				tracemokRcvr.nextConsumer.ConsumeTraces(ctx, generateTraces(tracemokRcvr.config.NumberOfTraces))
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (tracemokRcvr *tailtracerReceiver) Shutdown(ctx context.Context) error {
	tracemokRcvr.cancel()
	tracemokRcvr.logger.Info("I am done and ready to shutdown!")
	return nil
}


