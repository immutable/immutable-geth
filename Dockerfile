# Support setting various labels on the final image
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

# Build Geth in a stock Go builder container
FROM golang:1.21-alpine as builder

RUN apk add --no-cache gcc musl-dev linux-headers git

# Get dependencies - will also be cached if we won't change go.mod/go.sum
COPY go.mod /go-ethereum/
COPY go.sum /go-ethereum/
RUN cd /go-ethereum && go mod download
ADD . /go-ethereum
RUN cd /go-ethereum && go run build/ci.go install -static ./cmd/geth

# Pull Geth into a second stage deploy alpine container
FROM alpine:3.18.9

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-ethereum/build/bin/geth /usr/local/bin/

# Default config
RUN mkdir /etc/geth
ADD ./immutable/config/testnet.toml /etc/geth/testnet.toml
ADD ./immutable/config/mainnet.toml /etc/geth/mainnet.toml

# Public config
ADD ./immutable/config/testnet-public.toml /etc/geth/testnet-public.toml
ADD ./immutable/config/mainnet-public.toml /etc/geth/mainnet-public.toml

EXPOSE 8545 8546 30300 30300/udp
ENTRYPOINT ["geth"]

# Add some metadata labels to help programatic image consumption
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

LABEL commit="$COMMIT" version="$VERSION" buildnum="$BUILDNUM"
