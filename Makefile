.PHONY: protoc
protoc:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative chat/chat.proto

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
