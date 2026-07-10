#!/usr/bin/env bash
set -e

# Determine OS and Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "${OS}" in
    Linux*)     OS='Linux';;
    Darwin*)    OS='Darwin';;
    *)          echo "Unsupported OS: ${OS}"; exit 1;;
esac

case "${ARCH}" in
    x86_64)     ARCH='x86_64';;
    amd64)      ARCH='x86_64';;
    aarch64)    ARCH='arm64';;
    arm64)      ARCH='arm64';;
    *)          echo "Unsupported architecture: ${ARCH}"; exit 1;;
esac

REPO="mtrace-project/mtrace"
REQUESTED_VERSION="${1:-latest}"

if [ "${REQUESTED_VERSION}" = "latest" ]; then
    echo "Fetching latest release for ${REPO}..."
    LATEST_RELEASE=$(curl -s https://api.github.com/repos/${REPO}/releases/latest)
    VERSION=$(echo "${LATEST_RELEASE}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "${VERSION}" ]; then
        echo "Failed to get the latest release version."
        exit 1
    fi
    echo "Latest version is ${VERSION}"
else
    VERSION="${REQUESTED_VERSION}"
    if [[ "${VERSION}" != v* ]]; then
        VERSION="v${VERSION}"
    fi
    echo "Using requested version: ${VERSION}"
fi

# Determine file extension based on what goreleaser generates
if [ "${OS}" = "Linux" ]; then
    EXT="tar.zst"
else
    EXT="zip"
fi

FILENAME="mtrace_${VERSION#v}_${OS}_${ARCH}.${EXT}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"

echo "Downloading ${FILENAME} from ${DOWNLOAD_URL}..."
curl -sL "${DOWNLOAD_URL}" -o "${FILENAME}"

echo "Extracting..."
if [ "${EXT}" = "zip" ]; then
    unzip -q "${FILENAME}" -d mtrace_tmp
else
    mkdir -p mtrace_tmp
    tar -I zstd -xf "${FILENAME}" -C mtrace_tmp
fi

echo "Installing to /usr/local/bin/mtrace (requires sudo)..."
sudo mv mtrace_tmp/mtrace /usr/local/bin/mtrace
sudo chmod +x /usr/local/bin/mtrace

# Cleanup
rm -rf mtrace_tmp "${FILENAME}"

echo "mtrace installed successfully!"
