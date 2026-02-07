ENV ?= localhf
ENVFILE := $(ENV).env
GO_TAGS := -tags netgo
STATIC ?= 1
BIN_DIR := bin

GO_SRC := actions services devices
GO_FILES := $(shell find $(GO_SRC) -type f -name '*.go') main.go go.mod
BINARY := $(BIN_DIR)/turnkey

ifeq ($(STATIC), 1)
	KINDLEBT_LIB = -l:libkindlebt.a -lX11 -lace_bt -ldl -lpthread -llogc
else
	KINDLEBT_LIB = -lkindlebt -lX11
endif

GLOBALS := .globals.mk

$(GLOBALS): $(ENVFILE)
	@echo "# Auto-generated globals" > $@
	@while read line; do \
		case $$line in \
			*=*) \
				var=$${line%%=*} && \
				val=$${line#*=} && \
				val=$${val#\"} && \
				val=$${val%\"} && \
				val=$${val#\'} && \
				val=$${val%\'} && \
				val=$$(eval echo "$$val") && \
				abs=$$(realpath $$val) && \
				echo "$$var := $$abs" >> $@; \
				;; \
		esac; \
	done < $<

-include $(GLOBALS)

define GO_ENV
	export GOOS=linux && \
	export GOARCH=arm && \
	export GOARM=7 && \
	export CGO_ENABLED=1 && \
	export CC="$(CC_BIN)" && \
	export CGO_CFLAGS="-g -Wall -I. -I$(SYSROOT)/usr/include/ -I$(KINDLEBT_PATH)/include/ -I$(KINDLEBT_PATH)/subprojects/darkroot/acs_rtos_reference_ffs/ace/sdk/include/" && \
	export CGO_LDFLAGS="-L$(KINDLEBT_TARGET_PATH)/ $(KINDLEBT_LIB)"
endef

.PHONY: all
all: $(BINARY) build_daemons

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

$(GO_FILES): | $(BIN_DIR)

$(BINARY): $(GO_FILES)
	@$(GO_ENV) && go build $(GO_TAGS) -o $(BINARY) main.go

.PHONY: build
build: $(BINARY)

.PHONY: build_daemon
build_daemon:
	@$(GO_ENV) && \
	cd daemons && \
	go build $(GO_TAGS) -o $(realpath $(BIN_DIR))/$(GO_TARGET)_daemon ./$(GO_TARGET)

.PHONY: build_daemons
build_daemons:
	@$(MAKE) BIN_DIR=$(realpath $(BIN_DIR)) --no-print-directory -C daemons all

.PHONY: build_tests
build_tests:
	@$(GO_ENV) && go test $(GO_TAGS) -c -o $(BIN_DIR)/ ./...

.PHONY: vet
vet:
	@$(GO_ENV) && go vet $(GO_TAGS) ./...

.PHONY: clean
clean:
	rm -rf ./$(BIN_DIR)
	@$(MAKE) --no-print-directory -C daemons clean
