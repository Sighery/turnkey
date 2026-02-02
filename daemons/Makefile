PROTO_SRC := proto
PROTOGEN_OUT := protogen

PROTO_FILES := $(shell find $(PROTO_SRC) -name '*.proto')
PROTOGEN_FILES := $(patsubst $(PROTO_SRC)/%.proto,$(PROTOGEN_OUT)/%.pb.go,$(PROTO_FILES))

DAEMONS := kindlebt
BIN_DIR_REL := ../bin
BIN_DIR := $(realpath $(BIN_DIR_REL))
BINARIES := $(addsuffix _daemon,$(addprefix $(BIN_DIR)/,$(DAEMONS)))

.PHONY: all
all: protos $(BINARIES)

.PHONY: protos
protos: $(PROTOGEN_FILES)

$(PROTOGEN_OUT)/%.pb.go: $(PROTO_SRC)/%.proto
	@protoc \
		--proto_path=$(PROTO_SRC) \
		--go_out=$(PROTOGEN_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTOGEN_OUT) --go-grpc_opt=paths=source_relative \
		$<

$(BIN_DIR)/%_daemon: %/*.go $(PROTOGEN_FILES)
	@echo "Building daemon $*"
	@mkdir -p $(BIN_DIR)
	@$(MAKE) --no-print-directory -C .. GO_TARGET=$* build_daemon

.PHONY: clean
clean:
	rm -rf ./$(PROTOGEN_OUT)/*
