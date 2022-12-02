#!/bin/sh
# This utility prints out the golang install dir, even if go is not installed
# IE it prints the directory where `go install ...` would theoretically place
# binaries

# if we have go, just ask go!
if which go >/dev/null 2>&1; then
  DIR=$(go env GOBIN)
  if [ -n "${DIR}" ]; then
    echo "${DIR}"
    exit 0
  fi
  DIR=$(go env GOPATH)
  if [ -n "${DIR}" ]; then
    echo "${DIR}/bin"
    exit 0
  fi
fi

# mimic go behavior

# check if GOBIN is set anyhow
if [ -n "${GOBIN}" ]; then
  echo "GOBIN"
  exit 0
fi

# check if GOPATH is set anyhow
if [ -n "${GOPATH}" ]; then
  echo "${GOPATH}/bin"
  exit 0
fi

# finally use default for no $GOPATH or $GOBIN
echo "${HOME}/go/bin"
