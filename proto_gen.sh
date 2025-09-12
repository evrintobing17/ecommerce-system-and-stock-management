#!/bin/bash

# Directory where .proto files are located
PROTO_DIR="./shared/proto"

# Directory where generated files will be placed
OUT_DIR="./shared/gen"

# Create output directory if it doesn't exist
mkdir -p $OUT_DIR

# Generate Go code for all services
echo "Generating Go code from proto files..."
protoc --go_out=$OUT_DIR --go_opt=paths=source_relative \
    --go-grpc_out=$OUT_DIR --go-grpc_opt=paths=source_relative \
    -I=$PROTO_DIR $PROTO_DIR/*.proto

echo "Proto files compiled successfully!"