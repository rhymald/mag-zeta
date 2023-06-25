FROM golang:latest AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify
COPY ./ ./

RUN ls -la
RUN go build -o mag
RUN ls -la /app/

FROM ubuntu:latest 
COPY --from=builder /app/mag /mag
ENTRYPOINT ["./mag"]
