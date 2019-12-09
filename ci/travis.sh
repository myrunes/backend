#!/bin/bash

# install node.js & npm
curl -sL https://deb.nodesource.com/setup_13.x | bash -
apt-get install -y nodejs

# install vue CLI
npm i -g @vue/cli

# ensuring backend server dependencies
go mod tidy

# building backend
go build \
    -v -o ./bin/myrunes -ldflags "\
        -X github.com/myrunes/myrunes/internal/static.Release=TRUE" \
    ./cmd/server/*.go

cd ./web

# ensuring frontend dependencies
npm i 

# building web assets
npm run build