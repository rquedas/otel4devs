package tracemock

import (
	"time"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

const (
	typeStr = "tracemock"
	defaultInterval = 1 * time.Minute
)

func createDefaultConfig() config.Receiver {
	return &Config{
		ReceiverSettings:   config.NewReceiverSettings(config.NewComponentID(typeStr)),
		Interval: defaultInterval,
	}
}

// NewFactory creates a factory for tracemock receiver.
func NewFactory() component.ReceiverFactory {

}