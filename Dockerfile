FROM golang:1.17.2-alpine3.14 as builder
WORKDIR /build
COPY . .
RUN go build -o supermarket-api
FROM alpine:3.14.2
COPY --from=builder /build/supermarket-api /dist/supermarket-api
EXPOSE 6620
CMD ["/dist/supermarket-api"]
