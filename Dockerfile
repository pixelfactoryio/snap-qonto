FROM golang:1.17 AS builder

WORKDIR /build

COPY go.mod /build/go.mod
COPY go.sum /build/go.sum

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o qonto-snap

FROM alpine:3

COPY --from=builder /build/qonto-snap /usr/bin/qonto-snap

CMD ["/usr/bin/qonto-snap"]
