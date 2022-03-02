package tracemock

import (
	"context"
	"go.opentelemetry.io/collector/component"
)
type tracemock struct {
    host component.Host
	cancel context.CancelFunc
}

func (tracemokRcvr *tracemock) Start(ctx context.Context, host component.Host) error {
    tracemokRcvr.host = host
    ctx = context.Background()
	ctx, tracemokRcvr.cancel = context.WithCancel(ctx)
	return nil
}

func (tracemokRcvr *tracemock) Shutdown(ctx context.Context) error {
	tracemokRcvr.cancel()
	return nil
}
