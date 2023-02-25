FROM golang:1.17-alpine3.15 as builder
RUN apk add git openssh-client make curl bash


# COPY only the dep files for efficient caching

WORKDIR /go/src/github.com/lyft/flinkk8soperator



# COPY the rest of the source code
COPY . /go/src/github.com/lyft/flinkk8soperator/

# This 'linux_compile' target should compile binaries to the /artifacts directory
# The main entrypoint should be compiled to /artifacts/flinkk8soperator
RUN make linux_compile

# update the PATH to include the /artifacts directory
ENV PATH="/artifacts:${PATH}"

# This will eventually move to centurylink/ca-certs:latest for minimum possible image size
FROM alpine:3.15
RUN apk add --no-cache curl ca-certificates && update-ca-certificates
COPY --from=builder /artifacts /bin
CMD ["flinkoperator"]
