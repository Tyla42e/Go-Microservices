# syntax=docker/dockerfile:1

##
## Build the application from source
##

FROM golang:1.21 AS build-stage

WORKDIR /app

COPY ./utils/ ./utils
COPY ./eventtypes ./eventtypes
COPY ./services/posts ./services/posts

RUN cd /app/services/posts && go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -C ./services/posts .

##
## Deploy the application binary into a lean image
##

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/services/posts .

EXPOSE 4000
ENTRYPOINT ["./posts"]