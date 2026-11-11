ARG GO_VERSION=1.27

# BUILD STAGE
FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src

# Cache dependencies separately from source
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o /out/goddns .

# RUNTIME STAGE
# distroless/static: CA certificates included, runs as non-root (uid 65532)
FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=build /out/goddns /usr/local/bin/goddns

USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/goddns"]
