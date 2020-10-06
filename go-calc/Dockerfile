FROM golang:1.15.2 as builder

WORKDIR /build

COPY . /build
RUN make build

FROM gcr.io/distroless/static-debian10
COPY --from=builder /build/go-calc /

ENTRYPOINT ["/go-calc"]
CMD ["version"]
