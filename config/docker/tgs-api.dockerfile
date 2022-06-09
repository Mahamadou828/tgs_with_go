#Build the go binary of the service/api
FROM golang:1.18 as tgs_api
#Disable CGO to assure that the binary is not bind to anything
ENV CGO_ENABLED 0
#Backed to the main build variable
ARG BUILD_REF
ARG ENV
# Create app directory and use it as the working directory
RUN mkdir -p /service
#Copy the source code into the container
COPY . /service

#Build the service binary
WORKDIR /service/app/service/api
RUN go build -ldflags "-X main.build=${BUILD_REF}" -ldflags "-X main.env=${ENV}"


# Run the Go Binary in Alpine.
FROM alpine:3.15
ARG BUILD_DATE
ARG BUILD_REF
COPY --from=tgs_api --chown=sales:sales /service/app/service/api/api /service/api
WORKDIR /service
CMD ["./api"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="tgs-api" \
      org.opencontainers.image.authors="Mahamadou Samake <samaketech@gmail.com>" \
      org.opencontainers.image.source="https://github.com/Mahamadou828/tgs_with_go/app/api" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="Samaketech"