IMAGE        ?= ghcr.io/jchonig/docker-fetchbox
TAG          ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
GOIMAGE      := golang:1.26-alpine
GOIMAGE_TEST := golang:1.26
SRCDIR       := $(CURDIR)/src

.PHONY: all build lint test docker-build image-run-test image-test docker-push clean

all: build

build:
	docker run --rm \
		-v "$(SRCDIR):/src" -w /src \
		$(GOIMAGE) \
		go build -v -o /dev/null ./...

lint:
	docker run --rm \
		-v "$(SRCDIR):/src" -w /src \
		$(GOIMAGE) \
		sh -c 'go vet ./... && test -z "$$(gofmt -l .)"'

test:
	docker run --rm \
		-v "$(SRCDIR):/src" -w /src \
		$(GOIMAGE_TEST) \
		go test -v -race -count=1 ./...

docker-build:
	docker build -t $(IMAGE):$(TAG) .

# Run the smoke test against a pre-built image (used by CI after docker/build-push-action)
image-run-test:
	docker run --rm \
		--entrypoint /usr/local/bin/fetchbox \
		-v "$(CURDIR)/testdata/fetchbox.yml:/config/fetchbox.yml:ro" \
		$(IMAGE):$(TAG) \
		--config /config/fetchbox.yml

# Build the image locally then smoke-test it
image-test: docker-build image-run-test

docker-push: image-test
	docker push $(IMAGE):$(TAG)
	docker tag $(IMAGE):$(TAG) $(IMAGE):latest
	docker push $(IMAGE):latest

clean:
	rm -f src/fetchbox src/docker-fetchbox
