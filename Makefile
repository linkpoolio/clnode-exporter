.DEFAULT_GOAL := build
.PHONY: dep build install docker dockerpush

REPO=linkpoolio/clnode-exporter
LDFLAGS=-ldflags "-X github.com/linkpoolio/clnode-exporter/store.Sha=`git rev-parse HEAD`"

dep:
	@dep ensure

build: dep
	@go build $(LDFLAGS) -o clnode-exporter

install: dep
	@go install $(LDFLAGS)

docker:
	@docker build . -t $(REPO)

dockerpush:
	@docker push $(REPO)