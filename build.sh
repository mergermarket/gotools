#!/bin/bash -e

# This script can be run locally, as well as in Docker, but requires golangci-lint installed locally.
# See https://golangci-lint.run/welcome/install on how to install golang-lint, or look in ./Dockerfile

set -e
set -u

source "./build-functions.sh"

lintCode
fmtCode
buildCode
runTests
