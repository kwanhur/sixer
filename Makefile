# Copyright (c) 2022 kwanhur.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# init project path
WORKROOT := $(shell pwd)
OUTDIR   := $(WORKROOT)/output

# init environment variables
export PATH        := $(shell go env GOPATH)/bin:$(PATH)
export GO111MODULE := on

# init command params
GO           := go
GOBUILD      := $(GO) build
GOTEST       := $(GO) test
GOVET        := $(GO) vet
GOGET        := $(GO) get
GOGEN        := $(GO) generate
GOCLEAN      := $(GO) clean
GOINSTALL    := $(GO) install
GOFLAGS      := -race
STATICCHECK  := staticcheck
LICENSEEYE   := license-eye
MDLINT       := markdownlint
PIP          := pip3
PIPINSTALL   := $(PIP) install
NPM          := npm
NPMINSTALL   := $(NPM) install

# init arch
ARCH := $(shell getconf LONG_BIT)
ifeq ($(ARCH),64)
	GOTEST += $(GOFLAGS)
endif

# init git commit id
COMMIT ?= $(shell git rev-parse --short HEAD)

# go install package
# $(1) package name
# $(2) package address
define INSTALL_PKG
	@echo installing $(1)
	$(GOINSTALL) $(2)
	@echo $(1) installed
endef

define PIP_INSTALL_PKG
	@echo installing $(1)
	$(PIPINSTALL) $(1)
	@echo $(1) installed
endef

define NPM_INSTALL_PKG
	@echo installing $(1)
	$(NPMINSTALL) $(1)
	@echo $(1) installed
endef

# make, make all
all: compile

# make compile, go build
compile: build
build:
	$(GOBUILD) -ldflags "-X main.commit=$(COMMIT)" -o sixer

# make deps
deps:
	$(call NPM_INSTALL_PKG, markdownlint-cli)
	$(call PIP_INSTALL_PKG, pre-commit)
	$(call INSTALL_PKG, staticcheck, honnef.co/go/tools/cmd/staticcheck)
	$(call INSTALL_PKG, license-eye, github.com/apache/skywalking-eyes/cmd/license-eye@latest)

# make precommit, enable autoupdate and install with hooks
precommit:
	pre-commit autoupdate
	pre-commit install --install-hooks

# make check
check:
	$(STATICCHECK) ./...

# make license-check, check code file's license declaration
license-check:
	$(LICENSEEYE) header check

# make license-fix, fix code file's license declaration
license-fix:
	$(LICENSEEYE) header fix

# make clean
clean:
	$(GOCLEAN)
	rm -rf $(WORKROOT)/apisixer

# avoid filename conflict and speed up build 
.PHONY: all
