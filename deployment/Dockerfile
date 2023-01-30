FROM golang:1.17-alpine as builder
RUN apk add git

WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /main ./cmd/simplestforum/main.go


FROM alpine:3.15.2
ENV HTTP_PORT=8080
COPY --from=builder main /bin/main
COPY ./migrations ./migrations
EXPOSE 8080
ENTRYPOINT ["/bin/main"]