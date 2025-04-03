FROM golang:1.24.1-alpine3.21 AS build

COPY ./internal/ /app/internal/
COPY ./cmd/api/ /app/cmd/api/
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

RUN tree

RUN apk add --update gcc musl-dev
RUN go env -w CGO_ENABLED=1
RUN go mod download
RUN tree
RUN go build --tags "linux" /app/cmd/api/main.go 

FROM alpine:3.21 AS runtime
COPY --from=build /app/main /app/main

WORKDIR /app

RUN tree

EXPOSE 8080
CMD ["./main"]
