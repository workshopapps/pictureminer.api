# syntax=docker/dockerfile:1
FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /pictureminer
EXPOSE 9000
CMD [ "/pictureminer" ]r scan