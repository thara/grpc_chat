# Tool versoins
PROTOC_VERSION = 3.17.3
PROTOC_GEN_DOC_VERSION = 1.4.1
PROTOC_GEN_GO_VERSION = 1.26.0
PROTOC_GEN_GO_GRPC_VERSION = 1.1.0

.PHONY: protoc
protoc:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative chat/*.proto

.PHONY: fmt
fmt:
	@go fmt ./...
	@clang-format -i chat/chat.proto

.PHONY: runclient
runclient:
	@go run ./client

.PHONY: runserver
runserver:
	@go run ./server
