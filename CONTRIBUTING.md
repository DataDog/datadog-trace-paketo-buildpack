# Contributing

## Dependencies
For this buildpack to work locally, the following needs to be installed:

- go
- [pack](https://buildpacks.io/docs/tools/pack/)
- docker

## Running locally
To package this repository into a cloud native buildpack,
1. Clone this repository onto your local machine
2. Package this repository into a buildpack (of `.tgz` format) by calling the following in the root directory: 
   ```
    ./jam-darwin pack \
    --buildpack ./buildpack.toml \
    --version 1.0.0 \
    --output ./datadog-trace.tgz
   ```
   The [`jam-darwin` executable](https://github.com/paketo-buildpacks/packit) is a binary provided by Paketo to package its buildpacks.


To build a java application into an OCI image, we will use cloud native buildpack's CLI, `pack`. For `pack` to know that the user wants a Datadog trace agent to be contributed an OCI image, a [binding](https://paketo.io/docs/buildpacks/configuration/#bindings) needs to be provided along with the earlier package build pack. 

3. Create a binding using the following command
   ```
   mkdir binding
   echo 'DatadogTrace' > binding/type
   ```
   Within the directory `./binding`, additional files can be added to configure the agent (see [README](./README.md) for more details).
   
4. Build the image
   ```
    pack build <CONTAINER_NAME> \
    --buildpack paketo-buildpacks/java \
    --buildpack <PATH>/<TO>/datadog-trace.tgz \
    --volume "<PATH>/<TO>/binding:/platform/bindings/DatadogTrace"
   ```

When the running the docker image using,
```
docker run <CONTAINER_NAME>
```
the container shell should display logs from the datadog tracer, indicating that the agent is correctly attached.

## Testing
Run `go test ./...` to test run all unit tests for each of the directories.