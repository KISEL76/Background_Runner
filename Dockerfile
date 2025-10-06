FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . . 
RUN CGO_ENABLED=0 go build -o queue-svc ./cmd

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder app/queue-svc .
EXPOSE 8080
CMD [ "./queue-svc" ]