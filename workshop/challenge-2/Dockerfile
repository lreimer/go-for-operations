FROM golang:1.15.2 as builder

WORKDIR /build

COPY . /build
RUN go build -o k8s-logship-sidecar -ldflags="-s -w"

FROM gcr.io/distroless/static-debian10
COPY --from=builder /build/k8s-logship-sidecar /

ENV LOG_DIRECTORY=/logs
ENV LOG_FILE_PATTERN=.+.gz
ENV LOG_SCAN_INTERVAL=10

ENTRYPOINT ["/k8s-logship-sidecar"]
CMD [""]
