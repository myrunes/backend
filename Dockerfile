FROM golang:1.13

LABEL maintainer="zekro <contact@zekro.de>"

#### PREPARINGS #####

# install node.js
RUN curl -sL https://deb.nodesource.com/setup_13.x | bash - &&\
        apt-get install -y nodejs

# install vue CLI
RUN npm i -g @vue/cli

# set workdir to go application dir
WORKDIR /var/myrunes

# add nessecary repository files
ADD . . 

# ensure dependencies with go mod
RUN go mod tidy

# install web dependencies
RUN cd web &&\
    npm i

RUN mkdir ./bin

#### BUILD BACK END ####

RUN go build \
        -v -o ./bin/myrunes -ldflags "\
            -X github.com/zekroTJA/myrunes/internal/static.Release=TRUE \
            -X github.com/zekroTJA/myrunes/internal/static.AppVersion=$(git describe --tags)" \
        ./cmd/server/*.go

#### BUILD FRONT END ####

RUN cd ./web &&\
    npm run build &&\
    cd .. &&\
    mkdir -p ./bin/web &&\
    cp -r ./web/dist ./bin/web/dist


#### EXPOSE AND RUN SETTINGS ####

EXPOSE 8080

RUN mkdir -p /etc/myrunes

CMD ["./bin/myrunes", "-c", "/etc/myrunes/config.yml"]
