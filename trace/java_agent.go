/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package trace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/bindings"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type JavaAgent struct {
	BuildpackPath    string
	Context          libcnb.BuildContext
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

const (
	yamlExt        = ".yaml"
	propertiesExt  = ".properties"
	javaToolEnv    = "JAVA_TOOL_OPTIONS"
	jmxConfigEnv   = "DD_JMXFETCH_CONFIG"
	agentConfigEnv = "DD_TRACE_CONFIG"
)

func NewJavaAgent(buildpackPath string, dependency libpak.BuildpackDependency, cache libpak.DependencyCache,
	plan *libcnb.BuildpackPlan, context libcnb.BuildContext) JavaAgent {

	return JavaAgent{
		Context:          context,
		BuildpackPath:    buildpackPath,
		LayerContributor: libpak.NewDependencyLayerContributor(dependency, cache, plan),
	}
}

func (j JavaAgent) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Copying to agent to %s", layer.Path)

		//get the java agent and copy it into the image
		file := filepath.Join(layer.Path, filepath.Base(artifact.Name()))
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", artifact.Name(), file, err)
		}

		binding, ok, _ := bindings.ResolveOne(j.Context.Platform.Bindings, bindings.OfType("DatadogTrace"))

		//handle the bindings of `DatadogTrace` which agent configurations
		if ok {
			err := handleAgentProperties(binding, layer, j.Logger)
			if err != nil {
				return libcnb.Layer{}, fmt.Errorf("Unable to process datadog trace agent configuration from binding\n%w", err)
			}
		}

		layer.LaunchEnvironment.Appendf(javaToolEnv, " ", "-javaagent:"+file)

		return layer, nil
	}, libpak.LaunchLayer)
}

func (j JavaAgent) Name() string {
	return j.LayerContributor.LayerName()
}

//adds agent configuration if configuration files exists in the agent
func handleAgentProperties(binding libcnb.Binding, layer libcnb.Layer, logger bard.Logger) error {

	propertiesFileFound := false
	jmxFileFound := false

	err := filepath.Walk(binding.Path, func(path string, info os.FileInfo, err error) error {

		paths := strings.Split(path, "/")
		filename := paths[len(paths)-1]
		toPath := filepath.Join(layer.Path, filename)

		if strings.HasSuffix(path, propertiesExt) {
			logger.Info("Java agent configuration file %s detected", path)
			if propertiesFileFound {
				return fmt.Errorf("Unable to resolve binding configuration: More than 1 .properties supplied")
			}

			err := copyFile(path, toPath)
			if err != nil {
				return fmt.Errorf("Unable to resolve binding configuration\n%w", err)
			}
			layer.LaunchEnvironment.Appendf(agentConfigEnv, " ", toPath)
			propertiesFileFound = true

		} else if strings.HasSuffix(path, yamlExt) {
			logger.Info("JMX configuration file %s detected", path)
			if jmxFileFound {
				return fmt.Errorf("Unable to resolve binding configuration: More than 1 .yaml supplied")
			}

			err := copyFile(path, toPath)
			if err != nil {
				return fmt.Errorf("Unable to resolve binding configuration\n%w", err)
			}
			layer.LaunchEnvironment.Appendf(jmxConfigEnv, " ", toPath)
			jmxFileFound = true
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failed to resolve binding\n%w", err)
	}

	return nil
}

//copies a file into the built image
func copyFile(from string, to string) error {
	in, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("unable to open %s\n%w", from, err)
	}
	defer in.Close()

	if err := sherpa.CopyFile(in, to); err != nil {
		return fmt.Errorf("unable to copy %s to %s\n%w", in.Name(), to, err)
	}
	return nil
}
