#!/bin/sh

set -ex

MODE="${MODE:-release}"
TAGS="${TAGS:-}"
OUTPUT="${OUTPUT:-bin/kubetest2-tf}"
UPDATE_SUBMODULE="${UPDATE_SUBMODULE:-false}"
export CGO_ENABLED=0

git submodule update --init
if [[ "${UPDATE_SUBMODULE}" == "true" ]];then
	pushd data/data/k8s-ansible
	git checkout master
	git pull
	popd
fi

case "${MODE}" in
release)
	LDFLAGS="${LDFLAGS} -s -w"
	TAGS="${TAGS} release"
	if test "${SKIP_GENERATION}" != y
	then
		go generate ./data
	fi
	;;
dev)
  TAGS="dev"
	;;
*)
	echo "unrecognized mode: ${MODE}" >&2
	exit 1
esac

go build -tags "${TAGS}" -o "${OUTPUT}" ./kubetest2-tf
