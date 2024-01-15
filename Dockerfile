# syntax=docker/dockerfile:1

##
## Build the application from source
##

FROM golang:1.21 AS build-stage

WORKDIR /app

#COPY go.mod go.sum ./


COPY ./ ./
RUN cd /app/server
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -C /app/server  -o ./blog

##
## Run the tests in the container
##

FROM build-stage AS run-test-stage
RUN go test -v ./...

##
## Deploy the application binary into a lean image
##

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/server/blog ./blog

EXPOSE 8000 

USER nonroot:nonroot

ENTRYPOINT ["./blog"]