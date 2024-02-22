# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY internal/ internal/
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:latest

WORKDIR /

COPY --from=build-stage /server /server

ENTRYPOINT ["/server"]
