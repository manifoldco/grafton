VERSION?=$(shell git describe --tags --dirty | sed 's/^v//')
PKG=github.com/manifoldco/grafton

LD_FLAGS=-w -X $(PKG)/config.Version=$(VERSION)

GO_BUILD=CGO_ENABLED=0 packr build -i --ldflags="$(LD_FLAGS)"

MANIFOLD_VERSION=0.15.1
PROMULGATE_VERSION=0.0.9

rwildcard=$(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) \
	$(filter $(subst *,%,$2),$d))

HAS_GO_MOD=$(shell go help mod; echo $$?)
LINTERS=$(shell grep "// lint" tools.go | awk '{gsub(/\"/, "", $$1); print $$1}' | awk -F / '{print $$NF}') \
	gofmt \
	vet

all: ci
ci: $(LINTERS) cover build

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

ifeq ($(HAS_GO_MOD),0)
bootstrap:

vendor: go.sum
	GO111MODULE=on go mod vendor
else
bootstrap:
	which dep || go get -u github.com/golang/dep/cmd/dep

vendor: Gopkg.lock
	dep ensure -vendor-only
endif

mod-update:
	GO111MODULE=on go get -u -m
	GO111MODULE=on go mod tidy
	dep ensure -update

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

$(LINTERS): %: vendor/bin/gometalinter %-bin vendor generated
	PATH=`pwd`/vendor/bin:$$PATH gometalinter --tests --disable-all --vendor \
	     --deadline=5m -s data --skip generated --enable $@

.PHONY: cover $(LINTERS) $(COVER_TEST_PKGS:=-cover)

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
