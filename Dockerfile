# syntax=docker/dockerfile:1

ARG PNPM_VERSION=10.33.0
ARG GO_VERSION=1.26.1

## Node
################################################################################
# Use alpine image for base image
FROM node:lts-alpine3.23 AS base-node
WORKDIR /app/web

# Install pnpm.
RUN --mount=type=cache,target=/root/.npm \
    npm install -g pnpm@${PNPM_VERSION}

################################################################################
# Create a stage for installing production dependencies.
FROM base-node AS deps-node

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /root/.local/share/pnpm/store to speed up subsequent builds.
# Leverage bind mounts to package.json and pnpm-lock.yaml to avoid having to copy them
# into this layer.
RUN --mount=type=bind,source=web/package.json,target=package.json \
    --mount=type=bind,source=web/pnpm-lock.yaml,target=pnpm-lock.yaml \
    --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --prod --frozen-lockfile

################################################################################
# Create a stage for building the application.
FROM deps-node AS build-node

# Download additional development dependencies before building, as some projects require
# "devDependencies" to be installed to build. If you don't need this, remove this step.
RUN --mount=type=bind,source=web/package.json,target=package.json \
    --mount=type=bind,source=web/pnpm-lock.yaml,target=pnpm-lock.yaml \
    --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile

# Copy the rest of the source files into the image.
COPY web .

# Run the build script.
RUN pnpm run build


## Go
################################################################################
# Use alpine image as base image
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine3.23 AS base-go
WORKDIR /app

################################################################################
# Create a stage for building the application.
FROM base-go AS build-go
WORKDIR /app

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage bind mounts to go.sum and go.mod to avoid having to copy them into
# the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# This is the architecture you're building for, which is passed in by the builder.
# Placing it here allows the previous steps to be cached across architectures.
ARG TARGETARCH

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a (writable) bind mount to the current directory to avoid having to copy the
# source code into the container. (This mount must be writable because the "build-node" file mount will exist within it.)
# Leverage a bind mount to the files from the "build-node" stage.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,rw,target=. \
    --mount=type=bind,from=build-node,source=/app/web/build/,target=web/build/ \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .

## App
################################################################################
# Create a new stage for running the application that contains the minimal
# runtime dependencies for the application. This often uses a different base
# image from the build stage where the necessary files are copied from the build
# stage.
FROM alpine:latest
WORKDIR /app

RUN addgroup -g 1000 -S appuser && \
    adduser -S appuser -u 1000 -G appuser && \
    chown -R appuser:appuser /app

USER appuser

RUN mkdir .data

# Copy the executable from the "build-go" stage.
COPY --from=build-go /bin/server .

# Expose the port that the application listens on.
EXPOSE 3000

# What the container should run when it is started.
ENTRYPOINT [ "/app/server" ]
