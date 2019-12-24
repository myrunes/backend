#!/bin/bash

# ensuring backend server dependencies
go mod tidy

# building backend
go build \
    -v -o ./bin/myrunes -ldflags "\
        -X github.com/myrunes/backend/internal/static.Release=TRUE" \
    ./cmd/server/*.go
