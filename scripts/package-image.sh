#!/usr/bin/env bash

# check dependencies
if ! command -v pack &> /dev/null
then
    echo "command 'pack' could not be found"
fi
if ! command -v docker &> /dev/null
then
    echo "command 'docker' could not be found"
    exit
fi
if ! command -v go &> /dev/null
then
    echo "command 'go' could not be found"
    exit
fi

version="1.0.0"
packageName="datadog-trace.tgz"

case $(uname -s) in
    Linux*)     packageBinName=jam-linux;;
    Darwin*)    packageBinName=jam-darwin;;
    *) exit "this script only works on mac and linux"
esac

rm -rf build &> /dev/null
mkdir build

wget -O "./build/${packageBinName}" "https://github.com/paketo-buildpacks/packit/releases/latest/download/${packageBinName}"
chmod a+x "./build/${packageBinName}"

./build/jam-darwin pack --buildpack ./buildpack.toml --version ${version} --output "./build/${packageName}"

echo "[buildpack]
uri = './${packageName}'

[platform]
os = 'linux'" > ./build/package.toml
pack package-buildpack "datadog/datadog-trace:${version}" --config ./build/package.toml --format image
