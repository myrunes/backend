#!/bin/bash

# install dep dependnecy manager
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | bash

# install node.js & npm
curl -sL https://deb.nodesource.com/setup_12.x | bash -
apt-get install -y nodejs

# install vue CLI
npm i -g @vue/cli

# ensuring backend server dependencies
dep ensure -v

# building backend
go build \
    -v -o ./bin/myrunes -ldflags "\
        -X github.com/zekroTJA/myrunes/internal/static.Release=TRUE" \
    ./cmd/server/*.go

cd ./web

# ensuring frontend dependencies
npm i 

# building web assets
npm run build