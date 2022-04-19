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

func (tailtracerRcvr *tailtracerReceiver) Start(ctx context.Context, host component.Host) error {
    tailtracerRcvr.host = host
    ctx = context.Background()
	ctx, tailtracerRcvr.cancel = context.WithCancel(ctx)
 
	interval, _ := time.ParseDuration(tailtracerRcvr.config.Interval)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				tailtracerRcvr.logger.Info("I should start processing traces now!")
				tailtracerRcvr.nextConsumer.ConsumeTraces(ctx, generateTraces(tailtracerRcvr.config.NumberOfTraces))
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (tailtracerRcvr *tailtracerReceiver) Shutdown(ctx context.Context) error {
	tailtracerRcvr.cancel()
	tailtracerRcvr.logger.Info("I am done and ready to shutdown!")
	return nil
}


