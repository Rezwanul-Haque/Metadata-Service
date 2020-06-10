# Start from golang base image
FROM golang:1.14-alpine3.11 AS builder

#ENV GO111MODULE=on

# Add Maintainer info
LABEL maintainer="Rezwanul Haque <rezwanul.cse@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
RUN apk update && apk add --no-cache bash

# Configure the repo url
ENV REPO_URL=github.com/Rezwanul-Haque/Metadata-Service

# Setup our $GOPATH
ENV GOPATH=/app

# Set the current working directory inside the container
ENV APP_PATH=$GOPATH/src/$REPO_URL
ENV WORKPATH=$APP_PATH/src/

COPY src $WORKPATH

WORKDIR $WORKPATH

# Copy go mod and sum files
COPY ./src/go.mod .
COPY ./src/go.sum .

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

RUN ls .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o metadata-service .

##
#RUN git clone https://github.com/vishnubob/wait-for-it.git
#CMD ./wait-for-it/wait-for-it.sh --host=mysql --port=3306 --timeout=60 -- ./metadata-service

# Expose port 7001 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./metadata-service"]