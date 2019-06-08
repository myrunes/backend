FROM golang:1.11

LABEL maintainer="zekro <contact@zekro.de>"

#### PREPARINGS #####

# adding go binaries to path
ENV PATH="$GOPATH/bin:${PATH}"

# install node.js
RUN curl -sL https://deb.nodesource.com/setup_12.x | bash - &&\
    apt-get install -y nodejs

# install vue CLI
RUN npm i -g @vue/cli

# set workdir to go application dir
WORKDIR $GOPATH/src/github.com/zekroTJA/myrunes

# add nessecary repository files
ADD ./cmd ./cmd
ADD ./internal ./internal
ADD ./pkg ./pkg
ADD ./web ./web
ADD ./Gopkg.toml .

# install dep
RUN go get -u github.com/golang/dep/cmd/dep

# ensure go dependencies via dep
RUN dep ensure -v

# install web dependencies
RUN cd web &&\
    npm i

RUN mkdir ./bin

#### BUILD BACK END ####

RUN go build \
        -v -o ./bin/myrunes -ldflags "\
            -X github.com/zekroTJA/myrunes/internal/static.Release=TRUE" \
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
