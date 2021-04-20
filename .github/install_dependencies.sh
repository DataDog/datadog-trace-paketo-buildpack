#!/usr/bin/env bash
echo "Installing jam ${JAM_VERSION}"

mkdir -p "${HOME}"/bin

echo "${HOME}/bin" >> "${GITHUB_PATH}"
curl \
    --location \
    --show-error \
    --silent \
    "https://github.com/paketo-buildpacks/packit/releases/download/v${JAM_VERSION}/jam-linux" \
    --output ${HOME}/bin/jam

chmod a+x "${HOME}"/bin/jam

echo "Installing pack ${PACK_VERSION}"
curl \
    --location \
    --show-error \
    --silent \
    "https://github.com/buildpacks/pack/releases/download/v${PACK_VERSION}/pack-v${PACK_VERSION}-linux.tgz" \
| tar -C "${HOME}"/bin -xz pack