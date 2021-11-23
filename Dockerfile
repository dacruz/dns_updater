FROM golang:1.17.3 as builder

RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 go build -o dns_updater .

FROM alpine:3.14.3

RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/dns_updater /app/
CMD ["/app/dns_updater"]
