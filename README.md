<!-- # `gcr.io/paketo-buildpacks/azure-application-insights` -->
# Datadog Trace Agent Paketo Buildpack
The Paketo Datadog Trace Buildpack is a Cloud Native Buildpack that contributes the Data trace agent and configures it to connect to the service running on an OCI image.

## Behaviour
This buildpack will participate if all the following conditions are met:

* A [binding](https://paketo.io/docs/buildpacks/configuration/#bindings) exists with `type` of `DatadogTrace`

The buildpack will do the following for Java applications:

* Contributes a Java agent to a layer and configures `JAVA_TOOL_OPTIONS` to use it

## Configuring the Agent
The agent can be configured by adding the following files in the binding's directory:

| File Type        | Behaviour           |
| ------------- |:-------------:| 
| `.yaml`      |  Contributes the file into the image and assigns the env variable `DD_JMXFETCH_CONFIG` to the file's path in the image. | 
| `.properties`      | Contributes the file into the image and assigns the env variable `DD_TRACE_CONFIG` to the file's path in the image.      | 

Currently only one of each configuration file can be provided in the binding.

## Contributing
To contribute or use the buildpack locally, see [contributing](./CONTRIBUTING.md)

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
