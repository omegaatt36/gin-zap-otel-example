gin-zap-otel
===

## Dependency

### Jaeger

docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.6

# Quick Start

1. start service
```shell
go run main.go
```
2. send a request to service
```shell
curl http://localhost:8000
```
3. open browser and go to http://localhost:16686, see the trace result.