FROM golang:alpine AS builder
WORKDIR /server
COPY . .
RUN go build -o server ./app/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /server/server .
CMD [ "./server" ]