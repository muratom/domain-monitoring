.PHONY: all

gen-grpc:
	protoc -I . \
	--go_out ./emitter --go_opt paths=source_relative \
	--go-grpc_out ./emitter --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./emitter --grpc-gateway_opt paths=source_relative \
	emitter.proto
