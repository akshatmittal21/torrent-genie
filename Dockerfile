FROM golang:alpine AS builder
WORKDIR /go/src
COPY . .
RUN go build -o torrent-genie .

FROM alpine
WORKDIR /torrent-genie
COPY --from=builder /go/src/torrent-genie /torrent-genie/
EXPOSE 3600
CMD ["./torrent-genie"]