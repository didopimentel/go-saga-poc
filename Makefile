PG_ADDR ?= 'postgres://ps_user:ps_password@localhost:7002/go-saga-poc?sslmode=disable'



.PHONY: setup
setup:
	@echo "==> Setup: installing tools"
	go install github.com/bufbuild/buf/cmd/buf@v0.40.0
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-check-breaking@v0.40.0
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-check-lint@v0.40.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.3.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.3.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install github.com/golang/protobuf/protoc-gen-go@v1.5.1
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/matryer/moq@latest
	go install github.com/rakyll/gotest@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.14.1


################################################################################
## gRPC and protobuf
################################################################################

# generate go files from proto files using buf (http://buf.build)
.PHONY: generate

generate:
	@echo "==> Generate: creating files on 'protogen' dir"
	@command -v buf >/dev/null 2>&1 || { echo >&2 "Setup requires buf but it's not installed.  Aborting."; exit 1; }
	rm -rf protogen/orders
	buf generate --path ./proto/orders/api
	rm -rf protogen/payments
	buf generate --path ./proto/payments/api
	rm -rf protogen/delivery
	buf generate --path ./proto/delivery/api

################################################################################
## Migration and persistence make targets
################################################################################

.PHONY: migrations/up migrations/down migrations/down/yes db/dump-schema migrations/create

migrations/up:
	@command -v migrate >/dev/null 2>&1 || { echo >&2 "Setup requires migrate but it's not installed.  Aborting."; exit 1; }
	migrate -source="file:$$PWD/gateways/persistence/migrations" -database $(PG_ADDR) up

migrations/down:
	@command -v migrate >/dev/null 2>&1 || { echo >&2 "Setup requires migrate but it's not installed.  Aborting."; exit 1; }
	migrate -source="file:$$PWD/gateways/persistence/migrations" -database $(PG_ADDR) down
