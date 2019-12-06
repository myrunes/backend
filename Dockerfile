FROM golang:1.13 as build
LABEL maintainer="zekro <contact@zekro.de>"
ARG RELEASE=TRUE
RUN curl -sL https://deb.nodesource.com/setup_13.x | bash - &&\
        apt-get install -y nodejs
RUN npm i -g @vue/cli
WORKDIR /var/myrunes
ADD . . 
RUN go mod tidy
RUN cd web &&\
    npm ci
RUN mkdir ./bin
RUN go build \
        -v -o /app/myrunes -ldflags "\
            -X github.com/zekroTJA/myrunes/internal/static.Release=${RELEASE} \
            -X github.com/zekroTJA/myrunes/internal/static.AppVersion=$(git describe --tags)" \
        ./cmd/server/*.go
RUN cd ./web &&\
    npm run build &&\
    cd .. &&\
    mkdir -p /app/web &&\
    cp -r ./web/dist /app/web/dist


FROM debian:stretch-slim AS final
WORKDIR /app
COPY --from=build /app .

EXPOSE 8080
RUN mkdir -p /etc/myrunes
CMD ["/app/myrunes", "-c", "/etc/myrunes/config.yml", "-assets", "/app/web/dist"]
