FROM golang:1.13-alpine as build
ARG RELEASE=TRUE
WORKDIR /var/myrunes
ADD . .

RUN apk add git
RUN go mod download
RUN go build \
        -v -o /app/myrunes -ldflags "\
            -X github.com/myrunes/backend/internal/static.Release=${RELEASE} \
            -X github.com/myrunes/backend/internal/static.AppVersion=$(git describe --tags --abbrev=0)+$(git describe --tags | sed -n 's/^[0-9]\+\.[0-9]\+\.[0-9]\+-\([0-9]\+\)-.*$/\1/p')" \
        ./cmd/server/*.go

# ----------------------------------------------------------

FROM alpine:latest AS final
LABEL maintainer="zekro <contact@zekro.de>"
WORKDIR /app

RUN apk add ca-certificates
COPY --from=build /app .

EXPOSE 8080
RUN mkdir -p /etc/myrunes
CMD ["/app/myrunes", "-c", "/etc/myrunes/config.yml"]
