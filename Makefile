PROTO_DIR = proto 
OUT_DIR = .

PROTOC_GEN_GO = protoc-gen-go
PROTOC_GEN_GO_GRPC = protoc-gen-go-grpc 
PROTOC_FLAGS = --proto_path=$(PROTO_DIR) --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR)

PROTOS = ./proto/auth/auth.proto 

all: generate

generate:
	@echo "Generating gRPC code"
	@for file in $(PROTOS); do \
		protoc $(PROTOC_FLAGS) $$file; \
	done
	@echo "complete	"

clean:
	@echo "🧹 Cleaning generated files..."
	@find . -name "*.pb.go" -delete
	@echo "✅ Cleaned."

.PHONY: all generate clean