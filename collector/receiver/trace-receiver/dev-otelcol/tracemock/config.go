package tracemock

import (
  "fmt"
  "go.opentelemetry.io/collector/config"
)

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
   config.ReceiverSettings `mapstructure:",squash"`
   Interval    int `mapstructure:"interval"`
   NumberOfTraces int `mapstructure:"numberOfTraces"`
}


// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	if (cfg.NumberOfTraces <= 1){
	   return fmt.Errorf("numberOfTraces must be at least 1")
	}
	return nil
 }