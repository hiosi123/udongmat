FROM golang:1.21.0 AS builder

WORKDIR /udongmat

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o udongmatBin ./cmd

FROM alpine:latest

ENV TZ Asia/Seoul

WORKDIR /udongmat

COPY --from=builder /udongmat/udongmatBin .

CMD ["./udongmatBin"]

EXPOSE 8080
