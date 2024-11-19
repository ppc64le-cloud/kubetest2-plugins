# Common uses:
# installing a kubetest2-tf deployer: `make install-deployer-tf INSTALL_DIR=$HOME/go/bin`

# get the repo root and output path
REPO_ROOT:=$(shell pwd)
export REPO_ROOT
OUT_DIR=$(REPO_ROOT)/bin
# record the source commit in the binary, overridable
COMMIT?=$(shell git describe --tags --always 2>/dev/null)
INSTALL?=install
# make install will place binaries here
# the default path attempts to mimic go install
INSTALL_DIR?=$(shell $(REPO_ROOT)/hack/goinstalldir.sh)
# the output binary name, overridden when cross compiling
BINARY_NAME=kubetest2-tf
BINARY_PATH=./kubetest2-tf
BUILD_FLAGS=-trimpath -ldflags="-buildid= -X=github.com/ppc64le-cloud/kubetest2-plugins/kubetest2-tf/deployer.GitTag=$(COMMIT)"
# ==============================================================================

install-deployer-tf:
	git submodule update --init
	go build $(BUILD_FLAGS) -o $(OUT_DIR)/$(BINARY_NAME) $(BINARY_PATH)
	$(INSTALL) -d $(INSTALL_DIR)
	$(INSTALL) $(OUT_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
.PHONY: install-deployer-tf
