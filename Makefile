GO=go

.PHONY: all
all: clean build

.PHONY: build
build:
	@echo "build yutang service start >>>"
	@rm -rf go.sum
	${GO} get -u golang.org/x/lint/golint
	@echo "start lint >>>"
	@golint ./...
	@echo "lint completed >>>"
	${GO} get -u github.com/swaggo/swag/cmd/swag
	${GO} mod tidy
	@swag init --generalInfo=cmd/api/service.go --output=api/swagger-spec/api
	@echo "build mode: $(mode), bin name: $(o)"
	CGO_ENABLED=0 ${GO} build -o $(o) -ldflags "-X 'word/pkg/app.GinMode=$(mode)' -s -w" -tags doc cmd/main.go
	@echo ">>> build yutang service complete"

.PHONY: clean
clean:
	@echo "clean start >>>"
	@rm -rf bin/
	@echo ">>> clean complete"