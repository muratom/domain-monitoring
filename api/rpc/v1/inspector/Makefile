.PHONY: all

gen-http-server:
	oapi-codegen -config ./server.cfg.yaml ../inspector.yaml

gen-http-models:
	oapi-codegen -config ./models.cfg.yaml ../inspector.yaml

gen-http-client:
	oapi-codegen -config ./client.cfg.yaml ../inspector.yaml

gen-http: gen-http-server gen-http-models gen-http-client