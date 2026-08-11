FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Static, multi-arch-friendly build. Set TARGETARCH via buildx.
ARG TARGETARCH
RUN CGO_ENABLED=0 GOARCH=${TARGETARCH} go build -o /out/webhook ./cmd/webhook

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/webhook /webhook
USER nonroot:nonroot
ENTRYPOINT ["/webhook"]
