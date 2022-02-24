package tracemock

import (
	"fmt"
	"time"

	"go.opentelemetry.io/collector/config"
)

// Config represents the receiver config settings within the collector's config.yaml
type Config struct {
   config.ReceiverSettings `mapstructure:",squash"`
   Interval       string `mapstructure:"interval"`
   NumberOfTraces int `mapstructure:"number_of_traces"`
}


// Validate checks if the receiver configuration is valid
func (cfg *Config) Validate() error {
	if (cfg.NumberOfTraces < 1){
	   return fmt.Errorf("number_of_traces must be at least 1")
	}

    interval, _ := time.ParseDuration(cfg.Interval)
	if (interval.Minutes() < 1){
		return fmt.Errorf("when defined, the interval has to be set to at least 1 minute")
	 }
 
	return nil
 }