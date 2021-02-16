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

package trace_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/datadog-trace/trace"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"
)

func testJavaAgent(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		ctx    libcnb.BuildContext
		dep    = libpak.BuildpackDependency{
			URI:    "https://localhost/dd-java-agent.jar",
			SHA256: "799868a8196959d51f83ea7c5954c7ed6b29069b06286fae203b4d2e5d7ad53a",
		}
		//we are storing v0.70.0 dd trace agent in the cache
		dc = libpak.DependencyCache{CachePath: "testdata/dependencyCache"}
	)

	it.Before(func() {
		var err error

		ctx.Buildpack.Path, err = ioutil.TempDir("", "java-agent-buildpack")
		Expect(err).NotTo(HaveOccurred())

		ctx.Layers.Path, err = ioutil.TempDir("", "java-agent-layers")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Buildpack.Path)).To(Succeed())
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes Java agent as a part of the build pack", func() {

		j := trace.NewJavaAgent(ctx.Buildpack.Path, dep, dc, &libcnb.BuildpackPlan{}, ctx)
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = j.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.Launch).To(BeTrue())
		Expect(filepath.Join(layer.Path, "dd-java-agent.jar")).To(BeARegularFile())

		Expect(layer.LaunchEnvironment["JAVA_TOOL_OPTIONS.delim"]).To(Equal(" "))
		Expect(layer.LaunchEnvironment["JAVA_TOOL_OPTIONS.append"]).To(Equal(fmt.Sprintf("-javaagent:%s",
			filepath.Join(layer.Path, "dd-java-agent.jar"))))
	})

	it("creates correct flags for running agent given DatadogTrace binding", func() {

		binding, err := libcnb.NewBindingFromPath("testdata/binding")
		Expect(err).NotTo(HaveOccurred())

		fmt.Println(binding.Type)

		ctx.Platform.Bindings = libcnb.Bindings{
			binding,
		}

		j := trace.NewJavaAgent(ctx.Buildpack.Path, dep, dc, &libcnb.BuildpackPlan{}, ctx)
		layer, err := ctx.Layers.Layer("test-layer")

		layer, err = j.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())
		Expect(layer.Launch).To(BeTrue())

		Expect(layer.LaunchEnvironment["DD_TRACE_CONFIG.delim"]).To(Equal(" "))
		Expect(layer.LaunchEnvironment["DD_TRACE_CONFIG.append"]).To(Equal(fmt.Sprintf(filepath.Join(layer.Path, "agent.properties"))))

		Expect(layer.LaunchEnvironment["DD_JMXFETCH_CONFIG.delim"]).To(Equal(" "))
		Expect(layer.LaunchEnvironment["DD_JMXFETCH_CONFIG.append"]).To(Equal(fmt.Sprintf(filepath.Join(layer.Path, "conf.yaml"))))
	})
}
