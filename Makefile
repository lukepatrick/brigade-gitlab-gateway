# The Docker registry where images are pushed.
# Note that if you use an org (like on Quay and DockerHub), you should
# include that: quay.io/foo
DOCKER_REGISTRY    ?= lukepatrick
DOCKER_BUILD_FLAGS :=
LDFLAGS            :=

BINS        = brigade-gitlab-gateway
IMAGES      = brigade-gitlab-gateway
DOCKER_BINS = brigade-gitlab-gateway

GIT_TAG   = $(shell git describe --tags --always 2>/dev/null)
VERSION   ?= ${GIT_TAG}
IMAGE_TAG ?= ${VERSION}
LDFLAGS   += -X github.com/lukepatrick/brigade-gitlab-gateway/pkg/version.Version=$(VERSION)

# Build native binaries
.PHONY: build
build: $(BINS)

.PHONY: $(BINS)
$(BINS):
	go build -ldflags '$(LDFLAGS)' -o bin/$@ ./$@/cmd/$@

# Cross-compile for Docker+Linux
build-docker-bins: $(addsuffix -docker-bin,$(DOCKER_BINS))

%-docker-bin:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o ./$*/rootfs/$* ./$*/cmd/$*

# To use docker-build, you need to have Docker installed and configured. You should also set
# DOCKER_REGISTRY to your own personal registry if you are not pushing to the official upstream.
.PHONY: docker-build
docker-build: build-docker-bins
docker-build: $(addsuffix -image,$(IMAGES))

%-image:
	docker build $(DOCKER_BUILD_FLAGS) -t $(DOCKER_REGISTRY)/$*:$(IMAGE_TAG) $*

# You must be logged into DOCKER_REGISTRY before you can push.
.PHONY: docker-push
docker-push: $(addsuffix -push,$(IMAGES))

%-push:
	docker push $(DOCKER_REGISTRY)/$*:$(IMAGE_TAG)

.PRECIOUS: build-chart
.PHONY: build-chart
build-chart:
	helm package -d docs/ ./charts/brigade-gitlab-gateway
	helm repo index docs/

.PHONY: format
format:
	test -z "$$(find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -d {} + | tee /dev/stderr)" || \
	test -z "$$(find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -w {} + | tee /dev/stderr)"

HAS_GOMETALINTER := $(shell command -v gometalinter;)
HAS_DEP          := $(shell command -v dep;)
HAS_GIT          := $(shell command -v git;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_GIT
	$(error You must install git)
endif
ifndef HAS_GOMETALINTER
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
endif
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
ifndef HAS_BINDATA
	go get github.com/jteeuwen/go-bindata/...
endif
	dep ensure
