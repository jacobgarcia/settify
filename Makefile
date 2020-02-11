# Copyright 2016 The Kubernetes Authors.
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
#

ENV ?= DEV
PORT ?= 5000

# The binary to build (just the basename).
BIN := settify

# This repo's root import path (under GOPATH).
PKG := github.com/jacobgarcia/$(BIN)

ARCH ?= amd64
GOOS ?= linux

BASEIMAGE ?= alpine

# This version-strategy uses git tags to set the version string
VERSION := $(shell date +%Y%m%d%H%M%S)
RELEASE := $(shell git describe --tags --always)

PWD ?= $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# Information added to when building containers
GIT_BRANCH=$(shell git name-rev --name-only HEAD | sed "s/~.*//")
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_CREATOR=$(shell git log --format=format:%ae | head -n 1)

build: bin/$(ARCH)/$(BIN)

bin/$(ARCH)/$(BIN): build-dirs
	GOOS=$(GOOS)        \
	ARCH=$(ARCH)        \
	PKG=$(PKG)          \
	BIN=$(BIN)          \
	./scripts/build.sh

version:
	@echo $(VERSION)

tests:
	@if [ ! -d vendor ]; then $(MAKE) --no-print-directory update-vendors; fi
	@./scripts/test.sh

build-dirs:
	@mkdir -p bin/$(ARCH)
	@mkdir -p .go/src/$(PKG) .go/pkg .go/bin .go/std/$(ARCH)
	@if [ ! -d vendor ]; then $(MAKE) --no-print-directory update-vendors; fi

update-vendors:
	@dep ensure

clean: stop container-clean bin-clean files-clean

bin-clean:
	@rm -rf .go bin

files-clean:
	@rm -fr cpu-*.log mem-*.log block-*.log *.test $(BIN).log

build: build-dirs
	@go build -o bin/$(ARCH)/$(BIN) cmd/$(BIN)/main.go

GOLDEN_PKG ?= github.com/jacobgarcia/settify
update-golden-files:
	@go test $(GOLDEN_PKG) -update

up:
	@go run cmd/$(BIN)/main.go
	
include docker.mk
