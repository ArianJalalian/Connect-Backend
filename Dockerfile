FROM golang:latest AS builder

WORKDIR /app

# make a change to binaray, so it can be run on alpine container
ENV CGO_ENABLED=0 

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./
COPY pkg ./pkg
RUN mkdir -p storage/temp




FROM alpine:3.18.3 AS runner
WORKDIR /app
RUN chmod o+w /app

RUN adduser --system --uid 1001 application
RUN addgroup --system --gid 1001 application
USER application

WORKDIR /app


RUN touch .env
RUN mkdir -p storage/tmp

CMD ["./my-app", "serve"]
