.PHONY: all

gen-http-server:
	oapi-codegen -config ./internal/delivery/http/server.cfg.yaml ./api/openapi-spec/api.yaml

gen-http-models:
	oapi-codegen -config ./internal/delivery/http/models.cfg.yaml ./api/openapi-spec/api.yaml

gen-http: gen-http-server gen-http-models

gen-grpc:
	protoc -I ./api/proto \
	--go_out ./api/proto/gen/go --go_opt paths=source_relative \
	--go-grpc_out ./api/proto/gen/go --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./api/proto/gen/go --grpc-gateway_opt paths=source_relative \
	./api/proto/emitter/emitter.proto
