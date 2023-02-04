ARG GO_VERSION=1.16.6

# ----------------------- #
# Image to create our GO project
# ----------------------- #
# GO version to use
FROM golang:${GO_VERSION}-alpine AS builder
# No proxy to install depedencies
RUN go env -w GOPROXY=direct
# Install git to get dependencies defined
RUN apk add --no-cache git
# Install some security certificates to install application
RUN apk --no-cache add ca-certificates && update-ca-certificates
# Where application is going to be installed
WORKDIR /src
# Copy go.mod and go.sum files into source directory to install dependencies in next step
COPY ./go.mod ./go.sum ./
# Install dependencies copied
RUN go mod download
# Copy everything into source directory
COPY ./ ./
# Command to build application based on code
RUN CGO_ENABLED=0 go build \
		-installsuffix 'static' \
		-o /platzi-rest-ws .
# ----------------------- #
# ----------------------- #

# ----------------------- #
# Image to exceute our project file
# ----------------------- #
FROM scratch AS runner
# Copy certificates from builder image into this image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
# Copy our env file
COPY .env ./
# Copy image from builder to runner
COPY --from=builder /platzi-rest-ws /platzi-rest-ws
# Port to connect
EXPOSE 5050
# Command to be executed when the image is loaded
ENTRYPOINT [ "/platzi-rest-ws" ]
# ----------------------- #
# ----------------------- #
