FROM golang:latest AS builder
ADD . /prometheus-generator/
WORKDIR /prometheus-generator
RUN go get -d -v github.com/prometheus/client_golang/prometheus
RUN CGO_ENABLED=0 GOOS=linux go build -o prometheus-generator -a -tags netgo -ldflags '-w'

FROM busybox
EXPOSE 8080
COPY --from=builder /prometheus-generator/prometheus-generator ./
CMD ["/prometheus-generator"]
