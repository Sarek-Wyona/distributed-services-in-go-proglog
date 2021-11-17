.PHONY: protogen
## proto-generate: Generates golang code from the protos defined for v1
protogen:
	 protoc api/v1/*.proto  --go_out=. --go_opt=paths=source_relative --proto_path=.

.PHONY: grpcgen
grpcgen:
	protoc api/v1/*.proto --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --proto_path=.


.PHONY: test
## test: Runs tests for our code
test:
	go test -race ./...
