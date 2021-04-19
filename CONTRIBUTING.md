# Contributing

## Dependencies
For this buildpack to work locally, the following needs to be installed:

- go
- [pack](https://buildpacks.io/docs/tools/pack/)
- docker

## Running locally
### Packaging buildpack into an OCI image
1. Bundle this repository into `.tgz` file: 
   ```
    ./jam-darwin pack \
    --buildpack ./buildpack.toml \
    --version 1.0.0 \
    --output ./datadog-trace.tgz
   ```
   The `jam-darwin` executable is provided by Paketo to package its buildpacks (latest versions in the [link](https://github.com/paketo-buildpacks/packit)).

2. Create a `package.toml` in the project's root directory configure the build image
   ```
    echo "[package.toml]\nuri = "./datadog-trace.tgz" \n\n[platform]\nos = "linux"" > package.toml
   ```

3. Create the buildpack image using `pack`, referencing `package.toml`
   ```
    pack package-buildpack \
                "datadog/datadog-trace:1.0.0" \
                --config ./package.toml \
                --format image
   ```

The result would create docker image is visible with `docker image ls`.

### Packaging Java applications into an OCI image using buildpack
For `pack` to know that the user wants attach the Datadog trace agent to their application, a [binding](https://paketo.io/docs/buildpacks/configuration/#bindings) needs to be provided along with Datadog buildpack image was created prior. 

1. Create a binding using the following command
   ```
   mkdir binding
   echo 'DatadogTrace' > binding/type
   ```
   Within the directory `./binding`, additional files can be added to configure the agent (see [README](./README.md) for more details).
   
2. Build the image
   ```
    pack build <CONTAINER_NAME> \
    --buildpack paketo-buildpacks/java \
    --buildpack datadog/datadog-trace:1.0.0 \
    --volume "<PATH>/<TO>/binding:/platform/bindings/DatadogTrace"
   ```

When the running the docker image using `docker run <CONTAINER_NAME>`, the container shell should display logs from the Datadog tracer, indicating that the agent is correctly attached.

## Testing
Run `go test ./...` to test run all unit tests for each of the directories.

## Publishing
The buildpack is published by creating a release. When a release is created, a Github workflow will run in the background to package and publish the repo as a Docker image under Github Packages (see [publish.yml](.github/workflows/publish.yml)).

### Tagging Convention
The naming convention for tags should be period seperated numbers (e.g. `0.1.0`, instead of `v0.1.0`). Otherwise, the packager, `jam`, will not be able to parse the tag value.