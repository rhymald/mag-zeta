FROM golang:latest AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify
COPY ./ ./
RUN GCO_ENABLED=0 GOOS=linux go build -o mag

FROM ubuntu:mantic 
COPY --from=builder /app/mag /mag
EXPOSE 4917 9093
ENTRYPOINT ["./mag"]
