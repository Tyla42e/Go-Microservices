# syntax=docker/dockerfile:1

##
## Build the application from source
##

FROM golang:1.21 AS build-stage

WORKDIR /app

COPY ./utils/ ./utils
COPY ./eventtypes ./eventtypes
COPY ./services/comments ./services/comments

RUN cd /app/services/comments && go mod download

RUN ls -ls 
RUN CGO_ENABLED=0 GOOS=linux go build -C ./services/comments .

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

COPY --from=build-stage /app/services/comments .

EXPOSE 4001

#RUN ls -ls 
#USER nonroot:nonroot

#RUN chmod +x ./posts
ENTRYPOINT ["./comments"]