VERSION?=$(shell git describe --tags --dirty | sed 's/^v//')
GO_BUILD=CGO_ENABLED=0 go build -i --ldflags="-w"

rwildcard=$(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) \
	$(filter $(subst *,%,$2),$d))

LINTERS=\
	gofmt \
	golint \
	vet \
	misspell \
	ineffassign \
	deadcode

all: ci
ci: $(LINTERS) cover build

.PHONY: all ci

#################################################
# Bootstrapping for base golang package deps
#################################################

BOOTSTRAP=\
	github.com/golang/lint/golint \
	honnef.co/go/simple/cmd/gosimple \
	github.com/client9/misspell/cmd/misspell \
	github.com/gordonklaus/ineffassign \
	github.com/tsenart/deadcode \
	github.com/alecthomas/gometalinter \
	github.com/go-swagger/go-swagger/cmd/swagger

$(BOOTSTRAP):
	go get -u $@
bootstrap: $(BOOTSTRAP)
	glide -v || curl http://glide.sh/get | sh

vendor: glide.lock
	glide install

.PHONY: bootstrap $(BOOTSTRAP)

#################################################
# Test and linting
#################################################

test: vendor generated
	@CGO_ENABLED=0 go test -v $$(glide nv | grep -v generated)

comma:= ,
empty:=
space:= $(empty) $(empty)

COVER_TEST_PKGS:=$(shell find . -type f -name '*_test.go' | grep -v vendor | rev | cut -d "/" -f 2- | rev | grep -v generated | sort -u)
$(COVER_TEST_PKGS:=-cover): %-cover: all-cover.txt
	@CGO_ENABLED=0 go test -coverprofile=$@.out -covermode=atomic ./$*
	@if [ -f $@.out ]; then \
		grep -v "mode: atomic" < $@.out >> all-cover.txt; \
		rm $@.out; \
	fi

all-cover.txt:
	echo "mode: atomic" > all-cover.txt

cover: vendor generated all-cover.txt $(COVER_TEST_PKGS:=-cover)

METALINT=gometalinter --tests --disable-all --vendor --deadline=5m -s data \
	$$(glide nv | grep -v generated) --enable

$(LINTERS): vendor generated
	$(METALINT) $@

.PHONY: cover $(LINTERS) $(COVER_TEST_PKGS:=-cover)

#################################################
# Code generation
#################################################

generated/provider/client generated/provider/models: provider.yaml
	swagger generate client -f $< -A provider -t generated/provider
	touch generated/provider/client
	touch generated/provider/models

generated: generated/provider/client generated/provider/models

.PHONY: generated

#################################################
# Building
#################################################

build: bin/grafton

GRAFTON_DEPS=\
	vendor \
	$(wildcard *.go) \
	$(call rwildcard,acceptance,*.go) \
	$(call rwildcard,cmd,*.go) \
	$(call rwildcard,connector,*.go) \
	generated/provider/client \
	generated/provider/models

bin/grafton: $(GRAFTON_DEPS)
	$(GO_BUILD) -o bin/grafton ./cmd

.PHONY: build

#################################################
# Cleaning
#################################################

clean:
	rm bin/grafton
