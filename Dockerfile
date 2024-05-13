FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./gobank

FROM gcr.io/distroless/base-debian12

# WORKDIR /app
COPY --from=builder /build/gobank /app/gobank
CMD ["/app/gobank"]