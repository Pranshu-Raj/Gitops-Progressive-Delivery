IMAGE   ?= ghcr.io/pranshu-raj/gitops-progressive-delivery
VERSION ?= $(shell git rev-parse --short HEAD)

.PHONY: test build run image render tf-check

test:
	cd app && go vet ./... && go test -race -count=1 ./...

build:
	cd app && CGO_ENABLED=0 go build -ldflags="-X main.version=$(VERSION)" -o ../bin/app .

run: build
	./bin/app

image:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(VERSION) app

render:
	kubectl kustomize deploy/overlays/dev
	kubectl kustomize deploy/overlays/prod

tf-check:
	terraform -chdir=terraform fmt -check -recursive
	terraform -chdir=terraform validate
