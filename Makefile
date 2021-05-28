

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Set default version
$(shell echo $$(git fetch --tags && git tag -l --sort -version:refname | head -n 1)-dirty-$$(git rev-parse --short HEAD) > /tmp/newVersion)

all: manager

# Unfocus tests
unfocus:
	ginkgo unfocus controllers/

# Run tests
test: pack generate fmt vet manifests
	go test ./controllers/ -v

pack:
	pkger

override-release: current-version docker-build docker-push chart
	git tag -d $$(cat /tmp/newVersion)
	git push origin -d $$(cat /tmp/newVersion)
	git tag $$(cat /tmp/newVersion)
	git push origin $$(cat /tmp/newVersion)

patch-release: patch-version docker-build docker-push chart
	git tag $$(cat /tmp/newVersion);
	git push origin $$(cat /tmp/newVersion)

minor-release: minor-version docker-build docker-push chart
	git tag $$(cat /tmp/newVersion);
	git push origin $$(cat /tmp/newVersion)

major-release: major-version docker-build docker-push chart
	git tag $$(cat /tmp/newVersion);
	git push origin $$(cat /tmp/newVersion)

# Build manager binary
manager: pack generate fmt vet
	go build -ldflags "-X 'main.BuildVersion=$(shell cat /tmp/newVersion)'" -mod=readonly -o bin/cnvrg-operator main.go pkged.go
	$(shell if [ $$(echo $$(cat /tmp/newVersion) | grep dirty | wc -l) -eq "0" ]; then git tag $$(cat /tmp/newVersion); fi)



current-version:
	{ \
	set -e ;\
	currentVersion=$$(git fetch --tags && git tag -l --sort -version:refname | head -n 1) ;\
	echo $$currentVersion > /tmp/newVersion ;\
    }


patch-version:
	{ \
	set -e ;\
	currentVersion=$$(git fetch --tags && git tag -l --sort -version:refname | head -n 1) ;\
	patchVersion=$$(echo $$currentVersion | tr . " " | awk '{print $$3}') ;\
	patchVersion=$$(( $$patchVersion + 1 )) ;\
	newVersion=$$(echo $$currentVersion | tr . " " | awk -v pv=$$patchVersion '{print $$1"."$$2"."pv}') ;\
	echo $$newVersion > /tmp/newVersion ;\
    }

minor-version:
	{ \
	set -e ;\
	currentVersion=$$(git fetch --tags && git tag -l --sort -version:refname | head -n 1) ;\
	minorVersion=$$(echo $$currentVersion | tr . " " | awk '{print $$2}') ;\
	minorVersion=$$(( $$minorVersion + 1 )) ;\
	newVersion=$$(echo $$currentVersion | tr . " " | awk -v pv=$$minorVersion '{print $$1"."pv"."0}') ;\
	echo $$newVersion > /tmp/newVersion ;\
    }

major-version:
	{ \
	set -e ;\
	currentVersion=$$(git fetch --tags && git tag -l --sort -version:refname | head -n 1) ;\
	majorVersion=$$(echo $$currentVersion | tr . " " | awk '{print $$1}') ;\
	majorVersion=$$(( $$majorVersion + 1 )) ;\
	newVersion=$$(echo $$currentVersion | tr . " " | awk -v pv=$$majorVersion '{print pv".0.0"}') ;\
	echo $$newVersion > /tmp/newVersion ;\
    }

.PHONY: chart
chart:
	{ \
	helm repo add cnvrgv3 https://charts.v3.cnvrg.io ;\
	helm repo update ;\
	rm -fr /tmp/chart ;\
	cp -R chart /tmp/chart ;\
	VERSION=$$(cat /tmp/newVersion) envsubst < chart/Chart.yaml | tee tmp-file && mv tmp-file /tmp/chart/Chart.yaml ;\
	helm package /tmp/chart -d /tmp ;\
	if [ $$(echo $$(cat /tmp/newVersion) | grep dirty | wc -l) -eq "0" ]; then helm push /tmp/chart cnvrgv3 -u=${HELM_USER} -p=${HELM_PASS} --force; fi ;\
	}

chart-delete:
	curl -XDELETE https://$$HELM_USER:$$HELM_PASS@charts.v3.cnvrg.io/api/charts/cnvrg/${V}

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go operator run

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=docker.io/cnvrg/cnvrg-operator:$(shell cat /tmp/newVersion)
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=cnvrg-operator-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
	cp config/crd/bases/* pkg/controlplane/tmpl/crds
	cp config/crd/bases/* chart/crds

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: pack generate manifests
		docker build . -t docker.io/cnvrg/cnvrg-operator:$(shell cat /tmp/newVersion)

# Push the docker image
docker-push:
	docker push docker.io/cnvrg/cnvrg-operator:$(shell cat /tmp/newVersion)

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif