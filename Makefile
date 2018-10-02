VERSION?=$(shell git describe --tags --dirty | sed 's/^v//')
GO_BUILD=CGO_ENABLED=0 packr build -i --ldflags="-w"

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

CMD_PKGS=\
	github.com/golang/lint/golint \
	honnef.co/go/tools/cmd/gosimple \
	github.com/client9/misspell/cmd/misspell \
	github.com/gordonklaus/ineffassign \
	github.com/tsenart/deadcode \
	github.com/alecthomas/gometalinter \
	github.com/go-swagger/go-swagger/cmd/swagger

define VENDOR_BIN_TMPL
vendor/bin/$(notdir $(1)): vendor
	go build -o $$@ ./vendor/$(1)
VENDOR_BINS += vendor/bin/$(notdir $(1))
vendor/$(1): Gopkg.lock
	dep ensure -vendor-only
endef

$(foreach cmd_pkg,$(CMD_PKGS),$(eval $(call VENDOR_BIN_TMPL,$(cmd_pkg))))
$(patsubst %,%-bin,$(filter-out gofmt vet,$(LINTERS))): %-bin: vendor/bin/%
gofmt-bin vet-bin:

bootstrap:
	which dep || go get github.com/golang/dep/cmd/dep

vendor: Gopkg.lock
	dep ensure

.PHONY: bootstrap $(CMD_PKGS)

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
$(NO_WINDOWS:%=build/grafton_$(VERSION)_%.zip): build/grafton_$(VERSION)_%.zip: os-build/%/bin/grafton
	cd build/$*/bin; zip -r ../../grafton_$(VERSION)_$*.zip grafton

zips: $(OS_ARCH:%=build/grafton_$(VERSION)_%.zip)

.PHONY: zips $(OS_ARCH:%=os-build/%/bin/grafton)

#################################################
# Cleaning
#################################################

clean:
	rm -rf bin/grafton
	rm -rf build
