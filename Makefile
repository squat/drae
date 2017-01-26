.PHONY: all push container clean container-name container-latest push-latest 

BIN := drae
PROJECT := github.com/squat/$(BIN)
REGISTRY ?= index.docker.io
IMAGE ?= squat/$(BIN)

TAG := $(shell git describe --abbrev=0 --tags HEAD 2>/dev/null)
COMMIT := $(shell git rev-parse HEAD)
VERSION := $(COMMIT)
ifneq ($(TAG),)
    ifeq ($(COMMIT), $(shell git rev-list -n1 $(TAG)))
        VERSION := $(TAG)
    endif
endif
DIRTY := $(shell test -z "$$(git diff --shortstat 2>/dev/null)" || echo -dirty)
VERSION := $(VERSION)$(DIRTY)
LD_FLAGS := -ldflags \"-X $(PROJECT)/pkg/version.Version=$(VERSION)\"

BUILD_IMAGE ?= golang:1.7-alpine

all: build

build: bin/$(BIN)

bin:
	@mkdir -p bin

bin/$(BIN): bin cmd/$(BIN)/main.go glide.yaml
	@echo "building: $@"
	@docker run \
	    -u $$(id -u):$$(id -g) \
	    -v $$(pwd):/go/src/$(PROJECT) \
	    -v $$(pwd)/bin:/go/bin \
	    -w /go/src/$(PROJECT) \
	    $(BUILD_IMAGE) \
	    /bin/sh -c " \
	        GOOS=linux \
		CGO_ENABLED=0 \
		go build -o bin/$(BIN) \
		    $(LD_FLAGS) \
		    ./cmd/drae/... \
	    "

container: .container-$(VERSION) container-name
.container-$(VERSION): bin/$(BIN) Dockerfile
	@docker build -t $(IMAGE):$(VERSION) .
	@docker images -q $(IMAGE):$(VERSION) > $@

container-latest: .container-$(VERSION)
	@docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	@echo "container: $(IMAGE):latest"

container-name:
	@echo "container: $(IMAGE):$(VERSION)"

push: .push-$(VERSION) push-name
.push-$(VERSION): .container-$(VERSION)
	@docker push $(REGISTRY)/$(IMAGE):$(VERSION)
	@docker images -q $(IMAGE):$(VERSION) > $@

push-latest: container-latest
	@docker push $(REGISTRY)/$(IMAGE):latest
	@echo "pushed: $(IMAGE):latest"

push-name:
	@echo "pushed: $(IMAGE):$(VERSION)"

clean: container-clean bin-clean

container-clean:
	rm -rf .container-* .push-*

bin-clean:
	rm -rf bin
