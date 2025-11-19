# set specific Go version to build with
ARG GO_VERSION=1.22
 
# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build
 
RUN apk add --no-cache git
WORKDIR /src
COPY ./ ./

# Build the executable
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -installsuffix 'static' -o /app ./cmd/test-task-manager
 
# STAGE 2: build the container to run
FROM gcr.io/distroless/static AS final
 
LABEL maintainer="btcdirect"
USER nonroot:nonroot
 
# copy compiled app
COPY --from=build --chown=nonroot:nonroot /app /app
 
# run binary; use vector form
ENTRYPOINT ["/app"]
