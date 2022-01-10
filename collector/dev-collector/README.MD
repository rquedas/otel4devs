# HOW-TO: Build your own Collector instance

In order to build and debug your components, you are going to need your own Collector instance. That will allow you to launch and debug your Collector components directly within your favorite Golang IDE.

The other really cool aspect of approaching the component development this way is that you can use all the cool debugging features from your IDE (stack traces are great teachers!) to understand how the Collector itself interacts with your component code.

Start by creating a folder for your project and initializing a go module in it.

`mkdir otel-collector-dev`

`cd otel-collector-dev`

`go mod init github.com/rquedas/otel-collector-dev`

As any other Go based application, you will need a main.go file as the entry point for your Collector. So go ahead and create a main.go file with the following code.

> main.go

```go
package main

func main() {

}
```

## Instantiating your Collector

Everything you need to instantiate your Collector is within the go.opentelemetry.io/collector/service package.
Within the service package you will find the following function.

```go
func New(set CollectorSettings) (*Collector, error)
```

As you can see, The `New` function is responsible to instantiate a new `Collector` which is going to be configured by the given `CollectorSettings` 

Take a minute to read through the `CollectorSettings` struct declaration and documentation inside [services.go](https://github.com/open-telemetry/opentelemetry-collector/blob/v0.41.0/service/service.go)

&nbsp;
### Defining CollectorSettings

In order to properly initialize a `CollectorSettings` instance you will have to provide valid instances of BuildInfo, Factories and ConfigMapProvider structs.

&nbsp;

[BuildInfo](https://github.com/open-telemetry/opentelemetry-collector/blob/v0.41.0/component/build_info.go) 

It's meant to hold information about the collector instance. You can either pass your own information or use `component.NewDefaultBuildInfo()` to create a default one.

Here is what the main.go file looks like after adding the code to instantiate the BuildInfo

> main.go

```go
package main

import (
	"go.opentelemetry.io/collector/component"
)

func main() {
	info := component.BuildInfo{
		Command:  "otel-collector-dev",
		Description: "Custom Otel Collector for RQ Dev",
		Version:  "1.0.0",
	}
}
```

>### Reviewing the code
>
>- I imported the `go.opentelemetry.io/collector/component` package, which is where BuildInfo is declared.
>- I added a variable called `info` to hold the instance of the BuildInfo with my custom information.

&nbsp;

[Factories](https://github.com/open-telemetry/opentelemetry-collector/blob/v0.41.0/component/factories.go) 

If you take a look at the `go.opentelemetry.io/collector/component` package you will find out that every Collector's component have it's own dedicated Factory that allows the Collector to properly initialize the component.

Factories is a struct that is responsible to hold the references to all the factories of the components we want to have available within our Controller.

Here is what the main.go file looks like after adding the code to initialize Factories.

> main.go

```go
package main

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service/defaultcomponents"
	"log"
)

func main() {
	info := component.BuildInfo{
		Command:  "otel-collector-dev",
		Description: "Custom Otel Collector for RQ Dev",
		Version:  "1.0.0",
	}

    factories, err := defaultcomponents.Components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
    }
}
```

>### Reviewing the code
>
>- I imported the `go.opentelemetry.io/collector/service/defaultcomponents` package, so I can initialize my collector using the function defaultcomponents.Components()
> 
>- I added a variable called `factories` that will be holding all the factories for the components I want my Collector to be initialized with.
>
>	The defaultcomponents.Components() function will load the factories for the following components:
>	* Receivers: Otlp Receiver
>	* Processors: Batch and Memory Limiter Processors
>	* Exporters: Logging, Otlp and OtlpHttp Exporters
>
> &nbsp;
>
>- I imported `log` so I can add logging capabilities to my Collector entrypoint function.
>- I added an `if` statement to properly log an error in case something goes wrong with the defaultcomponents.Components() function.
> &nbsp;

&nbsp;

ConfigMapProvider

 An instance of a `configmapprovider.Provider`, which is meant to help developers access the Collector's component configurations independently of where they are stored.

For this tutorial, we will use a .yaml file to store the configuration of our Collector's components. Create a file named config.yaml within the `otel-collector-dev` folder and add the following content.

> config.yaml

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: localhost:55680

processors:

exporters:
  logging:

service:
  pipelines:
    traces:
      receivers:
        - otlp
      processors: []
      exporters:
        - logging
```


As of v0.41.0 of the Collector codebase, the `configmapprovider` package provides the following functions to generate `Provider` implementations:

```go
func NewFile(fileName string) Provider 
func NewInMemory(buf io.Reader) Provider
func NewProperties(properties []string) Provider 
````

Thanks to these functions, you can easily easily have a Provider to read the Collector's configurations from a .yaml file, key/value properties or in-memory.

Here is what my main.go looks like after adding the code to create the configmapprovider.Provider.

>main.go

```go
package main

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service/defaultcomponents"
	"go.opentelemetry.io/collector/config/configmapprovider"
)

func main(){
	info := component.BuildInfo{
		Command:  "otel-collector-dev",
		Description: "Custom Otel Collector for RQ Dev",
		Version:  "1.0.0",
	}

	factories, err := defaultcomponents.Components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
    }

	configMap := configmapprovider.NewFile("config.yaml")
	
}
```

>### Reviewing the code
>
>- I imported the `go.opentelemetry.io/collector/config/configmapprovider` package, so I can have access to the default functions to instantiate my `configmapprovider.Provider`
> 
>- I added a variable called `configMap` and initialized it using the `configmapprovider.NewFile` function which will load the contents of the `config.yaml` created earlier.
>
> &nbsp;


If you want to understand more about how the configmapprovider.Provider default implementations work, take a look at the [configmapprovider](https://github.com/open-telemetry/opentelemetry-collector/tree/v0.41.0/config/configmapprovider) folder inside the Collector's project.

## Creating the Collector service

We now have all the required elements to instantiate our Collector's build using the default collector components an their respective configurations within the config.yaml file.

Here is what my main.go looks like after adding the code to instantiate our Collector.

>main.go
```go
package main

import (
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmapprovider"
	"go.opentelemetry.io/collector/service/defaultcomponents"
	"go.opentelemetry.io/collector/service"

)

func main(){
	info := component.BuildInfo{
		Command:  "otel-collector-dev",
		Description: "Custom Otel Collector for RQ Dev",
		Version:  "1.0.0",
	}

	factories, err := defaultcomponents.Components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
    }

	configMap := configmapprovider.NewFile("config.yaml")
	
	collectorSettings := service.CollectorSettings{BuildInfo: info, Factories: factories, ConfigMapProvider: configMap}

	collectorApp, err := service.New(collectorSettings)
	
	if err != nil {
		log.Fatal("failed to construct the application: %w", err)
	}
}
````
>### Reviewing the code
>
>- I imported the `go.opentelemetry.io/collector/service` package, so I can have access to the functions to instantiate and run a Collector.
> 
>- I added a variable called `collectorSettings` and initialized it with a service.CollectorSettings instance that finally leverages all the code you wrote so.
> 
>- I added a variable called `collectorApp` and initialized it with the `service.New` function which will instantiate my Collector service with the proper components and configuration established in the give `collectorSettings`.
>
> &nbsp;

If you want to understand how the `service.New` function uses the CollectorSettings to instantiate the Collector, take a closer look at [collector.go](https://github.com/open-telemetry/opentelemetry-collector/blob/v0.41.0/service/collector.go
) inside the Collector project.


## Running the Collector
Once a Collector instance is created, you will use the following method to get your collector running:

```go
func (col *Collector) Run(ctx context.Context) error
```

Here is what my main.go looks like after adding the code to run our Collector instance.

>main.go
```go
package main

import (
	"context"
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmapprovider"
	"go.opentelemetry.io/collector/service"
	"go.opentelemetry.io/collector/service/defaultcomponents"
)

func main(){
	info := component.BuildInfo{
		Command:  "otel-collector-dev",
		Description: "Custom Otel Collector for RQ Dev",
		Version:  "1.0.0",
	}

	factories, err := defaultcomponents.Components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
    }

	configMap := configmapprovider.NewFile("config.yaml")
	
	collectorApp, err := service.New(collectorSettings)
	
	if err != nil {
		log.Fatal("failed to construct the application: %w", err)
	}

	err = collectorApp.Run(context.TODO())
	if err != nil {
		log.Fatal("application run finished with error: %w", err)
	}
}
````
>### Reviewing the code
>
>- I added a line to call `collectorApp.Run` with a `TODO context` and the necessary error handling code below it.
>
> &nbsp;

&nbsp;

If you want to understand how the `Run` method works, take a closer look at [collector.go](https://github.com/open-telemetry/opentelemetry-collector/blob/v0.41.0/service/collector.go
) inside the Collector project.

With this last code change, you are finally ready to run your collector within VSCode! 

Go ahead and select your main.go tab and then select `Run | Run Without Debugging` in the main menu to start your collector.

Here is what your DEBUG CONSOLE will look like if your collector is running properly.

![](images/dev-collector-img1.png?raw=true)

As you can see, the collector initialization generates an output showing what pipelines, exporters, receivers and extensions have been initialized based on the config.yaml we have defined earlier.

If there are any errors they will show up as part of the output too, but if everything went smoothly, the last line of your console should say "Everything is ready. Begin running and processing data.".

You have now a VSCode project that is bootstrapped and ready to support building and debugging components for the collector.



