GO_TIDY = go mod tidy

setup:
	go install github.com/vektra/mockery/v2@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@latest
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

tidy:
	$(GO_TIDY) && go mod vendor

fmt:
	find . -iname '*.go' -not -path '*/vendor/*' -print0 | xargs -0 gofmt -s -w

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...  && go tool cover -html=coverage.out

lint: fmt
	 golangci-lint run

win-env:
	go env -w GOOS=windows

pre-commit: win-env test coverage lint
