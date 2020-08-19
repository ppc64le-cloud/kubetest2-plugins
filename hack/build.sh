#!/bin/sh

set -ex

MODE="${MODE:-release}"
TAGS="${TAGS:-}"
OUTPUT="${OUTPUT:-bin/kubetest2-tf}"
export CGO_ENABLED=0

case "${MODE}" in
release)
	LDFLAGS="${LDFLAGS} -s -w"
	TAGS="${TAGS} release"
	if test "${SKIP_GENERATION}" != y
	then
		go generate -tags "${TAGS}" ./data
	fi
	;;
dev)
	;;
*)
	echo "unrecognized mode: ${MODE}" >&2
	exit 1
esac

go build -tags "${TAGS}" -o "${OUTPUT}" ./kubetest2-tf
