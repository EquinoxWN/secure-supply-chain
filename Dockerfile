# syntax=docker/dockerfile:1
# Base images are pinned by digest (ADR 0002); Dependabot proposes updates as reviewed PRs.

FROM golang:1.27-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
ARG VERSION=dev
ARG COMMIT=unknown
# Static, stripped, reproducible binary: no cgo, no build paths, fixed build ID.
RUN CGO_ENABLED=0 go build -trimpath -buildvcs=false \
      -ldflags "-s -w -buildid= -X main.version=${VERSION} -X main.commit=${COMMIT}" \
      -o /out/server ./cmd/secure-supply-chain

# distroless static: no shell, no package manager, runs as uid 65532.
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
COPY --from=build /out/server /server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/server"]
