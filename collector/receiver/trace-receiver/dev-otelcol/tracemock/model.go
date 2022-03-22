package tracemock

import (
	"math/rand"
	"strings"
	"time"

	"encoding/binary"

	crand "crypto/rand"

	"github.com/google/uuid"
	"go.opentelemetry.io/collector/model/pdata"
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
		 	newBackend.Endpoint = "api/v2.5/balance"
		case 2:
		  	newBackend.Endpoint = "api/v2.5/deposit"
		case 3:
			newBackend.Endpoint = "api/v2.5/withdrawn"

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

	atmInstLibray := appendAtmSystemInstrLibSpans(&resourceSpan)

	resourceSpan = traces.ResourceSpans().AppendEmpty()
	backendResource := resourceSpan.Resource()
	fillResourceWithBackendSystem(&backendResource, newBackendSystem)

	backendInstLibrary := appendAtmSystemInstrLibSpans(&resourceSpan)
	

	appendTraceSpans(&newBackendSystem, &backendInstLibrary, &atmInstLibray)

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

 func appendAtmSystemInstrLibSpans(resourceSpans *pdata.ResourceSpans) (pdata.InstrumentationLibrarySpans){
	iLibSpans := resourceSpans.InstrumentationLibrarySpans().AppendEmpty()
	iLibSpans.InstrumentationLibrary().SetName("atm-sytem")
	iLibSpans.InstrumentationLibrary().SetVersion("v1.0")
	return iLibSpans
}


func appendTraceSpans(backend *BackendSystem, backendInstrLbrSpans *pdata.InstrumentationLibrarySpans, atmInstrLbrSpans *pdata.InstrumentationLibrarySpans){
	traceId := NewTraceID()

	var atmOperationName string

	switch {
	case strings.Contains(backend.Endpoint, "balance"):
        atmOperationName = "Check Balance"
	case strings.Contains(backend.Endpoint, "deposit"):
		atmOperationName = "Make Deposit"
	case strings.Contains(backend.Endpoint, "withdraw"):
		atmOperationName = "Fast Cash"
	}

	atmSpanId := NewSpanID()
    atmSpanStartTime := time.Now()
    atmDuration, _ := time.ParseDuration("2s")
    atmSpanFinishTime := atmSpanStartTime.Add(atmDuration)


	atmSpan := atmInstrLbrSpans.Spans().AppendEmpty()
	atmSpan.SetTraceID(traceId)
	atmSpan.SetSpanID(atmSpanId)
	atmSpan.SetName(atmOperationName)
	atmSpan.SetKind(pdata.SpanKindClient)
	atmSpan.Status().SetCode(pdata.StatusCodeOk)
	atmSpan.SetStartTimestamp(pdata.NewTimestampFromTime(atmSpanStartTime))
	atmSpan.SetEndTimestamp(pdata.NewTimestampFromTime(atmSpanFinishTime))


	backendSpanId := NewSpanID()

	backendDuration, _ := time.ParseDuration("1s")
    backendSpanStartTime := atmSpanStartTime.Add(backendDuration)


	backendSpan := backendInstrLbrSpans.Spans().AppendEmpty()
	backendSpan.SetTraceID(atmSpan.TraceID())
	backendSpan.SetSpanID(backendSpanId)
	backendSpan.SetParentSpanID(atmSpan.SpanID())
	backendSpan.SetName(backend.Endpoint)
	backendSpan.SetKind(pdata.SpanKindServer)
	backendSpan.Status().SetCode(pdata.StatusCodeOk)
	backendSpan.SetStartTimestamp(pdata.NewTimestampFromTime(backendSpanStartTime))
	backendSpan.SetEndTimestamp(atmSpan.EndTimestamp())

}

func NewTraceID() pdata.TraceID{
	return pdata.NewTraceID(uuid.New())
}

func NewSpanID() pdata.SpanID {
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	randSource := rand.New(rand.NewSource(rngSeed))

	var sid [8]byte
	randSource.Read(sid[:])
    spanID := pdata.NewSpanID(sid)

	return spanID
}

