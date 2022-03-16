package tracemock

import (
	"go.opentelemetry.io/collector/model/pdata"
	"math/rand"
	"time"
	conventions "go.opentelemetry.io/collector/model/semconv/v1.8.0"
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
		OSType: "lnx",
		OSVersion: "4.16.10-300.fc28.x86_64",
		CloudProvider: "amzn",
		CloudRegion: "us-east-2",
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

	resourceSpan = traces.ResourceSpans().AppendEmpty()
	backendResource := resourceSpan.Resource()
	fillResourceWithBackendSystem(&backendResource, newBackendSystem)

	return traces
}

func fillResourceWithAtm(resource *pdata.Resource, atm Atm){
   atmAttrs := resource.Attributes()
   atmAttrs.InsertInt("atm.id", atm.ID)
   atmAttrs.InsertString("atm.stateid", atm.StateID)
   atmAttrs.InsertString("atm.ispnetwork", atm.ISPNetwork)
   atmAttrs.InsertString("atm.serialnumber", atm.SerialNumber)
   atmAttrs.InsertString(conventions.AttributeServiceName, atm.Name)
   atmAttrs.InsertString(conventions.AttributeServiceVersion, atm.Version)

}


func fillResourceWithBackendSystem(resource *pdata.Resource, backend BackendSystem){
	backendAttrs := resource.Attributes()
	var osType, cloudProvider string

	switch {
		case backend.CloudProvider == "amzn":
			cloudProvider = conventions.AttributeCloudProviderAWS
		case backend.OSType == "mcrsft":
			cloudProvider = conventions.AttributeCloudProviderAzure
		case backend.OSType == "gogl":
			cloudProvider = conventions.AttributeCloudProviderGCP		
	}

	backendAttrs.InsertString(conventions.AttributeCloudProvider, cloudProvider)
	backendAttrs.InsertString(conventions.AttributeCloudRegion, backend.CloudRegion)
	
	switch {
		case backend.OSType == "lnx":
			osType = conventions.AttributeOSTypeLinux
		case backend.OSType == "wndws":
			osType = conventions.AttributeOSTypeWindows
		case backend.OSType == "slrs":
			osType = conventions.AttributeOSTypeSolaris			
	}
	
	backendAttrs.InsertString(conventions.AttributeOSType, osType)
	backendAttrs.InsertString(conventions.AttributeOSVersion, backend.OSVersion)

	backendAttrs.InsertString(conventions.AttributeServiceName, backend.ProcessName)
	backendAttrs.InsertString(conventions.AttributeServiceVersion, backend.Version)

 }

