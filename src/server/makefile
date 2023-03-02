PROTOC = protoc
GO_OUT = --go_out=. --go_opt=paths=source_relative
GRPC_OUT = --go-grpc_out=. --go-grpc_opt=paths=source_relative

.PHONY: all
all: generate

.PHONY: generate
generate:
ifeq ($(PROTO_FILE),)
	$(foreach file,$(wildcard **/*.proto),$(PROTOC) $(GO_OUT) $(GRPC_OUT) $(file);)
else
	$(PROTOC) $(GO_OUT) $(GRPC_OUT) $(PROTO_FILE)
endif

.PHONY: clean
clean:
	rm -f ./**/*.pb.go

.PHONY: test
test:
	go test ./...
