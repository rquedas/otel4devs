package tailtracer

import (
	"math/rand"
	"strings"
	"time"

	"encoding/binary"

	crand "crypto/rand"

	"github.com/google/uuid"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/model/semconv/v1.9.0"
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

func generateTraces(numberOfTraces int) ptrace.Traces{
	traces := ptrace.NewTraces()

	for i := 0; i <= numberOfTraces; i++{
		newAtm := generateAtm()
		newBackendSystem := generateBackendSystem()
	
		resourceSpan := traces.ResourceSpans().AppendEmpty()
		atmResource := resourceSpan.Resource()
		fillResourceWithAtm(&atmResource, newAtm)
	
		atmInstScope := appendAtmSystemInstrScopeSpans(&resourceSpan)
	
		resourceSpan = traces.ResourceSpans().AppendEmpty()
		backendResource := resourceSpan.Resource()
		fillResourceWithBackendSystem(&backendResource, newBackendSystem)
	
		backendInstScope := appendAtmSystemInstrScopeSpans(&resourceSpan)
		
	
		appendTraceSpans(&newBackendSystem, &backendInstScope, &atmInstScope)
	} 

	return traces
}

func fillResourceWithAtm(resource *pcommon.Resource, atm Atm){
   atmAttrs := resource.Attributes()
   atmAttrs.InsertInt("atm.id", atm.ID)
   atmAttrs.InsertString("atm.stateid", atm.StateID)
   atmAttrs.InsertString("atm.ispnetwork", atm.ISPNetwork)
   atmAttrs.InsertString("atm.serialnumber", atm.SerialNumber)
   atmAttrs.InsertString(conventions.AttributeServiceName, atm.Name)
   atmAttrs.InsertString(conventions.AttributeServiceVersion, atm.Version)

}

func fillResourceWithBackendSystem(resource *pcommon.Resource, backend BackendSystem){
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

 func appendAtmSystemInstrScopeSpans(resourceSpans *ptrace.ResourceSpans) (ptrace.ScopeSpans){
	scopeSpans := resourceSpans.ScopeSpans().AppendEmpty()
	scopeSpans.Scope().SetName("atm-sytem")
	scopeSpans.Scope().SetVersion("v1.0")
	return scopeSpans
}


func appendTraceSpans(backend *BackendSystem, backendScopeSpans *ptrace.ScopeSpans, atmScopeSpans *ptrace.ScopeSpans){
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
    atmDuration, _ := time.ParseDuration("4s")
    atmSpanFinishTime := atmSpanStartTime.Add(atmDuration)


	atmSpan := atmScopeSpans.Spans().AppendEmpty()
	atmSpan.SetTraceID(traceId)
	atmSpan.SetSpanID(atmSpanId)
	atmSpan.SetName(atmOperationName)
	atmSpan.SetKind(ptrace.SpanKindClient)
	atmSpan.Status().SetCode(ptrace.StatusCodeOk)
	atmSpan.SetStartTimestamp(pcommon.NewTimestampFromTime(atmSpanStartTime))
	atmSpan.SetEndTimestamp(pcommon.NewTimestampFromTime(atmSpanFinishTime))


	backendSpanId := NewSpanID()

	backendDuration, _ := time.ParseDuration("2s")
    backendSpanStartTime := atmSpanStartTime.Add(backendDuration)


	backendSpan := backendScopeSpans.Spans().AppendEmpty()
	backendSpan.SetTraceID(atmSpan.TraceID())
	backendSpan.SetSpanID(backendSpanId)
	backendSpan.SetParentSpanID(atmSpan.SpanID())
	backendSpan.SetName(backend.Endpoint)
	backendSpan.SetKind(ptrace.SpanKindServer)
	backendSpan.Status().SetCode(ptrace.StatusCodeOk)
	backendSpan.SetStartTimestamp(pcommon.NewTimestampFromTime(backendSpanStartTime))
	backendSpan.SetEndTimestamp(atmSpan.EndTimestamp())

}

func NewTraceID() pcommon.TraceID{
	return pcommon.NewTraceID(uuid.New())
}

func NewSpanID() pcommon.SpanID {
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	randSource := rand.New(rand.NewSource(rngSeed))

	var sid [8]byte
	randSource.Read(sid[:])
    spanID := pcommon.NewSpanID(sid)

	return spanID
}

