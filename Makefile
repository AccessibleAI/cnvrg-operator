# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Set default version
# $(shell echo $$(git fetch --tags && git tag -l --sort -version:refname | head -n 1)-dirty-$$(git rev-parse --short HEAD) > /tmp/newVersion)

all: manager

# Unfocus tests
unfocus:
	ginkgo unfocus controllers/

bundle:
	kustomize build config/manifests | operator-sdk generate bundle -q --overwrite --version 4.3.16

# Run tests
test: generate fmt vet manifests
	rm -f ./controllers/test-report.html ./controllers/junit.xml
	CNVRG_OPERATOR_MAX_CONCURRENT_RECONCILES=1 go test ./controllers/ -v -timeout 40m

docker:
	docker buildx build --platform=linux/amd64 --load -t test .
test-report:
	docker run -v $$(pwd)/controllers:/tmp cnvrg/xunit-viewer xunit-viewer -r /tmp/junit.xml -o /tmp/test-report.html

override-release: current-version docker-build docker-push chart
	git tag -d $$(cat /tmp/newVersion)
	git push origin -d $$(cat /tmp/newVersion)
	git tag $$(cat /tmp/newVersion)
	git push origin $$(cat /tmp/newVersion)

rc-release: current-version docker-build docker-push chart

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
manager: generate fmt vet
	go build -ldflags "-X 'main.BuildVersion=$(shell cat /tmp/newVersion)'" -mod=readonly -o bin/cnvrg-operator main.go pkged.go
	$(shell if [ $$(echo $$(cat /tmp/newVersion) | grep dirty | wc -l) -eq "0" ]; then git tag $$(cat /tmp/newVersion); fi)



current-version:
	{ \
	set -e ;\
	currentBranch=$$(git rev-parse --abbrev-ref HEAD) ;\
	currentVersion=$$(git fetch --tags && git tag -l --sort -version:refname | head -n 1) ;\
 	if [[ $$currentBranch =~ .*"rc".* ]]; then currentVersion=$$currentBranch; fi ;\
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
	$(CONTROLLER_GEN) crd paths=./... output:artifacts:config=config/crd/bases
	cp config/crd/bases/* pkg/app/controlplane/tmpl/crds
	cp config/crd/bases/* charts/cnvrg-operator/crds
	sed 's/controller-gen$:/mlops.cnvrg.io\/default-loader: "true"\n    &/' pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgapps.yaml > tmp && cat tmp > pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgapps.yaml && rm tmp
	sed 's/controller-gen$:/mlops.cnvrg.io\/own: "false"\n    &/' pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgapps.yaml > tmp && cat tmp > pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgapps.yaml && rm tmp
	sed 's/controller-gen$:/mlops.cnvrg.io\/updatable: "true"\n    &/' pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgapps.yaml > tmp && cat tmp > pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgapps.yaml && rm tmp
	sed 's/controller-gen$:/mlops.cnvrg.io\/default-loader: "true"\n    &/' pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgthirdparties.yaml > tmp && cat tmp > pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgthirdparties.yaml && rm tmp
	sed 's/controller-gen$:/mlops.cnvrg.io\/own: "false"\n    &/' pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgthirdparties.yaml > tmp && cat tmp > pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgthirdparties.yaml && rm tmp
	sed 's/controller-gen$:/mlops.cnvrg.io\/updatable: "true"\n    &/' pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgthirdparties.yaml > tmp && cat tmp > pkg/app/controlplane/tmpl/crds/mlops.cnvrg.io_cnvrgthirdparties.yaml && rm tmp

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
docker-build: generate manifests
		docker build . -t docker.io/cnvrg/cnvrg-operator:$(shell cat /tmp/newVersion)

# Push the docker image
docker-push:
	docker push docker.io/cnvrg/cnvrg-operator:$(shell cat /tmp/newVersion)

build-copctl:
	go build -o bin/copctl cmd/copctl/*.go

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.9.2 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif




