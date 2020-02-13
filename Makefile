VERSION := $(shell git describe --tags 2>/dev/null || git rev-parse --abbrev-ref HEAD)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
REPO := $(shell cat go.mod | head -n 1 | awk '{print $$2}')

# Image URL to use all building/pushing image targets
IMG ?= kubevent:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:crdVersions=v1,trivialVersions=true"
# Use linker flags to provide version/build settings
LDFLAGS = -ldflags "-X=$(REPO)/cmd/kubevent/version.Version=$(VERSION) -X=$(REPO)/cmd/kubevent/version.Build=$(BUILD)"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: build

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

# Build manager binary
build: generate fmt vet
	go build $(LDFLAGS) -o bin ./cmd/...

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths="./..."

# Build the container image
container-build:
ifeq (, $(shell which podman))
	docker build . -t ${IMG}
else
	podman build . -t ${IMG}
endif

# Push the container image.
# If pushing to remote registry like docker.io, need to have username before the original image like <USERNAME>/<IMG> and has been authorized by login command.
container-push:
ifeq (, $(shell which podman))
	docker push $(IMG)
else
	podman push ${PROJECTNAME} docker-daemon:${IMG}
	podman push ${PROJECTNAME} docker://docker.io/${IMG}
endif

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.4 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
