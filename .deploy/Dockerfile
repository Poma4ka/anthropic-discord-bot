FROM golang:1.21-alpine AS build


WORKDIR /opt/build

COPY . .

RUN go mod tidy && \
    go mod download && \
    go build -ldflags="-s -w" -o ./dist/app ./src/main.go

FROM alpine:latest AS DEPLOY

COPY --from=build /opt/build/dist/app /opt/app

ENV GIN_MODE=release

ENTRYPOINT ["/opt/app"]
