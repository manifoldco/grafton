VERSION?=$(shell git describe --tags --dirty | sed 's/^v//')
PKG=github.com/manifoldco/grafton

LD_FLAGS=-w -X $(PKG)/config.Version=$(VERSION)

GO_BUILD=CGO_ENABLED=0 packr build -i --ldflags="$(LD_FLAGS)"

MANIFOLD_VERSION=0.15.1
PROMULGATE_VERSION=0.0.9

rwildcard=$(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) \
	$(filter $(subst *,%,$2),$d))

all: ci
ci: lint cover build

.PHONY: all ci

#################################################
# Bootstrapping for base golang package and tool deps
#################################################

CMD_PKGS=$(shell grep '	"' tools.go | awk -F '"' '{print $$2}')

define VENDOR_BIN_TMPL
vendor/bin/$(notdir $(1)): vendor/$(1) | vendor
	go build -a -o $$@ ./vendor/$(1)
VENDOR_BINS += vendor/bin/$(notdir $(1))
vendor/$(1): vendor
endef

$(foreach cmd_pkg,$(CMD_PKGS),$(eval $(call VENDOR_BIN_TMPL,$(cmd_pkg))))

$(patsubst %,%-bin,$(filter-out gofmt vet,$(LINTERS))): %-bin: vendor/bin/%
gofmt-bin vet-bin:

vendor: go.sum
	GO111MODULE=on go mod vendor

mod-update:
	GO111MODULE=on go get -u -m
	GO111MODULE=on go mod tidy

mod-tidy:
	GO111MODULE=on go mod tidy

.PHONY: $(CMD_PKGS)
.PHONY: mod-update mod-tidy

#################################################
# Test and linting
#################################################

test: vendor generated
	@CGO_ENABLED=0 go test -v $$(go list ./... | grep -v generated)

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

.golangci.gen.yml: .golangci.yml
	$(shell awk '/enable:/{y=1;next} y == 0 {print}' $< > $@)

LINTERS=$(filter-out megacheck,$(shell awk '/enable:/{y=1;next} y != 0 {print $$2}' .golangci.yml))

lint: vendor/bin/golangci-lint vendor .golangci.gen.yml 
	$< run -c .golangci.gen.yml $(LINTERS:%=-E %) ./...
	$< run -c .golangci.gen.yml -E megacheck ./...
#  Run imports separately because it can cause golangci to hang when it encounters build issues
#  holding up all other tests, it is also quite long, so the separation allows us to get speedier
#  validation through earlier completion of the lighter weight tests when running manually.
#  As of writing this I don't believe it is taking long enough to justify a separate travis job
	$< run -c .golangci.gen.yml -E goimports ./...

.PHONY: lint cover $(COVER_TEST_PKGS:=-cover)

#################################################
# Code generation
#################################################

generated/%/client: specs/%.yaml vendor/bin/swagger
	rm -rf generated/$*/
	vendor/bin/swagger generate client -f $< -t generated/$*
	touch generated/$*/client
	touch generated/$*/models

APIS=$(patsubst specs/%.yaml,%,$(wildcard specs/*.yaml))
API_CLIENTS=$(APIS:%=generated/%/client)
generated-clients: $(API_CLIENTS)
generated: generated-clients

.PHONY: generated-clients

#################################################
# Building
#################################################

PREFIX?=

SUFFIX=
ifeq ($(GOOS),windows)
SUFFIX=.exe
endif

build: $(PREFIX)bin/grafton

GRAFTON_DEPS=\
	vendor \
	$(wildcard *.go) \
	$(call rwildcard,acceptance,*.go) \
	$(call rwildcard,cmd,*.go) \
	$(call rwildcard,connector,*.go) \
	generated/provider/client \
	generated/provider/models

$(PREFIX)bin/grafton$(SUFFIX): $(GRAFTON_DEPS)
	go get github.com/gobuffalo/packr/...
	$(GO_BUILD) -o $(PREFIX)bin/grafton$(SUFFIX) ./cmd; packr clean

.PHONY: build

#################################################
# Test against sample provider
#################################################

sample-provider:
	GO111MODULE=on go get -u github.com/manifoldco/go-sample-provider/cmd/server
	GO111MODULE=on go build -i -o bin/sample-provider github.com/manifoldco/go-sample-provider/cmd/server
	./bin/grafton generate
	./bin/sample-provider --test --grafton-path="./bin/grafton"
	GO111MODULE=on go mod tidy

.PHONY: sample-provider

#################################################
# Releasing
#################################################

NO_WINDOWS= \
	darwin_amd64 \
	linux_amd64
OS_ARCH= \
	$(NO_WINDOWS) \
	windows_amd64

os=$(word 1,$(subst _, ,$1))
arch=$(word 2,$(subst _, ,$1))

os-build/windows_amd64/bin/grafton: os-build/%/bin/grafton:
	PREFIX=build/$*/ GOOS=$(call os,$*) GOARCH=$(call arch,$*) make build/$*/bin/grafton.exe
$(NO_WINDOWS:%=os-build/%/bin/grafton): os-build/%/bin/grafton:
	PREFIX=build/$*/ GOOS=$(call os,$*) GOARCH=$(call arch,$*) make build/$*/bin/grafton

build/grafton_$(VERSION)_windows_amd64.zip: build/grafton_$(VERSION)_%.zip: os-build/%/bin/grafton
	cd build/$*/bin; zip -r ../../grafton_$(VERSION)_$*.zip grafton.exe
$(NO_WINDOWS:%=build/grafton_$(VERSION)_%.tar.gz): build/grafton_$(VERSION)_%.tar.gz: os-build/%/bin/grafton
	cd build/$*/bin; tar -czf ../../grafton_$(VERSION)_$*.tar.gz grafton

zips: $(NO_WINDOWS:%=build/grafton_$(VERSION)_%.tar.gz) build/grafton_$(VERSION)_windows_amd64.zip

release: zips
	curl -LO https://releases.manifold.co/manifold-cli/$(MANIFOLD_VERSION)/manifold-cli_$(MANIFOLD_VERSION)_linux_amd64.tar.gz
	tar xvf manifold-cli_*
	curl -LO https://releases.manifold.co/promulgate/$(PROMULGATE_VERSION)/promulgate_$(PROMULGATE_VERSION)_linux_amd64.tar.gz
	tar xvf promulgate_*
	./manifold run -t manifold -p promulgate -- ./promulgate release v$(VERSION)

.PHONY: release zips $(OS_ARCH:%=os-build/%/bin/grafton)

#################################################
# Cleaning
#################################################

clean:
	rm -rf bin/grafton
	rm -rf build
