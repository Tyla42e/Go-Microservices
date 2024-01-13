# Use golang and the build env
FROM golang:latest AS build

WORKDIR /app

# download the required Go dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY / ./

# FAIL
RUN go build -o myrepo-test .

##########
# Deploy #
##########

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /myrepo-test /myrepo-test

USER nonroot:nonroot

ENTRYPOINT ["/myrepo-test"]

CMD ["/myrepo-test"]