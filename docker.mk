# Docker targets - dependent on basic makefile
#

REGISTRY ?= hub.docker.com/adrianforsius

IMAGE := $(REGISTRY)/go-service

BUILD_IMAGE ?= golang:1.10-alpine

docker-bin/$(ARCH)/$(BIN): build-dirs
	@echo "building: $@"
	@docker run                                                            \
		-ti                                                                \
		--rm                                                               \
		-u $$(id -u):$$(id -g)                                             \
		-v $(PWD)/.go:/go                                                 \
		-v $(PWD):/go/src/$(PKG)                                          \
		-v $(PWD)/bin/$(ARCH):/go/bin                                     \
		-v $(PWD)/bin/$(ARCH):/go/bin/linux_$(ARCH)                       \
		-v $(PWD)/.go/std/$(ARCH):/usr/local/go/pkg/linux_$(ARCH)_static  \
		-w /go/src/$(PKG)                                                  \
		$(BUILD_IMAGE)                                                     \
		/bin/sh -c "                                                       \
			GOOS=$(GOOS)                                                   \
			ARCH=$(ARCH)                                                   \
			PKG=$(PKG)                                                     \
			BIN=$(BIN)                                                     \
			./scripts/build.sh                                             \
		"

docker-tests: build-dirs
	@echo "Preparing environment for tests: "
	@docker run                                                            \
		-ti                                                                \
		--rm                                                               \
		-u $$(id -u):$$(id -g)                                             \
		-v $(PWD)/.go:/go                                                 \
		-v $(PWD):/go/src/$(PKG)                                          \
		-v $(PWD)/bin/$(ARCH):/go/bin                                     \
		-v $(PWD)/.go/std/$(ARCH):/usr/local/go/pkg/linux_$(ARCH)_static  \
		-w /go/src/$(PKG)                                                  \
		$(BUILD_IMAGE)                                                     \
		/bin/sh -c "                                                       \
			./scripts/test.sh                                              \
		"

docker-build:
	@$(MAKE) VERSION=latest container
	@docker-compose -f docker-compose.yml build

docker-up: docker-build
	@docker-compose -f docker-compose.yml up

docker-down:
	@docker-compose -f docker-compose.yml down

docker-restart: docker-down docker-up

docker-push: docker-login
	docker push $(IMAGE):$(VERSION)
	docker images -q $(IMAGE):$(VERSION)
	@echo "pushed: $(IMAGE):$(VERSION)"

docker-login:
	@docker login                              \
		--username="$$DOCKER_REGISTRY_USER"     \
		--password="$$DOCKER_REGISTRY_API_KEY"  \
		$$DOCKER_REGISTRY_URL

DOTFILE_IMAGE = $(subst :,_,$(subst /,_,$(IMAGE))-$(VERSION))

container: .container-$(DOTFILE_IMAGE)
.container-$(DOTFILE_IMAGE): docker-bin/$(ARCH)/$(BIN) Dockerfile.in
	@echo "Putting together container $(IMAGE):$(VERSION)"
	@sed \
		-e 's|ARG_BIN|$(BIN)|g' \
		-e 's|ARG_ARCH|$(ARCH)|g' \
		-e 's|ARG_FROM|$(BASEIMAGE)|g' \
		Dockerfile.in > .dockerfile-$(ARCH)
	@docker build -t $(IMAGE):$(VERSION) -f .dockerfile-$(ARCH) .
	@docker images -q $(IMAGE):$(VERSION) > $@

container-clean:
	rm -rf .container-* .dockerfile-* .push-*
