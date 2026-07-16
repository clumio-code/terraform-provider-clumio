#!/bin/bash

VERSION=$1
if [ -z "${VERSION}" ]; then
  echo "Usage: ./installer.sh <version>. An example of <version> would be v0.1.0."
  exit 1
fi
OS=$(uname | awk '{print tolower($0)}')
arch=$(uname -m)
case ${arch} in
  x86_64)
    ARCH="amd64"
    ;;
  i386|i686)
    ARCH="386"
    ;;
  arm64|armv8l|aarch64)
    ARCH="arm64"
    ;;
  armv7l|arm)
    ARCH="arm"
    ;;
  *)
    echo "Not supported architecture. Exiting..."
    exit 1
esac

# download fetches the URL in $1 and writes it to the file path in $2 using
# whichever of curl or wget is available.
download() {
  if command -v curl &> /dev/null; then
    curl -sSL "$1" -o "$2"
  elif command -v wget &> /dev/null; then
    wget -q "$1" -O "$2"
  else
    echo "wget or curl is required to download files. Exiting..."
    exit 1
  fi
}

VERSION_NUMBER=${VERSION:1}
TF_PLUGIN_DIR="${HOME}/.terraform.d/plugins/clumio.com/providers/clumio/${VERSION_NUMBER}/${OS}_${ARCH}"
if ! mkdir -p "${TF_PLUGIN_DIR}"; then
   echo "Error creating directory ${TF_PLUGIN_DIR}. Exiting..."
  exit 1
fi

PROVIDER_NAME=terraform-provider-clumio_${VERSION_NUMBER}_${OS}_${ARCH}
BINARY="https://github.com/clumio-code/terraform-provider-clumio/releases/download/${VERSION}/${PROVIDER_NAME}.zip"

if ! mkdir "${PROVIDER_NAME}"; then
  echo "Error creating directory ${PROVIDER_NAME}. Exiting..."
  exit 1
fi

if ! cd "${PROVIDER_NAME}"; then
  echo "Error changing directory to ${PROVIDER_NAME}. Exiting..."
  exit 1
fi

if ! download "${BINARY}" "${PROVIDER_NAME}.zip"; then
  echo "Error downloading ${BINARY}. Exiting..."
  exit 1
fi

# Verify the downloaded archive against the SHA256SUMS published with the
# release before unzipping/installing it, to detect tampered or corrupted
# downloads.
SHA256SUMS_NAME="terraform-provider-clumio_${VERSION_NUMBER}_SHA256SUMS"
SHA256SUMS_URL="https://github.com/clumio-code/terraform-provider-clumio/releases/download/${VERSION}/${SHA256SUMS_NAME}"
if ! download "${SHA256SUMS_URL}" "${SHA256SUMS_NAME}"; then
  echo "Error downloading ${SHA256SUMS_URL}. Exiting..."
  exit 1
fi

if command -v sha256sum &> /dev/null; then
  SHA256_CMD="sha256sum"
elif command -v shasum &> /dev/null; then
  SHA256_CMD="shasum -a 256"
else
  echo "sha256sum or shasum is required to verify the download. Exiting..."
  exit 1
fi

EXPECTED_SHA=$(awk -v file="${PROVIDER_NAME}.zip" '$2 == file {print $1}' "${SHA256SUMS_NAME}")
if [ -z "${EXPECTED_SHA}" ]; then
  echo "Could not find a checksum for ${PROVIDER_NAME}.zip in ${SHA256SUMS_NAME}. Exiting..."
  exit 1
fi
ACTUAL_SHA=$(${SHA256_CMD} "${PROVIDER_NAME}.zip" | awk '{print $1}')
if [ "${EXPECTED_SHA}" != "${ACTUAL_SHA}" ]; then
  echo "Checksum verification failed for ${PROVIDER_NAME}.zip. Exiting..."
  echo "  expected: ${EXPECTED_SHA}"
  echo "  actual:   ${ACTUAL_SHA}"
  exit 1
fi
echo "Checksum verified for ${PROVIDER_NAME}.zip."

if ! unzip -q "${PROVIDER_NAME}.zip"; then
  echo "Error unzipping ${PROVIDER_NAME}.zip. Exiting..."
  exit 1
fi

if ! cp "terraform-provider-clumio_${VERSION}" "${TF_PLUGIN_DIR}"; then
  echo "Error copying terraform-provider-clumio_${VERSION} to ${TF_PLUGIN_DIR}. Exiting..."
  exit 1
fi

if ! cd ..; then
  echo "cd .. failed. Exiting..."
  exit 1
fi

# cleanup downloaded files
if [ -d "${PROVIDER_NAME}" ]; then
  rm -rf "${PROVIDER_NAME}"
fi

if [ -f "${TF_PLUGIN_DIR}/terraform-provider-clumio_${VERSION}" ]; then
  echo "Clumio Terraform Provider installed successfully."
fi
