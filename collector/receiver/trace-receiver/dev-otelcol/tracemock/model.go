package tracemock

import (
	"go.opentelemetry.io/collector/model/pdata"
	"math/rand"
	"time"
)

type Atm struct{
    ID           int64
	Version      string
	Name         string
	StateID      string
	SerialNumber string
	ISPNetwork   string
}

type BackendSystem struct{
	Version       string
	ProcessName   string
	OSType        string
    OSVersion     string
	CloudProvider string
	CloudRegion   string	
	Endpoint      string
}

func generateAtm() Atm{
	i := getRandomNumber(1, 2)
    var newAtm Atm

	switch i {
		case 1:
			newAtm = Atm{
				ID: 111,
				Name: "ATM-111-IL",
				SerialNumber: "atmxph-2022-111",
				Version: "v1.0",
				ISPNetwork: "comcast-chicago",
				StateID: "IL",
		
			}
		
		case 2:
			newAtm = Atm{
				ID: 222,
				Name: "ATM-222-CA",
				SerialNumber: "atmxph-2022-222",
				Version: "v1.0",
				ISPNetwork: "comcast-sanfrancisco",
				StateID: "CA",
			}
	}

	return newAtm
}

func generateBackendSystem() BackendSystem{
    i := getRandomNumber(1, 3)

	newBackend := BackendSystem{
    	ProcessName: "accounts",
		Version: "v2.5",
	}

	switch i {
		case 1:
		 	newBackend.Endpoint = "accounts/api/v2.5/balance"
		case 2:
		  	newBackend.Endpoint = "accounts/api/v2.5/deposit"
		case 3:
			newBackend.Endpoint = "accounts/api/v2.5/withdrawn"

	}

	return newBackend
}

func getRandomNumber(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	i := (rand.Intn(max - min + 1) + min)
    return i	
} 

func generateTraces() pdata.Traces{
	traces := pdata.NewTraces()
	newAtm := generateAtm()
	newBackendSystem := generateBackendSystem()

	resourceSpan := traces.ResourceSpans().AppendEmpty()
	atmResource := resourceSpan.Resource()
	fillResourceWithAtm(&atmResource, newAtm)

	return traces
}

func fillResourceWithAtm(resource *pdata.Resource, atm Atm){
   atmAttrs := resource.Attributes()
   atmAttrs.InsertInt("atm.id", atm.ID)
   atmAttrs.InsertString("atm.stateid", atm.StateID)
   atmAttrs.InsertString("atm.ispnetwork", atm.ISPNetwork)
   atmAttrs.InsertString("atm.serialnumber", atm.SerialNumber)
}


