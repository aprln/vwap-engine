# ======
# Setup
# ======
setup-env:
	cp -i ./.env.example ./.env

# ======
# Build
# ======
build-local:
	go build -race -o bin/vwap main.go

build-docker:
	docker build . --rm -t april/vwap-engine

# ====
# Run
# ====
run-local: build-local
	./bin/vwap

run-docker: build-docker
	docker run -it april/vwap-engine bin/vwap

# =====
# Test
# =====
test-local:
	go test -race ./...

test-docker: build-docker
	docker run --rm -it april/vwap-engine go test -race /vwap-engine/...

# ==============
# Test coverage
# ==============
coverage-local:
	go test -race -covermode=atomic -coverprofile cover.out -coverpkg=./... ./...
	go tool cover -func cover.out

coverage-docker: build-docker
	docker run --rm -it april/vwap-engine go test -race -covermode=atomic -coverprofile cover.out -coverpkg=./... ./...
	docker run --rm -it april/vwap-engine go tool cover -func cover.out

# =====
# Lint
# =====
lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run -v