# Makefile

# Variables
PROTO_DIR = proto
PROTO_FILES = $(PROTO_DIR)/softphone.proto
PROTO_OUT_DIR = .
GO_OUT_DIR = internal
DOCKER_IMAGE_NAME = api-gateway
DOCKER_CONTAINER_NAME = api-gateway-container
CONFIG_FILE = internal/config/config.json

# Extract the server port from config.json
SERVER_PORT = $(shell jq -r '.server_port' $(CONFIG_FILE))

# Targets
.PHONY: all build proto test clean docker-build docker-run run

all: build

build: proto
	@echo "Building the application..."
	go build -o bin/api-gateway ./cmd/...

proto:
	@echo "Generating protobuf files..."
	protoc --go_out=$(PROTO_OUT_DIR) --go-grpc_out=$(PROTO_OUT_DIR) $(PROTO_FILES)

docker-build: build
	@echo "Building the Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME) .

docker-run: 
	@echo "Running the Docker container..."
	docker run -p $(SERVER_PORT):8080 --name $(DOCKER_CONTAINER_NAME) $(DOCKER_IMAGE_NAME)

run: build
	@echo "Running the application locally..."
	if [ ! -f $(CONFIG_FILE) ]; then echo "Config file missing"; exit 1; fi
	./bin/api-gateway

generate-mocks:
	@echo "Generating mocks..."
	mockgen -destination=test/integration/mocks/mock_softphoneserviceclient.go -package=mocks api-gateway/proto SoftPhoneServiceClient
	mockgen -source=internal/workermanager/manager.go -destination=test/mocks/mock_workermanager.go -package=mocks
	mockgen -source=internal/grpcclient/client.go -destination=test/mocks/mock_grpcclient.go -package=mocks

test-unit: generate-mocks
	@echo "Running unit tests..."
	go test ./test/unit/...

test-integration: generate-mocks
	@echo "Running integration tests..."
	go test -tags=integration ./test/integration/...

clean:
	@echo "Cleaning up..."
	rm -rf bin/*
	rm -rf $(PROTO_OUT_DIR)/*.pb.go
