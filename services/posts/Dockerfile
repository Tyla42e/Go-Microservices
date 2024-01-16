# syntax=docker/dockerfile:1

##
## Build the application from source
##

FROM golang:1.21 AS build-stage

WORKDIR /app

#COPY ./ ./

# RUN mkdir ./utils
# RUN mkdir ./eventtypes
# RUN mkdir ./services
# RUN mkdir ./services/posts

RUN ls -ls  ./

COPY ./utils/ ./utils
COPY ./eventtypes ./eventtypes
COPY ./services/posts ./services/posts


# RUN ls -ls  ./
# RUN ls -ls  ./utils
# RUN ls -ls  ./eventtypes
# RUN ls -ls  ./services/posts

RUN cd /app/services/posts && go mod download
# RUN ls -ls 
# RUN cd ./posts
#RUN ls -ls 
#RUN go mod download
RUN ls -ls 
RUN CGO_ENABLED=0 GOOS=linux go build -C ./services/posts .
RUN ls -ls  ./services/posts
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

COPY --from=build-stage /app/services/posts .

EXPOSE 4000

#RUN ls -ls 
#USER nonroot:nonroot

#RUN chmod +x ./posts
ENTRYPOINT ["./posts"]