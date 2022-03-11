package tracemock

import "go.opentelemetry.io/collector/model/pdata"

type Atm struct{
    ID           int64
	Version      string
	Name         string
	StateID      string
	SerialNumber string
	ISPNetwork   string
}

type Backendsystem struct{
	Version       string
	ProcessName   string
	OSType        string
    OSVersion     string
	CloudProvider string
	CloudRegion   string	
	Endpoint      string
}

func generateTraces() pdata.Traces{
	traces := pdata.NewTraces()
	
	return traces
}