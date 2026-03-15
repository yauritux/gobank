FROM golang:1.25.8 AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./gobank

FROM gcr.io/distroless/base-debian12

COPY --from=builder /build/gobank /app/gobank
CMD ["/app/gobank"]