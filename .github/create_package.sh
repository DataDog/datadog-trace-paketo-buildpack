#!/usr/bin/env bash
mkdir ${HOME}/build

echo "Package repo using jam"
jam pack --buildpack "${GITHUB_WORKSPACE}/buildpack.toml" --version ${BUILDPACK_VERSION} --output "${HOME}/build/${PACKAGE_NAME}"

echo "Building docker image using packaged tgz"
echo "[buildpack]
uri = '${HOME}/build/${PACKAGE_NAME}'
[platform]
os = 'linux'" > ${HOME}/build/package.toml
pack package-buildpack "${BUILDPACK_NAME}:${BUILDPACK_VERSION}" --config ${HOME}/build/package.toml --format image