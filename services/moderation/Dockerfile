# syntax=docker/dockerfile:1

##
## Build the application from source
##

FROM golang:1.21 AS build-stage

WORKDIR /app

COPY ./utils/ ./utils
COPY ./eventtypes ./eventtypes
COPY ./services/comments ./services/comments
COPY ./services/moderation ./services/moderation

RUN cd /app/services/moderation && go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -C ./services/moderation .

##
## Run the tests in the container
##

# FROM build-stage AS run-test-stage
# RUN go test -v ./...

##
## Deploy the application binary into a lean image
##

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/services/moderation .

EXPOSE 4003
ENTRYPOINT ["./moderation"]